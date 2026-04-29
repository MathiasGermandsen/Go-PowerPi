package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	Service string   `json:"svc"`
	Scopes  []string `json:"scopes"`
	jwt.RegisteredClaims
}

func GenerateToken(secret, service, audience string, scopes []string, expHours int) (tokenStr string, jti string, err error) {
	if secret == "" {
		return "", "", errors.New("JWT_SECRET must not be empty")
	}

	jti = uuid.NewString()
	now := time.Now().UTC()
	exp := now.Add(time.Duration(expHours) * time.Hour)

	claims := Claims{
		Service: service,
		Scopes:  scopes,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			Subject:   service,
			Audience:  jwt.ClaimStrings{audience},
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err = token.SignedString([]byte(secret))
	return tokenStr, jti, err
}

func ValidateToken(tokenStr, secret string) (*Claims, error) {
	if secret == "" {
		return nil, errors.New("JWT_SECRET must not be empty")
	}

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
