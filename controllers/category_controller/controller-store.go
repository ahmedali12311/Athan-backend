package category_controller

import (
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

// Test Cases -----------------------------------------------------------------

func testStore() []*test_utils.Test {
	testPng := config.GetRootPath("test_utils/test.png")
	form := map[string]string{
		"id":          TestID.String(),
		"img":         fmt.Sprintf("@%s", testPng),
		"name":        "Test Category",
		"is_disabled": "f",
		"is_featured": "t",
		"parent":      fmt.Sprintf("{%q:%q}", "id", category.CitySuperParent),
	}
	contentType, formData := test_utils.CreateForm(form)

	return []*test_utils.Test{
		{
			Name:        "categories/store:admin",
			ContentType: contentType,
			Path:        "/api/v3/categories",
			Method:      http.MethodPost,
			Code:        http.StatusCreated,
			Body:        formData,
			Params:      map[string]string{},
			TokenUserID: user.SuperAdminID,
			Cases: func(t *testing.T, res string) {
				var created category.Model
				test_utils.DecodeJSONString(t, &created, res)

				createdImgBasename := path.Base(*created.Img)
				createdThumbBasename := path.Base(*created.Thumb)

				createdImgPath := config.GetUploadsPath(
					fmt.Sprintf("categories/%s", createdImgBasename),
				)
				createdThumbPath := config.GetUploadsPath(
					fmt.Sprintf("categories/thumbs/%s", createdThumbBasename),
				)

				createdImgInfo, err := os.Stat(createdImgPath)
				if err != nil {
					t.Fatal(err)
				}

				createdThumbInfo, err := os.Stat(createdThumbPath)
				if err != nil {
					t.Fatal(err)
				}
				require.Equalf(
					t,
					"Test Category",
					created.Name,
					"admin categories store created name matches input",
				)
				require.Equalf(
					t,
					createdImgBasename,
					createdImgInfo.Name(),
					"admin categories created img",
				)
				require.Equalf(
					t,
					createdThumbBasename,
					createdThumbInfo.Name(),
					"admin categories created thumb",
				)

				require.Equalf(
					t,
					1,
					created.Depth,
					"admin categories store created depth = 1",
				)
				require.Equalf(
					t,
					0,
					created.Sort,
					"admin categories store created sort = 0",
				)
				require.Equalf(
					t,
					"e1cf1e12-3f2a-429e-af3b-a0b2ad092c13",
					created.Parent.ID.String(),
					"admin categories store created parent.id is valid",
				)
			},
		},
	}
}
