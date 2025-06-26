package wallet_controller

import (
	"errors"
	"net/http"

	"app/models/user"
	"app/models/wallet_transaction"
	"app/pkg/payment_gateway"

	setting "bitbucket.org/sadeemTechnology/backend-model-setting"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (c *ControllerBasic) TylncInitiate(ctx echo.Context) error {
	ctxUser := c.Utils.CtxUser(ctx)
	model := wallet_transaction.Model{
		ID:       uuid.New(),
		WalletID: ctxUser.ID,
		Type:     wallet_transaction.TypeCredit,
		User:     wallet_transaction.WalletUser{ID: &ctxUser.ID},
	}

	input := payment_gateway.TlyncRequest{
		WalletTransactionID: model.ID,
	}
	settings, err := c.Models.Setting.GetForTyrianAnt()
	if err != nil {
		return c.APIErr.Database(ctx, err, &setting.Model{})
	}

	v, err := c.GetValidator(ctx, model.ModelName())
	if err != nil {
		return err
	}
	var phone string

	dummyUser := user.Model{}
	dummyUser.MergePhone(v)
	if dummyUser.Phone != nil {
		phone = *dummyUser.Phone
	}

	if valid := model.MergeAndValidate(v); !valid {
		return c.APIErr.InputValidation(ctx, v)
	}
	input.Phone = phone
	input.Amount = model.Amount

	res, err := payment_gateway.TlyncInitiatePayment(settings, &input)
	if err != nil {
		return c.APIErr.ExternalRequestError(ctx, err)
	}

	model.TLyncURL = res.TLyncURL

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

func (c *ControllerBasic) TylncConfirm(ctx echo.Context) error {
	var result wallet_transaction.Model

	v, err := c.GetValidator(ctx, result.ModelName())
	if err != nil {
		return err
	}
	v.AssignUUID(
		"transaction_id",
		"id",
		"wallet_transactions",
		&result.ID,
		true,
	)

	if !v.Valid() {
		return c.APIErr.InputValidation(ctx, v)
	}

	ctxUser := c.Utils.CtxUser(ctx)
	if err := c.Models.Wallet.GetTransaction(&result, &ctxUser.ID); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}
	if result.IsConfirmed {
		err := errors.New(v.T.WalletTransactionAlreadyConfirmed())
		return c.APIErr.BadRequest(ctx, err)
	}

	settings, err := c.Models.Setting.GetForTyrianAnt()
	if err != nil {
		return c.APIErr.Database(ctx, err, &setting.Model{})
	}

	res, err := payment_gateway.TlyncTransactionReceipt(settings, result.ID)
	if err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}

	// Start transacting
	tx, err := c.Models.DB.Beginx()
	if err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	defer func() { _ = tx.Rollback() }()

	result.TLyncResponse = res.Response
	result.IsConfirmed = true
	result.PaymentMethod = res.PaymentMethod
	result.PaymentReference = res.PaymentReference

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

func (c *ControllerBasic) TylncAutoConfirm(ctx echo.Context) error {
	var result wallet_transaction.Model

	v, err := c.GetValidator(ctx, result.ModelName())
	if err != nil {
		return err
	}
	c.Utils.Logger.Info().Msgf(
		"received tlync backend-url callback values, encoded: %s",
		v.Data.Values.Encode(),
	)
	id, err := uuid.Parse(v.Data.Get("id"))
	if err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}
	result.ID = id
	if err := c.Models.Wallet.GetTransaction(&result, nil); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}
	settings, err := c.Models.Setting.GetForTyrianAnt()
	if err != nil {
		return c.APIErr.Database(ctx, err, &setting.Model{})
	}

	res, err := payment_gateway.TlyncTransactionReceipt(settings, result.ID)
	if err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}

	// Start transacting
	tx, err := c.Models.DB.Beginx()
	if err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	defer func() { _ = tx.Rollback() }()

	result.TLyncResponse = res.Response
	result.IsConfirmed = true
	result.PaymentMethod = res.PaymentMethod
	result.PaymentReference = res.PaymentReference

	if err := c.Models.Wallet.UpdateTransaction(&result, nil, tx); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}
	if err = tx.Commit(); err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	return ctx.JSON(http.StatusCreated, result)
}
