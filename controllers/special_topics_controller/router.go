package special_topics_controller

import (
	"app/controller"

	model "bitbucket.org/sadeemTechnology/backend-model"
)

func (m *Controllers) SetRoutes(
	d *controller.RouterDependencies,
) {
	b := d.E.Group("/special_topics")

	b.GET("", m.Basic.Index).Name = "special_topics:index:admin,public"
	b.GET("/:id", m.Basic.Show).Name = "special_topics:show:admin,public"

	r := d.Requires(model.ScopeAdmin)

	b.POST("", m.Basic.Store, r).Name = "special_topics:store:admin"
	b.PUT("/:id", m.Basic.Update, r).Name = "special_topics:update:admin"
	b.DELETE("/:id", m.Basic.Destroy, r).Name = "special_topics:destroy:admin"
}
