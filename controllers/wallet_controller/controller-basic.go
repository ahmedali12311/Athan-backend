package wallet_controller

import (
	"errors"
	"fmt"
	"net/http"
	"slices"

	"app/controller"
	"app/models/setting"
	"app/models/wallet"
	"app/models/wallet_transaction"
	"app/pkg/tlync"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
)

type ControllerBasic struct {
	*controller.Dependencies
}

// Scopes ---------------------------------------------------------------------

func (c *ControllerBasic) userScope(ctx echo.Context) *uuid.UUID {
	scopes := c.Utils.CtxScopes(ctx)
	if slices.Contains(scopes, "admin") {
		return nil
	}
	return &c.Utils.CtxUser(ctx).ID
}

// Actions --------------------------------------------------------------------

func (c *ControllerBasic) Index(ctx echo.Context) error {
	userID := c.userScope(ctx)
	indexResponse, err := c.Models.Wallet.GetAllTransactions(ctx, userID)
	if err != nil {
		return c.APIErr.Database(ctx, err, nil)
	}
	settings := tlync.Settings{}
	if err := c.Models.Setting.GetForTlync(&settings); err != nil {
		return c.APIErr.Database(ctx, err, &setting.Model{})
	}

	for i := 0; i < len(*indexResponse.Data) && i < 2; i++ {
		if !(*indexResponse.Data)[i].IsConfirmed {

			tlyncInput := tlync.ConfirmInput{
				StoreID:   settings.StoreID,
				CustomRef: (*indexResponse.Data)[i].ID.String(),
			}
			res, err := tlync.TransactionReceipt(&settings, &tlyncInput)
			if err != nil {
				c.APIErr.LoggedOnly(ctx, err)
				continue
			}
			if res.Result != "success" {
				err := fmt.Errorf("confirmation is not complete: %s", res.Message)
				c.APIErr.LoggedOnly(ctx, err)
				continue
			}
			// Start transacting
			tx, err := c.Models.DB.Beginx()
			if err != nil {
				c.APIErr.LoggedOnly(ctx, err)
				continue
			}
			defer func() { _ = tx.Rollback() }()

			// marshalling/unmarshalling res into wallet trx struct
			b, err := json.Marshal(res.Data)
			if err != nil {
				c.APIErr.LoggedOnly(ctx, fmt.Errorf("marshalling res.Data: %w", err))
				continue
			}
			if err := json.Unmarshal(b, &(*indexResponse.Data)[i].TLyncResponse); err != nil {
				c.APIErr.LoggedOnly(ctx, fmt.Errorf("unmarshalling res.Data: %w", err))
				continue
			}

			(*indexResponse.Data)[i].IsConfirmed = true
			(*indexResponse.Data)[i].PaymentMethod = &res.Data.Gateway
			(*indexResponse.Data)[i].PaymentReference = &res.Data.OrderID
			if err := c.Models.Wallet.UpdateTransaction((*indexResponse.Data)[i], nil, tx); err != nil {
				c.APIErr.LoggedOnly(ctx, err)
				continue
			}
			if err = tx.Commit(); err != nil {
				c.APIErr.LoggedOnly(ctx, err)
				continue
			}
		}
	}
	return ctx.JSON(http.StatusOK, indexResponse)
}

func (c *ControllerBasic) Show(ctx echo.Context) error {
	var result wallet_transaction.Model
	if err := c.Utils.ReadUUIDParam(&result.ID, ctx); err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}
	userID := c.userScope(ctx)
	if err := c.Models.Wallet.GetTransaction(&result, userID); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}
	return ctx.JSON(http.StatusOK, result)
}

func (c *ControllerBasic) Store(ctx echo.Context) error {
	var result wallet_transaction.Model

	v, err := c.GetValidator(ctx, result.ModelName())
	if err != nil {
		return err
	}
	userID := c.userScope(ctx)
	ctxUser := c.Utils.CtxUser(ctx)
	result.User.ID = userID
	result.RechargedBy.ID = &ctxUser.ID
	if valid := result.MergeAndValidate(v); !valid {
		return c.APIErr.InputValidation(ctx, v)
	}
	// Start transacting
	tx, err := c.Models.DB.Beginx()
	if err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	defer func() { _ = tx.Rollback() }()

	if err := c.Models.Wallet.CreateTransaction(&result, tx); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23514" {
			t := c.Utils.CtxT(ctx)
			return c.APIErr.BadRequest(ctx, errors.New(t.TransactionDeclined()))
		}
		return c.APIErr.Database(ctx, err, &result)
	}
	if err = tx.Commit(); err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusCreated, result)
}

//	func (c *ControllerAdmin) Update(ctx echo.Context) error {
//		var result wallet_transaction.Model
//		if err := c.Utils.ReadUUIDParam(&result.ID, ctx); err != nil {
//			return c.APIErr.BadRequest(ctx, err)
//		}
//		if err := c.Models.Wallet.GetTransaction(&result, nil); err != nil {
//			return c.APIErr.Database(ctx, err, "Wallet.GetTransaction", result.ModelName())
//		}
//
//		v := validator.NewValidator(
//			c.Utils.Logger,
//			c.Utils.CtxT(ctx),
//			c.Models.Wallet.DB,
//			c.Schemas.WalletTransaction,
//		)
//		if err := v.Parse(ctx.Request()); err != nil {
//			return c.APIErr.BadRequest(ctx, err)
//		}
//
//		// only admin allowed to insert different user wallet.
//		v.UnmarshalInto("user", &result.User)
//
//		if valid := result.MergeAndValidate(v); !valid {
//			return c.APIErr.InputValidation(ctx, v)
//		}
//		// Start transacting
//		tx, err := c.Models.Wallet.DB.Beginx()
//		if err != nil {
//			return c.APIErr.InternalServer(ctx, err)
//		}
//		defer func() { _ = tx.Rollback() }()
//
//		if err := c.Models.Wallet.UpdateTransaction(&result, nil, tx); err != nil {
//			return c.APIErr.Database(ctx, err, "Wallet.UpdateTransactions", result.ModelName())
//		}
//		if err = tx.Commit(); err != nil {
//			return c.APIErr.InternalServer(ctx, err)
//		}
//
//		return ctx.JSON(http.StatusOK, result)
//	}

func (c *ControllerBasic) Destroy(ctx echo.Context) error {
	var result wallet_transaction.Model
	if err := c.Utils.ReadUUIDParam(&result.ID, ctx); err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}
	// Start transacting
	tx, err := c.Models.DB.Beginx()
	if err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	defer func() { _ = tx.Rollback() }()

	if err := c.Models.Wallet.DestroyTransaction(&result, nil, tx); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}
	if err = tx.Commit(); err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, result)
}

func (c *ControllerBasic) Wallet(ctx echo.Context) error {
	ctxUser := c.Utils.CtxUser(ctx)

	result, err := c.Models.Wallet.GetOne(&ctxUser.ID)
	if err != nil {
		return c.APIErr.Database(ctx, err, &wallet.Model{})
	}
	return ctx.JSON(http.StatusOK, result)
}
