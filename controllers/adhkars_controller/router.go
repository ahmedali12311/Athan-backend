package adhkars_controller

import (
	"app/controller"

	model "bitbucket.org/sadeemTechnology/backend-model"
)

func (m *Controllers) SetRoutes(
	d *controller.RouterDependencies,
) {
	b := d.E.Group("/adhkars")

	b.GET("", m.Basic.Index).Name = "adhkars:index:admin,public"
	b.GET("/:id", m.Basic.Show).Name = "adhkars:show:admin,public"

	r := d.Requires(model.ScopeAdmin)

	b.POST("", m.Basic.Store, r).Name = "adhkars:store:admin"
	b.PUT("/:id", m.Basic.Update, r).Name = "adhkars:update:admin"
	b.DELETE("/:id", m.Basic.Destroy, r).Name = "adhkars:destroy:admin"
}
