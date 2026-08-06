package category

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/tommasogarofani/ingsw_25-26_group-11/server/domain"
)

// mockCategoryServiceForController implementa l'interfaccia categoryService
type mockCategoryServiceForController struct {
	mock.Mock
}

// CategoryService_Create implements [categoryService].
func (m *mockCategoryServiceForController) CategoryService_Create(ctx context.Context, c *domain.Categoria) error {
	args := m.Called(ctx, c)
	if args.Error(0) == nil {
		c.ID = 1 // simula l'assegnazione dell'ID da parte del repository
	}
	return args.Error(0)
}

// CategoryService_Delete implements [categoryService].
func (m *mockCategoryServiceForController) CategoryService_Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// CategoryService_GetAll implements [categoryService].
func (m *mockCategoryServiceForController) CategoryService_GetAll(ctx context.Context) ([]domain.Categoria, error) {
	args := m.Called(ctx)
	list, _ := args.Get(0).([]domain.Categoria)
	return list, args.Error(1)
}

// CategoryService_GetByID implements [categoryService].
func (m *mockCategoryServiceForController) CategoryService_GetByID(ctx context.Context, id int) (*domain.Categoria, error) {
	args := m.Called(ctx, id)
	c, _ := args.Get(0).(*domain.Categoria)
	return c, args.Error(1)
}

// CategoryService_Update implements [categoryService].
func (m *mockCategoryServiceForController) CategoryService_Update(ctx context.Context, c *domain.Categoria) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func TestCategoryController_Create(t *testing.T) {
	t.Run("mappa l'input e il risultato nel DTO", func(t *testing.T) {
		svc := new(mockCategoryServiceForController)
		svc.On("CategoryService_Create", mock.Anything, mock.MatchedBy(func(c *domain.Categoria) bool {
			return c.Nome == "Ferramenta" && c.Descrizione == "Minuteria metallica"
		})).Return(nil)

		c := NewCategoryController(svc)
		out, err := c.CategoryController_Create(context.Background(), CategoryInputDTO{Nome: "Ferramenta", Descrizione: "Minuteria metallica"})

		require.NoError(t, err)
		assert.Equal(t, 1, out.ID)
		assert.Equal(t, "Ferramenta", out.Nome)
		svc.AssertExpectations(t)
	})

	t.Run("propaga l'errore del service", func(t *testing.T) {
		svc := new(mockCategoryServiceForController)
		svc.On("CategoryService_Create", mock.Anything, mock.Anything).Return(errors.New("nome vuoto"))

		c := NewCategoryController(svc)
		_, err := c.CategoryController_Create(context.Background(), CategoryInputDTO{Nome: "   "})

		require.Error(t, err)
	})
}

func TestCategoryController_GetAll(t *testing.T) {
	svc := new(mockCategoryServiceForController)
	svc.On("CategoryService_GetAll", mock.Anything).Return([]domain.Categoria{
		{ID: 1, Nome: "Elettronica"},
		{ID: 2, Nome: "Ferramenta"},
	}, nil)

	c := NewCategoryController(svc)
	out, err := c.CategoryController_GetAll(context.Background())

	require.NoError(t, err)
	require.Len(t, out, 2)
	assert.Equal(t, "Elettronica", out[0].Nome)
}

func TestCategoryController_GetByID(t *testing.T) {
	t.Run("categoria trovata", func(t *testing.T) {
		svc := new(mockCategoryServiceForController)
		svc.On("CategoryService_GetByID", mock.Anything, 1).Return(&domain.Categoria{ID: 1, Nome: "Elettronica"}, nil)

		c := NewCategoryController(svc)
		out, err := c.CategoryController_GetByID(context.Background(), 1)

		require.NoError(t, err)
		assert.Equal(t, "Elettronica", out.Nome)
	})

	t.Run("categoria non trovata propaga l'errore", func(t *testing.T) {
		svc := new(mockCategoryServiceForController)
		svc.On("CategoryService_GetByID", mock.Anything, 99).Return(nil, domain.ErrCategoriaNonTrovata)

		c := NewCategoryController(svc)
		_, err := c.CategoryController_GetByID(context.Background(), 99)

		require.ErrorIs(t, err, domain.ErrCategoriaNonTrovata)
	})
}

func TestCategoryController_Update(t *testing.T) {
	svc := new(mockCategoryServiceForController)
	svc.On("CategoryService_Update", mock.Anything, mock.MatchedBy(func(c *domain.Categoria) bool {
		return c.ID == 5 && c.Nome == "Nuovo nome"
	})).Return(nil)

	c := NewCategoryController(svc)
	err := c.CategoryController_Update(context.Background(), 5, CategoryInputDTO{Nome: "Nuovo nome"})

	require.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestCategoryController_Delete(t *testing.T) {
	svc := new(mockCategoryServiceForController)
	svc.On("CategoryService_Delete", mock.Anything, 5).Return(nil)

	c := NewCategoryController(svc)
	err := c.CategoryController_Delete(context.Background(), 5)

	require.NoError(t, err)
	svc.AssertExpectations(t)
}
