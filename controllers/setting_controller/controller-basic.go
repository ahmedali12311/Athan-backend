package setting_controller

import (
	"database/sql"
	"errors"
	"net/http"
	"slices"

	"bitbucket.org/sadeemTechnology/backend-model-setting"

	"github.com/labstack/echo/v4"
)

// Actions --------------------------------------------------------------------

func (c *Controllers) Index(ctx echo.Context) error {
	ws := c.scope(ctx)
	indexResponse, err := c.Models.Setting.GetAll(ctx, ws)
	if err != nil {
		return c.APIErr.Database(ctx, err, nil)
	}
	return ctx.JSON(http.StatusOK, indexResponse)
}

func (c *Controllers) Show(ctx echo.Context) error {
	var result setting.Model
	if err := c.Utils.ReadIDParam(&result.ID, ctx); err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}
	ws := c.scope(ctx)
	if err := c.Models.Setting.GetOne(&result, ws); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}
	return ctx.JSON(http.StatusOK, result)
}

func (c *Controllers) Store(ctx echo.Context) error {
	var result setting.Model
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

	if err := c.Models.Setting.CreateOne(&result, tx); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}
	if err = tx.Commit(); err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	return ctx.JSON(http.StatusCreated, result)
}

func (c *Controllers) Update(ctx echo.Context) error {
	var result setting.Model
	if err := c.Utils.ReadIDParam(&result.ID, ctx); err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}
	if err := c.Models.Setting.GetOne(&result, nil); err != nil {
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

	ws := c.scope(ctx)
	if err := c.Models.Setting.UpdateOne(&result, ws, tx); err != nil {
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

func (c *Controllers) Destroy(ctx echo.Context) error {
	var result setting.Model
	if err := c.Utils.ReadIDParam(&result.ID, ctx); err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}
	ws := c.scope(ctx)
	if err := c.Models.Setting.GetOne(&result, ws); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}
	if slices.Contains(setting.CoreKeys, result.Key) {
		err := errors.New("not allowed to delete a core settings key")
		return c.APIErr.BadRequest(ctx, err)
	}
	// Start transacting
	tx, err := c.Models.DB.Beginx()
	if err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	defer func() { _ = tx.Rollback() }()

	if err := c.Models.Setting.DeleteOne(&result, ws, tx); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}
	if err := tx.Commit(); err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	return ctx.JSON(http.StatusOK, result)
}
