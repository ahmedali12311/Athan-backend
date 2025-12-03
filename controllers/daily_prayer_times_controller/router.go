package daily_prayer_times_controller

import (
	"app/controller"
	"app/models/daily_prayer_times"
)

func (m *Controllers) SetBasicRoutes(
	d *controller.RouterDependencies,
) {
	b := d.E.Group("/daily-prayer-times")

	b.GET("", m.Basic.Index).Name = "daily-prayer-times:index:admin,public"
	b.GET("/:id", m.Basic.Show).Name = "daily-prayer-times:show:admin,public"
	b.GET("/rolling", m.Basic.Rolling).Name = "daily-prayer-times:rolling:public"

	r := d.Requires(daily_prayer_times.ScopeAdmin)

	b.POST("", m.Basic.Store, r).Name = "daily-prayer-times:store:admin"
	b.PUT("/:id", m.Basic.Update, r).Name = "daily-prayer-times:update:admin"
	b.DELETE("/:id", m.Basic.Destroy, r).Name = "daily-prayer-times:destroy:admin"
	b.POST("/bulk-csv", m.Basic.BulkCSV, r).Name = "daily-prayer-times:bulk-csv:admin"

}
