package utils

import (
	"smartgas-payment/config"
	"smartgas-payment/internal/schemas"

	"github.com/golang-jwt/jwt/v4"
)

func ParseJwtToken(tokenString string, isRefreshToken bool) (*schemas.JwtClaims, error) {
	cfg := config.ConfigSettings

	secretKey := cfg.SecretKey
	if isRefreshToken {
		secretKey = cfg.SecretKeyRefresh
	}

	token, err := jwt.ParseWithClaims(tokenString, &schemas.JwtClaims{}, func(t *jwt.Token) (interface{}, error) {

		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims := token.Claims.(*schemas.JwtClaims)

	if err = claims.Valid(); err != nil {
		return nil, err
	}

	return claims, nil
}
