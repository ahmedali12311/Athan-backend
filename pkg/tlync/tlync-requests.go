package tlync

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	setting "bitbucket.org/sadeemTechnology/backend-model-setting"
	"github.com/goccy/go-json"
)

func InitiatePayment(
	settings *setting.TLync,
	input *InitiateInput,
) (*InitiateResponse, error) {
	endpoint := settings.Endpoint + "/payment/initiate"
	token := "Bearer " + settings.Token
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
	req.Header.Add("authorization", token)
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
	var result InitiateResponse
	if err := json.Unmarshal(bd, &result); err != nil {
		return nil, fmt.Errorf("can not unmarshal JSON: %w", err)
	}
	if result.Result == "" {
		return nil, errors.New("empty result: " + result.Message)
	}
	return &result, nil
}

func TransactionReceipt(
	settings *setting.TLync,
	input *ConfirmInput,
) (*ConfirmResponse, error) {
	endpoint := settings.Endpoint + "/receipt/transaction"
	token := "Bearer " + settings.Token
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
	req.Header.Add("authorization", token)
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
	var result ConfirmResponse
	if err := json.Unmarshal(bd, &result); err != nil {
		return nil, fmt.Errorf("can not unmarshal JSON: %w", err)
	}
	return &result, nil
}
