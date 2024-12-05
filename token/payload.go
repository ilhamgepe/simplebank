package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidTokenMethod = errors.New("unexpected signing method")
	ErrExpiredToken       = errors.New("token has expired")
	ErrInvalidToken       = errors.New("token is invalid")
	ErrParsingToken       = errors.New("failed to parse token")
)

type Payload struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	jwt.RegisteredClaims
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:       tokenID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID.String(),
			Subject:   username,
			Issuer:    "simple_bank",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(duration - time.Minute)),
		},
	}

	return payload, nil
}

func (payload *Payload) valid() error {
	if time.Now().After(payload.RegisteredClaims.ExpiresAt.Time) {
		return ErrExpiredToken
	}
	return nil
}
