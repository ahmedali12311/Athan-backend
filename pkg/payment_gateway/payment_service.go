package payment_gateway

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"bitbucket.org/sadeemTechnology/backend-model-setting"
	"github.com/goccy/go-json"
	"github.com/labstack/echo/v4"
)

func PaymentService(
	settings *setting.TyrianAnt,
	ctx echo.Context,
) (*PaymentServicesResponse, error) {
	query := url.Values{}

	if q := ctx.QueryParam("q"); q != "" {
		query.Set("q", q)
	}
	if filters := ctx.QueryParam("filters"); filters != "" {
		query.Set("filters", filters)
	}
	if sorts := ctx.QueryParam("sorts"); sorts != "" {
		query.Set("sorts", sorts)
	}

	query.Add("all", "")

	endpoint := settings.Endpoint + "/payment-services?" + query.Encode()
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		endpoint,
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-API-Key", settings.APIKey)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http.DefaultClient.Do: %w", err)
	}
	defer res.Body.Close()

	bd, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}
	var result PaymentServicesResponse

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(string(bd))
	}

	if err := json.Unmarshal(bd, &result); err != nil {
		return nil, fmt.Errorf("can not unmarshal JSON: %w", err)
	}

	return &result, nil
}
