package testhelper

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	mydb "github.com/aryyawijaya/simple-bank/db/sqlc"
	"github.com/aryyawijaya/simple-bank/middlewares"
	"github.com/aryyawijaya/simple-bank/modules/auth/password"
	"github.com/aryyawijaya/simple-bank/modules/auth/token"
	"github.com/aryyawijaya/simple-bank/util"
	"github.com/stretchr/testify/require"
)

var ph = password.NewPassHelper()

func RandomUser(t *testing.T) (user mydb.User, pass string) {
	pass = util.RandomString(8)
	hashedPass, err := ph.HashPassword(pass)
	require.NoError(t, err)

	user = mydb.User{
		Username:       util.RandomString(5),
		HashedPassword: hashedPass,
		FullName:       util.RandomString(5),
		Email:          util.RandomEmail(),
	}

	return
}

func RandomAccount(owner string) mydb.Account {
	return mydb.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    owner,
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}
}

func CustomAuthorization(
	t *testing.T,
	request *http.Request,
	tokenMaker token.Maker,
	authorizationType string,
	username string,
	duration time.Duration,
) {
	token, payload, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set(middlewares.AuthorizationHeaderKey, authorizationHeader)
}
