package setting_controller

import (
	"app/controller"

	"bitbucket.org/sadeemTechnology/backend-model-setting"
	"github.com/labstack/echo/v4"
)

type Controllers struct {
	*controller.Dependencies
}

func Get(d *controller.Dependencies) *Controllers {
	return &Controllers{d}
}

// scope ----------------------------------------------------------------------

func (c *Controllers) scope(ctx echo.Context) *setting.WhereScope {
	scopes := c.Utils.CtxScopes(ctx)
	var admin, public bool
	for _, v := range scopes {
		switch v {
		case "admin":
			admin = true
		case "public":
			public = true
		}
	}

	return &setting.WhereScope{
		IsAdmin:  admin,
		IsPublic: public && !admin,
	}
}
