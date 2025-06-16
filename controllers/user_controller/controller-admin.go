package user_controller

import (
	"fmt"
	"net/http"

	"app/controller"
	"app/models/role"
	"app/models/user"
	"bitbucket.org/sadeemTechnology/backend-config"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ControllerAdmin struct {
	*controller.Dependencies
}

// Scopes ---------------------------------------------------------------------
// handled by router, all actions are admin

// Actions --------------------------------------------------------------------

func (c *ControllerAdmin) Become(ctx echo.Context) error {
	var result user.Model
	if err := c.Utils.ReadUUIDParam(&result.ID, ctx); err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}
	if err := c.Models.User.GetOne(&result); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}
	tokenResponse, err := result.GenTokenResponse()
	if err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	cookie := result.GenCookie(
		tokenResponse.Token,
		config.TimeNow().Add(config.JwtExpiry),
	)
	ctx.SetCookie(&cookie)
	return ctx.JSON(http.StatusOK, tokenResponse)
}

func (c *ControllerAdmin) GrantRole(ctx echo.Context) error {
	var result role.Model
	v, err := c.GetValidator(ctx, result.ModelName())
	if err != nil {
		return err
	}

	var roleID int
	var userID uuid.UUID
	if valid := result.ValidateUserRole(v, &userID, &roleID); !valid {
		return c.APIErr.InputValidation(ctx, v)
	}

	// Start transacting
	tx, err := c.Models.DB.Beginx()
	if err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	defer func() { _ = tx.Rollback() }()
	if err := c.Models.User.GrantRole(&userID, &roleID, tx); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}
	if err := tx.Commit(); err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	responseBody := map[string]string{
		"message": fmt.Sprintf(
			"user %s granted role %d successfully",
			userID,
			roleID,
		),
	}
	return ctx.JSON(http.StatusOK, responseBody)
}

func (c *ControllerAdmin) RevokeRole(ctx echo.Context) error {
	var result role.Model
	v, err := c.GetValidator(ctx, result.ModelName())
	if err != nil {
		return err
	}

	var roleID int
	var userID uuid.UUID
	if valid := result.ValidateUserRole(v, &userID, &roleID); !valid {
		return c.APIErr.InputValidation(ctx, v)
	}
	// Start transacting
	tx, err := c.Models.DB.Beginx()
	if err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	defer func() { _ = tx.Rollback() }()
	// update user here
	if err := c.Models.User.RevokeRole(&userID, &roleID, tx); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}
	// Commit successful transaction
	if err := tx.Commit(); err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	responseBody := map[string]string{
		"message": fmt.Sprintf(
			"user %s revoked role %d successfully",
			userID,
			roleID,
		),
	}
	return ctx.JSON(http.StatusOK, responseBody)
}
