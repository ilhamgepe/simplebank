package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
}
type store struct {
	*Queries
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) Store {
	return &store{
		Queries: New(db),
		db:      db,
	}
}

func (s *store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	q := New(tx)

	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx Error: %v \n rollback error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
