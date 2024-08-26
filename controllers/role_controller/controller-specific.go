package role_controller

import (
	"net/http"

	"app/models/role"

	"github.com/labstack/echo/v4"
)

func (c *ControllerBasic) GrantAllPermissions(ctx echo.Context) error {
	var id int
	if err := c.Utils.ReadIDParam(&id, ctx); err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}
	affected, err := c.Models.Role.GrantAllPermissions(id)
	if err != nil {
		return c.APIErr.Database(ctx, err, &role.Model{})
	}
	response := map[string]any{
		"role_id":    id,
		"grant_type": "all permissions",
		"affected":   affected,
	}
	return ctx.JSON(http.StatusCreated, response)
}

func (c *ControllerBasic) GrantByScope(ctx echo.Context) error {
	var id int
	if err := c.Utils.ReadIDParam(&id, ctx); err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}
	v, err := c.GetValidator(ctx, "role")
	if err != nil {
		return err
	}
	scope := v.Data.Get("scope")
	affected, err := c.Models.Role.GrantByScope(id, scope)
	if err != nil {
		return c.APIErr.Database(ctx, err, &role.Model{})
	}
	response := map[string]any{
		"role_id":    id,
		"grant_type": "scope",
		"scope":      scope,
		"affected":   affected,
	}
	return ctx.JSON(http.StatusCreated, response)
}

func (c *ControllerBasic) RevokeAllPermissions(ctx echo.Context) error {
	var id int
	if err := c.Utils.ReadIDParam(&id, ctx); err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}
	affected, err := c.Models.Role.RevokeAllPermissions(id)
	if err != nil {
		return c.APIErr.Database(ctx, err, &role.Model{})
	}
	response := map[string]any{
		"role_id":     id,
		"revoke_type": "all permissions",
		"affected":    affected,
	}
	return ctx.JSON(http.StatusOK, response)
}
