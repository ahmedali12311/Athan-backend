package apierrors

import (
	"encoding/json"
	"errors"
	"net/http"

	"app/models/user"
	"app/utilities"

	"bitbucket.org/sadeemTechnology/backend-validator"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// Client-Side Status Codes
//
// 404 Not Found
// 401 Unauthorized
// 403 Forbidden
// 400 Bad Request
// 409 Conflict // duplicates
// 422 Unprocessable Entity // validation
// 429 Too Many Requests
//
// Server-Side Status Codes
//
// 500 Internal Server Error
// 502 Bad Gateway
// 503 Service Unavailable
// 504 Gateway Timed Out
// 501 Not Implemented

type ErrType string

const (
	Ok                 ErrType = "Ok" // 200
	TokenExpired       ErrType = "TokenExpired"
	ModelNotFound      ErrType = "ModelNotFound"
	FileValidation     ErrType = "FileValidation"
	UnhandledError     ErrType = "UnhandledError"
	InputValidation    ErrType = "InputValidation"
	Unauthenticated    ErrType = "Unauthenticated"
	UnauthorizedAccess ErrType = "UnauthorizedAccess"
	InvalidCredentials ErrType = "InvalidCredentials"

	BadRequest          ErrType = "BadRequest"          // 400
	Unauthorized        ErrType = "Unauthorized"        // 401
	Forbidden           ErrType = "Forbidden"           // 403
	NotFound            ErrType = "NotFound"            // 404
	MethodNotAllowed    ErrType = "MethodNotAllowed"    // 405
	NotAcceptable       ErrType = "NotAcceptable"       // 406
	Conflict            ErrType = "Conflict"            // 409
	Locked              ErrType = "Locked"              // 423
	InternalServerError ErrType = "InternalServerError" // 500
)

type Errors map[string][]string

type APIErrors struct {
	Utils  *utilities.Utils
	Logger *zerolog.Logger
}

// ResMessage is used to create a simple none-error message.
type ResMessage struct {
	Status    int    `json:"status"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}

type ErrMessage struct {
	Status    int     `json:"status"`
	Type      ErrType `json:"type"`
	Message   string  `json:"message"`
	RequestID string  `json:"request_id"`
	Errors    any     `json:"errors,omitempty"`
	Caller    any     `json:"caller,omitempty"`
}

func Get(utils *utilities.Utils, logger *zerolog.Logger) *APIErrors {
	return &APIErrors{
		Utils:  utils,
		Logger: logger,
	}
}

var (
	ErrDisabledAccount    = errors.New("ErrDisabledAccount")
	ErrDeletedAccount     = errors.New("ErrDeletedAccount")
	ErrUnauthorizedAccess = errors.New("ErrUnauthorizedAccess")
)

func (e *APIErrors) LoggedOnlyInfo(ctx echo.Context, msg string) {
	e.Logger.
		Info().
		Int("status", ctx.Response().Status).
		Str("Method", ctx.Request().Method).
		Str("Host", ctx.Request().Host).
		Str("URI", ctx.Request().RequestURI).
		Msg(msg)
}

func (e *APIErrors) LoggedOnly(ctx echo.Context, err error) {
	e.Logger.Error().
		Err(err).
		Int("status", ctx.Response().Status).
		Str("Method", ctx.Request().Method).
		Str("Host", ctx.Request().Host).
		Str("URI", ctx.Request().RequestURI).
		Msg("LogOnlyError")
}

func (e *APIErrors) logError(ctx echo.Context, err error) {
	authorizedID := "none"
	authUser := e.Utils.CtxUser(ctx)
	if authUser != nil {
		authorizedID = authUser.ID.String()
	}
	// FIX: body logger
	// body, _ := io.ReadAll(ctx.Request().Body)
	e.Logger.Error().
		Stack().
		Err(err).
		Int("status", ctx.Response().Status).
		Str("Method", ctx.Request().Method).
		Str("Host", ctx.Request().Host).
		Str("URI", ctx.Request().RequestURI).
		Str("authorized_id", authorizedID).
		// Str("request", string(body)).
		// Str("trace", string(debug.Stack())).
		Msg("internal server error")
}

func (e *APIErrors) respond(
	ctx echo.Context,
	status int,
	errType ErrType,
	message string,
	data any,
) error {
	errMsg := ErrMessage{
		Status:    status,
		Type:      errType,
		Message:   message,
		RequestID: ctx.Response().Header().Get(echo.HeaderXRequestID),
		Errors:    data,
	}
	ctx.Response().
		Header().
		Set(echo.HeaderContentType, "application/problem+json")
	return ctx.JSON(status, errMsg)
}

func (e *APIErrors) GlobalErrorHandler(err error, ctx echo.Context) {
	t := e.Utils.CtxT(ctx)
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok { //nolint:errorlint // uncomparable
		code = he.Code
	}
	var errData any
	var errType ErrType
	message := t.OutOfScopeError()

	switch code {
	case 400:
		errType = BadRequest
		message = t.BadRequest()
	case 404:
		errType = NotFound
		message = t.NotFound()
	case 405:
		errType = MethodNotAllowed
		message = t.MethodNotAllowed()
	case 409:
		errType = Conflict
		message = t.ConflictError()
	case 500:
		errType = InternalServerError
		message = t.InternalServerError()
		e.logError(ctx, err)
	default:
		errType = UnhandledError
		errData = err
		e.logError(ctx, err)
	}

	if err := e.respond(ctx, code, errType, message, errData); err != nil {
		e.logError(ctx, err)
	}
}

func (e *APIErrors) InternalServer(ctx echo.Context, err error) error {
	t := e.Utils.CtxT(ctx)
	message := t.InternalServerError()
	e.logError(ctx, err)
	return e.respond(
		ctx,
		http.StatusInternalServerError,
		InternalServerError,
		message,
		nil,
	)
}

func (e *APIErrors) NotFound(ctx echo.Context) error {
	t := e.Utils.CtxT(ctx)
	message := t.NotFound()
	return e.respond(ctx, http.StatusNotFound, NotFound, message, nil)
}

func (e *APIErrors) BadRequest(ctx echo.Context, err error) error {
	return e.respond(ctx, http.StatusBadRequest, BadRequest, err.Error(), nil)
}

func (e *APIErrors) InvalidCredentials(ctx echo.Context) error {
	t := e.Utils.CtxT(ctx)
	return e.respond(
		ctx,
		http.StatusUnauthorized,
		Unauthorized,
		t.InvalidCredentials(),
		nil,
	)
}

func (e *APIErrors) Unauthorized(ctx echo.Context) error {
	t := e.Utils.CtxT(ctx)
	return e.respond(
		ctx,
		http.StatusUnauthorized,
		Unauthorized,
		t.UnauthorizedAccess(),
		nil,
	)
}

func (e *APIErrors) DisabledModel(ctx echo.Context, modelName string) error {
	t := e.Utils.CtxT(ctx)
	return e.respond(
		ctx,
		http.StatusForbidden,
		Forbidden,
		t.ModelDisabled(t.ModelName(modelName)),
		nil,
	)
}

func (e *APIErrors) InputValidation(
	ctx echo.Context,
	v *validator.Validator,
) error {
	t := e.Utils.CtxT(ctx)
	errMap := v.GetErrorMap()
	return e.respond(
		ctx,
		http.StatusUnprocessableEntity,
		InputValidation,
		t.InputValidation(),
		errMap,
	)
}

func (e *APIErrors) InputValidationRequest(
	ctx echo.Context,
	errs any,
) error {
	t := e.Utils.CtxT(ctx)
	return e.respond(
		ctx,
		http.StatusUnprocessableEntity,
		InputValidation,
		t.InputValidation(),
		errs,
	)
}

func (e *APIErrors) Forbidden(ctx echo.Context, err error) error {
	t := e.Utils.CtxT(ctx)
	var message string
	if errors.Is(err, jwt.ErrTokenMalformed) { //nolint: gocritic //dw
		message = "malformed jwt token"
	} else if errors.Is(err, jwt.ErrTokenExpired) {
		// Token is either expired or not active yet
		message = t.JwtExpired()
	} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
		// Token is either expired or not active yet
		message = "jwt token is not valid yet"
	} else if errors.Is(err, jwt.ErrSignatureInvalid) {
		// Token is either expired or not active yet
		u := &user.Model{}
		cookie := u.InvalidCookie()
		ctx.SetCookie(&cookie)
		message = err.Error()
	} else if errors.Is(err, ErrDeletedAccount) {
		message = t.DeletedAccount()
	} else if errors.Is(err, ErrDisabledAccount) {
		message = t.DisabledAccount()
	} else if errors.Is(err, ErrUnauthorizedAccess) {
		message = t.UnauthorizedAccess()
	} else {
		message = err.Error()
	}

	return e.respond(ctx, http.StatusForbidden, Forbidden, message, nil)
}

func (e *APIErrors) Firebase(ctx echo.Context, err error) error {
	return e.respond(
		ctx,
		http.StatusUnauthorized,
		Unauthorized,
		"firebase authentication error",
		map[string]string{"message": err.Error()},
	)
}

func (e *APIErrors) ExternalRequestError(
	ctx echo.Context,
	err error,
) error {
	t := e.Utils.CtxT(ctx)
	var result map[string]any

	errJson := json.Unmarshal([]byte(err.Error()), &result)
	if errJson != nil {
		return e.respond(
			ctx,
			http.StatusBadRequest,
			BadRequest,
			t.ExternalRequestError(),
			err.Error(),
		)
	}
	return e.respond(
		ctx,
		http.StatusBadRequest,
		BadRequest,
		t.ExternalRequestError(),
		result,
	)
}
