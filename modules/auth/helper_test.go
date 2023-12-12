package auth_test

import (
	"testing"

	"github.com/aryyawijaya/simple-bank/modules/auth"
	"github.com/aryyawijaya/simple-bank/util"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

var ah = auth.NewAuthHelper()

func TestPassword(t *testing.T) {
	password := util.RandomString(8)

	hashedPass, err := ah.HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPass)

	// valid password
	err = ah.CheckPassword(hashedPass, password)
	require.NoError(t, err)

	// invalid password
	wrongPass := util.RandomString(8)
	err = ah.CheckPassword(hashedPass, wrongPass)
	require.Error(t, err)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	// rehash password
	hashedPass2, err := ah.HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPass2)
	require.NotEqual(t, hashedPass, hashedPass2)
}
