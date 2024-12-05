package token

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTMaker struct {
	secret string
}

const minSecretKeySize = 32

// NewJWTMaker returns a new JWTMaker instance with the given secret key.
// The provided secret key must be at least 32 characters long.
// If the secret key is too small, CreateToken will return an error.
func NewJWTMaker(secret string) (Maker, error) {
	if len(secret) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{
		secret: secret,
	}, nil
}

// CreateToken creates a new token for the given user with given duration.
// It's supposed to return a new, signed token that can be verified with the same maker.
// The function will return an error if the secret key is too small.
func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	return jwtToken.SignedString([]byte(maker.secret))
}

// VerifyToken checks if the given token is valid.
// It returns the payload of the token if the token is valid.
// If the token is invalid, it returns an error.
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, func(t *jwt.Token) (interface{}, error) {
		if v, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("unexpected signing method: %v", v)
			return nil, ErrInvalidTokenMethod
		}
		return []byte(maker.secret), nil
	})
	if err != nil {
		return nil, err
	}

	if !jwtToken.Valid {
		return nil, ErrInvalidToken
	}
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrParsingToken
	}

	if err := payload.valid(); err != nil {
		return nil, err
	}

	return payload, nil
}
