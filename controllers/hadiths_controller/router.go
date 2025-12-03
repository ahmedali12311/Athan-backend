package hadiths_controller

import (
	"app/controller"

	model "bitbucket.org/sadeemTechnology/backend-model"
)

func (m *Controllers) SetRoutes(
	d *controller.RouterDependencies,
) {
	b := d.E.Group("/hadiths")

	b.GET("", m.Basic.Index).Name = "hadiths:index:admin,public"
	b.GET("/:id", m.Basic.Show).Name = "hadiths:show:admin,public"

	r := d.Requires(model.ScopeAdmin)

	b.POST("", m.Basic.Store, r).Name = "hadiths:store:admin"
	b.PUT("/:id", m.Basic.Update, r).Name = "hadiths:update:admin"
	b.DELETE("/:id", m.Basic.Destroy, r).Name = "hadiths:destroy:admin"
}
