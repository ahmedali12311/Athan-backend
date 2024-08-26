package category_controller

// import (
// 	"fmt"
// 	"github.com/stretchr/testify/assert"
// 	"indigo-spirograph/cmd/api/category"
// 	"indigo-spirograph/cmd/config"
// 	"indigo-spirograph/cmd/test_utils"
// 	"net/http"
// 	"net/http/httptest"
// 	"os"
// 	"path"
// 	"testing"
// )
//
// func TestAdminCategoriesCreateUpdateDelete(t *testing.T) {
// 	route := "/api/v2/admin/categories"
// 	ts := httptest.NewServer(app.routes(true))
// 	defer ts.Close()
//
// 	token := test_utils.GetToken("admin", t, app.db)
//
// 	// --------------------------------------------------------------
// 	// 		Create new category
// 	// --------------------------------------------------------------
// 	var created category.Model
//
// 	form := map[string]string{
// 		"img":         "@../test_utils/test.png",
// 		"name":        "Test Category",
// 		"is_disabled": "f",
// 		"is_featured": "t",
// 		"parent_id":   "e1cf1e12-3f2a-429e-af3b-a0b2ad092c13", // name = items, depth = 1
// 	}
//
// 	contentType, formData, err := test_utils.CreateForm(form)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	requestTest := &test_utils.RequestTest{
// 		Path:        route,
// 		Method:      "POST",
// 		Token:       token,
// 		ContentType: contentType,
// 		Params:      nil,
// 		Body:        formData,
// 	}
//
// 	postResponse, postBody := test_utils.TestRequest(t, ts, requestTest)
// 	test_utils.DecodeJSONString(t, &created, postBody)
//
// 	createdImgBasename := path.Base(*created.Img)
// 	createdThumbBasename := path.Base(*created.Thumb)
//
// 	createdImgPath := config.GetUploadsPath(fmt.Sprintf("categories/%s", createdImgBasename))
// 	createdThumbPath := config.GetUploadsPath(fmt.Sprintf("categories/thumbs/%s", createdThumbBasename))
//
// 	createdImgInfo, err := os.Stat(createdImgPath)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	createdThumbInfo, err := os.Stat(createdThumbPath)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	assert.Equalf(t, http.StatusCreated, postResponse.StatusCode, "admin categories store is ok with token")
// 	assert.Equalf(t, "Test Category", created.Name, "admin categories store created name matches input")
// 	assert.Equalf(t, createdImgBasename, createdImgInfo.Name(), "admin categories created img")
// 	assert.Equalf(t, createdThumbBasename, createdThumbInfo.Name(), "admin categories created thumb")
// 	assert.Equalf(t, uint8(1), created.Depth, "admin categories store created depth = 1")
// 	assert.Equalf(t, uint16(0), created.Sort, "admin categories store created sort = 0")
// 	assert.Equalf(
// 		t,
// 		"e1cf1e12-3f2a-429e-af3b-a0b2ad092c13",
// 		created.Parent.ID.String(),
// 		"admin categories store created parent.id is valid",
// 	)
//
// 	// set path to created category id
// 	routeID := fmt.Sprintf("%s/%s", route, created.ID.String())
//
// 	// --------------------------------------------------------------
// 	// 		Update created category
// 	// --------------------------------------------------------------
//
// 	var updated category.Model
//
// 	formUpdate := map[string]string{
// 		"img":       "@../test_utils/test.png",
// 		"name":      "Test Category Updated",
// 		"parent_id": "7ba362b8-cfe7-425b-8882-9578151dff47", // name = sandwiches, depth = 2
// 	}
//
// 	contentTypeUpdate, formDataUpdate, errUpdate := test_utils.CreateForm(formUpdate)
// 	if errUpdate != nil {
// 		t.Fatal(err)
// 	}
//
// 	requestUpdateTest := &test_utils.RequestTest{
// 		Path:        routeID,
// 		Method:      "PUT",
// 		Token:       token,
// 		ContentType: contentTypeUpdate,
// 		Params:      nil,
// 		Body:        formDataUpdate,
// 	}
//
// 	updateResponse, updateBody := test_utils.TestRequest(t, ts, requestUpdateTest)
// 	test_utils.DecodeJSONString(t, &updated, updateBody)
//
// 	updatedImgBasename := path.Base(*updated.Img)
// 	updatedThumbBasename := path.Base(*updated.Thumb)
// 	updatedImgPath := config.GetUploadsPath(fmt.Sprintf("categories/%s", updatedImgBasename))
// 	updatedThumbPath := config.GetUploadsPath(fmt.Sprintf("categories/thumbs/%s", updatedThumbBasename))
//
// 	t.Log(updated)
// 	t.Log(updateBody)
// 	t.Log(updatedImgPath)
// 	t.Log(updatedThumbPath)
//
// 	// ensure that old files are removed
// 	_, errCreatedImgStillExist := os.Stat(createdImgPath)
//
// 	_, errCreatedThumbStillExist := os.Stat(createdThumbPath)
//
// 	// ensure new files are saved
// 	updatedImgInfo, err := os.Stat(updatedImgPath)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	updatedThumbInfo, err := os.Stat(updatedThumbPath)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	assert.Equalf(t, http.StatusOK, updateResponse.StatusCode, "admin categories store is ok with token")
// 	assert.Equalf(t, "Test Category Updated", updated.Name, "admin categories store updated name matches input")
// 	assert.Equalf(t, true, errCreatedImgStillExist != nil, "admin categories removed old img")
// 	assert.Equalf(t, true, errCreatedThumbStillExist != nil, "admin categories removed old thumb")
// 	assert.Equalf(t, updatedImgBasename, updatedImgInfo.Name(), "admin categories updated img")
// 	assert.Equalf(t, updatedThumbBasename, updatedThumbInfo.Name(), "admin categories updated thumb")
// 	assert.Equalf(t, uint8(2), updated.Depth, "admin categories store updated depth = 1")
// 	assert.Equalf(
// 		t,
// 		"7ba362b8-cfe7-425b-8882-9578151dff47",
// 		updated.Parent.ID.String(),
// 		"admin categories put parent.id is sandwiches",
// 	)
// 	assert.Equalf(
// 		t,
// 		"e1cf1e12-3f2a-429e-af3b-a0b2ad092c13",
// 		updated.SuperParent.ID.String(),
// 		"admin categories put super_parent.id is items",
// 	)
//
// 	//// --------------------------------------------------------------
// 	//// 		Delete the updated category
// 	//// --------------------------------------------------------------
//
// 	var deleted category.Model
//
// 	requestDeleteTest := &test_utils.RequestTest{
// 		Path:        routeID,
// 		Method:      "DELETE",
// 		Token:       token,
// 		ContentType: "application/json",
// 		Params:      nil,
// 		Body:        nil,
// 	}
// 	deleteResponse, deleteBody := test_utils.TestRequest(t, ts, requestDeleteTest)
// 	test_utils.DecodeJSONString(t, &deleted, deleteBody)
//
// 	// ensure that old files are removed
// 	_, errUpdatedImgStillExist := os.Stat(updatedImgPath)
// 	_, errUpdatedThumbStillExist := os.Stat(updatedThumbPath)
//
// 	assert.Equalf(t, http.StatusOK, deleteResponse.StatusCode, "admin categories delete is ok with token")
// 	assert.Equalf(t, "Test Category Updated", deleted.Name, "admin categories deleted name returns")
// 	assert.Equalf(t, true, errUpdatedImgStillExist != nil, "admin categories removed update img")
// 	assert.Equalf(t, true, errUpdatedThumbStillExist != nil, "admin categories removed update thumb")
// }
