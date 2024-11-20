package wallet_controller

import (
	"fmt"
	"net/http"

	"app/config"
	"app/models/setting"
	"app/models/user"
	"app/models/wallet_transaction"
	"app/pkg/tlync"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (c *ControllerBasic) Initiate(ctx echo.Context) error {
	ctxUser := c.Utils.CtxUser(ctx)
	model := wallet_transaction.Model{
		ID:       uuid.New(),
		WalletID: ctxUser.ID,
		Type:     wallet_transaction.TypeCredit,
		User: wallet_transaction.WalletUser{
			ID: &ctxUser.ID,
		},
		RechargedBy: wallet_transaction.WalletUser{
			ID: &ctxUser.ID,
		},
	}
	settings := tlync.Settings{}
	if err := c.Models.Setting.GetForTlync(&settings); err != nil {
		return c.APIErr.Database(ctx, err, &setting.Model{})
	}
	tlyncInput := tlync.InitiateInput{
		ID:          settings.StoreID,
		Amount:      fmt.Sprintf("%.2f", model.Amount),
		Phone:       "218910000000",
		Email:       "info@sample.com",
		BackendURL:  config.DOMAIN + "/api/v1/wallet-transactions/auto-confirm",
		FrontendURL: settings.FrontURL,
		CustomRef:   model.ID.String(),
	}
	v, err := c.GetValidator(ctx, model.ModelName())
	if err != nil {
		return err
	}
	var phone string
	email := v.Data.Get("email")

	dummyUser := user.Model{}
	dummyUser.MergePhone(v)
	if dummyUser.Phone != nil {
		phone = *dummyUser.Phone
	}
	if email != "" {
		tlyncInput.Email = email
	}
	if phone != "" {
		tlyncInput.Phone = phone
	}
	if valid := model.MergeAndValidate(v); !valid {
		return c.APIErr.InputValidation(ctx, v)
	}
	tlyncInput.Amount = fmt.Sprintf("%.2f", model.Amount)
	res, err := tlync.InitiatePayment(&settings, &tlyncInput)
	if err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}

	model.TLyncURL = &res.URL
	res.Amount = model.Amount
	// Start transacting
	tx, err := c.Models.DB.Beginx()
	if err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	defer func() { _ = tx.Rollback() }()

	if err := c.Models.Wallet.CreateTransaction(&model, tx); err != nil {
		return c.APIErr.Database(ctx, err, &model)
	}
	if err = tx.Commit(); err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	return ctx.JSON(http.StatusCreated, model)
}

func (c *ControllerBasic) Confirm(ctx echo.Context) error {
	var result wallet_transaction.Model
	if err := c.Utils.ReadUUIDParam(&result.ID, ctx); err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}
	ctxUser := c.Utils.CtxUser(ctx)
	if err := c.Models.Wallet.GetTransaction(&result, &ctxUser.ID); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}
	settings := tlync.Settings{}
	if err := c.Models.Setting.GetForTlync(&settings); err != nil {
		return c.APIErr.Database(ctx, err, &setting.Model{})
	}
	tlyncInput := tlync.ConfirmInput{
		StoreID:   settings.StoreID,
		CustomRef: result.ID.String(),
	}
	res, err := tlync.TransactionReceipt(&settings, &tlyncInput)
	if err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}
	if res.Result != "success" {
		err := fmt.Errorf("confirmation is not complete: %s", res.Message)
		return c.APIErr.BadRequest(ctx, err)
	}
	// Start transacting
	tx, err := c.Models.DB.Beginx()
	if err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	defer func() { _ = tx.Rollback() }()

	// marshalling/unmarshalling res into wallet trx struct
	b, err := json.Marshal(res.Data)
	if err != nil {
		return c.APIErr.BadRequest(
			ctx,
			fmt.Errorf("marshalling res.Data: %w", err),
		)
	}
	if err := json.Unmarshal(b, &result.TLyncResponse); err != nil {
		return c.APIErr.BadRequest(
			ctx,
			fmt.Errorf("unmarshalling res.Data: %w", err),
		)
	}

	result.IsConfirmed = true
	result.PaymentMethod = &res.Data.Gateway
	result.PaymentReference = &res.Data.OrderID

	if err := c.Models.Wallet.UpdateTransaction(
		&result,
		&ctxUser.ID,
		tx,
	); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}
	if err = tx.Commit(); err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	return ctx.JSON(http.StatusCreated, result)
}

func (c *ControllerBasic) AutoConfirm(ctx echo.Context) error {
	var result wallet_transaction.Model
	v, err := c.GetValidator(ctx, result.ModelName())
	if err != nil {
		return err
	}
	c.Utils.Logger.Info().Msgf(
		"received tlync backend-url callback values, encoded: %s",
		v.Data.Values.Encode(),
	)
	id, err := uuid.Parse(v.Data.Get("custom_ref"))
	if err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}
	result.ID = id
	if err := c.Models.Wallet.GetTransaction(&result, nil); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}
	settings := tlync.Settings{}
	if err := c.Models.Setting.GetForTlync(&settings); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}
	tlyncInput := tlync.ConfirmInput{
		StoreID:   settings.StoreID,
		CustomRef: result.ID.String(),
	}
	res, err := tlync.TransactionReceipt(&settings, &tlyncInput)
	if err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}
	if res.Result != "success" {
		err := fmt.Errorf("confirmation is not complete: %s", res.Message)
		return c.APIErr.BadRequest(ctx, err)
	}
	// Start transacting
	tx, err := c.Models.DB.Beginx()
	if err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	defer func() { _ = tx.Rollback() }()

	// marshalling/unmarshalling res into wallet trx struct
	b, err := json.Marshal(res.Data)
	if err != nil {
		return c.APIErr.BadRequest(
			ctx,
			fmt.Errorf("marshalling res.Data: %w", err),
		)
	}
	if err := json.Unmarshal(b, &result.TLyncResponse); err != nil {
		return c.APIErr.BadRequest(
			ctx,
			fmt.Errorf("unmarshalling res.Data: %w", err),
		)
	}

	result.IsConfirmed = true
	result.PaymentMethod = &res.Data.Gateway
	result.PaymentReference = &res.Data.OrderID
	if err := c.Models.Wallet.UpdateTransaction(&result, nil, tx); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}
	if err = tx.Commit(); err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	return ctx.JSON(http.StatusCreated, result)
}
