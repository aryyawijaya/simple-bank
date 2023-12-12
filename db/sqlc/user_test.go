package mydb

import (
	"context"
	"testing"

	"github.com/aryyawijaya/simple-bank/modules/auth"
	"github.com/aryyawijaya/simple-bank/util"
	"github.com/stretchr/testify/require"
)

var ah = auth.NewAuthHelper()

func createRandomUser(t *testing.T) User {
	hashedPass, err := ah.HashPassword(util.RandomString(8))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPass,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	createdUser, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, createdUser)

	require.Equal(t, arg.Username, createdUser.Username)
	require.Equal(t, arg.HashedPassword, createdUser.HashedPassword)
	require.Equal(t, arg.FullName, createdUser.FullName)
	require.Equal(t, arg.Email, createdUser.Email)

	require.True(t, createdUser.PasswordChangedAt.IsZero())
	require.NotZero(t, createdUser.CreatedAt)

	return createdUser
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	createdUser := createRandomUser(t)
	currUser, err := testQueries.GetUser(context.Background(), createdUser.Username)

	require.NoError(t, err)
	require.NotEmpty(t, currUser)

	require.Equal(t, createdUser, currUser)
}
