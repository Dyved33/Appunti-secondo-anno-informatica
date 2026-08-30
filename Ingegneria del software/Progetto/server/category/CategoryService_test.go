package category

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/tommasogarofani/ingsw_25-26_group-11/server/domain"
)

// mockCategoryRepository implementa l'interfaccia categoryRepository
type mockCategoryRepository struct {
	mock.Mock
}

func (m *mockCategoryRepository) CategoryRepository_Create(ctx context.Context, c *domain.Categoria) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *mockCategoryRepository) CategoryRepository_GetAll(ctx context.Context) ([]domain.Categoria, error) {
	args := m.Called(ctx)
	list, _ := args.Get(0).([]domain.Categoria)
	return list, args.Error(1)
}

func (m *mockCategoryRepository) CategoryRepository_GetByID(ctx context.Context, id int) (*domain.Categoria, error) {
	args := m.Called(ctx, id)
	c, _ := args.Get(0).(*domain.Categoria)
	return c, args.Error(1)
}

func (m *mockCategoryRepository) CategoryRepository_Update(ctx context.Context, c *domain.Categoria) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *mockCategoryRepository) CategoryRepository_Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type fakeTxRunner struct{}

func (fakeTxRunner) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func TestCategoryService_Create(t *testing.T) {
	tests := []struct {
		name      string
		nome      string
		mockSetup func(repo *mockCategoryRepository)
		wantErr   bool
	}{
		{
			name: "nome valido viene creato dopo il trim",
			nome: "  Ferramenta  ",
			mockSetup: func(repo *mockCategoryRepository) {
				repo.On("CategoryRepository_Create", mock.Anything, mock.MatchedBy(func(c *domain.Categoria) bool {
					return c.Nome == "Ferramenta"
				})).Return(nil)
			},
		},
		{
			name:    "nome vuoto viene rifiutato",
			nome:    "   ",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mockCategoryRepository)
			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}

			svc := NewCategoryService(repo, fakeTxRunner{})
			err := svc.CategoryService_Create(context.Background(), &domain.Categoria{Nome: tt.nome})

			if tt.wantErr {
				require.Error(t, err)
				repo.AssertNotCalled(t, "CategoryRepository_Create", mock.Anything, mock.Anything)
			} else {
				require.NoError(t, err)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestCategoryService_GetAll(t *testing.T) {
	repo := new(mockCategoryRepository)
	repo.On("CategoryRepository_GetAll", mock.Anything).Return([]domain.Categoria{{ID: 1, Nome: "Elettronica"}}, nil)

	svc := NewCategoryService(repo, fakeTxRunner{})
	got, err := svc.CategoryService_GetAll(context.Background())

	require.NoError(t, err)
	require.Len(t, got, 1)
}

func TestCategoryService_GetByID(t *testing.T) {
	repo := new(mockCategoryRepository)
	repo.On("CategoryRepository_GetByID", mock.Anything, 1).Return(&domain.Categoria{ID: 1, Nome: "Elettronica"}, nil)

	svc := NewCategoryService(repo, fakeTxRunner{})
	got, err := svc.CategoryService_GetByID(context.Background(), 1)

	require.NoError(t, err)
	require.Equal(t, "Elettronica", got.Nome)
}

func TestCategoryService_Update(t *testing.T) {
	t.Run("nome vuoto dopo il trim viene rifiutato", func(t *testing.T) {
		repo := new(mockCategoryRepository)
		svc := NewCategoryService(repo, fakeTxRunner{})

		err := svc.CategoryService_Update(context.Background(), &domain.Categoria{ID: 1, Nome: "   "})
		require.Error(t, err)
		repo.AssertNotCalled(t, "CategoryRepository_Update", mock.Anything, mock.Anything)
	})

	t.Run("nome valido aggiorna la categoria", func(t *testing.T) {
		repo := new(mockCategoryRepository)
		repo.On("CategoryRepository_Update", mock.Anything, mock.Anything).Return(nil)

		svc := NewCategoryService(repo, fakeTxRunner{})
		err := svc.CategoryService_Update(context.Background(), &domain.Categoria{ID: 1, Nome: "Elettronica"})

		require.NoError(t, err)
		repo.AssertExpectations(t)
	})
}

func TestCategoryService_Delete(t *testing.T) {
	repo := new(mockCategoryRepository)
	repo.On("CategoryRepository_Delete", mock.Anything, 1).Return(nil)

	svc := NewCategoryService(repo, fakeTxRunner{})
	err := svc.CategoryService_Delete(context.Background(), 1)

	require.NoError(t, err)
	repo.AssertExpectations(t)
}
