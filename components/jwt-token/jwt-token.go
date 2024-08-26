package jwt_token

import (
	"app/config"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// JWT structs

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
	Url   string `json:"url"`
}

type FileClaims struct {
	ID uuid.UUID `json:"id"`
	jwt.RegisteredClaims
}

type FuncMapClaims func() jwt.MapClaims

func GenFileToken(fn FuncMapClaims) (string, error) {
	claims := fn()
	claims["exp"] = jwt.NewNumericDate(time.Now().Add(time.Minute * 15))
	// claims := FileClaims{
	// 	id,
	// 	jwt.RegisteredClaims{
	// 		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
	// 	},
	// }

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString(config.JwtSecret)
	if err != nil {
		return "", nil
	}

	return signedString, nil
}

func ParseFileToken(bearer *string) (jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(*bearer,
		jwt.MapClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return config.JwtSecret, nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
