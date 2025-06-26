package sms

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type options struct {
	logger *zerolog.Logger
}

type Option func(options *options) error

func WithLogger(logger *zerolog.Logger) Option {
	return func(options *options) error {
		if logger != nil {
			options.logger = logger
			return nil
		}

		return errors.New("logger is nil")
	}
}

type SMSConfig struct {
	Url string
	Key string
	Jwt string
	Env string
}

type SMSInput struct {
	Phone   string `json:"phone"`
	Content string `json:"content"`
}

type SMS struct {
	config  SMSConfig
	options options
}

func NewSMS(config SMSConfig, opts ...Option) (*SMS, error) {
	var sms SMS
	for _, opt := range opts {
		if err := opt(&sms.options); err != nil {
			return nil, err
		}
	}

	sms.config = config

	if sms.config.Url == "" {
		return nil, errors.New("sms url can't be empty")
	}

	if sms.config.Jwt == "" {
		return nil, errors.New("sms jwt can't be empty")
	}

	return &sms, nil
}

func (s *SMS) SendMessage(phone, message string) error {
	sms_payload := SMSInput{
		Phone:   phone,
		Content: message,
	}

	token := "Bearer " + s.config.Jwt
	json_payload, err := json.Marshal(sms_payload)
	if err != nil {
		return err
	}

	// Todo: make it flexible for testing
	// and production
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		s.config.Url+"/api/v1/messages?test=",
		bytes.NewBuffer(json_payload),
	)
	if err != nil {
		return err
	}
	req.Header.Add("authorization", token)
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{
		Timeout: 60 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("http client.Do: %w", err)
	}
	defer res.Body.Close()

	bd, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("io.ReadAll: %w", err)
	}
	var response map[string]any
	if err := json.Unmarshal(bd, &response); err != nil {
		return fmt.Errorf("can not unmarshal JSON: %w", err)
	}

	if s.options.logger != nil {
		s.options.logger.Info().
			Str("form", fmt.Sprintf("%s - %s", phone, message)).
			RawJSON("otp response", bd).
			Msg("sms")
	}

	return nil
}
