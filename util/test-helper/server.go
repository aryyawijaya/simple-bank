package testhelper

import (
	"testing"
	"time"

	mydb "github.com/aryyawijaya/simple-bank/db/sqlc"
	"github.com/aryyawijaya/simple-bank/server-http-gin"
	"github.com/aryyawijaya/simple-bank/util"
	"github.com/stretchr/testify/require"
)

func NewTestServer(t *testing.T, store mydb.Store) *server.Server {
	config := &util.Config{
		TokenSymetricKey:    util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	s, err := server.NewServer(store, config)
	require.NoError(t, err)

	return s
}
