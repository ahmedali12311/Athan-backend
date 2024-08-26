package token_controller

import (
	"app/controller"
	"app/models/token"
)

func (m *Controllers) SetBasicRoutes(
	d *controller.RouterDependencies,
) {
	f := d.E.Group("/tokens")
	r := d.Requires(
		token.ScopeAdmin,
		token.ScopeOwn,
	)

	f.GET("", m.Basic.Index, r).Name = "tokens:index:admin,own"
	f.GET("/:id", m.Basic.Show, r).Name = "tokens:show:admin,own"
}
