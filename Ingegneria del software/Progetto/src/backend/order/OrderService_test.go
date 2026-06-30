package order

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/tommasogarofani/ingsw_25-26_group-11/src/backend/domain"
)

// mockOrderRepository implementa l'interfaccia orderRepository
type mockOrderRepository struct {
	mock.Mock
}

func (m *mockOrderRepository) OrderRepository_CreateOrdineEntrata(ctx context.Context, o *domain.OrdineEntrata) error {
	args := m.Called(ctx, o)
	return args.Error(0)
}
func (m *mockOrderRepository) OrderRepository_CreateOrdineUscita(ctx context.Context, o *domain.OrdineUscita) error {
	args := m.Called(ctx, o)
	return args.Error(0)
}
func (m *mockOrderRepository) OrderRepository_GetOrdineUscitaByID(ctx context.Context, id int) (*domain.OrdineUscita, error) {
	args := m.Called(ctx, id)
	o, _ := args.Get(0).(*domain.OrdineUscita)
	return o, args.Error(1)
}
func (m *mockOrderRepository) OrderRepository_GetOrdiniUscita(ctx context.Context) ([]domain.OrdineUscita, error) {
	args := m.Called(ctx)
	list, _ := args.Get(0).([]domain.OrdineUscita)
	return list, args.Error(1)
}
func (m *mockOrderRepository) OrderRepository_GetOrdiniUscitaByStato(ctx context.Context, stato domain.StatoOrdine) ([]domain.OrdineUscita, error) {
	args := m.Called(ctx, stato)
	list, _ := args.Get(0).([]domain.OrdineUscita)
	return list, args.Error(1)
}
func (m *mockOrderRepository) OrderRepository_UpdateOrdineUscita(ctx context.Context, o *domain.OrdineUscita) error {
	args := m.Called(ctx, o)
	return args.Error(0)
}
func (m *mockOrderRepository) OrderRepository_CreateDDT(ctx context.Context, d *domain.DDT) error {
	args := m.Called(ctx, d)
	return args.Error(0)
}
func (m *mockOrderRepository) OrderRepository_CountDDTByYear(ctx context.Context, year int) (int, error) {
	args := m.Called(ctx, year)
	return args.Int(0), args.Error(1)
}
func (m *mockOrderRepository) OrderRepository_CreateReso(ctx context.Context, re *domain.Reso, stato string) error {
	args := m.Called(ctx, re, stato)
	return args.Error(0)
}
func (m *mockOrderRepository) OrderRepository_GetResoByID(ctx context.Context, id int) (*domain.Reso, string, error) {
	args := m.Called(ctx, id)
	re, _ := args.Get(0).(*domain.Reso)
	return re, args.String(1), args.Error(2)
}
func (m *mockOrderRepository) OrderRepository_UpdateResoStato(ctx context.Context, id int, stato string) error {
	args := m.Called(ctx, id, stato)
	return args.Error(0)
}
func (m *mockOrderRepository) OrderRepository_GetAllResi(ctx context.Context) ([]map[string]interface{}, error) {
	args := m.Called(ctx)
	list, _ := args.Get(0).([]map[string]interface{})
	return list, args.Error(1)
}
func (m *mockOrderRepository) OrderRepository_CreateFattura(ctx context.Context, f *domain.Fattura) error {
	args := m.Called(ctx, f)
	return args.Error(0)
}
func (m *mockOrderRepository) OrderRepository_CountFattureByYearAndType(ctx context.Context, year int, tipo domain.TipoFattura) (int, error) {
	args := m.Called(ctx, year, tipo)
	return args.Int(0), args.Error(1)
}
func (m *mockOrderRepository) OrderRepository_GetAllFatture(ctx context.Context) ([]domain.Fattura, error) {
	args := m.Called(ctx)
	list, _ := args.Get(0).([]domain.Fattura)
	return list, args.Error(1)
}
func (m *mockOrderRepository) OrderRepository_DeleteFattura(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// mockProductRepoForOrders implementa l'interfaccia productRepoForOrders.
type mockProductRepoForOrders struct {
	mock.Mock
}

func (m *mockProductRepoForOrders) ProductRepository_GetByID(ctx context.Context, id int) (*domain.Prodotto, error) {
	args := m.Called(ctx, id)
	p, _ := args.Get(0).(*domain.Prodotto)
	return p, args.Error(1)
}
func (m *mockProductRepoForOrders) ProductRepository_GetGiacenza(ctx context.Context, prodID int) (*domain.Giacenza, error) {
	args := m.Called(ctx, prodID)
	g, _ := args.Get(0).(*domain.Giacenza)
	return g, args.Error(1)
}
func (m *mockProductRepoForOrders) ProductRepository_UpdateGiacenza(ctx context.Context, prodID int, delta int) error {
	args := m.Called(ctx, prodID, delta)
	return args.Error(0)
}

type fakeTxRunner struct{}

func (fakeTxRunner) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func newOrderServiceWithMocks() (*OrderService, *mockOrderRepository, *mockProductRepoForOrders) {
	orderRepo := new(mockOrderRepository)
	prodRepo := new(mockProductRepoForOrders)
	svc := NewOrderService(orderRepo, prodRepo, fakeTxRunner{})
	return svc, orderRepo, prodRepo
}

func TestOrderService_ConfirmOrdineEntrata(t *testing.T) {
	svc, orderRepo, prodRepo := newOrderServiceWithMocks()

	orderRepo.On("OrderRepository_CreateOrdineEntrata", mock.Anything, mock.MatchedBy(func(o *domain.OrdineEntrata) bool {
		return o.Stato == domain.StatoConfermato
	})).Return(nil)
	prodRepo.On("ProductRepoForOrders_UpdateGiacenza", mock.Anything, 5, 100).Return(nil)
	prodRepo.On("ProductRepoForOrders_UpdateGiacenza", mock.Anything, 6, 50).Return(nil)

	o := &domain.OrdineEntrata{
		Righe: []domain.RigaOrdineEntrata{
			{ProdottoID: 5, Quantita: 100},
			{ProdottoID: 6, Quantita: 50},
		},
	}

	err := svc.OrderService_ConfirmOrdineEntrata(context.Background(), o)
	require.NoError(t, err)
	assert.Equal(t, domain.StatoConfermato, o.Stato)
	orderRepo.AssertExpectations(t)
	prodRepo.AssertExpectations(t)
}

func TestOrderService_ConfirmOrdineEntrata_ErroreUpdateGiacenza(t *testing.T) {
	svc, orderRepo, prodRepo := newOrderServiceWithMocks()

	orderRepo.On("OrderRepository_CreateOrdineEntrata", mock.Anything, mock.Anything).Return(nil)
	prodRepo.On("ProductRepoForOrders_UpdateGiacenza", mock.Anything, 5, 100).Return(errors.New("db down"))

	o := &domain.OrdineEntrata{Righe: []domain.RigaOrdineEntrata{{ProdottoID: 5, Quantita: 100}}}
	err := svc.OrderService_ConfirmOrdineEntrata(context.Background(), o)

	require.Error(t, err)
}

func TestOrderService_ConfirmOrdineUscita(t *testing.T) {
	t.Run("giacenza sufficiente: decrementa e salva l'ordine", func(t *testing.T) {
		svc, orderRepo, prodRepo := newOrderServiceWithMocks()

		prodRepo.On("ProductRepoForOrders_GetGiacenza", mock.Anything, 5).Return(&domain.Giacenza{ProdottoID: 5, Quantita: 100}, nil)
		prodRepo.On("ProductRepoForOrders_UpdateGiacenza", mock.Anything, 5, -10).Return(nil)
		orderRepo.On("OrderRepository_CreateOrdineUscita", mock.Anything, mock.Anything).Return(nil)

		o := &domain.OrdineUscita{Righe: []domain.RigaOrdineUscita{{ProdottoID: 5, Quantita: 10}}}
		err := svc.OrderService_ConfirmOrdineUscita(context.Background(), o)

		require.NoError(t, err)
		assert.Equal(t, domain.StatoConfermato, o.Stato)
		orderRepo.AssertExpectations(t)
		prodRepo.AssertExpectations(t)
	})

	t.Run("giacenza insufficiente: restituisce errore dedicato e non salva l'ordine", func(t *testing.T) {
		svc, orderRepo, prodRepo := newOrderServiceWithMocks()

		prodRepo.On("ProductRepoForOrders_GetGiacenza", mock.Anything, 5).Return(&domain.Giacenza{ProdottoID: 5, Quantita: 3}, nil)
		prodRepo.On("ProductRepoForOrders_GetByID", mock.Anything, 5).Return(&domain.Prodotto{ID: 5, Nome: "Vite M6"}, nil)

		o := &domain.OrdineUscita{Righe: []domain.RigaOrdineUscita{{ProdottoID: 5, Quantita: 10}}}
		err := svc.OrderService_ConfirmOrdineUscita(context.Background(), o)

		require.Error(t, err)
		prodRepo.AssertNotCalled(t, "ProductRepoForOrders_UpdateGiacenza", mock.Anything, mock.Anything, mock.Anything)
		orderRepo.AssertNotCalled(t, "ProductRepoForOrders_CreateOrdineUscita", mock.Anything, mock.Anything)
	})

	t.Run("più righe: si ferma alla prima che fallisce e non decrementa le successive", func(t *testing.T) {
		svc, orderRepo, prodRepo := newOrderServiceWithMocks()

		prodRepo.On("ProductRepoForOrders_GetGiacenza", mock.Anything, 5).Return(&domain.Giacenza{ProdottoID: 5, Quantita: 100}, nil)
		prodRepo.On("ProductRepoForOrders_UpdateGiacenza", mock.Anything, 5, -10).Return(nil)
		prodRepo.On("ProductRepoForOrders_GetGiacenza", mock.Anything, 6).Return(&domain.Giacenza{ProdottoID: 6, Quantita: 0}, nil)
		prodRepo.On("ProductRepoForOrders_GetByID", mock.Anything, 6).Return(&domain.Prodotto{ID: 6, Nome: "Bullone"}, nil)

		o := &domain.OrdineUscita{Righe: []domain.RigaOrdineUscita{
			{ProdottoID: 5, Quantita: 10},
			{ProdottoID: 6, Quantita: 5},
		}}
		err := svc.OrderService_ConfirmOrdineUscita(context.Background(), o)

		require.Error(t, err)
		orderRepo.AssertNotCalled(t, "ProductRepoForOrders_CreateOrdineUscita", mock.Anything, mock.Anything)
	})
}

func TestOrderService_GeneraDDT(t *testing.T) {
	t.Run("ordine non ancora evaso: genera il DDT e marca l'ordine come evaso", func(t *testing.T) {
		svc, orderRepo, _ := newOrderServiceWithMocks()

		o := &domain.OrdineUscita{ID: 20, Stato: domain.StatoConfermato}
		orderRepo.On("ProductRepoForOrders_GetOrdineUscitaByID", mock.Anything, 20).Return(o, nil)
		orderRepo.On("ProductRepoForOrders_CountDDTByYear", mock.Anything, time.Now().Year()).Return(3, nil)
		orderRepo.On("ProductRepoForOrders_CreateDDT", mock.Anything, mock.MatchedBy(func(d *domain.DDT) bool {
			return d.OrdineUscitaID == 20 && d.Mittente == "Magazzino" && d.Destinatario == "Cliente Spa"
		})).Return(nil)
		orderRepo.On("ProductRepoForOrders_UpdateOrdineUscita", mock.Anything, mock.MatchedBy(func(o *domain.OrdineUscita) bool {
			return o.Stato == domain.StatoEvaso && o.DDTID != nil
		})).Return(nil)

		ddt, err := svc.OrderService_GeneraDDT(context.Background(), 20, "Magazzino", "Cliente Spa")

		require.NoError(t, err)
		require.NotNil(t, ddt)
		orderRepo.AssertExpectations(t)
	})

	t.Run("ordine già evaso: restituisce ErrDDTGiaGenerato", func(t *testing.T) {
		svc, orderRepo, _ := newOrderServiceWithMocks()

		orderRepo.On("ProductRepoForOrders_GetOrdineUscitaByID", mock.Anything, 20).
			Return(&domain.OrdineUscita{ID: 20, Stato: domain.StatoEvaso}, nil)

		_, err := svc.OrderService_GeneraDDT(context.Background(), 20, "Magazzino", "Cliente Spa")

		require.ErrorIs(t, err, domain.ErrDDTGiaGenerato)
		orderRepo.AssertNotCalled(t, "ProductRepoForOrders_CreateDDT", mock.Anything, mock.Anything)
	})
}

func TestOrderService_CreateReso(t *testing.T) {
	t.Run("ordine evaso: il reso viene registrato come da_verificare", func(t *testing.T) {
		svc, orderRepo, _ := newOrderServiceWithMocks()

		orderRepo.On("ProductRepoForOrders_GetOrdineUscitaByID", mock.Anything, 20).
			Return(&domain.OrdineUscita{ID: 20, Stato: domain.StatoEvaso}, nil)
		orderRepo.On("ProductRepoForOrders_CreateReso", mock.Anything, mock.Anything, "da_verificare").Return(nil)

		err := svc.OrderService_CreateReso(context.Background(), &domain.Reso{OrdineUscitaID: 20})
		require.NoError(t, err)
		orderRepo.AssertExpectations(t)
	})

	t.Run("ordine non evaso: viene rifiutato", func(t *testing.T) {
		svc, orderRepo, _ := newOrderServiceWithMocks()

		orderRepo.On("ProductRepoForOrders_GetOrdineUscitaByID", mock.Anything, 20).
			Return(&domain.OrdineUscita{ID: 20, Stato: domain.StatoConfermato}, nil)

		err := svc.OrderService_CreateReso(context.Background(), &domain.Reso{OrdineUscitaID: 20})
		require.Error(t, err)
		orderRepo.AssertNotCalled(t, "ProductRepoForOrders_CreateReso", mock.Anything, mock.Anything, mock.Anything)
	})
}

func TestOrderService_VerifyReso(t *testing.T) {
	t.Run("reso da_verificare: incrementa la giacenza e passa a verificato", func(t *testing.T) {
		svc, orderRepo, prodRepo := newOrderServiceWithMocks()

		orderRepo.On("ProductRepoForOrders_GetResoByID", mock.Anything, 3).
			Return(&domain.Reso{ID: 3, ProdottoID: 5, Quantita: 2}, "da_verificare", nil)
		prodRepo.On("ProductRepoForOrders_UpdateGiacenza", mock.Anything, 5, 2).Return(nil)
		orderRepo.On("ProductRepoForOrders_UpdateResoStato", mock.Anything, 3, "verificato").Return(nil)

		err := svc.OrderService_VerifyReso(context.Background(), 3)
		require.NoError(t, err)
		orderRepo.AssertExpectations(t)
		prodRepo.AssertExpectations(t)
	})

	t.Run("reso già elaborato: viene rifiutato", func(t *testing.T) {
		svc, orderRepo, prodRepo := newOrderServiceWithMocks()

		orderRepo.On("ProductRepoForOrders_GetResoByID", mock.Anything, 3).
			Return(&domain.Reso{ID: 3}, "verificato", nil)

		err := svc.OrderService_VerifyReso(context.Background(), 3)
		require.Error(t, err)
		prodRepo.AssertNotCalled(t, "ProductRepoForOrders_UpdateGiacenza", mock.Anything, mock.Anything, mock.Anything)
	})
}

func TestOrderService_ApproveReso(t *testing.T) {
	t.Run("reso verificato viene approvato", func(t *testing.T) {
		svc, orderRepo, _ := newOrderServiceWithMocks()

		orderRepo.On("ProductRepoForOrders_GetResoByID", mock.Anything, 3).Return(&domain.Reso{ID: 3}, "verificato", nil)
		orderRepo.On("ProductRepoForOrders_UpdateResoStato", mock.Anything, 3, "approvato").Return(nil)

		err := svc.OrderService_ApproveReso(context.Background(), 3)
		require.NoError(t, err)
		orderRepo.AssertExpectations(t)
	})

	t.Run("reso non ancora verificato viene rifiutato", func(t *testing.T) {
		svc, orderRepo, _ := newOrderServiceWithMocks()

		orderRepo.On("ProductRepoForOrders_GetResoByID", mock.Anything, 3).Return(&domain.Reso{ID: 3}, "da_verificare", nil)

		err := svc.OrderService_ApproveReso(context.Background(), 3)
		require.Error(t, err)
		orderRepo.AssertNotCalled(t, "ProductRepoForOrders_UpdateResoStato", mock.Anything, mock.Anything, mock.Anything)
	})
}

func TestOrderService_RifiutaReso(t *testing.T) {
	t.Run("reso verificato viene rifiutato e la giacenza torna indietro", func(t *testing.T) {
		svc, orderRepo, prodRepo := newOrderServiceWithMocks()

		orderRepo.On("ProductRepoForOrders_GetResoByID", mock.Anything, 3).
			Return(&domain.Reso{ID: 3, ProdottoID: 5, Quantita: 2}, "verificato", nil)
		prodRepo.On("ProductRepoForOrders_UpdateGiacenza", mock.Anything, 5, -2).Return(nil)
		orderRepo.On("ProductRepoForOrders_UpdateResoStato", mock.Anything, 3, "rifiutato").Return(nil)

		err := svc.OrderService_RifiutaReso(context.Background(), 3)
		require.NoError(t, err)
		orderRepo.AssertExpectations(t)
		prodRepo.AssertExpectations(t)
	})

	t.Run("reso non verificato non può essere rifiutato", func(t *testing.T) {
		svc, orderRepo, prodRepo := newOrderServiceWithMocks()

		orderRepo.On("GetResoByID", mock.Anything, 3).Return(&domain.Reso{ID: 3}, "da_verificare", nil)

		err := svc.OrderService_RifiutaReso(context.Background(), 3)
		require.Error(t, err)
		prodRepo.AssertNotCalled(t, "UpdateGiacenza", mock.Anything, mock.Anything, mock.Anything)
	})
}

func TestOrderService_GeneraFattura(t *testing.T) {
	t.Run("ordine evaso: calcola imponibile, IVA al 22% e totale", func(t *testing.T) {
		svc, orderRepo, _ := newOrderServiceWithMocks()

		o := &domain.OrdineUscita{
			ID:    20,
			Stato: domain.StatoEvaso,
			Righe: []domain.RigaOrdineUscita{
				{Quantita: 10, PrezzoUnitario: 5.0},  // 50
				{Quantita: 2, PrezzoUnitario: 25.0},  // 50
			},
		}
		orderRepo.On("GetOrdineUscitaByID", mock.Anything, 20).Return(o, nil)
		orderRepo.On("CountFattureByYearAndType", mock.Anything, time.Now().Year(), domain.TipoFatturaVendita).Return(0, nil)
		orderRepo.On("CreateFattura", mock.Anything, mock.MatchedBy(func(f *domain.Fattura) bool {
			return f.Imponibile == 100.0 &&
				assertAlmostEqual(f.IVA, 22.0) &&
				assertAlmostEqual(f.Totale, 122.0) &&
				f.Tipo == domain.TipoFatturaVendita
		})).Return(nil)

		f, err := svc.OrderService_GeneraFattura(context.Background(), 20)
		require.NoError(t, err)
		require.NotNil(t, f)
		orderRepo.AssertExpectations(t)
	})

	t.Run("ordine non evaso: la fattura viene rifiutata", func(t *testing.T) {
		svc, orderRepo, _ := newOrderServiceWithMocks()

		orderRepo.On("GetOrdineUscitaByID", mock.Anything, 20).
			Return(&domain.OrdineUscita{ID: 20, Stato: domain.StatoConfermato}, nil)

		_, err := svc.OrderService_GeneraFattura(context.Background(), 20)
		require.Error(t, err)
		orderRepo.AssertNotCalled(t, "CreateFattura", mock.Anything, mock.Anything)
	})
}

func TestOrderService_GeneraNotaCredito(t *testing.T) {
	svc, orderRepo, _ := newOrderServiceWithMocks()

	o := &domain.OrdineUscita{
		ID:    20,
		Stato: domain.StatoEvaso,
		Righe: []domain.RigaOrdineUscita{{Quantita: 1, PrezzoUnitario: 10.0}},
	}
	orderRepo.On("GetOrdineUscitaByID", mock.Anything, 20).Return(o, nil)
	orderRepo.On("CountFattureByYearAndType", mock.Anything, time.Now().Year(), domain.TipoNotaCredito).Return(0, nil)
	orderRepo.On("CreateFattura", mock.Anything, mock.MatchedBy(func(f *domain.Fattura) bool {
		return f.Tipo == domain.TipoNotaCredito
	})).Return(nil)

	f, err := svc.OrderService_GeneraNotaCredito(context.Background(), 20)
	require.NoError(t, err)
	assert.Equal(t, domain.TipoNotaCredito, f.Tipo)
}

func TestOrderService_GetNoteCredito(t *testing.T) {
	svc, orderRepo, _ := newOrderServiceWithMocks()

	orderRepo.On("GetAllFatture", mock.Anything).Return([]domain.Fattura{
		{ID: 1, Tipo: domain.TipoFatturaVendita},
		{ID: 2, Tipo: domain.TipoNotaCredito},
	}, nil)

	list, err := svc.OrderService_GetNoteCredito(context.Background())
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, domain.TipoNotaCredito, list[0].Tipo)
}

// assertAlmostEqual confronta due float con una piccola tolleranza per evitare
// falsi negativi dovuti ad arrotondamenti in virgola mobile.
func assertAlmostEqual(a, b float64) bool {
	const eps = 0.0001
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff < eps
}