package requester

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/goccy/go-json"
)

type Settings struct {
	TLSConfig *tls.Config
	Method    string
	URL       string
	Headers   map[string]string
}

func Request(settings *Settings, input any) (int, []byte, error) {
	var body io.Reader

	if input != nil {
		b, err := json.Marshal(input)
		if err != nil {
			return http.StatusServiceUnavailable, nil, err
		}
		body = bytes.NewBuffer(b)
	} else {
		body = http.NoBody
	}

	req, err := http.NewRequestWithContext(
		context.Background(),
		settings.Method,
		settings.URL,
		body,
	)
	if err != nil {
		return http.StatusServiceUnavailable, nil, fmt.Errorf(
			"request: %w",
			err,
		)
	}
	for k, v := range settings.Headers {
		req.Header.Add(k, v)
	}

	client := http.Client{
		Transport: &http.Transport{TLSClientConfig: settings.TLSConfig},
		Timeout:   60 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		return http.StatusServiceUnavailable, nil, fmt.Errorf(
			"request Do: %w",
			err,
		)
	}
	defer res.Body.Close()

	bd, err := io.ReadAll(res.Body)
	if err != nil {
		return res.StatusCode, nil, fmt.Errorf(
			"request body ReadAll: %w",
			err,
		)
	}

	return res.StatusCode, bd, nil
}
