package mydb

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/aryyawijaya/simple-bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	createdUser := createRandomUser(t)

	arg := CreateAccountParams{
		Owner:    createdUser.Username,
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	createdAccount := createRandomAccount(t)
	currAccount, err := testQueries.GetAccount(context.Background(), createdAccount.ID)

	require.NoError(t, err)
	require.NotEmpty(t, currAccount)

	require.Equal(t, createdAccount.ID, currAccount.ID)
	require.Equal(t, createdAccount.Owner, currAccount.Owner)
	require.Equal(t, createdAccount.Balance, currAccount.Balance)
	require.Equal(t, createdAccount.Currency, currAccount.Currency)
	require.WithinDuration(t, createdAccount.CreatedAt, currAccount.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	createdAccount := createRandomAccount(t)

	arg := UpdateAccountBalanceParams{
		ID:      createdAccount.ID,
		Balance: util.RandomBalance(),
	}

	updatedAccount, err := testQueries.UpdateAccountBalance(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount)

	require.Equal(t, createdAccount.ID, updatedAccount.ID)
	require.Equal(t, createdAccount.Owner, updatedAccount.Owner)
	require.Equal(t, arg.Balance, updatedAccount.Balance)
	require.Equal(t, createdAccount.Currency, updatedAccount.Currency)
	require.WithinDuration(t, createdAccount.CreatedAt, updatedAccount.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	createdAccount := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), createdAccount.ID)
	require.NoError(t, err)

	currAccount, err := testQueries.GetAccount(context.Background(), createdAccount.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, currAccount)
}

func TestListAccount(t *testing.T) {
	for i := 0; i < 5; i++ {
		createRandomAccount(t)
	}

	limit, offset := 3, 2
	arg := ListAccountsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, limit)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
