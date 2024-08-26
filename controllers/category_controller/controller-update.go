package category_controller

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"testing"

	"app/config"
	"app/models/category"
	"app/models/user"
	"app/test_utils"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
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
	ws.SortBeforeUpdate = result.Sort
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

// Test Cases -----------------------------------------------------------------

func testUpdate() []*test_utils.Test {
	testPng := config.GetRootPath("test_utils/test.png")
	form := map[string]string{
		"id":          TestID.String(),
		"img":         fmt.Sprintf("@%s", testPng),
		"name":        "Test Category Updated",
		"is_disabled": "t",
		"is_featured": "t",
		"parent":      fmt.Sprintf("{%q:%q}", "id", category.CitySuperParent),
	}
	contentType, formData := test_utils.CreateForm(form)

	return []*test_utils.Test{
		{
			Name:        "categories:update-admin-200",
			ContentType: contentType,
			Path:        "/api/v1/categories/" + TestID.String(),
			Method:      http.MethodPut,
			Code:        http.StatusOK,
			Body:        formData,
			Params:      map[string]string{},
			TokenUserID: user.SuperAdminID,
			Cases: func(t *testing.T, res string) {
				var updated category.Model
				test_utils.DecodeJSONString(t, &updated, res)

				updatedImgBasename := path.Base(*updated.Img)
				updatedThumbBasename := path.Base(*updated.Thumb)
				updatedImgPath := config.GetUploadsPath(
					fmt.Sprintf("categories/%s", updatedImgBasename),
				)
				updatedThumbPath := config.GetUploadsPath(
					fmt.Sprintf("categories/thumbs/%s", updatedThumbBasename),
				)

				// ensure new files are saved
				updatedImgInfo, err := os.Stat(updatedImgPath)
				if err != nil {
					t.Fatal(err)
				}

				updatedThumbInfo, err := os.Stat(updatedThumbPath)
				if err != nil {
					t.Fatal(err)
				}

				require.Equalf(
					t,
					"Test Category Updated",
					updated.Name,
					"admin categories updated name matches input",
				)
				require.Equalf(
					t,
					updatedImgBasename,
					updatedImgInfo.Name(),
					"admin categories updated img",
				)
				require.Equalf(
					t,
					updatedThumbBasename,
					updatedThumbInfo.Name(),
					"admin categories updated thumb",
				)

				require.Equalf(
					t,
					category.CitySuperParent,
					updated.Parent.ID.String(),
					"admin categories updated to new parent",
				)
				require.Equalf(
					t,
					category.CitySuperParent,
					updated.SuperParent.ID.String(),
					"admin categories updated to new super_parent",
				)
			},
		},
	}
}
