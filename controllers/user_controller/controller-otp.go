package user_controller

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"app/config"
	"app/controller"
	"app/models/user"
	"app/pkg/otp"

	"github.com/labstack/echo/v4"
)

type ControllerOTP struct {
	*controller.Dependencies
}

func (c *ControllerOTP) Request(ctx echo.Context) error {
	var result user.Model
	v, err := c.GetValidator(ctx, result.ModelName())
	if err != nil {
		return err
	}
	result.MergePhone(v)
	if !v.Valid() {
		return c.APIErr.InputValidation(ctx, v)
	}
	exists := true
	if err := c.Models.User.GetOne(&result); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exists = false
		} else {
			return c.APIErr.Database(ctx, err, &result)
		}
	}

	tx, err := c.Models.DB.Beginx()
	if err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	defer func() { _ = tx.Rollback() }()

	expires := time.Now().UTC().Add(5 * time.Minute)

	settings, err := c.Models.Setting.GetForOTP()
	if err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}
	input := &otp.Input{
		Phone: *result.Phone,
	}
	if !exists && result.Phone != nil {
		response, err := otp.Request(settings, input)
		if err != nil {
			// return c.APIErr.BadRequest(ctx, err)
			return ctx.JSON(response.Status, response)
		}
		result.MergeOTPCreate(v, &response.Pin, &expires)
		if err := c.Models.User.CreateOne(&result, tx); err != nil {
			return c.APIErr.Database(
				ctx,
				err,
				&result,
			)
		}
	}
	if result.PinExpiry != nil && exists {
		if time.Now().UTC().Before(*result.PinExpiry) {
			message := fmt.Sprintf(
				"otp still active, try submitting again in %.2f seconds",
				result.PinExpiry.Sub(time.Now().UTC()).Seconds(),
			)
			return ctx.JSON(http.StatusOK, map[string]any{
				"status":  "succuess",
				"exists":  exists,
				"message": message,
				// FIX: remove and replace by http response
				"pin": result.Pin,
			})
		}
	}

	response, err := otp.Request(settings, input)
	if err != nil {
		return ctx.JSON(response.Status, response)
	}
	result.Pin = &response.Pin
	result.PinExpiry = &expires
	if err := c.Models.User.UpdateOne(&result, tx); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}
	if err = tx.Commit(); err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}

	return ctx.JSON(
		response.Status,
		response,
	)
}

func (c *ControllerOTP) Login(ctx echo.Context) error {
	var result user.Model
	v, err := c.GetValidator(ctx, result.ModelName())
	if err != nil {
		return err
	}

	var pin string
	result.MergePhone(v)

	v.AssignString("pin", &pin)
	v.Check(pin != "", "pin", v.T.ValidateRequired())
	if !v.Valid() {
		return c.APIErr.InputValidation(ctx, v)
	}
	result.Pin = &pin

	if err := c.Models.User.GetOne(&result); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.APIErr.InvalidCredentials(ctx)
		}
		return c.APIErr.Database(ctx, err, &result)
	}

	if result.PinExpiry != nil {
		if !time.Now().UTC().Before(*result.PinExpiry) {
			message := fmt.Sprintf(
				"otp expired %.2f seconds ago, please request another code",
				result.PinExpiry.Sub(time.Now().UTC()).Abs().Seconds(),
			)
			return ctx.JSON(http.StatusOK, map[string]any{
				"status":  "error",
				"message": message,
				"pin":     result.Pin,
				// FIX: remove and replace by http response
			})
		}
	}
	if err := c.Models.User.Verify(&result.ID, nil); err != nil {
		c.APIErr.LoggedOnly(ctx, err)
	}
	tokenResponse, err := result.GenTokenResponse()
	if err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	cookie := result.GenCookie(
		tokenResponse.Token,
		time.Now().Add(config.JwtExpiry),
	)
	ctx.SetCookie(&cookie)
	return ctx.JSON(http.StatusOK, tokenResponse)
}
