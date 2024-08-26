package user_controller

import "app/controller"

type Controllers struct {
	OTP     *ControllerOTP
	Auth    *ControllerAuth
	Admin   *ControllerAdmin
	Basic   *ControllerBasic
	Profile *ControllerProfile
}

func Get(deps *controller.Dependencies) *Controllers {
	return &Controllers{
		OTP:     &ControllerOTP{deps},
		Auth:    &ControllerAuth{deps},
		Admin:   &ControllerAdmin{deps},
		Basic:   &ControllerBasic{deps},
		Profile: &ControllerProfile{deps},
	}
}
