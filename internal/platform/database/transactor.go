package database

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type contextKey string

const txKey contextKey = "tx"

// Transactor defines the interface for running operations within a transaction
type Transactor interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type transactor struct {
	db *sqlx.DB
}

// NewTransactor creates a new Transactor instance
func NewTransactor(db *sqlx.DB) Transactor {
	return &transactor{db: db}
}

// WithinTransaction executes a function within a database transaction.
// If the function returns an error, the transaction is rolled back.
// Otherwise, the transaction is committed.
func (t *transactor) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := t.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	// Inject transaction into context
	txCtx := context.WithValue(ctx, txKey, tx)

	if err := fn(txCtx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

// GetQueryer returns a sqlx.ExtContext from the context if it exists (transaction),
// otherwise it returns the provided *sqlx.DB.
func GetQueryer(ctx context.Context, db *sqlx.DB) sqlx.ExtContext {
	if tx, ok := ctx.Value(txKey).(*sqlx.Tx); ok {
		return tx
	}
	return db
}
