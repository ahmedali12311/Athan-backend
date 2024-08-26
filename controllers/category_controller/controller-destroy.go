package category_controller

import (
	"errors"
	"net/http"
	"testing"

	"app/models/category"
	"app/models/user"
	"app/test_utils"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
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

// Test Cases -----------------------------------------------------------------

var testDestroy = []*test_utils.Test{
	{
		Name:        "categories/destroy:admin",
		ContentType: "application/json",
		Path:        "/api/v3/categories/" + TestID.String(),
		Method:      http.MethodDelete,
		Code:        http.StatusOK,
		Body:        http.NoBody,
		TokenUserID: user.SuperAdminID,
		Cases: func(t *testing.T, body string) {
			var deleted category.Model
			test_utils.DecodeJSONString(t, &deleted, body)

			require.Equalf(
				t,
				"Test Category Updated",
				deleted.Name,
				"admin categories deleted category name matches updated",
			)

			require.Equalf(
				t,
				1,
				deleted.Depth,
				"admin categories destroyed category depth = 1",
			)
			require.Equalf(
				t,
				0,
				deleted.Sort,
				"admin categories destroyed category sort = 0",
			)
			require.Equalf(
				t,
				category.CitySuperParent,
				deleted.Parent.ID.String(),
				"admin categories destroyed category parent.id is valid",
			)
		},
	},
}
