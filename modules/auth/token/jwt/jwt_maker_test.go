package jwt_test

import (
	"testing"
	"time"

	"github.com/aryyawijaya/simple-bank/modules/auth/token"
	"github.com/aryyawijaya/simple-bank/modules/auth/token/jwt"
	"github.com/aryyawijaya/simple-bank/util"
	jwtv5 "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := jwt.NewJWTMaker(util.RandomString(32))
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

func TestExpiredJWTToken(t *testing.T) {
	maker, err := jwt.NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	jwtToken, payload, err := maker.CreateToken(util.RandomString(5), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, jwtToken)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(jwtToken)
	require.Error(t, err)
	require.EqualError(t, err, token.ErrTokenExpired.Error())
	require.Nil(t, payload)
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	payload, err := token.NewPayload(util.RandomOwner(), time.Minute)
	require.NoError(t, err)

	claims := jwt.JWTPayload{
		*payload,
		jwtv5.RegisteredClaims{
			ID:        payload.ID.String(),
			IssuedAt:  &jwtv5.NumericDate{Time: payload.IssuedAt},
			ExpiresAt: &jwtv5.NumericDate{Time: payload.ExpiredAt},
		},
	}

	jwtToken := jwtv5.NewWithClaims(jwtv5.SigningMethodNone, claims)
	resultToken, err := jwtToken.SignedString(jwtv5.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := jwt.NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(resultToken)
	require.Error(t, err)
	require.EqualError(t, err, token.ErrTokenInvalid.Error())
	require.Nil(t, payload)
}
