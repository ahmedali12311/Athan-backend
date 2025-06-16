package category_controller

import (
	"net/http"

	"bitbucket.org/sadeemTechnology/backend-model-category"
	"github.com/labstack/echo/v4"
)

func (c *Controllers) Show(ctx echo.Context) error {
	var result category.Model
	if err := c.Utils.ReadUUIDParam(&result.ID, ctx); err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}
	ws := c.scope(ctx)
	if err := c.Models.Category.GetOne(&result, ws); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}
	return ctx.JSON(http.StatusOK, result)
}
