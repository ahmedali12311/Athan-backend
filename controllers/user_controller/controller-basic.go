package user_controller

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"app/controller"
	"app/models/user"
	"app/pkg/gis"

	"github.com/labstack/echo/v4"
)

type ControllerBasic struct {
	*controller.Dependencies
}

// Scopes ---------------------------------------------------------------------

// func (c *ControllerBasic) userScope(ctx echo.Context) *uuid.UUID {
// 	scopes := c.Utils.CtxScopes(ctx)
// 	if slices.Contains(scopes, "admin") {
// 		return nil
// 	}
// 	return &c.Utils.CtxUser(ctx).ID
// }

// Actions --------------------------------------------------------------------

func (c *ControllerBasic) Index(ctx echo.Context) error {
	indexResponse, err := c.Models.User.GetAll(ctx)
	if err != nil {
		return c.APIErr.Database(ctx, err, nil)
	}
	return ctx.JSON(http.StatusOK, indexResponse)
}

func (c *ControllerBasic) Show(ctx echo.Context) error {
	var result user.Model
	if err := c.Utils.ReadUUIDParam(&result.ID, ctx); err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}
	if err := c.Models.User.GetOne(&result); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}
	return ctx.JSON(http.StatusOK, result)
}

func (c *ControllerBasic) Store(ctx echo.Context) error {
	result := user.Model{
		Location:  gis.EmptyPoint,
		CreatedAt: time.Time{},
	}
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

	if err := c.Models.User.CreateOne(&result, tx); err != nil {
		defer v.DeleteNewPicture()
		return c.APIErr.Database(ctx, err, &result)
	}
	if err = tx.Commit(); err != nil {
		defer v.DeleteNewPicture()
		return c.APIErr.InternalServer(ctx, err)
	}
	if err := c.Models.User.GetRoles(&result); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}

	return ctx.JSON(http.StatusCreated, result)
}

func (c *ControllerBasic) Update(ctx echo.Context) error {
	var result user.Model
	v, err := c.GetValidator(ctx, result.ModelName())
	if err != nil {
		return err
	}
	if err := c.Utils.ReadUUIDParam(&result.ID, ctx); err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}
	if err := c.Models.User.GetOne(&result); err != nil {
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

	if err := c.Models.User.UpdateOne(&result, tx); err != nil {
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
	if err := c.Models.User.GetRoles(&result); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}

	return ctx.JSON(http.StatusOK, result)
}

func (c *ControllerBasic) Clear(ctx echo.Context) error {
	var result user.Model
	if err := c.Utils.ReadUUIDParam(&result.ID, ctx); err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}
	if err := c.Models.User.GetOne(&result); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}

	// Start transacting
	tx, err := c.Models.DB.Beginx()
	if err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	defer func() { _ = tx.Rollback() }()

	if err := c.Models.User.ClearOne(&result.ID, tx); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}

	if err := tx.Commit(); err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	// Delete only if commit succeeds
	v, err := c.GetValidator(ctx, result.ModelName())
	if err != nil {
		return err
	}
	v.SaveOldImgThumbDists(&result)
	v.DeleteOldPicture()

	return ctx.JSON(http.StatusOK, result)
}
