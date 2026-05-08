package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// devSecret is used when no JWT_SECRET env var is set.
// TODO: replace with a secret manager or required env var before production.
const devSecret = "triangleoflove-dev-secret-not-for-production"

// Claims are the JWT payload fields.
type Claims struct {
	AccountID string `json:"accountId"`
	Role      string `json:"role"`
	jwt.RegisteredClaims
}

// SignToken issues a signed JWT for the given account ID and role.
func SignToken(accountID, role string) (string, error) {
	claims := Claims{
		AccountID: accountID,
		Role:      role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(devSecret))
}

// ParseToken validates a JWT string and returns its claims.
func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(devSecret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
