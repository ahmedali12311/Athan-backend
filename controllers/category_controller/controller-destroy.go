package category_controller

import (
	"errors"
	"net/http"

	"app/models/category"

	"github.com/labstack/echo/v4"
)

func (c *Controllers) Destroy(ctx echo.Context) error {
	var result category.Model
	if err := c.Utils.ReadUUIDParam(&result.ID, ctx); err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}
	ws := c.scope(ctx)
	if err := c.Models.Category.GetOne(&result, ws); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}
	hasChildren, err := c.Models.Category.HasChildren(&result)
	if err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}
	v, err := c.GetValidator(ctx, result.ModelName())
	if err != nil {
		return err
	}
	if hasChildren {
		err := errors.New(v.T.UnDestroyableCategory())
		return c.APIErr.Forbidden(ctx, err)
	}

	// Start transacting
	tx, err := c.Models.DB.Beginx()
	if err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	defer func() { _ = tx.Rollback() }()

	if err := c.Models.Category.DeleteOne(&result, ws, tx); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}

	if err := tx.Commit(); err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	v.SaveOldImgThumbDists(&result)
	v.DeleteOldPicture()

	return ctx.JSON(http.StatusOK, result)
}
