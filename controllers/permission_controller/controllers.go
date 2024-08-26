package permission_controller

import "app/controller"

type Controllers struct {
	Basic *ControllerBasic
}

func Get(deps *controller.Dependencies) *Controllers {
	return &Controllers{
		Basic: &ControllerBasic{deps},
	}
}
