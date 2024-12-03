package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(6)
	hashed, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashed)

	err = CheckPasswordHash(password, hashed)
	require.NoError(t, err)

	wrongPassword := RandomString(6)
	err = CheckPasswordHash(wrongPassword, hashed)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashed2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashed)
	require.NotEqual(t, hashed, hashed2)
}
