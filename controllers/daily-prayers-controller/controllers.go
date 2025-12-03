package daily_prayer_times_controller

import "app/controller"

type Controllers struct {
	Basic *ControllerBasic
}

func Get(d *controller.Dependencies) *Controllers {
	return &Controllers{
		Basic: &ControllerBasic{d},
	}
}
