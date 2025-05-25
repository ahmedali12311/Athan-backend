package wallet_controller

import (
	"errors"
	"fmt"
	"net/http"

	"app/models/setting"
	"app/models/wallet_transaction"
	"app/pkg/payment_gateway"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (c *ControllerBasic) MasaratInitiate(ctx echo.Context) error {
	ctxUser := c.Utils.CtxUser(ctx)
	model := wallet_transaction.Model{
		ID:       uuid.New(),
		WalletID: ctxUser.ID,
		Type:     wallet_transaction.TypeCredit,
		User:     wallet_transaction.WalletUser{ID: &ctxUser.ID},
	}

	input := payment_gateway.MasaratInitiateRequest{
		WalletTransactionID: model.ID,
	}
	settings := payment_gateway.Settings{}
	if err := c.Models.Setting.GetForPaymentGateway(&settings); err != nil {
		return c.APIErr.Database(ctx, err, &setting.Model{})
	}

	v, err := c.GetValidator(ctx, model.ModelName())
	if err != nil {
		return err
	}

	v.AssignFloat("amount", &model.Amount)
	v.Check(model.Amount > 0, "amount", v.T.ValidateMustBeGtZero())

	v.AssignString("identity_card", &input.IdentityCard, 0, 50)
	v.Check(input.IdentityCard != "", "identity_card", v.T.ValidateRequired())

	if input.PaymentServiceID, err = uuid.Parse(v.Data.Get("payment_service_id")); err != nil {
		v.Check(false, "payment_service_id", v.T.ValidateUUID())
	}
	if !v.Valid() {
		return c.APIErr.InputValidation(ctx, v)
	}

	input.Amount = model.Amount

	res, err := payment_gateway.MasaratInitiatePayment(&settings, &input)
	if err != nil {
		return c.APIErr.ExternalRequestError(ctx, err)
	}
	model.PaymentReference = res.PaymentReference

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

func (c *ControllerBasic) MasaratConfirm(ctx echo.Context) error {
	var result wallet_transaction.Model

	ctxUser := c.Utils.CtxUser(ctx)

	v, err := c.GetValidator(ctx, result.ModelName())
	if err != nil {
		return err
	}

	var input payment_gateway.MasaratConfirmRequest

	v.AssignString("pin", &input.Pin, 4, 8)
	v.AssignUUID(
		"transaction_id",
		"id",
		"wallet_transactions",
		&result.ID,
		true,
	)
	input.WalletTransactionID = result.ID

	if !v.Valid() {
		return c.APIErr.InputValidation(ctx, v)
	}

	if err := c.Models.Wallet.GetTransaction(&result, &ctxUser.ID); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}

	if result.IsConfirmed {
		err := errors.New(v.T.WalletTransactionAlreadyConfirmed())
		return c.APIErr.BadRequest(ctx, err)
	}

	settings := payment_gateway.Settings{}
	if err := c.Models.Setting.GetForPaymentGateway(&settings); err != nil {
		return c.APIErr.Database(ctx, err, &setting.Model{})
	}

	res, err := payment_gateway.MasaratTransactionConfirm(&settings, &input)
	if err != nil {
		return c.APIErr.ExternalRequestError(ctx, err)
	}

	// Start transacting
	tx, err := c.Models.DB.Beginx()
	if err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	defer func() { _ = tx.Rollback() }()

	result.TLyncResponse = res.Response
	result.IsConfirmed = true
	result.Notes = res.Notes
	result.PaymentMethod = res.PaymentMethod
	if result.PaymentMethod != nil {
		note := fmt.Sprintf(
			"تعبئة المحفظة بخدمة %s",
			*result.PaymentMethod,
		)
		result.Notes = &note
	}

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
