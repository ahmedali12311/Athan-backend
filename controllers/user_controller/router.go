//nolint:lll
package user_controller

import (
	"app/controller"
	"app/models/user"
)

func (m *Controllers) SetOTPRoutes(
	d *controller.RouterDependencies,
) {
	f := d.E.Group("/otp")
	f.POST("/login", m.OTP.Login).Name = "otp:login:public"
	f.POST("/request", m.OTP.Request).Name = "otp:request:public"
}

func (m *Controllers) SetAuthRoutes(
	d *controller.RouterDependencies,
) {
	d.E.POST("/login", m.Auth.Login).Name = "auth:login:public"
	// e.POST("/phone-login", m.Public.PhoneLogin).Name = "public:phone-login"
	// e.POST(
	//        "/register-unverified",
	//        m.Public.RegisterUnverified,
	//    ).Name = "public:register-unverified"

	// requires jwt
	d.E.GET("/logout", m.Auth.Logout).Name = "auth:logout:public"

	d.E.POST("/forget-my-password", m.Auth.ForgetMyPassword).Name = "auth:forget-my-password:public"
	d.E.POST("/reset-password", m.Auth.ResetPassword).Name = "auth:reset-password:public"
}

func (m *Controllers) SetProfileRoutes(
	d *controller.RouterDependencies,
) {
	f := d.E.Group("/me")

	f.GET("", m.Profile.Me).Name = "user:me:public"

	r := d.Requires(
		user.ScopeOwn,
	)

	f.PUT("", m.Profile.Update, r).Name = "user:me-update:own"
	// f.DELETE("", m.Profile.Clear, r).Name = "user:me-clear:own"
}

func (m *Controllers) SetAdminRoutes(
	d *controller.RouterDependencies,
) {
	f := d.E.Group("/users")
	r := d.Requires(
		user.ScopeAdmin,
	)

	f.GET("", m.Basic.Index, r).Name = "users:index:admin"
	f.POST("", m.Basic.Store, r).Name = "users:store:admin"
	f.GET("/:id", m.Basic.Show, r).Name = "users:show:admin"
	f.PUT("/:id", m.Basic.Update, r).Name = "users:update:admin"
	f.DELETE("/:id", m.Basic.Clear, r).Name = "users:clear:admin"

	f.GET("/:id/become", m.Admin.Become, r).Name = "users:become:admin"
	f.POST("/grant-role", m.Admin.GrantRole, r).Name = "users:grant-role:admin"
	f.POST("/revoke-role", m.Admin.RevokeRole, r).Name = "users:revoke-role:admin"
}
