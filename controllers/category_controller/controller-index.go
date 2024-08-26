package category_controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (c *Controllers) Index(ctx echo.Context) error {
	ws := c.scope(ctx)
	indexResponse, err := c.Models.Category.GetAll(ctx, ws)
	if err != nil {
		return c.APIErr.Database(ctx, err, nil)
	}
	return ctx.JSON(http.StatusOK, indexResponse)
}
