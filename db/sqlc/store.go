package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	*Queries
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(q *Queries) error) error {
	tx, err := store.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = fn(store.Queries.WithTx(tx))
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

/*
* TransferTxparams adalah sebuah input / DTO untuk TraTransferTx
 */
type TransferTxparams struct {
	FromAccountID int64   `json:"from_account_id"`
	ToAccountID   int64   `json:"to_account_id"`
	Amount        int64   `json:"amount"`
	Description   *string `json:"description"`
}

/*
TraTransferTxResult adalah hasil dari traTransferTx

	Transfer	"ini adalah hasil dari transfer"
	FromAccount	"ini adalah account pengirim"
	ToAccount	"ini adalah account penerima"
	* semua catatan uang akan berada di entries table
	FromEntry	"ini adalah entry pengirim (negative amount)"
	ToEntry		"ini adalah entry penerima (positive amount)"
*/
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

/*
* transfertx melakukan transfer uang dari akun 1 ke akun 2
* akan membuat transfer record, add acount entries,dan update account balance
* dalam 1 database transaction
 */
func (store *Store) TransferTx(ctx context.Context, arg *TransferTxparams) (TransferTxResult, error) {
	var result TransferTxResult
	// membuat transaction dengan menggunakan execTx
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// membuat transfer
		result.Transfer, err = q.CreateTransfer(ctx, &CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		// menambahkan entry
		result.FromEntry, err = q.CreateEntry(ctx, &CreateEntryParams{
			AccountID:   arg.FromAccountID,
			Amount:      -arg.Amount,
			Description: arg.Description,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, &CreateEntryParams{
			AccountID:   arg.ToAccountID,
			Amount:      arg.Amount,
			Description: arg.Description,
		})
		if err != nil {
			return err
		}

		// untuk menghindari deadlock bisa dengan cara mengurutkan id/identifier apapun
		// jadi tidak ada yang namanya id 1 tf ke id 2 tp ada id 2 tf ke id 1, karna yang di lakukan dari yang kecil dulu
		// sehingga menghindari deadlock
		if arg.ToAccountID < arg.FromAccountID {
			result.ToAccount, err = q.UpdateAccountBalance(ctx, &UpdateAccountBalanceParams{
				Amount: arg.Amount,
				ID:     arg.ToAccountID,
			})
			if err != nil {
				return err
			}

			result.FromAccount, err = q.UpdateAccountBalance(ctx, &UpdateAccountBalanceParams{
				Amount: -arg.Amount,
				ID:     arg.FromAccountID,
			})
			if err != nil {
				return err
			}

			return nil
		}

		// update account
		// update sender account balance -> sender account.Balance - arg.Amount
		result.FromAccount, err = q.UpdateAccountBalance(ctx, &UpdateAccountBalanceParams{
			Amount: -arg.Amount,
			ID:     arg.FromAccountID,
		})
		if err != nil {
			return err
		}

		// get account 2
		// update account 1 balance -> account1.Balance - arg.Amount
		result.ToAccount, err = q.UpdateAccountBalance(ctx, &UpdateAccountBalanceParams{
			Amount: arg.Amount,
			ID:     arg.ToAccountID,
		})
		if err != nil {
			return err
		}
		return nil
	})

	return result, err
}
