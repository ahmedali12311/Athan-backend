package otp

import (
	"app/pkg/requester"
	"fmt"
	"net/http"

	"github.com/goccy/go-json"
)

func Request(result *Response, settings *Settings, input *Input) error {
	path := settings.URL + EndpointV1Pins + "?key=" + settings.Key
	if settings.Env != "production" {
		path += "&test"
	}
	s := &requester.Settings{
		Method: http.MethodPost,
		URL:    path,
		Headers: map[string]string{
			"Authorization": Bearer + settings.JWT,
			"Content-Type":  "application/json",
		},
	}
	statusCode, body, err := requester.Request(s, input)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("can not unmarshal JSON: %w", err)
	}
	result.Status = statusCode
	if statusCode != http.StatusCreated {
		return fmt.Errorf(
			"service returned %d, %s, %s",
			statusCode,
			result.Type,
			result.Message,
		)
	}
	return nil
}
