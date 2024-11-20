package wallet_controller

//
// import (
// 	"errors"
// 	"fmt"
// 	"net/http"
//
// 	"app/cmd/api/controller"
// 	"app/cmd/api/validator"
// 	"app/internal/database/user_account"
// 	"app/internal/database/wallet_transaction"
// 	"app/internal/generics"
// 	"app/internal/sasv4"
// 	"app/internal/tlync"
//
// 	"github.com/goccy/go-json"
// 	"github.com/google/uuid"
// 	"github.com/labstack/echo/v4"
// )
//
// type ControllerProfile struct {
// 	controller.Dependencies
// }
//
// func (c *ControllerProfile) Index(ctx echo.Context) error {
// 	ctxUser := c.Utils.CtxUser(ctx)
// 	indexResponse, err := c.Queries.Wallet.GetAllTransactions(ctx, &ctxUser.ID)
// 	if err != nil {
// 		return c.APIErr.Database(ctx, err, "Wallet.GetAllTransactions", "")
// 	}
// 	return ctx.JSON(http.StatusOK, indexResponse)
// }
//
// func (c *ControllerProfile) Show(ctx echo.Context) error {
// 	var result wallet_transaction.Model
// 	if err := c.Utils.ReadUUIDParam(&result.ID, ctx); err != nil {
// 		return c.APIErr.BadRequest(ctx, err)
// 	}
//
// 	ctxUser := c.Utils.CtxUser(ctx)
// 	if err := c.Queries.Wallet.GetTransaction(&result, &ctxUser.ID); err != nil {
// 		return c.APIErr.Database(ctx, err, "Wallet.GetTransaction", result.ModelName())
// 	}
//
// 	return ctx.JSON(http.StatusOK, result)
// }
//
// func (c *ControllerProfile) Store(ctx echo.Context) error {
// 	v := validator.NewValidator(
// 		c.Utils.Logger,
// 		c.Utils.CtxT(ctx),
// 		c.Queries.Wallet.DB,
// 		c.Schemas.WalletTransaction,
// 	)
// 	if err := v.Parse(ctx.Request()); err != nil {
// 		return c.APIErr.BadRequest(ctx, err)
// 	}
// 	ctxUser := c.Utils.CtxUser(ctx)
// 	model := wallet_transaction.Model{
// 		ID:       uuid.New(),
// 		WalletID: ctxUser.ID,
// 		Type:     wallet_transaction.TypeCredit,
// 		User: wallet_transaction.WalletUser{
// 			ID: &ctxUser.ID,
// 		},
// 		RechargedBy: wallet_transaction.WalletUser{
// 			ID: &ctxUser.ID,
// 		},
// 	}
// 	settings := tlync.Settings{}
// 	if err := c.Queries.Setting.GetForTlync(&settings); err != nil {
// 		return c.APIErr.Database(ctx, err, "Setting.GetForTlync", "setting")
// 	}
// 	tlyncInput := tlync.InitiateInput{
// 		ID:          settings.StoreID,
// 		Amount:      fmt.Sprintf("%.2f", model.Amount),
// 		Phone:       "218910001234",
// 		Email:       "info@flexnet.ly",
// 		BackendURL:  "https://flexnet.sadeem-lab.com/api/v1/meta",
// 		FrontendURL: settings.FrontURL,
// 		CustomRef:   model.ID.String(),
// 	}
// 	email := v.Data.Get("email")
// 	phone := v.Data.Get("phone")
// 	if email != "" {
// 		tlyncInput.Email = email
// 	}
// 	if phone != "" {
// 		tlyncInput.Phone = phone
// 	}
// 	if valid := model.MergeAndValidate(v); !valid {
// 		return c.APIErr.InputValidation(ctx, v)
// 	}
// 	tlyncInput.Amount = fmt.Sprintf("%.2f", model.Amount)
// 	res, err := tlync.InitiatePayment(&settings, &tlyncInput)
// 	if err != nil {
// 		return c.APIErr.BadRequest(ctx, err)
// 	}
//
// 	model.TLyncURL = &res.URL
// 	res.Amount = model.Amount
// 	// Start transacting
// 	tx, err := c.Queries.Wallet.DB.Beginx()
// 	if err != nil {
// 		return c.APIErr.InternalServer(ctx, err)
// 	}
// 	defer func() { _ = tx.Rollback() }()
//
// 	if err := c.Queries.Wallet.CreateTransaction(&model, &ctxUser.ID, tx); err != nil {
// 		return c.APIErr.Database(ctx, err, "Wallet.CreateTransaction", model.ModelName())
// 	}
// 	if err = tx.Commit(); err != nil {
// 		return c.APIErr.InternalServer(ctx, err)
// 	}
// 	return ctx.JSON(http.StatusCreated, model)
// }
//
// func (c *ControllerProfile) Confirm(ctx echo.Context) error {
// 	var result wallet_transaction.Model
// 	if err := c.Utils.ReadUUIDParam(&result.ID, ctx); err != nil {
// 		return c.APIErr.BadRequest(ctx, err)
// 	}
// 	ctxUser := c.Utils.CtxUser(ctx)
// 	if err := c.Queries.Wallet.GetTransaction(&result, &ctxUser.ID); err != nil {
// 		return c.APIErr.Database(ctx, err, "Wallet.GetTransaction", result.ModelName())
// 	}
// 	settings := tlync.Settings{}
// 	if err := c.Queries.Setting.GetForTlync(&settings); err != nil {
// 		return c.APIErr.Database(ctx, err, "Setting.GetForTlync", "setting")
// 	}
// 	tlyncInput := tlync.ConfirmInput{
// 		StoreID:   settings.StoreID,
// 		CustomRef: result.ID.String(),
// 	}
// 	res, err := tlync.TransactionReceipt(&settings, &tlyncInput)
// 	if err != nil {
// 		return c.APIErr.BadRequest(ctx, err)
// 	}
// 	if res.Result != "success" {
// 		err := fmt.Errorf("confirmation is not complete: %s", res.Message)
// 		return c.APIErr.BadRequest(ctx, err)
// 	}
// 	// Start transacting
// 	tx, err := c.Queries.Wallet.DB.Beginx()
// 	if err != nil {
// 		return c.APIErr.InternalServer(ctx, err)
// 	}
// 	defer func() { _ = tx.Rollback() }()
//
// 	// marshalling/unmarshalling res into wallet trx struct
// 	b, err := json.Marshal(res.Data)
// 	if err != nil {
// 		return c.APIErr.BadRequest(ctx, fmt.Errorf("marshalling res.Data: %w", err))
// 	}
// 	if err := json.Unmarshal(b, &result.TLyncResponse); err != nil {
// 		return c.APIErr.BadRequest(ctx, fmt.Errorf("unmarshalling res.Data: %w", err))
// 	}
//
// 	result.IsConfirmed = true
// 	result.PaymentMethod = &res.Data.Gateway
// 	result.PaymentReference = &res.Data.OrderID
// 	if err := c.Queries.Wallet.UpdateTransaction(&result, &ctxUser.ID, tx); err != nil {
// 		return c.APIErr.Database(ctx, err, "Wallet.UpdateTransaction", result.ModelName())
// 	}
// 	if err = tx.Commit(); err != nil {
// 		return c.APIErr.InternalServer(ctx, err)
// 	}
// 	return ctx.JSON(http.StatusCreated, result)
// }
//
// // Transfer moves funds from app wallet to sasv4 wallet
// func (c *ControllerProfile) Transfer(ctx echo.Context) error {
// 	v := validator.NewValidator(
// 		c.Utils.Logger,
// 		c.Utils.CtxT(ctx),
// 		c.Queries.Wallet.DB,
// 		c.Schemas.WalletTransaction,
// 	)
// 	if err := v.Parse(ctx.Request()); err != nil {
// 		return c.APIErr.BadRequest(ctx, err)
// 	}
// 	ctxUser := c.Utils.CtxUser(ctx)
//
// 	amount := v.Data.GetFloat("amount")
// 	accountID, err := v.Data.GetUUID("account_id")
// 	if err != nil || accountID == nil {
// 		v.Check(false, "account_id", v.T.ValidateUUID())
// 	}
// 	v.Check(amount > 0, "amount", v.T.ValidateMustBeGtZero())
// 	v.Exists(accountID, "account_id", "user_accounts", true)
// 	if !v.Valid() {
// 		return c.APIErr.InputValidation(ctx, v)
// 	}
// 	if ctxUser.Wallet.Credit <= amount {
// 		err := errors.New(v.T.CartOrderNotEnoughCredit(amount))
// 		return c.APIErr.BadRequest(ctx, err)
// 	}
//
// 	userAccount := user_account.Model{
// 		ID: *accountID,
// 	}
// 	if err := c.Queries.UserAccount.GetOne(&userAccount, &ctxUser.ID); err != nil {
// 		return c.APIErr.Database(ctx, err, "UserAccount.GetOne", userAccount.ModelName())
// 	}
//
// 	sasv4TrxUUID := uuid.New().String()
// 	walletTRX := wallet_transaction.Model{
// 		ID:               uuid.New(),
// 		WalletID:         ctxUser.ID,
// 		Amount:           amount,
// 		IsConfirmed:      true,
// 		Type:             wallet_transaction.TypeDebit,
// 		PaymentMethod:    generics.Ptr("transfer"),
// 		PaymentReference: &sasv4TrxUUID,
// 		Notes:            generics.Ptr("wallet to sasv4 transfer"),
// 		User: wallet_transaction.WalletUser{
// 			ID: &ctxUser.ID,
// 		},
// 		RechargedBy: wallet_transaction.WalletUser{
// 			ID: nil,
// 		},
// 	}
//
// 	// Start transacting
// 	tx, err := c.Queries.Wallet.DB.Beginx()
// 	if err != nil {
// 		return c.APIErr.InternalServer(ctx, err)
// 	}
// 	defer func() { _ = tx.Rollback() }()
//
// 	if err := c.Queries.Wallet.CreateTransaction(&walletTRX, &ctxUser.ID, tx); err != nil {
// 		return c.APIErr.Database(ctx, err, "Wallet.CreateTransaction", walletTRX.ModelName())
// 	}
//
// 	// Transfer operation start after the creation of local transaction
// 	// 1. get sas credentials and login admin
// 	creds, err := c.Queries.Setting.GetSASv4()
// 	if err != nil {
// 		return c.APIErr.Database(ctx, err, "Setting.GetSASv4", "setting")
// 	}
// 	inputadmin := sasv4.LoginInput{
// 		Username: creds.Username,
// 		Password: creds.Password,
// 		Language: "en",
// 	}
// 	resAdmin, err := sasv4.AdminLogin(&inputadmin)
// 	if err != nil {
// 		return c.APIErr.BadRequest(ctx, fmt.Errorf("SASv4.AdminLogin: %w", err))
// 	}
//
// 	input := sasv4.UserDepositInput{
// 		UserID:        userAccount.SasID,
// 		UserUsername:  userAccount.Username,
// 		Amount:        amount,
// 		Comment:       "user transferred by flexnet/sadeem app",
// 		TransactionID: sasv4TrxUUID,
// 	}
// 	res, err := sasv4.UserDeposit(resAdmin.Token, &input)
// 	if err != nil {
// 		return c.APIErr.BadRequest(ctx, fmt.Errorf("SASv4.UserDeposit: %w", err))
// 	}
// 	if res.Status != 200 {
// 		return c.APIErr.BadRequest(ctx, fmt.Errorf("SASv4.UserDeposit not OK"))
// 	}
//
// 	if err = tx.Commit(); err != nil {
// 		return c.APIErr.InternalServer(ctx, err)
// 	}
// 	return ctx.JSON(http.StatusCreated, walletTRX)
// }
