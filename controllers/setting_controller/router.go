package setting_controller

import (
	"app/controller"
	"app/models/setting"
)

func (m *Controllers) SetBasicRoutes(
	d *controller.RouterDependencies,
) {
	f := d.E.Group("/settings")

	f.GET("", m.Basic.Index).Name = "settings:index:public,admin"
	f.GET("/:id", m.Basic.Show).Name = "settings:show:public,admin"

	r := d.Requires(setting.ScopeAdmin)

	f.POST("", m.Basic.Store, r).Name = "settings:store:admin"
	f.PUT("/:id", m.Basic.Update, r).Name = "settings:update:admin"
	f.DELETE("/:id", m.Basic.Destroy, r).Name = "settings:destroy:admin"
}
