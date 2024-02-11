package mydb_test

import (
	"context"
	"fmt"
	"testing"

	mydb "github.com/aryyawijaya/simple-bank/db/sqlc"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := mydb.NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Printf(">> before: %v %v\n", account1.Balance, account2.Balance)

	// run n concurrent transfer transaction
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan mydb.TransferTxResult)

	for i := 0; i < n; i++ {
		// txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			// ctx := context.WithValue(context.Background(), txKey, txName)
			ctx := context.Background()
			result, err := store.TransferTx(ctx, mydb.TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	// check results
	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		// blocking code, wait data from Go routines
		err := <-errs
		require.NoError(t, err)

		// blocking code, wait data from Go routines
		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -1*amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check balance
		// check account frist
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		fmt.Printf(">> tx: %v %v\n", fromAccount.Balance, toAccount.Balance)

		// check account's balance
		// outFromAcc1 --> amount money going out from account1
		outFromAcc1 := account1.Balance - fromAccount.Balance
		// inToAcc2 --> amount money going in to account2
		inToAcc2 := toAccount.Balance - account2.Balance
		// outFromAcc1 should be same as inToAcc2
		require.Equal(t, outFromAcc1, inToAcc2)
		require.True(t, outFromAcc1 > 0)
		// outFromAcc1 should be divisible by amount
		// --> either 1, 2, ..., n times transfer
		require.True(t, outFromAcc1%amount == 0)

		// define k should be distinct from 1-n
		k := int(outFromAcc1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)

		existed[k] = true
	}

	// check final balance
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Printf(">> after: %v %v\n", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
}

func TestTransferTxDeadlockByOrderQuery(t *testing.T) {
	store := mydb.NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Printf(">> before: %v %v\n", account1.Balance, account2.Balance)

	// run n concurrent transfer transaction
	/*
		use n = 10
		5 from account 1 to account 2
		5 from account 2 to account 1
	*/
	n := 10
	amount := int64(10)

	errs := make(chan error)
	// delete result, because just need care about deadlock error

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {
			ctx := context.Background()
			_, err := store.TransferTx(ctx, mydb.TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		// blocking code, wait data from Go routines
		err := <-errs
		require.NoError(t, err)
	}

	// check final balance
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Printf(">> after: %v %v\n", updatedAccount1.Balance, updatedAccount2.Balance)

	// this case balance account 1 & 2 should be equal before transaction
	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
