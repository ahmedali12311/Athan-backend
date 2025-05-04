package otp

import "github.com/google/uuid"

const (
	Bearer = "Bearer "

	KeySadeemOTPURL = "sadeem_otp_url"
	KeySadeemOTPKey = "sadeem_otp_key"
	KeySadeemOTPJWT = "sadeem_otp_jwt"
	KeySadeemOTPEnv = "sadeem_otp_env"
	// EndpointV1Pins appended to the Settins.URL and sent as POST
	EndpointV1Pins = "/api/v1/pins"
)

type Settings struct {
	Key string
	JWT string
	URL string
	// Env Should be development or production
	Env string
}

// 1 request: send otp

type Input struct {
	Phone string `json:"phone"`
}

// Response the only success code is 201 created, else it is an error
type Response struct {
	ID *uuid.UUID `json:"id,omitempty"`
	// Pin the 6 random generated numbers
	Pin string `json:"pin,omitempty"`
	// Code the phone international key eg: 218
	Code string `json:"code,omitempty"`
	// Region the iso international code eg: LY
	Region string `json:"region,omitempty"`
	// Number the national phone number eg: 910001234
	Number string `json:"number,omitempty"`
	// Content the sms message content containint pin eg: Your OTP is 001100
	Content string `json:"content,omitempty"`

	// if any of these exists it is an error
	Status    int    `json:"status,omitempty"`
	Type      string `json:"type,omitempty"`
	Message   string `json:"message,omitempty"`
	RequestID string `json:"request_id,omitempty"`

	Errors map[string][]string `json:"errors,omitempty"`
}
