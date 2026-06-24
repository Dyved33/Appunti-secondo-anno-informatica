package db

import (
	"context"
	"database/sql"
)

type contextKey string

const txKey contextKey = "tx"

// QueryExecutor defines the subset of sql.DB and sql.Tx methods needed to run queries.
type QueryExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// GetExecutor returns the transaction from the context if present, otherwise the default database connection.
func GetExecutor(ctx context.Context, defaultDB *sql.DB) QueryExecutor {
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		return tx
	}
	return defaultDB
}

// TransactionHelper facilitates running operations within a transaction.
type TransactionHelper struct {
	DB *sql.DB
}

// NewTransactionHelper creates a new TransactionHelper.
func NewTransactionHelper(db *sql.DB) *TransactionHelper {
	return &TransactionHelper{DB: db}
}

// WithinTransaction executes the given function inside a transaction.
// If the context already contains a transaction, it executes the function directly.
func (th *TransactionHelper) WithinTransaction(ctx context.Context, fn func(context.Context) error) error {
	if _, ok := ctx.Value(txKey).(*sql.Tx); ok {
		return fn(ctx)
	}

	tx, err := th.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	txCtx := context.WithValue(ctx, txKey, tx)
	err = fn(txCtx)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
