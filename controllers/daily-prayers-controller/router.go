package daily_prayer_times_controller

import (
	"app/controller"
	"app/models/city"
)

func (m *Controllers) SetBasicRoutes(
	d *controller.RouterDependencies,
) {
	f := d.E.Group("/daily-times")

	f.GET("", m.Basic.Index).Name = "dailytimes:index:admin,public"
	f.GET("/:id", m.Basic.Show).Name = "dailytimes:show:admin,public"

	r := d.Requires(city.ScopeAdmin)
	f.POST("", m.Basic.Store, r).Name = "dailytimes:store:admin"
	f.DELETE("/:id", m.Basic.Destroy, r).Name = "dailytimes:destroy:admin"
	f.PUT("/:id", m.Basic.Update, r).Name = "dailytimes:update:admin"
}
