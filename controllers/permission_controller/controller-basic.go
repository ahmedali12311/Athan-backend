package permission_controller

import (
	"net/http"

	"app/controller"
	"app/models/permission"

	"github.com/labstack/echo/v4"
)

type ControllerBasic struct {
	*controller.Dependencies
}

func (c *ControllerBasic) Index(ctx echo.Context) error {
	indexResponse, err := c.Models.Permission.GetAll(ctx)
	if err != nil {
		return c.APIErr.Database(ctx, err, nil)
	}
	return ctx.JSON(http.StatusOK, indexResponse)
}

func (c *ControllerBasic) Show(ctx echo.Context) error {
	var result permission.Model
	if err := c.Utils.ReadIDParam(&result.ID, ctx); err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}
	if err := c.Models.Permission.GetOne(&result); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}
	return ctx.JSON(http.StatusOK, result)
}
