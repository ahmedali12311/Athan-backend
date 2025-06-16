package category_controller

import (
	"database/sql"
	"errors"
	"net/http"

	"bitbucket.org/sadeemTechnology/backend-model-category"
	"github.com/labstack/echo/v4"
)

func (c *Controllers) Update(ctx echo.Context) error {
	var result category.Model
	v, err := c.GetValidator(ctx, result.ModelName())
	if err != nil {
		return err
	}
	if err := c.Utils.ReadUUIDParam(&result.ID, ctx); err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}
	ws := c.scope(ctx)
	if err := c.Models.Category.GetOne(&result, ws); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}
	if valid := result.MergeAndValidate(v); !valid {
		defer v.DeleteNewPicture()
		return c.APIErr.InputValidation(ctx, v)
	}
	// Start transacting
	tx, err := c.Models.DB.Beginx()
	if err != nil {
		defer v.DeleteNewPicture()
		return c.APIErr.InternalServer(ctx, err)
	}
	defer func() { _ = tx.Rollback() }()

	if err := c.Models.Category.UpdateOne(&result, ws, tx); err != nil {
		defer v.DeleteNewPicture()
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
		defer v.DeleteNewPicture()
		return c.APIErr.InternalServer(ctx, err)
	}
	return ctx.JSON(http.StatusOK, result)
}
