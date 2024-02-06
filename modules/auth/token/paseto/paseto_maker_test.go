package paseto_test

import (
	"testing"
	"time"

	"github.com/aryyawijaya/simple-bank/modules/auth/token"
	"github.com/aryyawijaya/simple-bank/modules/auth/token/paseto"
	"github.com/aryyawijaya/simple-bank/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := paseto.NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomString(5)
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := paseto.NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	pasetoToken, payload, err := maker.CreateToken(util.RandomString(5), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, pasetoToken)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(pasetoToken)
	require.Error(t, err)
	require.EqualError(t, err, token.ErrTokenExpired.Error())
	require.Nil(t, payload)
}
