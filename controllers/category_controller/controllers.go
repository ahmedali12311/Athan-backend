package category_controller

import (
	"app/controller"
	"app/models/category"

	"github.com/labstack/echo/v4"
)

type Controllers struct {
	*controller.Dependencies
}

func Get(d *controller.Dependencies) *Controllers {
	return &Controllers{d}
}

// scope ----------------------------------------------------------------------
func (c *Controllers) scope(ctx echo.Context) *category.WhereScope {
	scopes := c.Utils.CtxScopes(ctx)
	var admin bool
	for _, v := range scopes {
		switch v { //nolint: gocritic // dw
		case "admin":
			admin = true
		}
	}
	// showArchived := ctx.QueryParams().Has("archived")

	return &category.WhereScope{
		IsPublic: !admin,
	}
}
