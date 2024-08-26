package token_controller

import (
	"net/http"
	"slices"

	"app/controller"
	"app/models/token"

	"github.com/google/uuid"
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
	indexResponse, err := c.Models.Token.GetAll(ctx, userID)
	if err != nil {
		return c.APIErr.Database(ctx, err, &token.Model{})
	}
	return ctx.JSON(http.StatusOK, indexResponse)
}

func (c *ControllerBasic) Show(ctx echo.Context) error {
	var result token.Model
	if err := c.Utils.ReadIDParam(&result.ID, ctx); err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}
	userID := c.userScope(ctx)
	if err := c.Models.Token.GetOne(&result, userID); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}
	return ctx.JSON(http.StatusOK, result)
}
