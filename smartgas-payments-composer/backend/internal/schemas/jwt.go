package schemas

import (
	"smartgas-payment/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type JwtClaims struct {
	Sub uuid.UUID `json:"sub"`
	jwt.RegisteredClaims
}

func (j JwtClaims) ClaimToken() (string, error) {
	cfg := config.ConfigSettings

	if j.ExpiresAt == nil {
		j.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(cfg.JwtExpMinutes)))
	}

	j.IssuedAt = jwt.NewNumericDate(time.Now())
	j.Issuer = "Smart Gas"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, j)

	return token.SignedString([]byte(cfg.SecretKey))
}

func (j JwtClaims) ClaimRefreshToken() (string, error) {
	cfg := config.ConfigSettings

	if j.ExpiresAt == nil {
		j.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * time.Duration(cfg.JwtRefreshExpDays)))
	}

	j.IssuedAt = jwt.NewNumericDate(time.Now())
	j.Issuer = "Smart Gas"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, j)

	return token.SignedString([]byte(cfg.SecretKeyRefresh))
}
