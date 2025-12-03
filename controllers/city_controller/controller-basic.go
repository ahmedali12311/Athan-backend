package city_controller

import (
	"database/sql"
	"errors"
	"net/http"

	"app/controller"
	"app/models/city"

	"github.com/labstack/echo/v4"
)

type ControllerBasic struct {
	*controller.Dependencies
}

// Scopes ---------------------------------------------------------------------

func (c *ControllerBasic) scope(ctx echo.Context) *city.WhereScope {
	scopes := c.Utils.CtxScopes(ctx)

	var admin, marketer bool
	for _, v := range scopes {
		switch v {
		case "admin":
			admin = true
		}
	}
	ws := &city.WhereScope{
		IsPublic:   !admin,
		IsMarketer: marketer,
	}
	return ws
}

func (c *ControllerBasic) Index(ctx echo.Context) error {
	ws := c.scope(ctx)
	indexResponse, err := c.Models.City.GetAll(ctx, ws)
	if err != nil {
		return c.APIErr.Database(ctx, err, nil)
	}
	return ctx.JSON(http.StatusOK, indexResponse)
}

func (c *ControllerBasic) Show(ctx echo.Context) error {
	var result city.Model
	if err := c.Utils.ReadUUIDParam(&result.ID, ctx); err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}
	ws := c.scope(ctx)
	if err := c.Models.City.GetOne(&result, ws); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}
	return ctx.JSON(http.StatusOK, result)
}

func (c *ControllerBasic) Store(ctx echo.Context) error {
	var result city.Model

	v, err := c.GetValidator(ctx, result.ModelName())
	if err != nil {
		return err
	}
	if valid := result.MergeAndValidate(v); !valid {
		return c.APIErr.InputValidation(ctx, v)
	}

	if !v.Valid() {
		return c.APIErr.InputValidation(ctx, v)
	}
	// Start transacting
	tx, err := c.Models.DB.Beginx()
	if err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	defer func() { _ = tx.Rollback() }()

	if err := c.Models.City.CreateOne(&result, tx); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}
	if err = tx.Commit(); err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	return ctx.JSON(http.StatusCreated, result)
}

func (c *ControllerBasic) Destroy(ctx echo.Context) error {
	var result city.Model
	if err := c.Utils.ReadUUIDParam(&result.ID, ctx); err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}
	// Start transacting
	tx, err := c.Models.DB.Beginx()
	if err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	defer func() { _ = tx.Rollback() }()

	ws := c.scope(ctx)

	if err := c.Models.City.DeleteOne(&result, ws, tx); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	} else {
		if err := tx.Commit(); err != nil {
			return c.APIErr.InternalServer(ctx, err)
		}
	}
	return ctx.JSON(http.StatusOK, result)
}

func (c *ControllerBasic) Update(ctx echo.Context) error {
	t := c.Utils.CtxT(ctx)
	var result city.Model

	if err := c.Utils.ReadUUIDParam(&result.ID, ctx); err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}
	ws := c.scope(ctx)
	if err := c.Models.City.GetOne(&result, ws); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}

	v, err := c.GetValidator(ctx, result.ModelName())
	if err != nil {
		return err
	}
	if valid := result.MergeAndValidate(v); !valid {
		return c.APIErr.InputValidation(ctx, v)
	}
	// Start transacting
	tx, err := c.Models.DB.Beginx()
	if err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	defer func() { _ = tx.Rollback() }()

	if err := c.Models.City.UpdateOne(&result, ws, tx); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return c.APIErr.Database(
				ctx,
				errors.New(t.ConflictError()),
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
