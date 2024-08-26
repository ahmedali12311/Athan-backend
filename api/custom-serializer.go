package api

import (
	"fmt"
	"net/http"

	"github.com/goccy/go-json"
	"github.com/labstack/echo/v4"
)

// CustomJSONSerializer implements JSON encoding using encoding/json.
type CustomJSONSerializer struct{}

// Serialize converts an interface into a json and writes it to the response.
// You can optionally use the indent parameter to produce pretty JSONs.
func (d CustomJSONSerializer) Serialize(
	c echo.Context,
	i any,
	indent string,
) error {
	enc := json.NewEncoder(c.Response())
	if indent != "" {
		enc.SetIndent("", indent)
	}
	return enc.Encode(i)
}

// Deserialize reads a JSON from a request body and converts it into an
// interface.
func (d CustomJSONSerializer) Deserialize(c echo.Context, i any) error {
	err := json.NewDecoder(c.Request().Body).Decode(i)

	ute, ok := err.(*json.UnmarshalTypeError) //nolint: errorlint
	if ok {
		msg := fmt.Sprintf(
			"Unmarshal type error: expected=%v, got=%v, field=%v, offset=%v",
			ute.Type,
			ute.Value,
			ute.Field,
			ute.Offset,
		)
		return echo.
			NewHTTPError(http.StatusBadRequest, msg).
			SetInternal(err)
	}

	se, ok := err.(*json.SyntaxError) //nolint: errorlint
	if ok {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			fmt.Sprintf(
				"Syntax error: offset=%v, error=%v",
				se.Offset,
				se.Error(),
			)).SetInternal(err)
	}
	return err
}
