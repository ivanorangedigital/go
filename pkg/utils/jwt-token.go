package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

// create jwt with user id
func NewJWT(userId, secret string, maxAge uint64) (string, error) {
	now := time.Now()
	claims := &UserClaims{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(maxAge) * time.Second)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// return userId
func ValidateJWT(token, secret string) (string, error) {
	claims := new(UserClaims)
	// parser and token verify
	t, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		// alghoritm verify
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid signing method")
		}
		// return secret to sign
		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}

	if !t.Valid {
		return "", fmt.Errorf("Token invalid")
	}

	return claims.UserID, nil
}
