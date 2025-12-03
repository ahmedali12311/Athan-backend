package daily_prayer_times_controller

import (
	"app/controller"
	"app/models/daily_prayer_times"
)

func (m *Controllers) SetBasicRoutes(
	d *controller.RouterDependencies,
) {
	b := d.E.Group("/daily_prayer_times")

	b.GET("", m.Basic.Index).Name = "daily_prayer_times:index:admin,public"
	b.GET("/:id", m.Basic.Show).Name = "daily_prayer_times:show:admin,public"

	r := d.Requires(daily_prayer_times.ScopeAdmin)

	b.POST("", m.Basic.Store, r).Name = "daily_prayer_times:store:admin"
	b.PUT("/:id", m.Basic.Update, r).Name = "daily_prayer_times:update:admin"
	b.DELETE("/:id", m.Basic.Destroy, r).Name = "daily_prayer_times:destroy:admin"
	b.POST("/bulk_csv", m.Basic.BulkCSV, r).Name = "daily_prayer_times:csv-bulk:admin"

}
