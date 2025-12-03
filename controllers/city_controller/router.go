package city_controller

import (
	"app/controller"
	"app/models/city"
)

func (m *Controllers) SetBasicRoutes(
	d *controller.RouterDependencies,
) {
	f := d.E.Group("/cities")

	f.GET("", m.Basic.Index).Name = "cities:index:admin,public"
	f.GET("/:id", m.Basic.Show).Name = "cities:show:admin,public"

	r := d.Requires(city.ScopeAdmin)

	f.POST("", m.Basic.Store, r).Name = "cities:store:admin"
	f.DELETE("/:id", m.Basic.Destroy, r).Name = "cities:destroy:admin"
	f.PUT("/:id", m.Basic.Update).Name = "cities:update:admin"

}
