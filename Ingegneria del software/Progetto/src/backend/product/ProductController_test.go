package product

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/tommasogarofani/ingsw_25-26_group-11/src/backend/domain"
)

// mockProductServiceForController implementa l'interfaccia productService
type mockProductServiceForController struct {
	mock.Mock
}

func (m *mockProductServiceForController) ProductService_Create(ctx context.Context, p *domain.Prodotto) error {
	args := m.Called(ctx, p)
	if args.Get(0) == nil {
		p.ID = 1 // simula l'assegnazione dell'ID da parte del repository
		return nil
	}
	return args.Error(0)
}

func (m *mockProductServiceForController) ProductService_GetProducts(ctx context.Context) ([]domain.Prodotto, error) {
	args := m.Called(ctx)
	list, _ := args.Get(0).([]domain.Prodotto)
	return list, args.Error(1)
}

func (m *mockProductServiceForController) ProductService_GetByBarcode(ctx context.Context, barcode string) (*domain.Prodotto, error) {
	args := m.Called(ctx, barcode)
	p, _ := args.Get(0).(*domain.Prodotto)
	return p, args.Error(1)
}

func (m *mockProductServiceForController) ProductService_Update(ctx context.Context, p *domain.Prodotto) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *mockProductServiceForController) ProductService_Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockProductServiceForController) ProductService_UpdateListinoPrices(ctx context.Context, productID int, acqPrice, vendPrice float64) error {
	args := m.Called(ctx, productID, acqPrice, vendPrice)
	return args.Error(0)
}

func (m *mockProductServiceForController) ProductService_GetValoreMagazzino(ctx context.Context) (int, float64, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Get(1).(float64), args.Error(2)
}

func (m *mockProductServiceForController) ProductService_GetInventarioReport(ctx context.Context) (*domain.InventoryReport, error) {
	args := m.Called(ctx)
	r, _ := args.Get(0).(*domain.InventoryReport)
	return r, args.Error(1)
}

func (m *mockProductServiceForController) ProductService_GetProdottiSottoSoglia(ctx context.Context) ([]domain.InventoryItem, error) {
	args := m.Called(ctx)
	list, _ := args.Get(0).([]domain.InventoryItem)
	return list, args.Error(1)
}

func TestProductController_AddProduct(t *testing.T) {
	t.Run("mappa correttamente l'input nel domain object e il risultato nel DTO", func(t *testing.T) {
		svc := new(mockProductServiceForController)
		svc.On("ProductService_Create", mock.Anything, mock.MatchedBy(func(p *domain.Prodotto) bool {
			return p.Nome == "Vite" && p.Barcode == "8001234567890" && p.CategoriaID == 3
		})).Return(nil)

		c := NewProductController(svc)
		out, err := c.ProductController_AddProduct(context.Background(), ProductInputDTO{
			Codice:      "P001",
			Barcode:     "8001234567890",
			Nome:        "Vite",
			CategoriaID: 3,
			PrezzoAcq:   0.1,
			PrezzoVend:  0.25,
			SogliaMin:   100,
		})

		require.NoError(t, err)
		assert.Equal(t, 1, out.ID) // assegnato dal mock per simulare il repository
		assert.Equal(t, "Vite", out.Nome)
		svc.AssertExpectations(t)
	})

	t.Run("propaga l'errore del service", func(t *testing.T) {
		svc := new(mockProductServiceForController)
		svc.On("ProductService_Create", mock.Anything, mock.Anything).Return(domain.ErrBarcodeNonValido)

		c := NewProductController(svc)
		_, err := c.ProductController_AddProduct(context.Background(), ProductInputDTO{Barcode: "non valido"})

		require.ErrorIs(t, err, domain.ErrBarcodeNonValido)
	})
}

func TestProductController_GetProducts(t *testing.T) {
	svc := new(mockProductServiceForController)
	svc.On("ProductService_GetProducts", mock.Anything).Return([]domain.Prodotto{
		{ID: 1, Nome: "Vite"},
		{ID: 2, Nome: "Bullone"},
	}, nil)

	c := NewProductController(svc)
	out, err := c.ProductController_GetProducts(context.Background())

	require.NoError(t, err)
	require.Len(t, out, 2)
	assert.Equal(t, "Vite", out[0].Nome)
}

func TestProductController_GetProductByBarcode(t *testing.T) {
	t.Run("rimuove gli spazi prima di chiamare il service", func(t *testing.T) {
		svc := new(mockProductServiceForController)
		svc.On("ProductService_GetByBarcode", mock.Anything, "8001234567890").
			Return(&domain.Prodotto{Nome: "Vite", Barcode: "8001234567890"}, nil)

		c := NewProductController(svc)
		out, err := c.ProductController_GetProductByBarcode(context.Background(), "  8001234567890  ")

		require.NoError(t, err)
		assert.Equal(t, "Vite", out.Nome)
		svc.AssertExpectations(t)
	})

	t.Run("propaga l'errore prodotto non trovato", func(t *testing.T) {
		svc := new(mockProductServiceForController)
		svc.On("ProductService_GetByBarcode", mock.Anything, "8000000000000").
			Return(nil, domain.ErrProdottoNonTrovato)

		c := NewProductController(svc)
		_, err := c.ProductController_GetProductByBarcode(context.Background(), "8000000000000")

		require.ErrorIs(t, err, domain.ErrProdottoNonTrovato)
	})
}

func TestProductController_EditProduct(t *testing.T) {
	svc := new(mockProductServiceForController)
	svc.On("ProductService_Update", mock.Anything, mock.MatchedBy(func(p *domain.Prodotto) bool {
		return p.ID == 7 && p.Nome == "Vite aggiornata"
	})).Return(nil)

	c := NewProductController(svc)
	err := c.ProductController_EditProduct(context.Background(), 7, ProductInputDTO{Nome: "Vite aggiornata"})

	require.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestProductController_RemoveProduct(t *testing.T) {
	svc := new(mockProductServiceForController)
	svc.On("ProductService_Delete", mock.Anything, 4).Return(nil)

	c := NewProductController(svc)
	err := c.ProductController_RemoveProduct(context.Background(), 4)

	require.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestProductController_UpdatePrices(t *testing.T) {
	svc := new(mockProductServiceForController)
	svc.On("ProductService_UpdateListinoPrices", mock.Anything, 1, 1.5, 3.0).Return(nil)

	c := NewProductController(svc)
	err := c.ProductController_UpdatePrices(context.Background(), 1, PriceUpdateInputDTO{PrezzoAcq: 1.5, PrezzoVend: 3.0})

	require.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestProductController_GetInventoryValue(t *testing.T) {
	svc := new(mockProductServiceForController)
	svc.On("ProductService_GetValoreMagazzino", mock.Anything).Return(150, 999.99, nil)

	c := NewProductController(svc)
	qty, value, err := c.ProductController_GetInventoryValue(context.Background())

	require.NoError(t, err)
	assert.Equal(t, 150, qty)
	assert.Equal(t, 999.99, value)
}

func TestProductController_GetInventoryReport(t *testing.T) {
	want := &domain.InventoryReport{TotaleQuantita: 10}
	svc := new(mockProductServiceForController)
	svc.On("ProductService_GetInventarioReport", mock.Anything).Return(want, nil)

	c := NewProductController(svc)
	got, err := c.ProductController_GetInventoryReport(context.Background())

	require.NoError(t, err)
	assert.Same(t, want, got)
}

func TestProductController_GetLowStockProducts(t *testing.T) {
	svc := new(mockProductServiceForController)
	svc.On("ProductService_GetProdottiSottoSoglia", mock.Anything).Return([]domain.InventoryItem{
		{ProdottoID: 1, SottoSoglia: true},
	}, nil)

	c := NewProductController(svc)
	got, err := c.ProductController_GetLowStockProducts(context.Background())

	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.True(t, got[0].SottoSoglia)
}

func TestProductController_GetLowStockProducts_PropagaErrore(t *testing.T) {
	svc := new(mockProductServiceForController)
	svc.On("ProductService_GetProdottiSottoSoglia", mock.Anything).Return(nil, errors.New("db down"))

	c := NewProductController(svc)
	_, err := c.ProductController_GetLowStockProducts(context.Background())

	require.Error(t, err)
}
