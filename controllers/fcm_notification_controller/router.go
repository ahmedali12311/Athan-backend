package fcm_notification_controller

import (
	"app/controller"
	"app/models/fcm_notification"
)

func (m *Controllers) SetBasicRoutes(
	d *controller.RouterDependencies,
) {
	f := d.E.Group("/fcm-notifications")

	r := d.Requires(fcm_notification.ScopeAdmin)

	f.GET("", m.Basic.Index, r).Name = "fcm-notifications:index:admin"
	f.POST("", m.Basic.Store, r).Name = "fcm-notifications:store:admin"
	f.POST("/notify-user", m.Basic.NotifyUser, r).Name = "fcm-notifications.notify-user:store:admin"
	f.GET("/:id", m.Basic.Show, r).Name = "fcm-notifications:show:admin"
	f.PUT("/:id", m.Basic.Update, r).Name = "fcm-notifications:update:admin"
	f.DELETE("/:id", m.Basic.Destroy, r).Name = "fcm-notifications:destroy:admin"
}
