package category_controller

import (
	"net/http"

	"bitbucket.org/sadeemTechnology/backend-model-category"
	"github.com/labstack/echo/v4"
)

func (c *Controllers) Store(ctx echo.Context) error {
	var result category.Model
	v, err := c.GetValidator(ctx, result.ModelName())
	if err != nil {
		return err
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

	if err := c.Models.Category.CreateOne(&result, tx); err != nil {
		defer v.DeleteNewPicture()
		return c.APIErr.Database(ctx, err, &result)
	}
	if err = tx.Commit(); err != nil {
		defer v.DeleteNewPicture()
		return c.APIErr.InternalServer(ctx, err)
	}
	return ctx.JSON(http.StatusCreated, result)
}
