package user

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tommasogarofani/ingsw_25-26_group-11/server/domain"
)

// Query usata nei test: deve corrispondere a quella in Repository.GetByUsername.
// sqlmock confronta per regex per default, quindi basta una sottostringa stabile.
const getByUsernameQuery = `SELECT id, username, password_hash, ruolo, domanda_sicurezza, risposta_hash, tentativi_falliti, bloccato, created_at`

func newMockRepo(t *testing.T) (*UserRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	repo := NewUserRepository(db)
	return repo, mock, func() { db.Close() }
}

func TestUserRepository_GetByUsername(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		username  string
		mockSetup func(mock sqlmock.Sqlmock)
		wantErr   error // se non nil, ci aspettiamo errors.Is(err, wantErr)
		wantUser  bool
	}{
		{
			name:     "utente trovato",
			username: "mario",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "username", "password_hash", "ruolo",
					"domanda_sicurezza", "risposta_hash", "tentativi_falliti",
					"bloccato", "created_at",
				}).AddRow(1, "mario", "hash123", "magazziniere", "Il tuo colore preferito?", "rispostaHash", 0, false, now)

				mock.ExpectQuery(getByUsernameQuery).
					WithArgs("mario").
					WillReturnRows(rows)
			},
			wantUser: true,
		},
		{
			name:     "utente non trovato -> credenziali errate",
			username: "inesistente",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(getByUsernameQuery).
					WithArgs("inesistente").
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: domain.ErrCredenzaliErrate,
		},
		{
			name:     "errore generico del database viene propagato",
			username: "mario",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(getByUsernameQuery).
					WithArgs("mario").
					WillReturnError(errors.New("connessione al database persa"))
			},
			wantErr: errors.New("connessione al database persa"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, mock, closeDB := newMockRepo(t)
			defer closeDB()

			tt.mockSetup(mock)

			got, err := repo.UserRepository_GetByUsername(context.Background(), tt.username)

			if tt.wantErr != nil {
				require.Error(t, err)
				if errors.Is(tt.wantErr, domain.ErrCredenzaliErrate) {
					assert.ErrorIs(t, err, domain.ErrCredenzaliErrate)
				}
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
			}

			if tt.wantUser {
				require.NotNil(t, got)
				assert.Equal(t, tt.username, got.Username)
				assert.Equal(t, domain.Ruolo("magazziniere"), got.Ruolo)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_Update(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(mock sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "update riuscito",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE Utenti SET password_hash = \?, tentativi_falliti = \?, bloccato = \? WHERE id = \?`).
					WithArgs("nuovoHash", 1, false, 42).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
		{
			name: "errore del database durante l'update",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE Utenti SET password_hash = \?, tentativi_falliti = \?, bloccato = \? WHERE id = \?`).
					WithArgs("nuovoHash", 1, false, 42).
					WillReturnError(errors.New("db irraggiungibile"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, mock, closeDB := newMockRepo(t)
			defer closeDB()

			tt.mockSetup(mock)

			u := &domain.Utente{
				ID:               42,
				PasswordHash:     "nuovoHash",
				TentativiFalliti: 1,
				Bloccato:         false,
			}
			err := repo.UserRepository_Update(context.Background(), u)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := NewUserRepository(db)

	mock.ExpectExec(`INSERT INTO Utenti`).
		WithArgs("giulia", "hash", "contabile", "domanda", "rispostaHash", 0, false, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(7, 1))

	u := &domain.Utente{
		Username:         "giulia",
		PasswordHash:     "hash",
		Ruolo:            domain.Ruolo("contabile"),
		DomandaSicurezza: "domanda",
		RispostaHash:     "rispostaHash",
		CreatedAt:        time.Now(),
	}

	err = repo.UserRepository_Create(context.Background(), u)
	require.NoError(t, err)
	assert.Equal(t, 7, u.ID, "l'ID auto-generato dovrebbe essere assegnato all'utente")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Delete(t *testing.T) {
	// NB: richiede che la query in Repository.Delete usi backtick invece di apici
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := NewUserRepository(db)

	mock.ExpectExec(`DELETE FROM Utenti WHERE id = \?`).
		WithArgs(5).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.UserRepository_Delete(context.Background(), &domain.Utente{ID: 5})
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
