package category_controller

import (
	"net/http"
	"strings"
	"testing"

	"app/models/category"
	"app/models/user"
	"app/test_utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var TestID = uuid.New()

func TestCases() []*test_utils.Test {
	cases := []*test_utils.Test{
		{
			Name:        "categories/index:admin",
			ContentType: "application/json",
			Path:        "/api/v1/categories",
			Method:      http.MethodGet,
			Body:        http.NoBody,
			Code:        http.StatusOK,
			TokenUserID: user.SuperAdminID,
			Params: map[string]string{
				"page":     "1",
				"paginate": "3",
				"sorts":    "-created_at",
			},
			Cases: func(t *testing.T, res string) {
				require.NotNil(t, res)
				require.True(t, strings.Contains(res, "\"meta\""))
			},
		},
		{
			Name:        "categories/index:public",
			ContentType: "application/json",
			Path:        "/api/v1/categories",
			Method:      http.MethodGet,
			Body:        http.NoBody,
			Code:        http.StatusOK,
			TokenUserID: user.SuperAdminID,
			Params: map[string]string{
				"page":     "1",
				"paginate": "3",
				"sorts":    "-created_at",
			},
		},
		{
			Name:        "categories/show",
			ContentType: "application/json",
			Path:        "/api/v1/categories/" + category.CitySuperParent,
			Method:      http.MethodGet,
			Body:        http.NoBody,
			Code:        http.StatusOK,
			TokenUserID: user.SuperAdminID,
			Params:      map[string]string{},
		},
	}

	cases = append(cases, testStore()...)
	cases = append(cases, testUpdate()...)
	cases = append(cases, testDestroy...)
	return cases
}
