package hadiths_controller

import (
	"app/controller"
)

type Controllers struct {
	*controller.Dependencies
	Basic *ControllerBasic
}

func Get(d *controller.Dependencies) *Controllers {
	return &Controllers{
		Dependencies: d,
		Basic:        &ControllerBasic{Dependencies: d},
	}
}
