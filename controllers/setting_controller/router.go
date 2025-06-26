package setting_controller

import (
	"app/controller"

	"bitbucket.org/sadeemTechnology/backend-model"
)

func (m *Controllers) SetBasicRoutes(
	d *controller.RouterDependencies,
) {
	f := d.E.Group("/settings")

	f.GET("", m.Index).Name = "settings:index:public,admin"
	f.GET("/:id", m.Show).Name = "settings:show:public,admin"

	r := d.Requires(model.ScopeAdmin)

	f.POST("", m.Store, r).Name = "settings:store:admin"
	f.PUT("/:id", m.Update, r).Name = "settings:update:admin"
	f.DELETE("/:id", m.Destroy, r).Name = "settings:destroy:admin"
}
