package auth

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newMockSessionRepo(t *testing.T) (*SessionRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	return NewSessionRepository(db), mock, func() { db.Close() }
}

func TestSessionRepository_Create(t *testing.T) {
	repo, mock, closeDB := newMockSessionRepo(t)
	defer closeDB()

	now := time.Now()
	mock.ExpectExec("INSERT INTO Sessioni").
		WithArgs(1, "hash123", now).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.SessionRepository_Create(context.Background(), 1, "hash123", now)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSessionRepository_IsValid(t *testing.T) {
	repo, mock, closeDB := newMockSessionRepo(t)
	defer closeDB()

	t.Run("sessione valida", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM Sessioni WHERE token_hash = \\? AND expires_at > NOW\\(\\)").
			WithArgs("hash123").
			WillReturnRows(rows)

		valid, err := repo.SessionRepository_IsValid(context.Background(), "hash123")
		assert.NoError(t, err)
		assert.True(t, valid)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("sessione scaduta o inesistente", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM Sessioni WHERE token_hash = \\? AND expires_at > NOW\\(\\)").
			WithArgs("hash_scaduto").
			WillReturnRows(rows)

		valid, err := repo.SessionRepository_IsValid(context.Background(), "hash_scaduto")
		assert.NoError(t, err)
		assert.False(t, valid)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSessionRepository_DeleteByTokenHash(t *testing.T) {
	repo, mock, closeDB := newMockSessionRepo(t)
	defer closeDB()

	mock.ExpectExec("DELETE FROM Sessioni WHERE token_hash = \\?").
		WithArgs("hash123").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.SessionRepository_DeleteByTokenHash(context.Background(), "hash123")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
