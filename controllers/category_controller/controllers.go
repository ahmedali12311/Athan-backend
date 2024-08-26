package category_controller

import (
	"slices"

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

// Scopes ---------------------------------------------------------------------

// isPublic
func (c *Controllers) isPublic(scopes []string) bool {
	return !slices.Contains(scopes, category.ScopeAdmin)
}

// scope ----------------------------------------------------------------------
func (c *Controllers) scope(ctx echo.Context) *category.WhereScope {
	scopes := c.Utils.CtxScopes(ctx)
	public := c.isPublic(scopes)
	// showArchived := ctx.QueryParams().Has("archived")

	return &category.WhereScope{
		IsPublic: public,
	}
}
