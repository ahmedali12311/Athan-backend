package auth

import (
	"time"

	"app/config"

	"github.com/golang-jwt/jwt/v4"
)

// GenJWT generates a jwt with passed subjectID and custom expiry
func GenJWT(
	subjectID string,
	expiry time.Duration,
) (string, error) {
	return jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:   config.DOMAIN,
			Subject:  subjectID,
			IssuedAt: jwt.NewNumericDate(config.TimeNow()),
			ExpiresAt: jwt.NewNumericDate(
				config.TimeNow().Add(expiry),
			),
		},
	).SignedString(config.JwtSecret)
}
