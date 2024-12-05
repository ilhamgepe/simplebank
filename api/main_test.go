package api

import (
	"os"
	"testing"
	"time"

	db "github.com/ilhamgepe/simplebank/db/sqlc"
	"github.com/ilhamgepe/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := utils.Config{
		TokenSymmetricKey:   utils.RandomString(32),
		AccessTokenDuration: time.Minute,
	}
	server, err := NewServer(store, config)
	require.NoError(t, err)
	return server
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
