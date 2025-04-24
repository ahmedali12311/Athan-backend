package payment_gateway

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/goccy/go-json"
)

func PaymentService(
	settings *Settings,
) (*PaymentServicesResponse, error) {
	endpoint := settings.Endpoint + "/payment-services"
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
