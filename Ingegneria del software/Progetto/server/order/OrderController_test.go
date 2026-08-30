package order

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/tommasogarofani/ingsw_25-26_group-11/server/domain"
)

// mockOrderServiceForController implementa l'interfaccia orderService.
type mockOrderServiceForController struct {
	mock.Mock
}

// OrderService_GeneraNotaCreditoById implements [orderService].
func (m *mockOrderServiceForController) OrderService_GeneraNotaCreditoById(ctx context.Context, orderID int) (*domain.Fattura, error) {
	args := m.Called(ctx, orderID)
	f, _ := args.Get(0).(*domain.Fattura)
	return f, args.Error(1)
}

// OrderService_GetAllFatture implements [orderService].
func (m *mockOrderServiceForController) OrderService_GetAllFatture(ctx context.Context) ([]domain.Fattura, error) {
	args := m.Called(ctx)
	list, _ := args.Get(0).([]domain.Fattura)
	return list, args.Error(1)
}

// OrderService_GetAllResi implements [orderService].
func (m *mockOrderServiceForController) OrderService_GetAllResi(ctx context.Context) ([]ResoDTO, error) {
	args := m.Called(ctx)
	list, _ := args.Get(0).([]ResoDTO)
	return list, args.Error(1)
}

// OrderService_GetDDTByOrdineUscita implements [orderService].
func (m *mockOrderServiceForController) OrderService_GetDDTByOrdineUscita(ctx context.Context, orderID int) (*domain.DDT, error) {
	args := m.Called(ctx, orderID)
	d, _ := args.Get(0).(*domain.DDT)
	return d, args.Error(1)
}

// OrderService_GetDDTByYear implements [orderService].
func (m *mockOrderServiceForController) OrderService_GetDDTByYear(ctx context.Context, year int) ([]domain.DDT, error) {
	args := m.Called(ctx, year)
	list, _ := args.Get(0).([]domain.DDT)
	return list, args.Error(1)
}

// OrderService_GetOrdiniEntrata implements [orderService].
func (m *mockOrderServiceForController) OrderService_GetOrdiniEntrata(ctx context.Context) ([]domain.OrdineEntrata, error) {
	args := m.Called(ctx)
	list, _ := args.Get(0).([]domain.OrdineEntrata)
	return list, args.Error(1)
}

// OrderService_GetOrdiniEntrataById implements [orderService].
func (m *mockOrderServiceForController) OrderService_GetOrdiniEntrataById(ctx context.Context, id int) (*domain.OrdineEntrata, error) {
	args := m.Called(ctx, id)
	o, _ := args.Get(0).(*domain.OrdineEntrata)
	return o, args.Error(1)
}

func (m *mockOrderServiceForController) OrderService_DeleteOrdineEntrata(ctx context.Context, o *domain.OrdineEntrata) error {
	args := m.Called(ctx, o)
	return args.Error(0)
}

// OrderService_GetOrdiniUscitaById implements [orderService].
func (m *mockOrderServiceForController) OrderService_GetOrdiniUscitaById(ctx context.Context, id int) (*domain.OrdineUscita, error) {
	args := m.Called(ctx, id)
	o, _ := args.Get(0).(*domain.OrdineUscita)
	return o, args.Error(1)
}

// OrderService_GetOrdiniUscitaNonEvasi implements [orderService].
func (m *mockOrderServiceForController) OrderService_GetOrdiniUscitaNonEvasi(ctx context.Context) ([]domain.OrdineUscita, error) {
	args := m.Called(ctx)
	list, _ := args.Get(0).([]domain.OrdineUscita)
	return list, args.Error(1)
}

func (m *mockOrderServiceForController) OrderService_ConfirmOrdineEntrata(ctx context.Context, o *domain.OrdineEntrata) error {
	args := m.Called(ctx, o)
	return args.Error(0)
}

func (m *mockOrderServiceForController) OrderService_ConfirmOrdineUscita(ctx context.Context, o *domain.OrdineUscita) error {
	args := m.Called(ctx, o)
	return args.Error(0)
}

func (m *mockOrderServiceForController) OrderService_GetOrdiniUscita(ctx context.Context) ([]domain.OrdineUscita, error) {
	args := m.Called(ctx)
	list, _ := args.Get(0).([]domain.OrdineUscita)
	return list, args.Error(1)
}

func (m *mockOrderServiceForController) OrderService_GetOrdiniUscitaByStato(ctx context.Context, stato domain.StatoOrdineUscita) ([]domain.OrdineUscita, error) {
	args := m.Called(ctx, stato)
	list, _ := args.Get(0).([]domain.OrdineUscita)
	return list, args.Error(1)
}

func (m *mockOrderServiceForController) OrderService_GetDettaglioOrdineUscita(ctx context.Context, id int) (*domain.OrdineUscita, error) {
	args := m.Called(ctx, id)
	o, _ := args.Get(0).(*domain.OrdineUscita)
	return o, args.Error(1)
}

func (m *mockOrderServiceForController) OrderService_UpdateOrdineUscita(ctx context.Context, o *domain.OrdineUscita) error {
	args := m.Called(ctx, o)
	return args.Error(0)
}

func (m *mockOrderServiceForController) OrderService_DeleteOrdineUscita(ctx context.Context, o *domain.OrdineUscita) error {
	args := m.Called(ctx, o)
	return args.Error(0)
}

func (m *mockOrderServiceForController) OrderService_GeneraDDT(ctx context.Context, orderID int, mittente, destinatario string) (*domain.DDT, error) {
	args := m.Called(ctx, orderID, mittente, destinatario)
	d, _ := args.Get(0).(*domain.DDT)
	return d, args.Error(1)
}

func (m *mockOrderServiceForController) OrderService_CreateReso(ctx context.Context, r *domain.Reso) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}

func (m *mockOrderServiceForController) OrderService_VerifyReso(ctx context.Context, resoID int) error {
	args := m.Called(ctx, resoID)
	return args.Error(0)
}

func (m *mockOrderServiceForController) OrderService_ApproveReso(ctx context.Context, resoID int) error {
	args := m.Called(ctx, resoID)
	return args.Error(0)
}

func (m *mockOrderServiceForController) OrderService_RifiutaReso(ctx context.Context, resoID int) error {
	args := m.Called(ctx, resoID)
	return args.Error(0)
}

func (m *mockOrderServiceForController) OrderService_DeleteReso(ctx context.Context, resoID int) error {
	args := m.Called(ctx, resoID)
	return args.Error(0)
}

func (m *mockOrderServiceForController) OrderService_GeneraFattura(ctx context.Context, orderID int) (*domain.Fattura, error) {
	args := m.Called(ctx, orderID)
	f, _ := args.Get(0).(*domain.Fattura)
	return f, args.Error(1)
}

func (m *mockOrderServiceForController) OrderService_GeneraNotaCredito(ctx context.Context, orderID int) (*domain.Fattura, error) {
	args := m.Called(ctx, orderID)
	f, _ := args.Get(0).(*domain.Fattura)
	return f, args.Error(1)
}

func (m *mockOrderServiceForController) OrderService_GetFatture(ctx context.Context) ([]domain.Fattura, error) {
	args := m.Called(ctx)
	list, _ := args.Get(0).([]domain.Fattura)
	return list, args.Error(1)
}

func (m *mockOrderServiceForController) OrderService_DeleteFattura(ctx context.Context, fatturaID int) error {
	args := m.Called(ctx, fatturaID)
	return args.Error(0)
}

func (m *mockOrderServiceForController) OrderService_GetNoteCredito(ctx context.Context) ([]domain.Fattura, error) {
	args := m.Called(ctx)
	list, _ := args.Get(0).([]domain.Fattura)
	return list, args.Error(1)
}

func TestOrderController_RegisterEntryOrder(t *testing.T) {
	svc := new(mockOrderServiceForController)
	svc.On("OrderService_ConfirmOrdineEntrata", mock.Anything, mock.MatchedBy(func(o *domain.OrdineEntrata) bool {
		return o.Numero == "OE-001" && o.Fornitore == "ACME Srl" && len(o.Righe) == 1 && o.Righe[0].ProdottoID == 5
	})).Return(nil)

	c := NewOrderController(svc)
	err := c.OrderController_RegisterEntryOrder(context.Background(), OrdineEntrataDTO{
		Numero:    "OE-001",
		Fornitore: "ACME Srl",
		Righe: []RigaOrdineEntrataDTO{
			{ProdottoID: 5, Quantita: 100, PrezzoUnitario: 0.5},
		},
	})

	require.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestOrderController_RegisterExitOrder(t *testing.T) {
	t.Run("mappa correttamente l'input nel domain object", func(t *testing.T) {
		svc := new(mockOrderServiceForController)
		svc.On("OrderService_ConfirmOrdineUscita", mock.Anything, mock.MatchedBy(func(o *domain.OrdineUscita) bool {
			return o.Numero == "OU-001" && o.Cliente == "Cliente Spa" && len(o.Righe) == 1
		})).Return(nil)

		c := NewOrderController(svc)
		err := c.OrderController_RegisterExitOrder(context.Background(), OrdineUscitaDTO{
			Numero:  "OU-001",
			Cliente: "Cliente Spa",
			Righe: []RigaOrdineUscitaDTO{
				{ProdottoID: 5, Quantita: 10, PrezzoUnitario: 1.0},
			},
		})

		require.NoError(t, err)
		svc.AssertExpectations(t)
	})

	t.Run("propaga l'errore di giacenza insufficiente", func(t *testing.T) {
		svc := new(mockOrderServiceForController)
		svc.On("OrderService_ConfirmOrdineUscita", mock.Anything, mock.Anything).
			Return(domain.NewGiacenzaInsufficienteError(5, "Vite M6", 3, 10))

		c := NewOrderController(svc)
		err := c.OrderController_RegisterExitOrder(context.Background(), OrdineUscitaDTO{})

		require.Error(t, err)
	})
}

func TestOrderController_GetOrdiniEntrata(t *testing.T) {
	svc := new(mockOrderServiceForController)
	svc.On("OrderService_GetOrdiniEntrata", mock.Anything).Return([]domain.OrdineEntrata{
		{ID: 1, Numero: "OE-001", Fornitore: "ACME Srl"},
	}, nil)
	svc.On("OrderService_GetOrdiniEntrataById", mock.Anything, 1).Return(&domain.OrdineEntrata{
		ID:     1,
		Numero: "OE-001",
		Righe:  []domain.RigaOrdineEntrata{{ID: 10, ProdottoID: 5, Quantita: 100, PrezzoUnitario: 0.5}},
	}, nil)

	c := NewOrderController(svc)
	out, err := c.OrderController_GetOrdiniEntrata(context.Background())

	require.NoError(t, err)
	require.Len(t, out, 1)
	require.Len(t, out[0].Righe, 1)
	assert.Equal(t, 5, out[0].Righe[0].ProdottoID)
	assert.Equal(t, 100, out[0].Righe[0].Quantita)
	assert.Equal(t, 0.5, out[0].Righe[0].PrezzoUnitario)
}

func TestOrderController_GetOrdiniEntrataById(t *testing.T) {
	svc := new(mockOrderServiceForController)
	svc.On("OrderService_GetOrdiniEntrataById", mock.Anything, 1).Return(&domain.OrdineEntrata{ID: 1, Numero: "OE-001", Fornitore: "ACME Srl"}, nil)

	c := NewOrderController(svc)
	out, err := c.OrderController_GetOrdiniEntrataById(context.Background(), 1)

	require.NoError(t, err)
	require.NotNil(t, out)
}

func TestOrderController_DeleteOrdineEntrata(t *testing.T) {
	svc := new(mockOrderServiceForController)

	order := &domain.OrdineEntrata{
		ID:     1,
		Numero: "OE-001",
	}

	svc.On(
		"OrderService_GetOrdiniEntrataById",
		mock.Anything,
		1,
	).Return(order, nil)

	svc.On(
		"OrderService_DeleteOrdineEntrata",
		mock.Anything,
		order,
	).Return(nil)

	c := NewOrderController(svc)

	err := c.OrderController_DeleteOrdineEntrata(
		context.Background(),
		1,
	)

	require.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestOrderController_GetOrdiniUscita(t *testing.T) {
	svc := new(mockOrderServiceForController)
	ddtID := 7
	svc.On("OrderService_GetOrdiniUscita", mock.Anything).Return([]domain.OrdineUscita{
		{ID: 1, Numero: "OU-001", Cliente: "Cliente A", Stato: domain.StatoEvaso, DDTID: &ddtID},
	}, nil)
	svc.On("OrderService_GetOrdiniUscitaById", mock.Anything, 1).Return(&domain.OrdineUscita{
		ID:    1,
		Righe: []domain.RigaOrdineUscita{{ID: 10, ProdottoID: 5, Quantita: 2, PrezzoUnitario: 1.5}},
	}, nil)

	c := NewOrderController(svc)
	out, err := c.OrderController_GetOrdiniUscita(context.Background())

	require.NoError(t, err)
	require.Len(t, out, 1)
	assert.Equal(t, "evaso", out[0].Stato)
	require.NotNil(t, out[0].DDTID)
	assert.Equal(t, 7, *out[0].DDTID)
}

func TestOrderController_GetOrdersByStatus(t *testing.T) {
	svc := new(mockOrderServiceForController)
	svc.On("OrderService_GetOrdiniUscitaByStato", mock.Anything, domain.StatoOrdineUscita("confermato")).
		Return([]domain.OrdineUscita{{ID: 1, Numero: "OU-001", Stato: domain.StatoOrdineUscita("confermato")}}, nil)

	c := NewOrderController(svc)
	out, err := c.OrderController_GetOrdersByStatus(context.Background(), "confermato")

	require.NoError(t, err)
	require.Len(t, out, 1)
	assert.Equal(t, "confermato", out[0].Stato)
}

func TestOrderController_UpdateOrdineUscita(t *testing.T) {
	svc := new(mockOrderServiceForController)
	svc.On("OrderService_GetOrdiniUscitaById", mock.Anything, 1).Return(&domain.OrdineUscita{ID: 1, Numero: "OU-001", Cliente: "Cliente A"}, nil)
	svc.On("OrderService_UpdateOrdineUscita", mock.Anything, mock.Anything).Return(nil)

	c := NewOrderController(svc)
	err := c.OrderController_UpdateOrdineUscita(context.Background(), 1, OrdineUscitaDTO{
		Numero:  "OU-002",
		Cliente: "Cliente B",
	})

	require.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestOrderController_DeleteOrdineUscita(t *testing.T) {
	svc := new(mockOrderServiceForController)

	order := &domain.OrdineUscita{
		ID:     1,
		Numero: "OU-001",
	}

	svc.On(
		"OrderService_GetOrdiniUscitaById",
		mock.Anything,
		1,
	).Return(order, nil)

	svc.On(
		"OrderService_DeleteOrdineUscita",
		mock.Anything,
		order,
	).Return(nil)

	c := NewOrderController(svc)

	err := c.OrderController_DeleteOrdineUscita(
		context.Background(),
		1,
	)

	require.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestOrderController_GeneraDDT(t *testing.T) {
	t.Run("genera il DDT e lo mappa nel DTO", func(t *testing.T) {
		svc := new(mockOrderServiceForController)
		svc.On("OrderService_GeneraDDT", mock.Anything, 20, "Magazzino", "Cliente Spa").
			Return(&domain.DDT{ID: 5, Numero: "DDT-2026-0001", OrdineUscitaID: 20}, nil)

		c := NewOrderController(svc)
		out, err := c.OrderController_GeneraDDT(context.Background(), 20, DDTEntrataDTO{Mittente: "Magazzino", Destinatario: "Cliente Spa"})

		require.NoError(t, err)
		assert.Equal(t, "DDT-2026-0001", out.Numero)
	})

	t.Run("propaga l'errore se il DDT è già stato generato", func(t *testing.T) {
		svc := new(mockOrderServiceForController)
		svc.On("OrderService_GeneraDDT", mock.Anything, 20, "Magazzino", "Cliente Spa").
			Return(nil, domain.ErrDDTGiaGenerato)

		c := NewOrderController(svc)
		_, err := c.OrderController_GeneraDDT(context.Background(), 20, DDTEntrataDTO{Mittente: "Magazzino", Destinatario: "Cliente Spa"})

		require.ErrorIs(t, err, domain.ErrDDTGiaGenerato)
	})
}

func TestOrderController_CreateReturn(t *testing.T) {
	svc := new(mockOrderServiceForController)
	svc.On("OrderService_CreateReso", mock.Anything, mock.MatchedBy(func(r *domain.Reso) bool {
		return r.OrdineUscitaID == 20 && r.ProdottoID == 5 && r.Quantita == 2 && r.Motivo == "Difettoso"
	})).Return(nil)

	c := NewOrderController(svc)
	err := c.OrderController_CreateReturn(context.Background(), ResoDTO{
		OrdineUscitaID: 20,
		ProdottoID:     5,
		Quantita:       2,
		Motivo:         "Difettoso",
		CreatedByID:    1,
	})

	require.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestOrderController_UpdateStatoReso(t *testing.T) {

	t.Run("stato da_verificare chiama VerifyReso", func(t *testing.T) {
		service := new(mockOrderServiceForController)

		controller := &OrderController{
			Service: service,
		}

		service.On(
			"OrderService_VerifyReso",
			mock.Anything,
			3,
		).Return(nil)

		err := controller.OrderController_UpdateStatoReso(
			context.Background(),
			3,
			"da_verificare",
		)

		require.NoError(t, err)
		service.AssertExpectations(t)
	})

	t.Run("stato approvato chiama ApproveReso", func(t *testing.T) {
		service := new(mockOrderServiceForController)

		controller := &OrderController{
			Service: service,
		}

		service.On(
			"OrderService_ApproveReso",
			mock.Anything,
			3,
		).Return(nil)

		err := controller.OrderController_UpdateStatoReso(
			context.Background(),
			3,
			"approvato",
		)

		require.NoError(t, err)
		service.AssertExpectations(t)
	})

	t.Run("stato rifiutato chiama RifiutaReso", func(t *testing.T) {
		service := new(mockOrderServiceForController)

		controller := &OrderController{
			Service: service,
		}

		service.On(
			"OrderService_RifiutaReso",
			mock.Anything,
			3,
		).Return(nil)

		err := controller.OrderController_UpdateStatoReso(
			context.Background(),
			3,
			"rifiutato",
		)

		require.NoError(t, err)
		service.AssertExpectations(t)
	})

	t.Run("stato non valido non chiama nessun servizio", func(t *testing.T) {
		service := new(mockOrderServiceForController)

		controller := &OrderController{
			Service: service,
		}

		err := controller.OrderController_UpdateStatoReso(
			context.Background(),
			3,
			"stato_non_valido",
		)

		require.NoError(t, err)

		service.AssertNotCalled(
			t,
			"OrderService_VerifyReso",
			mock.Anything,
			mock.Anything,
		)
		service.AssertNotCalled(
			t,
			"OrderService_ApproveReso",
			mock.Anything,
			mock.Anything,
		)
		service.AssertNotCalled(
			t,
			"OrderService_RifiutaReso",
			mock.Anything,
			mock.Anything,
		)
	})
}

func TestOrderController_GetReturns(t *testing.T) {
	svc := new(mockOrderServiceForController)
	svc.On("OrderService_GetAllResi", mock.Anything).Return([]ResoDTO{
		{ID: 1, OrdineUscitaID: 20, ProdottoID: 5, Quantita: 2, Motivo: "Difettoso"},
	}, nil)

	c := NewOrderController(svc)
	out, err := c.OrderController_GetReturns(context.Background())

	require.NoError(t, err)
	require.Len(t, out, 1)
	assert.Equal(t, 20, out[0].OrdineUscitaID)
	assert.Equal(t, "Difettoso", out[0].Motivo)
}

func TestOrderController_DeleteReso(t *testing.T) {
	svc := new(mockOrderServiceForController)
	svc.On("OrderService_DeleteReso", mock.Anything, 3).Return(nil)

	c := NewOrderController(svc)
	err := c.OrderController_DeleteReso(context.Background(), 3)

	require.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestOrderController_GeneraFattura(t *testing.T) {
	t.Run("emette la fattura e la mappa nel DTO", func(t *testing.T) {
		svc := new(mockOrderServiceForController)
		svc.On("OrderService_GeneraFattura", mock.Anything, 20).
			Return(&domain.Fattura{ID: 1, Numero: "FT-2026-0001", Tipo: domain.TipoFatturaVendita, Totale: 122.0}, nil)

		c := NewOrderController(svc)
		out, err := c.OrderController_GeneraFattura(context.Background(), 20)

		require.NoError(t, err)
		assert.Equal(t, "FT-2026-0001", out.Numero)
		assert.Equal(t, 122.0, out.Totale)
	})

	t.Run("propaga l'errore se l'ordine non è evaso", func(t *testing.T) {
		svc := new(mockOrderServiceForController)
		svc.On("OrderService_GeneraFattura", mock.Anything, 20).
			Return(nil, errors.New("impossibile emettere fattura per un ordine non ancora evaso"))

		c := NewOrderController(svc)
		_, err := c.OrderController_GeneraFattura(context.Background(), 20)

		require.Error(t, err)
	})
}

func TestOrderController_DeleteFattura(t *testing.T) {
	svc := new(mockOrderServiceForController)
	svc.On("OrderService_DeleteFattura", mock.Anything, 1).Return(nil)

	c := NewOrderController(svc)
	err := c.OrderController_DeleteFattura(context.Background(), 1)

	require.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestOrderController_GeneraNotaCredito(t *testing.T) {
	svc := new(mockOrderServiceForController)
	svc.On("OrderService_GeneraNotaCreditoById", mock.Anything, 20).
		Return(&domain.Fattura{ID: 2, Numero: "NC-2026-0001", Tipo: domain.TipoNotaCredito}, nil)

	c := NewOrderController(svc)
	out, err := c.OrderController_GeneraNotaCredito(context.Background(), 20)

	require.NoError(t, err)
	assert.Equal(t, "NC-2026-0001", out.Numero)
}

func TestOrderController_GetInvoices(t *testing.T) {
	svc := new(mockOrderServiceForController)
	svc.On("OrderService_GetAllFatture", mock.Anything).Return([]domain.Fattura{
		{ID: 1, Numero: "FT-2026-0001"},
		{ID: 2, Numero: "FT-2026-0002"},
	}, nil)

	c := NewOrderController(svc)
	out, err := c.OrderController_GetInvoices(context.Background())

	require.NoError(t, err)
	require.Len(t, out, 2)
}

func TestOrderController_GetCreditNotes(t *testing.T) {
	svc := new(mockOrderServiceForController)
	svc.On("OrderService_GetNoteCredito", mock.Anything).Return([]domain.Fattura{
		{ID: 3, Numero: "NC-2026-0001"},
	}, nil)

	c := NewOrderController(svc)
	out, err := c.OrderController_GetCreditNotes(context.Background())

	require.NoError(t, err)
	require.Len(t, out, 1)
	assert.Equal(t, "NC-2026-0001", out[0].Numero)
}

func TestOrderController_FatturaToDTO(t *testing.T) {
	f := domain.Fattura{
		ID:             1,
		Numero:         "FT-2026-0001",
		Tipo:           domain.TipoFatturaVendita,
		OrdineUscitaID: 20,
		Imponibile:     100.0,
		IVA:            22.0,
		Totale:         122.0,
	}

	c := NewOrderController(nil)
	dto := c.OrderController_FatturaToDTO(f)

	assert.Equal(t, f.ID, dto.ID)
	assert.Equal(t, f.Numero, dto.Numero)
	assert.Equal(t, string(f.Tipo), dto.Tipo)
	assert.Equal(t, f.OrdineUscitaID, dto.OrdineUscitaID)
	assert.Equal(t, f.Imponibile, dto.Imponibile)
	assert.Equal(t, f.IVA, dto.IVA)
	assert.Equal(t, f.Totale, dto.Totale)
}
