package user_controller

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"app/apierrors"
	"app/controller"
	"app/models/user"
	"app/pkg/generics"
	"bitbucket.org/sadeemTechnology/backend-config"

	"github.com/labstack/echo/v4"
	"golang.org/x/exp/rand"
)

type ControllerAuth struct {
	*controller.Dependencies
}

func (c *ControllerAuth) Login(ctx echo.Context) error {
	var result user.Model
	comparePassword := false
	v, err := c.GetValidator(ctx, result.ModelName())
	if err != nil {
		return err
	}
	valid, err := result.MergeLogin(v, nil, &comparePassword)
	if err != nil {
		return c.APIErr.Firebase(ctx, err)
	}
	if !valid {
		return c.APIErr.InputValidation(ctx, v)
	}
	if err := c.Models.User.GetOne(&result); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.APIErr.InvalidCredentials(ctx)
		}
		return c.APIErr.Database(ctx, err, &result)
	}
	// does not compare password if firebase_id_token is provided or is
	// registration
	if comparePassword {
		if ok, err := result.Password.
			Match(result.PasswordHash); err != nil || !ok {
			return c.APIErr.InvalidCredentials(ctx)
		}
	}
	tokenResponse, err := result.GenTokenResponse()
	if err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	cookie := result.GenCookie(
		tokenResponse.Token,
		config.TimeNow().Add(config.JwtExpiry),
	)
	ctx.Response().Writer.Header().Add("HX-Redirect", "/me")
	ctx.SetCookie(&cookie)
	return ctx.JSON(http.StatusOK, tokenResponse)
}

//	func (c *ControllerPublic) PhoneLogin(ctx echo.Context) error {
//		result := user.Model{
//			Location: gis.EmptyPoint,
//		}
//		v := validator.NewValidator(
//			c.Utils.Logger,
//			c.Utils.CtxT(ctx),
//			c.Queries.User.DB,
//			c.Schemas.User,
//		)
//		if err := v.Parse(ctx.Request()); err != nil {
//			return c.APIErr.BadRequest(ctx, err)
//		}
//		firebaseIDToken := v.Data.Values.Get("firebase_id_token")
//		if firebaseIDToken == "" {
//			return c.APIErr.InvalidCredentials(ctx)
//		}
//		if err := result.
// VerifyIDToken(firebaseIDToken, c.Utils.FB); err != nil {
//			return c.APIErr.Firebase(ctx, err)
//		}
//		isRegister := false
//		if err := c.Queries.User.GetOne(&result); err != nil {
//			if errors.Is(err, sql.ErrNoRows) {
//				// flag for registration
//				isRegister = true
//			} else {
//				return c.APIErr.
//                  BadRequest(ctx, fmt.Errorf("error getting user: %w", err))
//			}
//		}
//		if isRegister {
//			newRegisterdUser := user.Model{
//				Gender:     nil,
//				Phone:      result.Phone,
//				IsVerified: true,
//				Location:   gis.EmptyPoint,
//			}
//			if valid := newRegisterdUser.MergeAndValidate(v); !valid {
//				return c.APIErr.InputValidation(ctx, v)
//			}
//			if err := c.register(ctx, &newRegisterdUser); err != nil {
// return c.APIErr.
//              BadRequest(ctx, fmt.Errorf("error registering user: %w", err))
//			}
//			result = newRegisterdUser
//		}
//		tokenResponse, err := result.GenTokenResponse()
//		if err != nil {
//			return c.APIErr.InternalServer(ctx, err)
//		}
//		cookie := result.
// GenCookie(tokenResponse.Token, config.TimeNow().Add(config.JwtExpiry))
//		ctx.SetCookie(&cookie)
//		return ctx.JSON(http.StatusOK, tokenResponse)
//	}
//
// func (c *ControllerPublic) register(
//     ctx echo.Context,
//     newUser *user.Model,
// ) error {
//		// Start transacting
//		tx, err := c.Queries.DB.Beginx()
//		if err != nil {
//			return err
//		}
//		defer func() { _ = tx.Rollback() }()
//		// create user
//		if err := c.Queries.User.CreateOne(newUser, tx); err != nil {
//			return err
//		}
//		if err = tx.Commit(); err != nil {
//			return err
//		}
//		return nil
//	}

func (c *ControllerAuth) Logout(ctx echo.Context) error {
	var model user.Model
	var token user.Token
	var accessToken string
	var message string
	var resType apierrors.ErrType
	var status int

	reqCookie, err := ctx.Cookie("accessToken")
	if err != nil {
		accessToken = ""
	} else {
		accessToken = reqCookie.Value
	}
	if accessToken != "" {
		message = "logged out"
		status = http.StatusOK
		resType = apierrors.Ok
	} else {
		message = "not logged in..."
		status = http.StatusUnauthorized
		resType = apierrors.Unauthorized
	}
	cookie := model.GenCookie(
		token,
		time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
	)
	ctx.SetCookie(&cookie)
	ctx.Set("scopes", []string{})
	ctx.Response().Writer.Header().Add("HX-Redirect", "/")
	response := apierrors.ErrMessage{
		Status:    status,
		Type:      resType,
		Message:   message,
		RequestID: ctx.Response().Header().Get(echo.HeaderXRequestID),
		Errors:    nil,
	}
	return ctx.JSON(status, response)
}

// func (c *ControllerPublic) RegisterUnverified(ctx echo.Context) error {
// 	v := validator.NewValidator(
// 		c.Utils.Logger,
// 		c.Utils.CtxT(ctx),
// 		c.Queries.User.DB,
// 		c.Schemas.User,
// 	)
// 	if err := v.Parse(ctx.Request()); err != nil {
// 		return c.APIErr.BadRequest(ctx, err)
// 	}
// 	result := user.Model{}
// 	if valid := result.MergeRegisterUnverified(v); !valid {
// 		defer v.DeleteNewPicture()
// 		return c.APIErr.InputValidation(ctx, v)
// 	}
// 	// Start transacting required for user
// 	tx, err := c.Queries.DB.Beginx()
// 	if err != nil {
// 		return c.APIErr.InternalServer(ctx, err)
// 	}
// 	defer func() { _ = tx.Rollback() }()
//
// 	if err := c.Queries.User.CreateOne(&result, tx); err != nil {
// 		defer v.DeleteNewPicture()
//	usedPhone := strings.Contains(err.Error(), "users_phone_unique_nullable")
//	usedEmail := strings.Contains(err.Error(), "users_email_unique_nullable")
//	if usedPhone || usedEmail {
// 			return c.APIErr.InvalidCredentials(ctx)
// 		}
// return c.APIErr.Database(ctx, err, "User.CreateOne", result.ModelName())
// 	}
// 	if err = tx.Commit(); err != nil {
// 		defer v.DeleteNewPicture()
// 		return c.APIErr.InternalServer(ctx, err)
// 	}
// 	tokenResponse, err := result.GenTokenResponse()
// 	if err != nil {
// 		return c.APIErr.InternalServer(ctx, err)
// 	}
// 	cookie := result.
// GenCookie(tokenResponse.Token, config.TimeNow().Add(config.JwtExpiry))
// 	ctx.SetCookie(&cookie)
// 	return ctx.JSON(http.StatusOK, tokenResponse)
// }

func (c *ControllerAuth) ForgetMyPassword(ctx echo.Context) error {
	var result user.Model

	v, err := c.GetValidator(ctx, result.ModelName())
	if err != nil {
		return err
	}

	// Delete the password field from the data to ensure it doesn't get updated
	v.Data.Del("password")

	valid, err := result.MergeForgetPassword(v)
	if err != nil {
		return c.APIErr.Firebase(ctx, err)
	}
	if !valid {
		return c.APIErr.InputValidation(ctx, v)
	}
	if err := c.Models.User.GetOne(&result); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}

	expires := time.Now().UTC().Add(5 * time.Minute)

	if result.PinExpiry != nil {
		if time.Now().UTC().Before(*result.PinExpiry) {
			message := fmt.Sprintf(
				"pin still active, try submitting again in %.2f seconds",
				result.PinExpiry.Sub(time.Now().UTC()).Seconds(),
			)
			return ctx.JSON(http.StatusOK, map[string]any{
				"message": message,
			})
		}
	}

	// TODO: check this overflow warning

	rng := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	result.Pin = generics.Ptr(fmt.Sprintf("%d", rng.Intn(900000)+100000))
	result.PinExpiry = &expires

	if err := c.Models.User.UpdateOne(&result, c.Models.DB); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}

	// email way from pale-skull
	// if err := c.Models.User.SendPasswordResetEmail(&result); err != nil {
	// 	return c.APIErr.BadRequest(ctx, err)
	// }

	msg := "A password reset PIN has been sent to your email." +
		" Please check your inbox."

	return ctx.JSON(
		http.StatusOK,
		map[string]string{
			"status":  "success",
			"message": msg,
		},
	)
}

func (c *ControllerAuth) ResetPassword(ctx echo.Context) error {
	var result user.Model
	v, err := c.GetValidator(ctx, result.ModelName())
	if err != nil {
		return err
	}

	valid, err := result.MergeForgetPassword(v)
	if err != nil {
		return c.APIErr.Firebase(ctx, err)
	}
	if !valid {
		return c.APIErr.InputValidation(ctx, v)
	}

	var pin string

	v.AssignString("pin", &pin, 6, 6)
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
				"pin expired %.2f seconds ago, please request another code",
				result.PinExpiry.Sub(time.Now().UTC()).Abs().Seconds(),
			)
			return ctx.JSON(http.StatusOK, map[string]any{
				"status":  "error",
				"message": message,
			})
		}
	}
	result.Pin = nil
	result.PinExpiry = nil
	if err := c.Models.User.UpdateOne(&result, c.Models.DB); err != nil {
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

// OTP WAY from royal-helix
// func (c *ControllerAuth) ForgetMyPassword(ctx echo.Context) error {
// 	var result user.Model

// 	v, err := c.GetValidator(ctx, result.ModelName())
// 	if err != nil {
// 		return err
// 	}
// 	valid, err := result.MergeForgetPassword(v, false)
// 	if err != nil {
// 		return c.APIErr.Firebase(ctx, err)
// 	}
// 	if !valid {
// 		return c.APIErr.InputValidation(ctx, v)
// 	}
// 	if err := c.Models.User.GetOne(&result); err != nil {
// 		return c.APIErr.Database(ctx, err, &result)
// 	}

// 	tx, err := c.Models.DB.Beginx()
// 	if err != nil {
// 		return c.APIErr.InternalServer(ctx, err)
// 	}
// 	defer func() { _ = tx.Rollback() }()

// 	expires := time.Now().UTC().Add(5 * time.Minute)

// 	settings, err := c.Models.Setting.GetForOTP()
// 	if err != nil {
// 		return c.APIErr.Database(ctx, err, &result)
// 	}
// 	input := &otp.Input{
// 		Phone: *result.Phone,
// 	}
// 	if result.PinExpiry != nil {
// 		if time.Now().UTC().Before(*result.PinExpiry) {
// 			message := fmt.Sprintf(
// 				"otp still active, try submitting again in %.2f seconds",
// 				result.PinExpiry.Sub(time.Now().UTC()).Seconds(),
// 			)
// 			return ctx.JSON(http.StatusOK, map[string]any{
// 				"status":  "succuess",
// 				"message": message,
// 				// FIX: remove and replace by http response
// 				"pin": result.Pin,
// 			})
// 		}
// 	}

// 	response, err := otp.Request(settings, input)
// 	if err != nil {
// 		return ctx.JSON(response.Status, response)
// 	}
// 	result.Pin = &response.Pin
// 	result.PinExpiry = &expires
// 	if err := c.Models.User.UpdateOne(&result, tx); err != nil {
// 		return c.APIErr.Database(ctx, err, &result)
// 	}
// 	if err = tx.Commit(); err != nil {
// 		return c.APIErr.InternalServer(ctx, err)
// 	}

// 	return ctx.JSON(
// 		response.Status,
// 		response,
// 	)
// }
