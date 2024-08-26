package utilities

import (
	"errors"
	"strconv"

	"app/config"
	"app/models/user"
	"app/translations"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type Utils struct {
	Version        *string
	CommitDescribe *string
	CommitCount    *string
	Settings       *config.Settings
	Logger         *zerolog.Logger
	FB             *firebase.App
	FBM            *messaging.Client
}

func Get(
	version,
	commitDescribe,
	commitCount *string,
	settings *config.Settings,
	logger *zerolog.Logger,
	fb *firebase.App,
	fbm *messaging.Client,
) *Utils {
	return &Utils{
		version,
		commitDescribe,
		commitCount,
		settings,
		logger,
		fb,
		fbm,
	}
}

// ReadUUIDParam parses and validates uuid id parameters.
func (u *Utils) ReadUUIDParam(id *uuid.UUID, ctx echo.Context) error {
	paramID := ctx.Param("id")
	parsed, err := uuid.Parse(paramID)
	if err != nil || parsed == uuid.Nil {
		return err
	}
	if id == nil {
		id = &parsed
	}
	*id = parsed
	return nil
}

func (u *Utils) ReadUUIDParamArg(id *uuid.UUID, param string, ctx echo.Context) error {
	paramID := ctx.Param(param)
	parsed, err := uuid.Parse(paramID)
	if err != nil || parsed == uuid.Nil {
		return err
	}
	if id == nil {
		id = &parsed
	}
	*id = parsed
	return nil
}

// ReadIDParam parses and validates integer id parameters.
func (u *Utils) ReadIDParam(id *int, ctx echo.Context) error {
	paramID := ctx.Param("id")
	parsed, err := strconv.ParseInt(paramID, 10, 64)
	if err != nil || parsed < 1 {
		return errors.New("invalid id parameter")
	}
	*id = int(parsed)
	return nil
}

// Context Variables ----------------------------------------------------------

func (u *Utils) CtxLang(ctx echo.Context) *string {
	ctxLang, ok := ctx.Get("lang").(string)
	if !ok {
		lang := "ar"
		ctxLang = lang
	}
	return &ctxLang
}

// CtxT is container for translated messages.
func (u *Utils) CtxT(ctx echo.Context) *translations.Translations {
	ctxLocalizer, ok := ctx.Get("t").(*translations.Translations)
	if !ok {
		ctxLocalizer = nil
	}
	return ctxLocalizer
}

func (u *Utils) CtxUser(ctx echo.Context) *user.Model {
	ctxUser, ok := ctx.Get("user").(*user.Model)
	if !ok {
		ctxUser = nil
	}
	return ctxUser
}

func (u *Utils) CtxScopes(ctx echo.Context) []string {
	ctxScopes, ok := ctx.Get("scopes").([]string)
	if !ok {
		ctxScopes = []string{}
	}
	return ctxScopes
}

// User Notification + Firebase -----------------------------------------------

// func (u *Utils) SendUserFCM(
// 	n *user_notification.Model,
// 	token *string,
// ) error {
// 	if n.User.ID == nil {
// 		return errors.New("user id not provided for notification")
// 	}
// 	if token != nil {
// 		message := firebaseutils.BuildTokenMessage(
// 			n.Title,
// 			n.Body,
// 			*token,
// 			n.Data,
// 		)
// 		response, err := u.FBM.Send(context.Background(), message)
// 		if err != nil {
// 			return err
// 		}
// 		n.Response = &response
// 	}
//
// 	return nil
// }
