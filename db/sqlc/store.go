package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

type Store struct {
	*Queries
	db *pgx.Conn
}

func NewStore(db *pgx.Conn) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = fn(store.Queries.WithTx(tx))
	log.Println("err execTx", err)
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

		// TODO: update account balance

		return nil
	})

	return result, err
}
