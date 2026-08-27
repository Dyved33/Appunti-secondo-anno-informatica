package domain

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetExecutor_NoTxInContext(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	ctx := context.Background()
	exec := GetExecutor(ctx, db)
	assert.Same(t, db, exec)
}

func TestGetExecutor_WithTxInContext(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	tx, err := db.BeginTx(context.Background(), nil)
	require.NoError(t, err)

	ctx := context.WithValue(context.Background(), txKey, tx)
	exec := GetExecutor(ctx, db)
	assert.Same(t, tx, exec)

	_ = tx.Rollback()
}

func TestNewTransactionHelper(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	th := NewTransactionHelper(db)
	assert.NotNil(t, th)
	assert.Same(t, db, th.DB)
}

func TestWithinTransaction_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectCommit()

	th := NewTransactionHelper(db)
	err = th.WithinTransaction(context.Background(), func(ctx context.Context) error {
		exec := GetExecutor(ctx, db)
		_, ok := exec.(*sql.Tx)
		assert.True(t, ok, "GetExecutor should return the *sql.Tx inside transaction")
		return nil
	})

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestWithinTransaction_RollbackOnError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sentinelErr := errors.New("business logic failed")
	mock.ExpectBegin()
	mock.ExpectRollback()

	th := NewTransactionHelper(db)
	err = th.WithinTransaction(context.Background(), func(ctx context.Context) error {
		return sentinelErr
	})

	require.Error(t, err)
	assert.Equal(t, sentinelErr, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestWithinTransaction_BeginTxFails(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin().WillReturnError(errors.New("connection refused"))

	th := NewTransactionHelper(db)
	err = th.WithinTransaction(context.Background(), func(ctx context.Context) error {
		t.Fatal("fn should never be called when BeginTx fails")
		return nil
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "connection refused")
}

func TestWithinTransaction_NestedTx(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	tx, err := db.BeginTx(context.Background(), nil)
	require.NoError(t, err)

	outerCtx := context.WithValue(context.Background(), txKey, tx)

	th := NewTransactionHelper(db)
	err = th.WithinTransaction(outerCtx, func(ctx context.Context) error {
		exec := GetExecutor(ctx, db)
		assert.Same(t, tx, exec, "should reuse existing tx, not create a new one")
		return nil
	})

	require.NoError(t, err)
	_ = tx.Rollback()
	_ = db.Close()
}
