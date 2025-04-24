package payment_gateway

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/goccy/go-json"
)

func MasaratInitiatePayment(
	settings *Settings,
	input *MasaratInitiateRequest,
) (*Response, error) {
	endpoint := settings.Endpoint + "/payment-gateways/masarat/initiate"
	jsonInput, err := json.Marshal(*input)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		endpoint,
		bytes.NewBuffer(jsonInput),
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
	var result Response

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(string(bd))
	}

	if err := json.Unmarshal(bd, &result); err != nil {
		return nil, fmt.Errorf("can not unmarshal JSON: %w", err)
	}

	return &result, nil
}

func MasaratTransactionConfirm(
	settings *Settings,
	input *MasaratConfirmRequest,
) (*Response, error) {
	endpoint := settings.Endpoint + "/payment-gateways/masarat/" +
		input.WalletTransactionID.String() + "/confirm"

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		endpoint,
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-API-Key", settings.APIKey)
	req.Header.Add("accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http.DefaultClient.Do: %w", err)
	}
	defer res.Body.Close()

	bd, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(string(bd))
	}

	var result Response
	if err := json.Unmarshal(bd, &result); err != nil {
		return nil, fmt.Errorf("can not unmarshal JSON: %w", err)
	}
	return &result, nil
}
