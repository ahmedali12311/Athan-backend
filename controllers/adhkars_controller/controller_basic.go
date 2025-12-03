package adhkars_controller

import (
	"database/sql"
	"errors"
	"net/http"

	"app/controller"
	"app/models/adhkars"

	"github.com/labstack/echo/v4"
)

type ControllerBasic struct {
	*controller.Dependencies
}

// Scopes ---------------------------------------------------------------------

func (c *ControllerBasic) scope(
	ctx echo.Context,
) *adhkars.WhereScope {
	scopes := c.Utils.CtxScopes(ctx)
	ctxUser := c.Utils.CtxUser(ctx)

	var admin, public bool
	for _, v := range scopes {
		switch v {
		case "admin":
			admin = true
		case "public":
			public = true
		}
	}

	ws := &adhkars.WhereScope{
		IsAdmin:     admin,
		IsPublic:    public && !admin,
		QueryParams: ctx.QueryParams(),
	}

	if ctxUser != nil {
		ws.UserID = &ctxUser.ID
	}

	return ws
}

// Actions --------------------------------------------------------------------

func (c *ControllerBasic) Index(ctx echo.Context) error {
	ws := c.scope(ctx)
	indexResponse, err := c.Models.Adhkars.GetAll(ctx, ws)
	if err != nil {
		return c.APIErr.Database(ctx, err, nil)
	}
	return ctx.JSON(http.StatusOK, indexResponse)
}

func (c *ControllerBasic) Show(ctx echo.Context) error {
	var result adhkars.Model
	if err := c.Utils.ReadUUIDParam(&result.ID, ctx); err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}
	ws := c.scope(ctx)
	if err := c.Models.Adhkars.GetOne(&result, ws); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}
	return ctx.JSON(http.StatusOK, result)
}

func (c *ControllerBasic) Store(ctx echo.Context) error {
	var result adhkars.Model

	v, err := c.GetValidator(ctx, result.ModelName())
	if err != nil {
		return err
	}

	if valid := result.MergeAndValidate(v); !valid {
		return c.APIErr.InputValidation(ctx, v)
	}
	ws := c.scope(ctx)

	if ws.UserID != nil {
		result.CreatedByID = *ws.UserID
	}
	tx, err := c.Models.DB.Beginx()
	if err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if err := c.Models.Adhkars.CreateOne(&result, tx); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}

	if err = tx.Commit(); err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	return ctx.JSON(http.StatusCreated, result)
}

func (c *ControllerBasic) Update(ctx echo.Context) error {
	var result adhkars.Model
	if err := c.Utils.ReadUUIDParam(&result.ID, ctx); err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}
	ws := c.scope(ctx)
	if err := c.Models.Adhkars.GetOne(&result, ws); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}

	v, err := c.GetValidator(ctx, result.ModelName())
	if err != nil {
		return err
	}

	if valid := result.MergeAndValidate(v); !valid {
		return c.APIErr.InputValidation(ctx, v)
	}

	tx, err := c.Models.DB.Beginx()
	if err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if err := c.Models.Adhkars.UpdateOne(&result, ws, tx); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return c.APIErr.Database(
				ctx,
				errors.New(v.T.ConflictError()),
				&result,
			)
		default:
			return c.APIErr.Database(ctx, err, &result)
		}
	}

	if err := tx.Commit(); err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	return ctx.JSON(http.StatusOK, result)
}

func (c *ControllerBasic) Destroy(ctx echo.Context) error {
	var result adhkars.Model
	if err := c.Utils.ReadUUIDParam(&result.ID, ctx); err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}
	ws := c.scope(ctx)

	tx, err := c.Models.DB.Beginx()
	if err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if err := c.Models.Adhkars.DeleteOne(&result, ws, tx); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}

	if err := tx.Commit(); err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	return ctx.JSON(http.StatusOK, result)
}
