package adhkars_controller

import (
	"app/controller"
)

type Controllers struct {
	*controller.Dependencies
	Basic *ControllerBasic // Make sure to add this so it can be used in router.go
}

func Get(d *controller.Dependencies) *Controllers {
	return &Controllers{
		Dependencies: d,
		Basic:        &ControllerBasic{Dependencies: d}, // Initialize ControllerBasic
	}
}
