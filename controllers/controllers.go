package controllers

import (
	"app/controller"
	"app/controllers/category_controller"
	"app/controllers/meta_controller"
	"app/controllers/permission_controller"
	"app/controllers/role_controller"
	"app/controllers/setting_controller"
	"app/controllers/token_controller"
	"app/controllers/user_controller"
)

type Controllers struct {
	// API ---------------------------------------------------------------------

	Category *category_controller.Controllers
	Meta     *meta_controller.Controllers

	Permission *permission_controller.Controllers

	Role    *role_controller.Controllers
	Setting *setting_controller.Controllers
	Token   *token_controller.Controllers
	User    *user_controller.Controllers
}

func Setup(d *controller.Dependencies) *Controllers {
	return &Controllers{
		Category:   category_controller.Get(d),
		Meta:       meta_controller.Get(d),
		Permission: permission_controller.Get(d),
		Role:       role_controller.Get(d),
		Setting:    setting_controller.Get(d),
		Token:      token_controller.Get(d),
		User:       user_controller.Get(d),
	}
}
