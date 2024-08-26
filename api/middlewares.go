package api

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"app/config"
	"app/models/user"
	"app/translations"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// Function to sanitize request body (exclude sensitive information)
func sanitizeRequestBody(body []byte) string {
	// Convert byte slice to string for easier manipulation
	bodyStr := string(body)

	// Define patterns for sanitization
	passwordPatterns := `(?m)"password":\s*".*?"`
	tokenPattern := `"value":".*?"`

	// Sanitize password
	bodyStr = regexp.MustCompile(passwordPatterns).ReplaceAllString(bodyStr, `name="password"$1\r\n\r\n****\r\n`)

	// Sanitize token value
	sanitizedBody := regexp.MustCompile(tokenPattern).ReplaceAllString(bodyStr, `"value": "***"`)

	return sanitizedBody
}

func (app *Application) SetupMiddlewares(e *echo.Echo, isTest bool) {
	// e.Mount("/debug", middleware.Profiler())
	e.HTTPErrorHandler = app.APIErrors.GlobalErrorHandler
	e.JSONSerializer = CustomJSONSerializer{}
	// Locale must be the first middleware
	e.Use(app.locale())
	// e.Use(middleware.Recover())
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		Skipper:           middleware.DefaultSkipper,
		StackSize:         4 << 10, // 4 KB
		DisableStackAll:   false,
		DisablePrintStack: false,
		LogLevel:          log.ERROR,
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			c.Logger().Errorf("[PANIC RECOVER] %v", err)
			fmt.Println(string(stack)) //nolint: forbidigo //dw
			return app.APIErrors.InternalServer(c, err)
		},
	}))
	e.Use(middleware.Secure())

	// --------------------------------------------------------------
	// 		Standard Middlewares
	// --------------------------------------------------------------
	// CORS
	app.cors(e)

	// Request body limit
	e.Use(middleware.BodyLimit(config.MaxFormMemory))

	// X-Request-Id header
	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		TargetHeader: echo.HeaderXRequestID,
		Generator:    uuid.NewString,
	}))

	if !isTest {
		// Rate Limiter
		app.rateLimiter(e)
		// Gzip
		e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
			Level: 5,
			Skipper: func(ctx echo.Context) bool {
				// skip gzip on routes with this string
				if strings.Contains(ctx.Request().URL.Path, "swagger") {
					return true
				}
				if strings.Contains(ctx.Request().URL.Path, "metrics") {
					return true
				}
				return false
			},
		}))
	}

	// Register the logger middleware
	e.Use(middleware.BodyDump(func(ctx echo.Context, reqBody, resBody []byte) {
		var authorizedID string
		authUser := app.Utils.CtxUser(ctx)
		if authUser != nil {
			authorizedID = authUser.ID.String()
		}

		// Log the request and response
		app.Logger.Info().
			Str("Method", ctx.Request().Method).
			Str("Host", ctx.Request().Host).
			Str("URI", ctx.Path()).
			Str("authorized_id", authorizedID).
			Int("status", ctx.Response().Status).
			Str("request", sanitizeRequestBody(reqBody)).
			RawJSON("response", []byte(sanitizeRequestBody(resBody))).
			Msg("request")
	}))

	// ------------------------------------------------------------------------
	// 		Static Files
	// ------------------------------------------------------------------------
	e.Static(
		"/uploads/*",
		config.GetUploadsPath(""),
	).Name = "dir:uploads:public"

	e.Static("/*", config.GetRootPath("public")).Name = "dir:public:public"

	// Custom Middlewares
	e.Use(app.jwt())
	e.Use(app.scope())
}

func (app *Application) rateLimiter(e *echo.Echo) {
	rateLimitConfig := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{
				Rate:      10,
				Burst:     30,
				ExpiresIn: 1 * time.Minute,
			},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(
			ctx echo.Context,
			err error,
		) error {
			return ctx.JSON(http.StatusForbidden, nil)
		},
		DenyHandler: func(
			ctx echo.Context,
			identifier string,
			err error,
		) error {
			return ctx.JSON(http.StatusTooManyRequests, nil)
		},
	}
	e.Use(middleware.RateLimiterWithConfig(rateLimitConfig))
}

func (app *Application) cors(e *echo.Echo) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: app.Config.AllowedOrigins,
		AllowHeaders: []string{
			echo.HeaderAccept,
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderCacheControl,
			echo.HeaderAuthorization,
			echo.HeaderXRequestedWith,
			echo.HeaderXContentTypeOptions,
			echo.HeaderAccessControlAllowOrigin,
			echo.HeaderAccessControlAllowHeaders,
			echo.HeaderAccessControlAllowCredentials,
		},
		AllowCredentials: true,
	}))
}

func (app *Application) locale() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			lang := ctx.QueryParam("lang")
			switch lang {
			case "ar":
				app.Lang = "ar"
			case "en":
				app.Lang = "en"
			default:
				app.Lang = "ar"
			}
			ctx.Set("lang", app.Lang)
			localizer := i18n.NewLocalizer(app.LangBundle, app.Lang)
			t := &translations.Translations{Localizer: localizer}
			t.TranslateModels() // important to translate all api models
			ctx.Set("t", t)
			return next(ctx)
		}
	}
}

func (app *Application) jwt() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			var accessToken string
			var authorizedUser user.Model
			cookie, err := ctx.Cookie("accessToken")
			if err != nil {
				accessToken = ""
			} else {
				accessToken = strings.
					Replace(cookie.Value, "accessToken=", "", 1)
			}
			bearer := ctx.Request().Header.Get("Authorization")
			if bearer == "" && accessToken == "" {
				// no values are provided go to the next context with empty user
				app.CtxUser = nil
				return next(ctx)
			}
			if bearer != "" {
				splitBearer := strings.Split(bearer, " ")
				if len(splitBearer) == 2 {
					bearer = splitBearer[1]
				}
			}
			if accessToken != "" {
				bearer = accessToken
			}
			if bearer != "" {
				token, claims, err := authorizedUser.ParseToken(&bearer)
				if err != nil {
					return app.APIErrors.Forbidden(ctx, err)
				}
				if token.Valid {
					userUUID, err := uuid.Parse(claims.Subject)
					if err != nil {
						return app.APIErrors.Forbidden(ctx, err)
					}
					authorizedUser.ID = userUUID
					if err := app.
						Models.User.GetOne(&authorizedUser); err != nil {
						return app.APIErrors.Database(
							ctx,
							err,
							&authorizedUser,
						)
					}
					// go to the next context with found user must be of
					// type *user.Model
					app.CtxUser = &authorizedUser
					ctx.Set("user", app.CtxUser)
					return next(ctx)
				} else {
					app.CtxUser = nil
					ctx.Set("user", app.CtxUser)
					return app.APIErrors.Forbidden(ctx, err)
				}
			}
			app.CtxUser = nil
			ctx.Set("user", app.CtxUser)
			return next(ctx)
		}
	}
}
