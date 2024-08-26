package test_utils

import (
	"bytes"
	"context"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"app/models/user"

	"github.com/google/uuid"
)

type IndexTest struct {
	Params        map[string]string
	CanBeDisabled bool
}

type Test struct {
	Name        string
	ContentType string
	Method      string
	Path        string
	Body        io.Reader
	Code        int
	Params      map[string]string
	// TokenUserID used to generate a TokenValue if provided
	TokenUserID string
	// TokenValue is assigned by the user id if provided, do not assign this
	// value manually
	TokenValue string
	Cases      func(*testing.T, string)
}

// CreateImage for form data file upload mocking
func CreateImage(name string) *image.RGBA {
	width := 512
	height := 256

	upLeft := image.Point{}
	lowRight := image.Point{X: width, Y: height}

	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})

	cyan := color.RGBA{R: 100, G: 200, B: 200, A: 0xff}

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			switch {
			case x < width/2 && y < height/2: // upper left quadrant
				img.Set(x, y, cyan)
			case x >= width/2 && y >= height/2: // lower right quadrant
				img.Set(x, y, color.White)
			default:
				// Use zero value.
			}
		}
	}

	f, err := os.Create(name) //nolint: gosec // dont worry
	if err != nil {
		return nil
	}
	if err := png.Encode(f, img); err != nil {
		return nil
	}
	return img
}

// SetParams adds key, value params to a request
func SetParams(req *http.Request, params map[string]string) {
	if params != nil {
		q := req.URL.Query()
		for k, v := range params {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}
}

// DecodeJSONString decodes response body into a struct
func DecodeJSONString(t *testing.T, result any, body string) {
	x := json.NewDecoder(strings.NewReader(body))
	if err := x.Decode(&result); err != nil {
		t.Fatalf(
			"result didn't decode into type: %s, error: %s",
			reflect.TypeOf(result),
			err.Error(),
		)
	}
}

// GetToken returns a token string generated for the user
// created in the seed with admin role
func GetToken(userID string, t *testing.T, userQueries *user.Queries) string {
	if userID != "" {
		usr := &user.Model{
			ID: uuid.MustParse(userID),
		}
		if err := userQueries.GetOne(usr); err != nil {
			t.Fatal(err)
		}
		tokenResponse, err := usr.GenTokenResponse()
		if err != nil {
			t.Fatal(err)
		}
		return tokenResponse.Token.Value
	}
	return ""
}

// CreateForm builds a multipart/form-data with file input and string fields
func CreateForm(form map[string]string) (string, io.Reader) {
	body := new(bytes.Buffer)
	fd := multipart.NewWriter(body)

	defer fd.Close()

	for key, val := range form {
		if strings.HasPrefix(val, "@") {
			val = val[1:]
			img := CreateImage(val)

			part, err := fd.CreateFormFile(key, "test.png")
			if err != nil {
				panic(err.Error())
			}

			if err := png.Encode(part, img); err != nil {
				panic(err.Error())
			}

		} else {
			if err := fd.WriteField(key, val); err != nil {
				panic(err.Error())
			}
		}
	}
	return fd.FormDataContentType(), body
}

// TestRequest tests a single request, returns response and the body in a string
func TestRequest(
	t *testing.T,
	ts *httptest.Server,
	tc *Test,
) (*http.Response, string) {
	req, err := http.NewRequestWithContext(
		context.Background(),
		tc.Method,
		ts.URL+tc.Path,
		tc.Body,
	)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	req.Header.Set("Content-Type", tc.ContentType)
	if tc.TokenValue != "" {
		req.Header.Set("Authorization", tc.TokenValue)
	}

	SetParams(req, tc.Params)

	resp, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			return
		}
	}()

	return resp, string(respBody)
}

// // ModelGetRequestsTest performs a test for all index, show requests in the
// // provided board accepts board and token to authorize based on board
// func ModelGetRequestsTest[T finder.Model](
// 	t *testing.T,
// 	routes http.Handler,
// 	test *IndexTest,
// 	scope, token string,
// ) {
// 	target := &finder.IndexResponse[T]{}
// 	var m T
// 	table := strings.ReplaceAll(m.TableName(), "_", "-")
// 	name := strings.ReplaceAll(m.ModelName(), "_", "-")
//
// 	requestTest := &Test{
// 		Method:      "GET",
// 		Path:        fmt.Sprintf("/api/v1/%s", table),
// 		ContentType: "application/json",
// 		Params:      test.Params,
// 		Body:        http.NoBody,
// 	}
//
// 	ts := httptest.NewServer(routes)
// 	defer ts.Close()
//
// 	response, body := TestRequest(t, ts, requestTest)
// 	defer response.Body.Close()
// 	DecodeJSONString(t, target, body)
// 	t.Log(target.Data)
// 	results := *target.Data
//
// 	assert.Equalf(
// 		t,
// 		http.StatusOK,
// 		response.StatusCode,
// 		fmt.Sprintf("%s/%s index is ok", scope, table),
// 	)
//
// 	assert.LessOrEqual(
// 		t,
// 		uint64(len(results)),
// 		target.Meta.Paginate,
// 		fmt.Sprintf("%s/%s index within pagination limit", scope, table),
// 	)
//
// 	if len(results) > 0 {
//
// 		var modelMap map[string]interface{}
// 		data, err := json.Marshal(results[0])
// 		if err != nil {
// 			t.Fatal(err)
// 		}
//
// 		if err := json.Unmarshal(data, &modelMap); err != nil {
// 			t.Fatal(err)
// 		}
//
// 		// assert that a public result have is_disabled=false
// 		if test.CanBeDisabled && scope == "public" {
// 			isDisabled := modelMap["is_disabled"]
//
// 			assert.Equalf(
// 				t,
// 				false,
// 				isDisabled,
// 				fmt.Sprintf(
// 					"/%s for scope %s results must not be disabled",
// 					scope,
// 					name,
// 				),
// 			)
// 		}
//
// 		// Test show for the first result if it exists
// 		idType := reflect.TypeOf(modelMap["id"]).Kind()
//
// 		var ID string
//
// 		if idType == reflect.String {
// 			ID = modelMap["id"].(string)
// 		} else {
// 			// its a float64
// 			floatID := modelMap["id"].(float64)
// 			intID := int64(floatID)
// 			ID = fmt.Sprintf("%d", intID)
// 		}
//
// 		var pathString string
//
// 		if reflect.TypeOf(modelMap["id"]).Kind() == reflect.String {
// 			pathString = fmt.Sprintf("/api/v1/%s/%s", table, ID)
// 		} else {
// 			intID, err := strconv.ParseUint(ID, 10, 64)
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			ID = strconv.FormatUint(intID, 10)
// 			pathString = fmt.Sprintf("/api/v1/%s/%s", table, ID)
// 		}
//
// 		requestShowTest := &Test{
// 			Method:      "GET",
// 			Path:        pathString,
// 			Token:       token,
// 			ContentType: "application/json",
// 			Params:      test.Params,
// 			Body:        http.NoBody,
// 		}
//
// 		response, body := TestRequest(t, ts, requestShowTest)
// 		DecodeJSONString(t, &target, body)
// 		defer response.Body.Close()
//
// 		assert.Equalf(
// 			t,
// 			http.StatusOK,
// 			response.StatusCode,
// 			pathString,
// 		)
// 	}
// }
