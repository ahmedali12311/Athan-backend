package wallet_controller

import (
	"net/http"

	"app/pkg/payment_gateway"

	"bitbucket.org/sadeemTechnology/backend-model-setting"
	"github.com/labstack/echo/v4"
)

func (c *ControllerBasic) PaymentService(ctx echo.Context) error {
	settings, err := c.Models.Setting.GetForTyrianAnt()
	if err != nil {
		return c.APIErr.Database(ctx, err, &setting.Model{})
	}
	res, err := payment_gateway.PaymentService(settings, ctx)
	if err != nil {
		return c.APIErr.ExternalRequestError(ctx, err)
	}
	return ctx.JSON(http.StatusCreated, res)
}
