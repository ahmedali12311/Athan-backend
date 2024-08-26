package api

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"app/apierrors"

	"github.com/labstack/echo/v4"
)

func (app *Application) scope() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			scopes := []string{"public"}

			// This guard should protect all further requests
			// from null pointer dereference to ctxUser
			// it is set in the jwt middleware
			// no need to check it after this point
			if app.CtxUser == nil {
				ctx.Set("scopes", scopes)
				return next(ctx)
			}
			if app.CtxUser.IsDeleted {
				return app.APIErrors.Forbidden(ctx, apierrors.ErrDeletedAccount)
			}
			if app.CtxUser.IsDisabled {
				return app.APIErrors.
					Forbidden(ctx, apierrors.ErrDisabledAccount)
			}
			if app.CtxUser.Permissions == nil {
				ctx.Set("scopes", scopes)
				return next(ctx)
			}
			if len(*app.CtxUser.Permissions) == 0 {
				ctx.Set("scopes", scopes)
				return next(ctx)
			}
			method := ctx.Request().Method
			if method == "" {
				method = "GET"
			}

			// NOTE: need to check any paramaterized values and replaced them
			// here, should returns path as router:
			//  /api/v1/me
			//  /api/v1/categories/:id
			// paramID := ctx.Param("id")
			// path := ctx.Request().URL.Path

			// use registered path
			path := ctx.Path()

			// NOTE: this returns when requesting view components in html
			if !strings.Contains(path, "/api") {
				ctx.Set("scopes", scopes)
				return next(ctx)
			}

			// if paramID != "" {
			// 	segments := strings.Split(path, "/")
			// 	for i := range segments {
			// 		if segments[i] == paramID {
			// 			segments[i] = ":id"
			// 		}
			// 	}
			// 	path = strings.Join(segments, "/")
			// }
			for _, p := range app.Permissions {
				if p.Method == method && p.Path == path {
					for _, up := range *app.CtxUser.Permissions {
						if up == p.ID {
							scopes = append(scopes, p.Scope)
						}
					}
				}
			}
			ctx.Set("scopes", scopes)
			return next(ctx)
		}
	}
}

func (app *Application) requires(scopes ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			ctxScopes := app.Utils.CtxScopes(ctx)
			matches := 0
			for i := range scopes {
				if slices.Contains(ctxScopes, scopes[i]) {
					matches += 1
				}
			}
			if matches == 0 {
				msg := fmt.Sprintf(
					"insufficient scopes: [ %s ] requires: [ %s ]",
					strings.Join(ctxScopes, ","),
					strings.Join(scopes, ","),
				)
				err := errors.New(msg)
				return app.APIErrors.Forbidden(ctx, err)
			}
			return next(ctx)
		}
	}
}
