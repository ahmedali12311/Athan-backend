package user_controller

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"app/controller"
	"app/models/token"

	"github.com/labstack/echo/v4"
)

// Scopes ---------------------------------------------------------------------
// handled by router, all actions are own

// Actions --------------------------------------------------------------------

type ControllerProfile struct {
	*controller.Dependencies
}

// for this project it will be different
func (c *ControllerProfile) Me(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, c.Utils.CtxUser(ctx))
}

func (c *ControllerProfile) Update(ctx echo.Context) error {
	ctxUser := c.Utils.CtxUser(ctx)
	v, err := c.GetValidator(ctx, ctxUser.ModelName())
	if err != nil {
		return err
	}
	if valid := ctxUser.MergeAndValidate(v); !valid {
		defer v.DeleteNewPicture()
		return c.APIErr.InputValidation(ctx, v)
	}
	// Start transacting
	tx, err := c.Models.DB.Beginx()
	if err != nil {
		defer v.DeleteNewPicture()
		return c.APIErr.InternalServer(ctx, err)
	}
	defer func() { _ = tx.Rollback() }()

	fcmToken := v.Data.Get("fcm_token")
	tokenType := token.TypeFCM

	if fcmToken != "" {
		userID := ctxUser.ID.String()
		if err := c.Models.User.
			SetFCMToken(&userID, &tokenType, &fcmToken, tx); err != nil {
			defer v.DeleteNewPicture()
			return c.APIErr.Database(
				ctx,
				fmt.Errorf("error setting fcm token: %w", err),
				ctxUser,
			)
		}
	}

	if err := c.Models.User.UpdateOne(ctxUser, tx); err != nil {
		defer v.DeleteNewPicture()
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return c.APIErr.Database(
				ctx,
				errors.New(v.T.ConflictError()),
				ctxUser,
			)
		default:
			return c.APIErr.Database(ctx, err, ctxUser)
		}
	}

	if err := tx.Commit(); err != nil {
		defer v.DeleteNewPicture()
		return c.APIErr.InternalServer(ctx, err)
	}
	return ctx.JSON(http.StatusOK, ctxUser)
}

func (c *ControllerProfile) Clear(ctx echo.Context) error {
	ctxUser := c.Utils.CtxUser(ctx)

	// Start transacting
	tx, err := c.Models.DB.Beginx()
	if err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	defer func() { _ = tx.Rollback() }()

	if err := c.Models.User.ClearOne(&ctxUser.ID, tx); err != nil {
		return c.APIErr.Database(ctx, err, ctxUser)
	}

	if err := tx.Commit(); err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	// Delete only if commit succeeds
	v, err := c.GetValidator(ctx, ctxUser.ModelName())
	if err != nil {
		return err
	}
	v.SaveOldImgThumbDists(ctxUser)
	v.DeleteOldPicture()
	return ctx.JSON(http.StatusOK, ctxUser)
}
