package password_test

import (
	"testing"

	"github.com/aryyawijaya/simple-bank/entity"
	"github.com/aryyawijaya/simple-bank/modules/auth/password"
	"github.com/aryyawijaya/simple-bank/util"
	"github.com/stretchr/testify/require"
)

var ph = password.NewPassHelper()

func TestPassword(t *testing.T) {
	password := util.RandomString(8)

	hashedPass, err := ph.HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPass)

	// valid password
	err = ph.CheckPassword(hashedPass, password)
	require.NoError(t, err)

	// invalid password
	wrongPass := util.RandomString(8)
	err = ph.CheckPassword(hashedPass, wrongPass)
	require.Error(t, err)
	require.EqualError(t, err, entity.ErrPasswordInvalid.Error())

	// rehash password
	hashedPass2, err := ph.HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPass2)
	require.NotEqual(t, hashedPass, hashedPass2)
}
