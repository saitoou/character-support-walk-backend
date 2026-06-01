package database

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type txKey struct{}

func GetExecutor(ctx context.Context, db DBTX) DBTX {
	tx, ok := ctx.Value(txKey{}).(*sql.Tx)
	if ok {
		return tx
	}
	return db
}

type Transaction struct {
	db *sql.DB
}

func NewTransaction(db *sql.DB) *Transaction {
	return &Transaction{db: db}
}

func (t *Transaction) RunInTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	txCtx := context.WithValue(ctx, txKey{}, tx)

	if err := fn(txCtx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("rollback failed: %w: original error: %v", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
