package order

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

// mockOrderRepository implementa l'interfaccia orderRepository
type mockOrderRepository struct {
	mock.Mock
}

// OrderRepository_GetOrdineEntrataByID implements [orderRepository].
func (m *mockOrderRepository) OrderRepository_GetOrdineEntrataByID(ctx context.Context, id int) (*domain.OrdineEntrata, error) {
	args := m.Called(ctx, id)
	o, _ := args.Get(0).(*domain.OrdineEntrata)
	return o, args.Error(1)
}

// OrderRepository_GetOrdiniEntrata implements [orderRepository].
func (m *mockOrderRepository) OrderRepository_GetOrdiniEntrata(ctx context.Context) ([]domain.OrdineEntrata, error) {
	args := m.Called(ctx)
	list, _ := args.Get(0).([]domain.OrdineEntrata)
	return list, args.Error(1)
}

func (m *mockOrderRepository) OrderRepository_CreateOrdineEntrata(ctx context.Context, o *domain.OrdineEntrata) error {
	args := m.Called(ctx, o)
	return args.Error(0)
}

func (m *mockOrderRepository) OrderRepository_DeleteOrdineEntrata(ctx context.Context, o *domain.OrdineEntrata) error {
	args := m.Called(ctx, o)
	return args.Error(0)
}

//==========================================================================================

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

func (m *mockOrderRepository) OrderRepository_GetOrdiniUscitaByStato(ctx context.Context, stato domain.StatoOrdineUscita) ([]domain.OrdineUscita, error) {
	args := m.Called(ctx, stato)
	list, _ := args.Get(0).([]domain.OrdineUscita)
	return list, args.Error(1)
}

func (m *mockOrderRepository) OrderRepository_UpdateOrdineUscita(ctx context.Context, o *domain.OrdineUscita) error {
	args := m.Called(ctx, o)
	return args.Error(0)
}

func (m *mockOrderRepository) OrderRepository_DeleteOrdineUscita(ctx context.Context, o *domain.OrdineUscita) error {
	args := m.Called(ctx, o)
	return args.Error(0)
}

//==========================================================================================

func (m *mockOrderRepository) OrderRepository_CreateDDT(ctx context.Context, d *domain.DDT) error {
	args := m.Called(ctx, d)
	return args.Error(0)
}

func (m *mockOrderRepository) OrderRepository_GetDDTByYear(ctx context.Context, year int) ([]domain.DDT, error) {
	args := m.Called(ctx, year)
	list, _ := args.Get(0).([]domain.DDT)
	return list, args.Error(1)
}

func (m *mockOrderRepository) OrderRepository_GetDDTByOrdineUscita(ctx context.Context, orderID int) (*domain.DDT, error) {
	args := m.Called(ctx, orderID)
	return args.Get(0).(*domain.DDT), args.Error(1)
}

//==========================================================================================

func (m *mockOrderRepository) OrderRepository_CreateReso(ctx context.Context, re *domain.Reso) error {
	args := m.Called(ctx, re)
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

func (m *mockOrderRepository) OrderRepository_GetAllResi(ctx context.Context) ([]domain.Reso, error) {
	args := m.Called(ctx)
	list, _ := args.Get(0).([]domain.Reso)
	return list, args.Error(1)
}

func (m *mockOrderRepository) OrderRepository_DeleteReso(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

//==========================================================================================

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

//==========================================================================================

// mockProductRepoForOrders implementa l'interfaccia productRepoForOrders.
type mockProductRepoForOrders struct {
	mock.Mock
}

func (m *mockProductRepoForOrders) ProductRepository_GetByID(ctx context.Context, id int) (*domain.Prodotto, error) {
	args := m.Called(ctx, id)
	p, _ := args.Get(0).(*domain.Prodotto)
	return p, args.Error(1)
}
func (m *mockProductRepoForOrders) ProductRepository_GetGiacenzaById(ctx context.Context, prodID int) (*domain.Giacenza, error) {
	args := m.Called(ctx, prodID)
	g, _ := args.Get(0).(*domain.Giacenza)
	return g, args.Error(1)
}
func (m *mockProductRepoForOrders) ProductRepository_UpdateGiacenza(ctx context.Context, g *domain.Giacenza) error {
	args := m.Called(ctx, g)
	return args.Error(0)
}

//==========================================================================================

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
		return o.ID == 0 && len(o.Righe) == 2
	})).Return(nil)
	prodRepo.On("ProductRepository_UpdateGiacenza", mock.Anything, mock.MatchedBy(func(g *domain.Giacenza) bool {
		return g.ProdottoID == 5 && g.Quantita == 100
	})).Return(nil)
	prodRepo.On("ProductRepository_UpdateGiacenza", mock.Anything, mock.MatchedBy(func(g *domain.Giacenza) bool {
		return g.ProdottoID == 6 && g.Quantita == 50
	})).Return(nil)

	o := &domain.OrdineEntrata{
		Righe: []domain.RigaOrdineEntrata{
			{ProdottoID: 5, Quantita: 100},
			{ProdottoID: 6, Quantita: 50},
		},
	}

	err := svc.OrderService_ConfirmOrdineEntrata(context.Background(), o)
	require.NoError(t, err)
	orderRepo.AssertExpectations(t)
	prodRepo.AssertExpectations(t)
}

func TestOrderService_ConfirmOrdineEntrata_ErroreUpdateGiacenza(t *testing.T) {
	svc, orderRepo, prodRepo := newOrderServiceWithMocks()

	orderRepo.On("OrderRepository_CreateOrdineEntrata", mock.Anything, mock.Anything).Return(nil)
	prodRepo.On("ProductRepository_UpdateGiacenza", mock.Anything, &domain.Giacenza{ProdottoID: 5, Quantita: 100}).Return(errors.New("db down"))

	o := &domain.OrdineEntrata{Righe: []domain.RigaOrdineEntrata{{ProdottoID: 5, Quantita: 100}}}
	err := svc.OrderService_ConfirmOrdineEntrata(context.Background(), o)

	require.Error(t, err)
}

func TestOrderService_DeleteOrdineEntrata(t *testing.T) {
	svc, orderRepo, _ := newOrderServiceWithMocks()

	o := &domain.OrdineEntrata{ID: 1, Numero: "OE-001"}
	orderRepo.On("OrderRepository_DeleteOrdineEntrata", mock.Anything, o).Return(nil)

	err := svc.OrderService_DeleteOrdineEntrata(context.Background(), o)
	require.NoError(t, err)
	orderRepo.AssertExpectations(t)
}

func TestOrderService_DeleteOrdineEntrata_Errore(t *testing.T) {
	svc, orderRepo, _ := newOrderServiceWithMocks()

	o := &domain.OrdineEntrata{ID: 1}
	orderRepo.On("OrderRepository_DeleteOrdineEntrata", mock.Anything, o).Return(errors.New("db error"))

	err := svc.OrderService_DeleteOrdineEntrata(context.Background(), o)
	require.Error(t, err)
}

func TestOrderService_DeleteOrdineUscita(t *testing.T) {
	svc, orderRepo, _ := newOrderServiceWithMocks()

	o := &domain.OrdineUscita{ID: 1, Numero: "OU-001"}
	orderRepo.On("OrderRepository_DeleteOrdineUscita", mock.Anything, o).Return(nil)

	err := svc.OrderService_DeleteOrdineUscita(context.Background(), o)
	require.NoError(t, err)
	orderRepo.AssertExpectations(t)
}

func TestOrderService_DeleteOrdineUscita_Errore(t *testing.T) {
	svc, orderRepo, _ := newOrderServiceWithMocks()

	o := &domain.OrdineUscita{ID: 1}
	orderRepo.On("OrderRepository_DeleteOrdineUscita", mock.Anything, o).Return(errors.New("db error"))

	err := svc.OrderService_DeleteOrdineUscita(context.Background(), o)
	require.Error(t, err)
}

func TestOrderService_DeleteFattura(t *testing.T) {
	svc, orderRepo, _ := newOrderServiceWithMocks()

	orderRepo.On("OrderRepository_DeleteFattura", mock.Anything, 5).Return(nil)

	err := svc.OrderService_DeleteFattura(context.Background(), 5)
	require.NoError(t, err)
	orderRepo.AssertExpectations(t)
}

func TestOrderService_DeleteFattura_Errore(t *testing.T) {
	svc, orderRepo, _ := newOrderServiceWithMocks()

	orderRepo.On("OrderRepository_DeleteFattura", mock.Anything, 5).Return(errors.New("db error"))

	err := svc.OrderService_DeleteFattura(context.Background(), 5)
	require.Error(t, err)
}

func TestOrderService_GetOrdiniEntrata(t *testing.T) {
	svc, orderRepo, _ := newOrderServiceWithMocks()

	expected := []domain.OrdineEntrata{
		{ID: 1},
		{ID: 2},
	}
	orderRepo.On("OrderRepository_GetOrdiniEntrata", mock.Anything).Return(expected, nil)

	list, err := svc.OrderService_GetOrdiniEntrata(context.Background())
	require.NoError(t, err)
	assert.Equal(t, expected, list)
}

func TestOrderService_GetOrdiniEntrataById(t *testing.T) {
	svc, orderRepo, _ := newOrderServiceWithMocks()

	expected := domain.OrdineEntrata{ID: 1}
	orderRepo.On("OrderRepository_GetOrdineEntrataByID", mock.Anything, 1).Return(&expected, nil)

	o, err := svc.OrderService_GetOrdiniEntrataById(context.Background(), 1)
	require.NoError(t, err)
	assert.Equal(t, &expected, o)
}

//==========================================================================================

func TestOrderService_ConfirmOrdineUscita(t *testing.T) {
	t.Run("giacenza sufficiente: conferma l'ordine e scarica la merce", func(t *testing.T) {
		svc, orderRepo, prodRepo := newOrderServiceWithMocks()

		prodRepo.On("ProductRepository_GetGiacenzaById", mock.Anything, 5).Return(&domain.Giacenza{ProdottoID: 5, Quantita: 100}, nil)
		prodRepo.On("ProductRepository_UpdateGiacenza", mock.Anything, &domain.Giacenza{ProdottoID: 5, Quantita: -10}).Return(nil)
		orderRepo.On("OrderRepository_CreateOrdineUscita", mock.Anything, mock.MatchedBy(func(o *domain.OrdineUscita) bool {
			return o.Stato == domain.StatoConfermato
		})).Return(nil)

		o := &domain.OrdineUscita{
			Righe: []domain.RigaOrdineUscita{
				{ProdottoID: 5, Quantita: 10},
			},
		}

		err := svc.OrderService_ConfirmOrdineUscita(context.Background(), o)
		require.NoError(t, err)
		assert.Equal(t, domain.StatoConfermato, o.Stato)
		orderRepo.AssertExpectations(t)
		prodRepo.AssertExpectations(t)
	})

	t.Run("giacenza insufficiente: restituisce errore e non aggiorna nulla", func(t *testing.T) {
		svc, orderRepo, prodRepo := newOrderServiceWithMocks()

		prodRepo.On("ProductRepository_GetGiacenzaById", mock.Anything, 5).Return(&domain.Giacenza{ProdottoID: 5, Quantita: 3}, nil)
		prodRepo.On("ProductRepository_GetByID", mock.Anything, 5).Return(&domain.Prodotto{ID: 5, Nome: "Vite M6"}, nil)

		o := &domain.OrdineUscita{
			Righe: []domain.RigaOrdineUscita{
				{ProdottoID: 5, Quantita: 10},
			},
		}

		err := svc.OrderService_ConfirmOrdineUscita(context.Background(), o)
		require.Error(t, err)
		var giacErr *domain.GiacenzaInsufficienteError
		require.ErrorAs(t, err, &giacErr)
		assert.Equal(t, 5, giacErr.ProdottoID)

		prodRepo.AssertNotCalled(t, "ProductRepository_UpdateGiacenza", mock.Anything, mock.Anything)
		orderRepo.AssertNotCalled(t, "OrderRepository_CreateOrdineUscita", mock.Anything, mock.Anything)
	})

	t.Run("prodotto con giacenza non esistente (0 nel db): errore giacenza insufficiente", func(t *testing.T) {
		svc, orderRepo, prodRepo := newOrderServiceWithMocks()

		prodRepo.On("ProductRepository_GetGiacenzaById", mock.Anything, 5).Return(&domain.Giacenza{ProdottoID: 5, Quantita: 100}, nil)
		prodRepo.On("ProductRepository_UpdateGiacenza", mock.Anything, &domain.Giacenza{ProdottoID: 5, Quantita: -10}).Return(nil)
		prodRepo.On("ProductRepository_GetGiacenzaById", mock.Anything, 6).Return(&domain.Giacenza{ProdottoID: 6, Quantita: 0}, nil)
		prodRepo.On("ProductRepository_GetByID", mock.Anything, 6).Return(&domain.Prodotto{ID: 6, Nome: "Bullone"}, nil)

		o := &domain.OrdineUscita{
			Righe: []domain.RigaOrdineUscita{
				{ProdottoID: 5, Quantita: 10},
				{ProdottoID: 6, Quantita: 5},
			},
		}

		err := svc.OrderService_ConfirmOrdineUscita(context.Background(), o)
		require.Error(t, err)

		orderRepo.AssertNotCalled(t, "OrderRepository_CreateOrdineUscita", mock.Anything, mock.Anything)
	})
}

func TestOrderService_GetOrdiniUscita(t *testing.T) {
	svc, orderRepo, _ := newOrderServiceWithMocks()

	expected := []domain.OrdineUscita{
		{ID: 1, Stato: domain.StatoConfermato},
		{ID: 2, Stato: domain.StatoConfermato},
	}
	orderRepo.On("OrderRepository_GetOrdiniUscita", mock.Anything).Return(expected, nil)

	list, err := svc.OrderService_GetOrdiniUscita(context.Background())
	require.NoError(t, err)
	assert.Equal(t, expected, list)
}

func TestOrderService_GetOrdiniUscitaById(t *testing.T) {
	svc, orderRepo, _ := newOrderServiceWithMocks()

	expected := &domain.OrdineUscita{
		ID:    1,
		Stato: domain.StatoConfermato,
	}

	orderRepo.On(
		"OrderRepository_GetOrdineUscitaByID",
		mock.Anything,
		1,
	).Return(expected, nil)

	o, err := svc.OrderService_GetOrdiniUscitaById(
		context.Background(),
		1,
	)

	require.NoError(t, err)
	assert.Equal(t, expected, o)
}

func TestOrderService_GetOrdiniUscitaByStato(t *testing.T) {
	svc, orderRepo, _ := newOrderServiceWithMocks()

	expected := []domain.OrdineUscita{
		{ID: 1, Stato: domain.StatoConfermato},
		{ID: 2, Stato: domain.StatoConfermato},
	}
	orderRepo.On("OrderRepository_GetOrdiniUscitaByStato", mock.Anything, domain.StatoConfermato).Return(expected, nil)

	list, err := svc.OrderService_GetOrdiniUscitaByStato(context.Background(), domain.StatoConfermato)
	require.NoError(t, err)
	assert.Equal(t, expected, list)
}

func TestOrderService_GetOrdiniUscitaNonEvasi(t *testing.T) {
	svc, orderRepo, _ := newOrderServiceWithMocks()

	expected := []domain.OrdineUscita{
		{ID: 1, Stato: domain.StatoConfermato},
		{ID: 2, Stato: domain.StatoConfermato},
	}
	orderRepo.On("OrderRepository_GetOrdiniUscitaByStato", mock.Anything, domain.StatoConfermato).Return(expected, nil)

	list, err := svc.OrderService_GetOrdiniUscitaNonEvasi(context.Background())
	require.NoError(t, err)
	assert.Equal(t, expected, list)
}

func TestOrderService_UpdateOrdineUscita(t *testing.T) {
	t.Run("ordine senza righe: aggiorna direttamente", func(t *testing.T) {
		svc, orderRepo, _ := newOrderServiceWithMocks()

		o := &domain.OrdineUscita{ID: 1, Stato: domain.StatoConfermato}
		orderRepo.On("OrderRepository_UpdateOrdineUscita", mock.Anything, o).Return(nil)

		err := svc.OrderService_UpdateOrdineUscita(context.Background(), o)
		require.NoError(t, err)
		orderRepo.AssertExpectations(t)
	})

	t.Run("ordine con righe: scarica giacenza e aggiorna", func(t *testing.T) {
		svc, orderRepo, prodRepo := newOrderServiceWithMocks()

		o := &domain.OrdineUscita{
			ID:    1,
			Stato: domain.StatoConfermato,
			Righe: []domain.RigaOrdineUscita{
				{ProdottoID: 5, Quantita: 10},
			},
		}
		prodRepo.On("ProductRepository_GetGiacenzaById", mock.Anything, 5).Return(&domain.Giacenza{ProdottoID: 5, Quantita: 100}, nil)
		prodRepo.On("ProductRepository_UpdateGiacenza", mock.Anything, &domain.Giacenza{ProdottoID: 5, Quantita: -10}).Return(nil)
		orderRepo.On("OrderRepository_UpdateOrdineUscita", mock.Anything, o).Return(nil)

		err := svc.OrderService_UpdateOrdineUscita(context.Background(), o)
		require.NoError(t, err)
		orderRepo.AssertExpectations(t)
		prodRepo.AssertExpectations(t)
	})

	t.Run("giacenza insufficiente: errore", func(t *testing.T) {
		svc, orderRepo, prodRepo := newOrderServiceWithMocks()

		o := &domain.OrdineUscita{
			ID: 1,
			Righe: []domain.RigaOrdineUscita{
				{ProdottoID: 5, Quantita: 100},
			},
		}
		prodRepo.On("ProductRepository_GetGiacenzaById", mock.Anything, 5).Return(&domain.Giacenza{ProdottoID: 5, Quantita: 5}, nil)
		prodRepo.On("ProductRepository_GetByID", mock.Anything, 5).Return(&domain.Prodotto{ID: 5, Nome: "Vite"}, nil)

		err := svc.OrderService_UpdateOrdineUscita(context.Background(), o)
		require.Error(t, err)
		var giacErr *domain.GiacenzaInsufficienteError
		require.ErrorAs(t, err, &giacErr)
		orderRepo.AssertNotCalled(t, "OrderRepository_UpdateOrdineUscita", mock.Anything, mock.Anything)
	})

	t.Run("errore aggiornamento giacenza: propagato", func(t *testing.T) {
		svc, _, prodRepo := newOrderServiceWithMocks()

		o := &domain.OrdineUscita{
			ID: 1,
			Righe: []domain.RigaOrdineUscita{
				{ProdottoID: 5, Quantita: 10},
			},
		}
		prodRepo.On("ProductRepository_GetGiacenzaById", mock.Anything, 5).Return(&domain.Giacenza{ProdottoID: 5, Quantita: 100}, nil)
		prodRepo.On("ProductRepository_UpdateGiacenza", mock.Anything, mock.Anything).Return(errors.New("db error"))

		err := svc.OrderService_UpdateOrdineUscita(context.Background(), o)
		require.Error(t, err)
	})
}

//==========================================================================================

func TestOrderService_GeneraDDT(t *testing.T) {
	t.Run("ordine evaso correttamente: genera DDT", func(t *testing.T) {
		svc, orderRepo, _ := newOrderServiceWithMocks()

		o := &domain.OrdineUscita{
			ID:     20,
			Numero: "OU-2026-0001",
			Stato:  domain.StatoConfermato,
		}

		orderRepo.On(
			"OrderRepository_GetOrdineUscitaByID",
			mock.Anything,
			20,
		).Return(o, nil)

		orderRepo.On(
			"OrderRepository_GetDDTByYear",
			mock.Anything,
			time.Now().Year(),
		).Return([]domain.DDT{}, nil)

		orderRepo.On(
			"OrderRepository_CreateDDT",
			mock.Anything,
			mock.MatchedBy(func(d *domain.DDT) bool {
				return d.OrdineUscitaID == 20 &&
					d.Numero != "" &&
					d.Mittente == "Magazzino A" &&
					d.Destinatario == "Cliente B"
			}),
		).Return(nil)

		orderRepo.On(
			"OrderRepository_UpdateOrdineUscita",
			mock.Anything,
			mock.MatchedBy(func(o *domain.OrdineUscita) bool {
				return o.Stato == domain.StatoEvaso &&
					o.DDTID != nil
			}),
		).Return(nil)

		ddt, err := svc.OrderService_GeneraDDT(
			context.Background(),
			20,
			"Magazzino A",
			"Cliente B",
		)

		require.NoError(t, err)
		require.NotNil(t, ddt)
		assert.Equal(t, 20, ddt.OrdineUscitaID)

		orderRepo.AssertExpectations(t)
	})
}

func TestOrderService_GetDDTByOrdineUscita(t *testing.T) {
	svc, orderRepo, _ := newOrderServiceWithMocks()

	expected := &domain.DDT{ID: 1, OrdineUscitaID: 20, Numero: "DDT-2026-0001"}
	orderRepo.On("OrderRepository_GetDDTByOrdineUscita", mock.Anything, 20).Return(expected, nil)

	ddt, err := svc.OrderService_GetDDTByOrdineUscita(context.Background(), 20)
	require.NoError(t, err)
	assert.Equal(t, expected, ddt)
}

func TestOrderService_GetDDTByYear(t *testing.T) {
	svc, orderRepo, _ := newOrderServiceWithMocks()

	expected := []domain.DDT{
		{ID: 1, OrdineUscitaID: 20, Numero: "DDT-2026-0001"},
		{ID: 2, OrdineUscitaID: 20, Numero: "DDT-2026-0002"},
	}
	orderRepo.On("OrderRepository_GetDDTByYear", mock.Anything, time.Now().Year()).Return(expected, nil)

	list, err := svc.OrderService_GetDDTByYear(context.Background(), time.Now().Year())
	require.NoError(t, err)
	assert.Equal(t, expected, list)
}

//==========================================================================================

func TestOrderService_CreateReso(t *testing.T) {
	t.Run("crea un reso correttamente in stato da_verificare", func(t *testing.T) {
		svc, orderRepo, _ := newOrderServiceWithMocks()

		o := &domain.OrdineUscita{ID: 20, Stato: domain.StatoEvaso, Righe: []domain.RigaOrdineUscita{{ProdottoID: 5, Quantita: 10}}}
		orderRepo.On("OrderRepository_GetOrdineUscitaByID", mock.Anything, 20).Return(o, nil)
		orderRepo.On("OrderRepository_CreateReso", mock.Anything, mock.Anything).Return(nil)

		r := &domain.Reso{OrdineUscitaID: 20, ProdottoID: 5, Quantita: 2}
		err := svc.OrderService_CreateReso(context.Background(), r)
		require.NoError(t, err)
	})

	t.Run("ordine non evaso: errore reso non consentito", func(t *testing.T) {
		svc, orderRepo, _ := newOrderServiceWithMocks()

		o := &domain.OrdineUscita{ID: 20, Stato: domain.StatoConfermato}
		orderRepo.On("OrderRepository_GetOrdineUscitaByID", mock.Anything, 20).Return(o, nil)

		r := &domain.Reso{OrdineUscitaID: 20, ProdottoID: 5, Quantita: 2}
		err := svc.OrderService_CreateReso(context.Background(), r)
		require.Error(t, err)
		orderRepo.AssertNotCalled(t, "OrderRepository_CreateReso", mock.Anything, mock.Anything, mock.Anything)
	})
}

func TestOrderService_VerifyReso(t *testing.T) {
	t.Run("reso da verificare viene verificato e la giacenza aggiornata", func(t *testing.T) {
		svc, orderRepo, prodRepo := newOrderServiceWithMocks()

		orderRepo.On(
			"OrderRepository_GetResoByID",
			mock.Anything,
			3,
		).Return(
			&domain.Reso{
				ID:         3,
				ProdottoID: 5,
				Quantita:   2,
			},
			"da_verificare",
			nil,
		)

		prodRepo.On(
			"ProductRepository_UpdateGiacenza",
			mock.Anything,
			mock.MatchedBy(func(g *domain.Giacenza) bool {
				return g.ProdottoID == 5 && g.Quantita == 2
			}),
		).Return(nil)

		orderRepo.On(
			"OrderRepository_UpdateResoStato",
			mock.Anything,
			3,
			"accettato",
		).Return(nil)

		err := svc.OrderService_VerifyReso(context.Background(), 3)

		require.NoError(t, err)
		orderRepo.AssertExpectations(t)
		prodRepo.AssertExpectations(t)
	})
}
func TestOrderService_ApproveReso(t *testing.T) {
	t.Run("reso verificato viene approvato", func(t *testing.T) {
		svc, orderRepo, _ := newOrderServiceWithMocks()

		orderRepo.On("OrderRepository_GetResoByID", mock.Anything, 3).Return(&domain.Reso{ID: 3}, "accettato", nil)
		orderRepo.On("OrderRepository_UpdateResoStato", mock.Anything, 3, "accettato").Return(nil)

		err := svc.OrderService_ApproveReso(context.Background(), 3)
		require.NoError(t, err)
		orderRepo.AssertExpectations(t)
	})

	t.Run("reso non ancora verificato viene rifiutato", func(t *testing.T) {
		svc, orderRepo, _ := newOrderServiceWithMocks()

		orderRepo.On("OrderRepository_GetResoByID", mock.Anything, 3).Return(&domain.Reso{ID: 3}, "da_verificare", nil)

		err := svc.OrderService_ApproveReso(context.Background(), 3)
		require.Error(t, err)
		orderRepo.AssertNotCalled(t, "OrderRepository_UpdateResoStato", mock.Anything, mock.Anything, mock.Anything)
	})
}

func TestOrderService_RifiutaReso(t *testing.T) {
	t.Run("reso verificato viene rifiutato e la giacenza torna indietro", func(t *testing.T) {
		svc, orderRepo, prodRepo := newOrderServiceWithMocks()

		orderRepo.On("OrderRepository_GetResoByID", mock.Anything, 3).
			Return(&domain.Reso{
				ID:         3,
				ProdottoID: 5,
				Quantita:   2,
			}, "accettato", nil)
		prodRepo.On("ProductRepository_UpdateGiacenza", mock.Anything, mock.MatchedBy(func(g *domain.Giacenza) bool {
			return g.ProdottoID == 5 && g.Quantita == -2
		})).Return(nil)
		orderRepo.On("OrderRepository_UpdateResoStato", mock.Anything, 3, "rifiutato").Return(nil)

		err := svc.OrderService_RifiutaReso(context.Background(), 3)
		require.NoError(t, err)
		orderRepo.AssertExpectations(t)
		prodRepo.AssertExpectations(t)
	})

	t.Run("reso non verificato non può essere rifiutato", func(t *testing.T) {
		svc, orderRepo, prodRepo := newOrderServiceWithMocks()

		orderRepo.On("OrderRepository_GetResoByID", mock.Anything, 3).Return(&domain.Reso{ID: 3}, "da_verificare", nil)

		err := svc.OrderService_RifiutaReso(context.Background(), 3)
		require.Error(t, err)
		prodRepo.AssertNotCalled(t, "ProductRepository_UpdateGiacenza", mock.Anything, mock.Anything)
	})
}

func TestOrderService_GetAllResi(t *testing.T) {
	svc, orderRepo, _ := newOrderServiceWithMocks()

	expectedRepo := []domain.Reso{
		{ID: 1, OrdineUscitaID: 20, ProdottoID: 5, Quantita: 2},
		{ID: 2, OrdineUscitaID: 21, ProdottoID: 6, Quantita: 1},
	}

	expectedDTO := []ResoDTO{
		{ID: 1, OrdineUscitaID: 20, ProdottoID: 5, Quantita: 2},
		{ID: 2, OrdineUscitaID: 21, ProdottoID: 6, Quantita: 1},
	}

	orderRepo.On("OrderRepository_GetAllResi", mock.Anything).
		Return(expectedRepo, nil)

	list, err := svc.OrderService_GetAllResi(context.Background())

	require.NoError(t, err)
	assert.Equal(t, expectedDTO, list)
}

func TestOrderService_DeleteReso(t *testing.T) {
	svc, orderRepo, _ := newOrderServiceWithMocks()

	orderRepo.On("OrderRepository_DeleteReso", mock.Anything, 3).Return(nil)

	err := svc.OrderService_DeleteReso(context.Background(), 3)
	require.NoError(t, err)
	orderRepo.AssertExpectations(t)
}

//==========================================================================================

func TestOrderService_GeneraFattura(t *testing.T) {
	t.Run("ordine evaso: calcola imponibile, IVA al 22% e totale", func(t *testing.T) {
		svc, orderRepo, _ := newOrderServiceWithMocks()

		o := &domain.OrdineUscita{
			ID:    20,
			Stato: domain.StatoEvaso,
			Righe: []domain.RigaOrdineUscita{
				{Quantita: 10, PrezzoUnitario: 5.0}, // 50
				{Quantita: 2, PrezzoUnitario: 25.0}, // 50
			},
		}
		orderRepo.On("OrderRepository_GetOrdineUscitaByID", mock.Anything, 20).Return(o, nil)
		orderRepo.On("OrderRepository_CountFattureByYearAndType", mock.Anything, time.Now().Year(), domain.TipoFatturaVendita).Return(0, nil)
		orderRepo.On("OrderRepository_CreateFattura", mock.Anything, mock.MatchedBy(func(f *domain.Fattura) bool {
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

		orderRepo.On("OrderRepository_GetOrdineUscitaByID", mock.Anything, 20).
			Return(&domain.OrdineUscita{ID: 20, Stato: domain.StatoConfermato}, nil)

		_, err := svc.OrderService_GeneraFattura(context.Background(), 20)
		require.Error(t, err)
		orderRepo.AssertNotCalled(t, "OrderRepository_CreateFattura", mock.Anything, mock.Anything)
	})
}

func TestOrderService_GetAllFatture(t *testing.T) {
	svc, orderRepo, _ := newOrderServiceWithMocks()

	orderRepo.On("OrderRepository_GetAllFatture", mock.Anything).Return([]domain.Fattura{
		{ID: 1, Tipo: domain.TipoFatturaVendita},
		{ID: 2, Tipo: domain.TipoNotaCredito},
	}, nil)

	list, err := svc.OrderService_GetAllFatture(context.Background())
	require.NoError(t, err)
	require.Len(t, list, 2)
}

func TestOrderService_GeneraNotaCreditoById(t *testing.T) {
	svc, orderRepo, _ := newOrderServiceWithMocks()

	o := &domain.OrdineUscita{
		ID:    20,
		Stato: domain.StatoEvaso,
		Righe: []domain.RigaOrdineUscita{{Quantita: 1, PrezzoUnitario: 10.0}},
	}
	orderRepo.On("OrderRepository_GetOrdineUscitaByID", mock.Anything, 20).Return(o, nil)
	orderRepo.On("OrderRepository_CountFattureByYearAndType", mock.Anything, time.Now().Year(), domain.TipoNotaCredito).Return(0, nil)
	orderRepo.On("OrderRepository_CreateFattura", mock.Anything, mock.MatchedBy(func(f *domain.Fattura) bool {
		return f.Tipo == domain.TipoNotaCredito
	})).Return(nil)

	f, err := svc.OrderService_GeneraNotaCreditoById(context.Background(), 20)
	require.NoError(t, err)
	assert.Equal(t, domain.TipoNotaCredito, f.Tipo)
}

func TestOrderService_GetNoteCredito(t *testing.T) {
	svc, orderRepo, _ := newOrderServiceWithMocks()

	orderRepo.On("OrderRepository_GetAllFatture", mock.Anything).Return([]domain.Fattura{
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
