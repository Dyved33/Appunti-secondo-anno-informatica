package category

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tommasogarofani/ingsw_25-26_group-11/server/domain"
)

func newMockCategoryRepo(t *testing.T) (*CategoryRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	return NewCategoryRepository(db), mock, func() { db.Close() }
}

func TestCategoryRepository_Create(t *testing.T) {
	repo, mock, closeDB := newMockCategoryRepo(t)
	defer closeDB()

	mock.ExpectExec(`INSERT INTO Categorie \(nome, descrizione\) VALUES \(\?, \?\)`).
		WithArgs("Ferramenta", "Minuteria metallica").
		WillReturnResult(sqlmock.NewResult(3, 1))

	c := &domain.Categoria{Nome: "Ferramenta", Descrizione: "Minuteria metallica"}
	err := repo.CategoryRepository_Create(context.Background(), c)

	require.NoError(t, err)
	assert.Equal(t, 3, c.ID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCategoryRepository_GetAll(t *testing.T) {
	repo, mock, closeDB := newMockCategoryRepo(t)
	defer closeDB()

	rows := sqlmock.NewRows([]string{"id", "nome", "descrizione"}).
		AddRow(1, "Elettronica", "").
		AddRow(2, "Ferramenta", "Minuteria metallica")

	mock.ExpectQuery(`SELECT id, nome, descrizione FROM Categorie ORDER BY nome`).
		WillReturnRows(rows)

	got, err := repo.CategoryRepository_GetAll(context.Background())
	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, "Elettronica", got[0].Nome)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCategoryRepository_GetByID(t *testing.T) {
	tests := []struct {
		name      string
		id        int
		mockSetup func(mock sqlmock.Sqlmock)
		wantErr   error
	}{
		{
			name: "categoria trovata",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "nome", "descrizione"}).AddRow(1, "Elettronica", "")
				mock.ExpectQuery(`SELECT id, nome, descrizione FROM Categorie WHERE id = \?`).
					WithArgs(1).
					WillReturnRows(rows)
			},
		},
		{
			name: "categoria non trovata",
			id:   99,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, nome, descrizione FROM Categorie WHERE id = \?`).
					WithArgs(99).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: domain.ErrCategoriaNonTrovata,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, mock, closeDB := newMockCategoryRepo(t)
			defer closeDB()
			tt.mockSetup(mock)

			got, err := repo.CategoryRepository_GetByID(context.Background(), tt.id)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCategoryRepository_Update(t *testing.T) {
	repo, mock, closeDB := newMockCategoryRepo(t)
	defer closeDB()

	mock.ExpectExec(`UPDATE Categorie SET nome = \?, descrizione = \? WHERE id = \?`).
		WithArgs("Elettronica", "Componenti elettronici", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.CategoryRepository_Update(context.Background(), &domain.Categoria{ID: 1, Nome: "Elettronica", Descrizione: "Componenti elettronici"})
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCategoryRepository_Delete(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(mock sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "delete riuscita",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM Categorie WHERE id = \?`).
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
		{
			name: "errore propagato dal database",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM Categorie WHERE id = \?`).
					WithArgs(1).
					WillReturnError(errors.New("violazione vincolo di chiave esterna"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, mock, closeDB := newMockCategoryRepo(t)
			defer closeDB()
			tt.mockSetup(mock)

			err := repo.CategoryRepository_Delete(context.Background(), 1)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
