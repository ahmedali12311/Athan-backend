package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"app/config"
	"app/pkg/validator"

	firebase "firebase.google.com/go/v4"
	"github.com/golang-jwt/jwt/v4"
)

// jwt structs

type Token struct {
	Value     string `json:"value"`
	Type      string `json:"type"`
	ExpiresIn string `json:"expires_in"`
}

type TokenResponse struct {
	Token Token `json:"token"`
}

type FileTokenResponse struct {
	Token Token  `json:"token"`
	URL   string `json:"url"`
}

type CustomOrderFileClaims struct {
	OrderID string `json:"order_id"`
	jwt.RegisteredClaims
}

// jwt stuff

func (m *Model) ParseToken(
	bearer *string,
) (*jwt.Token, *jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(
		*bearer,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (any, error) {
			return config.JwtSecret, nil
		},
	)
	if err != nil {
		return &jwt.Token{}, &jwt.RegisteredClaims{}, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return token, claims, nil
	}
	return &jwt.Token{}, &jwt.RegisteredClaims{}, err
}

func (m *Model) GenTokenResponse() (TokenResponse, error) {
	// Create the claims
	claims := jwt.RegisteredClaims{
		Issuer:    config.DOMAIN,
		Subject:   m.ID.String(),
		IssuedAt:  jwt.NewNumericDate(config.TimeNow()),
		ExpiresAt: jwt.NewNumericDate(config.TimeNow().Add(config.JwtExpiry)),
		// NotBefore: jwt.NewNumericDate(config.TimeNow()),
		// ID:        m.ID.String(),
		// Audience:  []string{config.DOMAIN},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedString, err := token.SignedString(config.JwtSecret)
	if err != nil {
		return TokenResponse{}, err
	}
	tokenResponse := TokenResponse{
		Token: Token{
			Value:     signedString,
			Type:      "bearer",
			ExpiresIn: fmt.Sprintf("%dh", int(config.JwtExpiry.Hours())),
		},
	}
	return tokenResponse, nil
}

// GenCookie generates an http only cookie for the token given
// with expires time in the future for valid tokens and
// in the past for invalidating tokens (logging out).
func (m *Model) GenCookie(token Token, expires time.Time) http.Cookie {
	return http.Cookie{
		Name:     "accessToken",
		Value:    token.Value,
		Path:     "/",
		MaxAge:   0,
		Expires:  expires,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
}

// VerifyIDToken sets the user phone to the value of phone claims
// or returns an error.
func (m *Model) VerifyIDToken(firebaseIDToken string, fb *firebase.App) error {
	// firebase id token is set use it to get phone and ignore email + password
	auth, err := fb.Auth(context.Background())
	if err != nil {
		return err
	}
	token, err := auth.VerifyIDToken(context.Background(), firebaseIDToken)
	if err != nil {
		return err
	}
	// assertion 1: phone exists in claims
	tokenPhones, found := token.Firebase.Identities["phone"]
	if found {
		// assertion 2: phone claims type assert to array of interface
		tokenPhones, ok := tokenPhones.([]any)
		if !ok {
			err = errors.New("invalid phone claims")
			return err
		}
		// assertion 3: phone array of interface is type string on first element
		tokenPhoneValue, ok := tokenPhones[0].(string)
		if !ok {
			return errors.New("invalid phone claims")
		}
		phone := strings.Replace(tokenPhoneValue, "+", "", 1)
		m.Phone = &phone
		return nil
	}
	return errors.New("no phone claims")
}

// Model Utilities ------------------------------------------------------------

// MergeLogin returns valid boolen or firebase error
func (m *Model) MergeLogin(
	v *validator.Validator,
	fb *firebase.App,
	comparePassword *bool,
) (bool, error) {
	var email, password string
	firebaseIDToken := v.Data.Values.Get("firebase_id_token")
	if firebaseIDToken != "" {
		if err := m.VerifyIDToken(firebaseIDToken, fb); err != nil {
			v.Check(false, "firebase_id_token", err.Error())
			return false, err
		}
	} else {
		*comparePassword = true
		email = v.Data.Values.Get("email")
		password = v.Data.Values.Get("password")
		if password == "" {
			v.Check(false, "password", v.T.ValidateRequired())
		}
		if email == "" {
			v.Check(false, "email", v.T.ValidateRequired())
		}
		m.Email = &email
		m.Password.Plaintext = &password
	}
	return v.Valid(), nil
}

func (m *Model) MergeForgetPassword(v *validator.Validator) (bool, error) {
	m.MergeEmailPassword(v, false, true)

	return true, nil
}

func (m *Model) InvalidCookie() http.Cookie {
	return http.Cookie{
		Name:     "accessToken",
		Value:    "",
		Path:     "/",
		MaxAge:   0,
		Expires:  time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
}
