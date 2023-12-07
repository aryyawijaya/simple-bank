package mydb

import (
	"context"
	"database/sql"
	"fmt"
)

// Store interface provide all function Queries and + transaction
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

/*
SQLStore struct (implement Store interface)
provide all function Queries and + transaction
*/
type SQLStore struct {
	*Queries
	db *sql.DB
}

// create new Store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// executes function within a database transaction
func (store *SQLStore) exectTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// param for transfer money (transaction)
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// result the transfer money (transaction)
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// var txKey = struct{}{}

// performs transfer money
/*
contains create transfer record, add account entries,
and update accounts balance within 1 database transaction
*/
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.exectTx(ctx, func(q *Queries) error {
		var err error

		// txName := ctx.Value(txKey)

		// fmt.Printf("%v create transfer\n", txName)
		// create transfer record
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(arg))
		if err != nil {
			return err
		}

		// fmt.Printf("%v create entry 1\n", txName)
		// add account entries
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -1 * arg.Amount,
		})
		if err != nil {
			return err
		}
		// fmt.Printf("%v create entry 2\n", txName)
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// TO-DO:
		// update accounts balance (due to need locking & prevent potential deadlock)

		// fmt.Printf("%v get account 1 for update\n", txName)
		// get FromAccount balance and subtract with Amount
		/* account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}

		// fmt.Printf("%v update account 1 balance\n", txName)
		result.FromAccount, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
			ID:      arg.FromAccountID,
			Balance: account1.Balance - arg.Amount,
		})
		if err != nil {
			return err
		} */

		// better implementation to update account balance
		// result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		// 	ID:     arg.FromAccountID,
		// 	Amount: -1 * arg.Amount,
		// })
		// if err != nil {
		// 	return err
		// }

		// fmt.Printf("%v get account 2 for update\n", txName)
		// get ToAccount balance and add with Amount
		/* account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}

		// fmt.Printf("%v update account 2 balance\n", txName)
		result.ToAccount, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
			ID:      arg.ToAccountID,
			Balance: account2.Balance + arg.Amount,
		})
		if err != nil {
			return err
		} */

		// better implementation to update account balance
		// result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		// 	ID:     arg.ToAccountID,
		// 	Amount: arg.Amount,
		// })
		// if err != nil {
		// 	return err
		// }

		// acquire locks in consistent order to avoid deadlock
		if arg.FromAccountID < arg.ToAccountID {
			// result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			// 	ID:     arg.FromAccountID,
			// 	Amount: -1 * arg.Amount,
			// })
			// if err != nil {
			// 	return err
			// }
			// result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			// 	ID:     arg.ToAccountID,
			// 	Amount: arg.Amount,
			// })
			// if err != nil {
			// 	return err
			// }
			result.FromAccount, result.ToAccount, err = addMoney(
				ctx,
				q,
				arg.FromAccountID,
				-1*arg.Amount,
				arg.ToAccountID,
				arg.Amount,
			)
			if err != nil {
				return err
			}
		} else {
			// result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			// 	ID:     arg.ToAccountID,
			// 	Amount: arg.Amount,
			// })
			// if err != nil {
			// 	return err
			// }
			// result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			// 	ID:     arg.FromAccountID,
			// 	Amount: -1 * arg.Amount,
			// })
			// if err != nil {
			// 	return err
			// }
			result.ToAccount, result.FromAccount, err = addMoney(
				ctx,
				q,
				arg.ToAccountID,
				arg.Amount,
				arg.FromAccountID,
				-1*arg.Amount,
			)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return // equal to return account1, account2, err (Go feature)
	}
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	return
}
