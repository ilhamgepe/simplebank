package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomInt64(t *testing.T) {
	assert.NotNil(t, RandomInt64(10, 20))
}

func TestRandomString(t *testing.T) {
	assert.NotNil(t, RandomString(10))
}

func TestRandomOwner(t *testing.T) {
	assert.NotNil(t, RandomOwner())
}

func TestRandomBalance(t *testing.T) {
	assert.NotNil(t, RandomBalance())
}

func TestRandomEntriesAmount(t *testing.T) {
	assert.NotNil(t, RandomEntriesAmount())
}

func TestRandomCurrency(t *testing.T) {
	assert.NotNil(t, RandomCurrency())
}
