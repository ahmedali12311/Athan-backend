package category_controller

import (
	"app/controller"
)

func (m *Controllers) SetRoutes(
	d *controller.RouterDependencies,
) {
	f := d.E.Group("/categories")

	f.GET("", m.Index).Name = "categories:index:public,admin"
	f.GET("/:id", m.Show).Name = "categories:show:public,admin"

	r := d.Requires("admin")

	f.POST("", m.Store, r).Name = "categories:store:admin"
	f.PUT("/:id", m.Update, r).Name = "categories:update:admin"
	f.DELETE("/:id", m.Destroy, r).Name = "categories:destroy:admin"
}
