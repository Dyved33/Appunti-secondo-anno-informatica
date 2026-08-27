package product

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/tommasogarofani/ingsw_25-26_group-11/server/domain"
)

// mockProductRepository implementa l'interfaccia productRepository
type mockProductRepository struct {
	mock.Mock
}

// ProductRepository_GetGiacenzaById implements [productRepository].
func (m *mockProductRepository) ProductRepository_GetGiacenzaById(ctx context.Context, prodID int) (*domain.Giacenza, error) {
	args := m.Called(ctx, prodID)
	g, _ := args.Get(0).(*domain.Giacenza)
	return g, args.Error(1)
}

// ProductRepository_UpdateGiacenza implements [productRepository].
func (m *mockProductRepository) ProductRepository_UpdateGiacenza(ctx context.Context, g *domain.Giacenza) error {
	args := m.Called(ctx, g)
	return args.Error(0)
}

func (m *mockProductRepository) ProductRepository_GetAllProducts(ctx context.Context) ([]domain.Prodotto, error) {
	args := m.Called(ctx)
	list, _ := args.Get(0).([]domain.Prodotto)
	return list, args.Error(1)
}

func (m *mockProductRepository) ProductRepository_GetByID(ctx context.Context, id int) (*domain.Prodotto, error) {
	args := m.Called(ctx, id)
	p, _ := args.Get(0).(*domain.Prodotto)
	return p, args.Error(1)
}

func (m *mockProductRepository) ProductRepository_GetByBarcode(ctx context.Context, barcode string) (*domain.Prodotto, error) {
	args := m.Called(ctx, barcode)
	p, _ := args.Get(0).(*domain.Prodotto)
	return p, args.Error(1)
}

func (m *mockProductRepository) ProductRepository_Create(ctx context.Context, p *domain.Prodotto) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *mockProductRepository) ProductRepository_Update(ctx context.Context, p *domain.Prodotto) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *mockProductRepository) ProductRepository_Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockProductRepository) ProductRepository_GetAllGiacenze(ctx context.Context) ([]domain.Giacenza, error) {
	args := m.Called(ctx)
	list, _ := args.Get(0).([]domain.Giacenza)
	return list, args.Error(1)
}

func (m *mockProductRepository) ProductRepository_GetGiacenzaByID(ctx context.Context, productID int) (*domain.Giacenza, error) {
	args := m.Called(ctx, productID)
	g, _ := args.Get(0).(*domain.Giacenza)
	return g, args.Error(1)
}

// fakeTxRunner esegue subito la funzione passata, senza una vera transazione:
// nei test di Service vogliamo isolare la logica di business, non testare
// la gestione delle transazioni (che appartiene ai test di domain.TransactionHelper).
type fakeTxRunner struct{}

func (fakeTxRunner) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

const validEAN13 = "3614273613422" // barcode EAN-13 sintatticamente valido usato nei test

func TestProductService_Create(t *testing.T) {
	tests := []struct {
		name      string
		input     *domain.Prodotto
		mockSetup func(repo *mockProductRepository)
		wantErr   error
	}{
		{
			name:  "barcode valido viene creato",
			input: &domain.Prodotto{Nome: "Vite", Barcode: "  " + validEAN13 + "  "},
			mockSetup: func(repo *mockProductRepository) {
				repo.On("ProductRepository_Create", mock.Anything, mock.MatchedBy(func(p *domain.Prodotto) bool {
					return p.Barcode == validEAN13 // verifica che il trim sia avvenuto prima della Create
				})).Return(nil)
			},
		},
		{
			name:    "barcode non valido viene rifiutato senza toccare il repository",
			input:   &domain.Prodotto{Nome: "Vite", Barcode: "123"},
			wantErr: domain.ErrBarcodeNonValido,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mockProductRepository)
			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}

			svc := NewProductService(repo, fakeTxRunner{})
			err := svc.ProductService_Create(context.Background(), tt.input)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestProductService_GetGiacenzaByID(t *testing.T) {
	t.Run("giacenza esistente", func(t *testing.T) {
		repo := new(mockProductRepository)
		repo.On("ProductRepository_GetGiacenzaById", mock.Anything, 1).Return(&domain.Giacenza{ProdottoID: 1, Quantita: 50}, nil)

		svc := NewProductService(repo, fakeTxRunner{})
		g, err := svc.ProductService_GetGiacenzaByID(context.Background(), 1)

		require.NoError(t, err)
		assert.Equal(t, 50, g.Quantita)
	})

	t.Run("giacenza inesistente restituisce errore", func(t *testing.T) {
		repo := new(mockProductRepository)
		repo.On("ProductRepository_GetGiacenzaById", mock.Anything, 99).Return(nil, domain.ErrProdottoNonTrovato)

		svc := NewProductService(repo, fakeTxRunner{})
		_, err := svc.ProductService_GetGiacenzaByID(context.Background(), 99)

		require.ErrorIs(t, err, domain.ErrProdottoNonTrovato)
	})
}

func TestProductService_Update(t *testing.T) {
	t.Run("barcode non valido viene rifiutato", func(t *testing.T) {
		repo := new(mockProductRepository)
		svc := NewProductService(repo, fakeTxRunner{})

		err := svc.ProductService_Update(context.Background(), &domain.Prodotto{Barcode: "non-valido"})
		require.ErrorIs(t, err, domain.ErrBarcodeNonValido)
		repo.AssertNotCalled(t, "ProductRepository_Update", mock.Anything, mock.Anything)
	})

	t.Run("barcode valido aggiorna il prodotto", func(t *testing.T) {
		repo := new(mockProductRepository)
		repo.On("ProductRepository_Update", mock.Anything, mock.Anything).Return(nil)

		svc := NewProductService(repo, fakeTxRunner{})
		err := svc.ProductService_Update(context.Background(), &domain.Prodotto{Barcode: validEAN13})

		require.NoError(t, err)
		repo.AssertExpectations(t)
	})
}

func TestProductService_GetProductByBarcode(t *testing.T) {
	t.Run("barcode non valido viene rifiutato senza interrogare il repository", func(t *testing.T) {
		repo := new(mockProductRepository)
		svc := NewProductService(repo, fakeTxRunner{})

		_, err := svc.ProductService_GetByBarcode(context.Background(), "abc")
		require.ErrorIs(t, err, domain.ErrBarcodeNonValido)
		repo.AssertNotCalled(t, "ProductRepository_GetByBarcode", mock.Anything, mock.Anything)
	})

	t.Run("barcode valido interroga il repository dopo il trim", func(t *testing.T) {
		repo := new(mockProductRepository)
		repo.On("ProductRepository_GetByBarcode", mock.Anything, validEAN13).
			Return(&domain.Prodotto{Barcode: validEAN13, Nome: "Vite"}, nil)

		svc := NewProductService(repo, fakeTxRunner{})
		p, err := svc.ProductService_GetByBarcode(context.Background(), "  "+validEAN13+"  ")

		require.NoError(t, err)
		assert.Equal(t, "Vite", p.Nome)
		repo.AssertExpectations(t)
	})
}

func TestProductService_Delete(t *testing.T) {
	repo := new(mockProductRepository)
	repo.On("ProductRepository_Delete", mock.Anything, 5).Return(nil)

	svc := NewProductService(repo, fakeTxRunner{})
	err := svc.ProductService_Delete(context.Background(), 5)

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestProductService_GetAllProducts(t *testing.T) {
	repo := new(mockProductRepository)
	expected := []domain.Prodotto{
		{ID: 1, Nome: "Vite"},
		{ID: 2, Nome: "Bullone"},
	}
	repo.On("ProductRepository_GetAllProducts", mock.Anything).Return(expected, nil)

	svc := NewProductService(repo, fakeTxRunner{})
	list, err := svc.ProductService_GetAllProducts(context.Background())

	require.NoError(t, err)
	assert.Len(t, list, 2)
	assert.Equal(t, "Vite", list[0].Nome)
	repo.AssertExpectations(t)
}

func TestProductService_GetAllProducts_Errore(t *testing.T) {
	repo := new(mockProductRepository)
	repo.On("ProductRepository_GetAllProducts", mock.Anything).Return(nil, errors.New("db down"))

	svc := NewProductService(repo, fakeTxRunner{})
	_, err := svc.ProductService_GetAllProducts(context.Background())
	require.Error(t, err)
}

func TestProductService_GetByID(t *testing.T) {
	repo := new(mockProductRepository)
	expected := &domain.Prodotto{ID: 5, Nome: "Vite M6"}
	repo.On("ProductRepository_GetByID", mock.Anything, 5).Return(expected, nil)

	svc := NewProductService(repo, fakeTxRunner{})
	p, err := svc.ProductService_GetByID(context.Background(), 5)

	require.NoError(t, err)
	assert.Equal(t, "Vite M6", p.Nome)
	repo.AssertExpectations(t)
}

func TestProductService_GetByID_NonTrovato(t *testing.T) {
	repo := new(mockProductRepository)
	repo.On("ProductRepository_GetByID", mock.Anything, 999).Return(nil, errors.New("not found"))

	svc := NewProductService(repo, fakeTxRunner{})
	_, err := svc.ProductService_GetByID(context.Background(), 999)
	require.Error(t, err)
}

func TestProductService_GetValoreMagazzino(t *testing.T) {
	repo := new(mockProductRepository)
	cat := &domain.Categoria{Nome: "Ferramenta"}
	products := []domain.Prodotto{
		{ID: 1, Nome: "Vite", SogliaMin: 100, PrezzoVend: 0.25, Categoria: cat},
		{ID: 2, Nome: "Bullone", SogliaMin: 50, PrezzoVend: 0.40, Categoria: cat},
	}
	giacenze := []domain.Giacenza{
		{ProdottoID: 1, Quantita: 30},
		{ProdottoID: 2, Quantita: 200},
	}

	repo.On("ProductRepository_GetAllProducts", mock.Anything).Return(products, nil)
	repo.On("ProductRepository_GetAllGiacenze", mock.Anything).Return(giacenze, nil)

	svc := NewProductService(repo, fakeTxRunner{})
	qty, valore, err := svc.ProductService_GetValoreMagazzino(context.Background())

	require.NoError(t, err)
	assert.Equal(t, 230, qty)
	assert.InDelta(t, 30*0.25+200*0.40, valore, 0.0001)
	repo.AssertExpectations(t)
}

func TestProductService_GetValoreMagazzino_Errore(t *testing.T) {
	repo := new(mockProductRepository)
	repo.On("ProductRepository_GetAllProducts", mock.Anything).Return(nil, errors.New("db down"))

	svc := NewProductService(repo, fakeTxRunner{})
	qty, valore, err := svc.ProductService_GetValoreMagazzino(context.Background())
	require.Error(t, err)
	assert.Equal(t, 0, qty)
	assert.Equal(t, 0.0, valore)
}

func TestProductService_GetAllGiacenze_ErroreGiacenze(t *testing.T) {
	repo := new(mockProductRepository)
	repo.On("ProductRepository_GetAllProducts", mock.Anything).Return([]domain.Prodotto{}, nil)
	repo.On("ProductRepository_GetAllGiacenze", mock.Anything).Return(nil, errors.New("db down"))

	svc := NewProductService(repo, fakeTxRunner{})
	_, err := svc.ProductService_GetAllGiacenze(context.Background())
	require.Error(t, err)
}

//====================================

func TestProductService_UpdateListinoPrices(t *testing.T) {
	t.Run("aggiorna i prezzi mantenendo gli altri campi", func(t *testing.T) {
		repo := new(mockProductRepository)
		existing := &domain.Prodotto{ID: 1, Nome: "Vite", PrezzoAcq: 1.0, PrezzoVend: 2.0}
		repo.On("ProductRepository_GetByID", mock.Anything, 1).Return(existing, nil)
		repo.On("ProductRepository_Update", mock.Anything, mock.MatchedBy(func(p *domain.Prodotto) bool {
			return p.PrezzoAcq == 1.5 && p.PrezzoVend == 3.0 && p.Nome == "Vite"
		})).Return(nil)

		svc := NewProductService(repo, fakeTxRunner{})
		err := svc.ProductService_UpdateListinoPrices(context.Background(), 1, 1.5, 3.0)

		require.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("prodotto inesistente propaga l'errore senza chiamare Update", func(t *testing.T) {
		repo := new(mockProductRepository)
		repo.On("ProductRepository_GetByID", mock.Anything, 99).Return(nil, domain.ErrProdottoNonTrovato)

		svc := NewProductService(repo, fakeTxRunner{})
		err := svc.ProductService_UpdateListinoPrices(context.Background(), 99, 1.0, 2.0)

		require.ErrorIs(t, err, domain.ErrProdottoNonTrovato)
		repo.AssertNotCalled(t, "ProductRepository_Update", mock.Anything, mock.Anything)
	})
}

func TestProductService_GetAllGiacenze(t *testing.T) {
	repo := new(mockProductRepository)

	categoria := &domain.Categoria{Nome: "Ferramenta"}
	products := []domain.Prodotto{
		{ID: 1, Codice: "P001", Nome: "Vite", SogliaMin: 100, PrezzoVend: 0.25, Categoria: categoria},
		{ID: 2, Codice: "P002", Nome: "Bullone", SogliaMin: 50, PrezzoVend: 0.40},
		{ID: 3, Codice: "P003", Nome: "Senza giacenza", SogliaMin: 10, PrezzoVend: 1.0}, // non comparirà nel report
	}
	giacenze := []domain.Giacenza{
		{ProdottoID: 1, Quantita: 30},  // sotto soglia (100)
		{ProdottoID: 2, Quantita: 200}, // sopra soglia
	}

	repo.On("ProductRepository_GetAllProducts", mock.Anything).Return(products, nil)
	repo.On("ProductRepository_GetAllGiacenze", mock.Anything).Return(giacenze, nil)

	svc := NewProductService(repo, fakeTxRunner{})
	report, err := svc.ProductService_GetAllGiacenze(context.Background())

	require.NoError(t, err)
	require.Len(t, report.Articoli, 2, "deve includere solo i prodotti per cui esiste una giacenza")
	assert.Equal(t, 230, report.TotaleQuantita)
	assert.InDelta(t, 30*0.25+200*0.40, report.ValoreTotale, 0.0001)

	require.Len(t, report.ArticoliSottoSoglia, 1)
	assert.Equal(t, 1, report.ArticoliSottoSoglia[0].ProdottoID)
	assert.Equal(t, "Ferramenta", report.ArticoliSottoSoglia[0].Categoria)

	repo.AssertExpectations(t)
}

func TestProductService_GetProdottiSottoSoglia(t *testing.T) {
	t.Run("propaga l'errore del report", func(t *testing.T) {
		repo := new(mockProductRepository)
		repo.On("ProductRepository_GetAllProducts", mock.Anything).Return(nil, errors.New("db down"))

		svc := NewProductService(repo, fakeTxRunner{})
		_, err := svc.ProductService_GetProdottiSottoSoglia(context.Background())

		require.Error(t, err)
	})

	t.Run("restituisce solo gli articoli sotto soglia", func(t *testing.T) {
		repo := new(mockProductRepository)
		repo.On("ProductRepository_GetAllProducts", mock.Anything).Return([]domain.Prodotto{
			{ID: 1, SogliaMin: 100},
		}, nil)
		repo.On("ProductRepository_GetAllGiacenze", mock.Anything).Return([]domain.Giacenza{
			{ProdottoID: 1, Quantita: 5, UpdatedAt: time.Now()},
		}, nil)

		svc := NewProductService(repo, fakeTxRunner{})
		items, err := svc.ProductService_GetProdottiSottoSoglia(context.Background())

		require.NoError(t, err)
		require.Len(t, items, 1)
		assert.True(t, items[0].SottoSoglia)
	})
}