package wallet_controller

import (
	"app/models/setting"
	"app/pkg/payment_gateway"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (c *ControllerBasic) PaymentService(ctx echo.Context) error {
	settings := payment_gateway.Settings{}
	if err := c.Models.Setting.GetForPaymentGateway(&settings); err != nil {
		return c.APIErr.Database(ctx, err, &setting.Model{})
	}
	res, err := payment_gateway.PaymentService(&settings)
	if err != nil {
		return c.APIErr.ExternalRequestError(ctx, err)
	}
	return ctx.JSON(http.StatusCreated, res)
}
