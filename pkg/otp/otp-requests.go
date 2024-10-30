package otp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/goccy/go-json"
)

func Request(settings *Settings, input *Input) (*Response, error) {
	token := "Bearer " + settings.JWT
	jsonInput, err := json.Marshal(*input)
	if err != nil {
		return nil, err
	}
	path := settings.URL + EndpointV1Pins + "?key=" + settings.Key
	if settings.Env != "production" {
		path += "&test"
	}
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		path,
		bytes.NewBuffer(jsonInput),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Add("authorization", token)
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{
		Timeout: 60 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http client.Do: %w", err)
	}
	defer res.Body.Close()

	bd, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}
	var result Response
	if err := json.Unmarshal(bd, &result); err != nil {
		return nil, fmt.Errorf("can not unmarshal JSON: %w", err)
	}
	if res.StatusCode != http.StatusCreated {
		return &result, fmt.Errorf(
			"service returned %d, %s, %s",
			res.StatusCode,
			result.Type,
			result.Message,
		)
	}
	result.Status = res.StatusCode
	return &result, nil
}
