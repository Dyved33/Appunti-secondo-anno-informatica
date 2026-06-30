package order

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/tommasogarofani/ingsw_25-26_group-11/src/backend/domain"
)

// mockOrderServiceForController implementa l'interfaccia orderService
// (vedi 01_PATCH_product_category_order.md).
type mockOrderServiceForController struct {
	mock.Mock
}

// OrderService_ApproveReso implements [orderService].
func (m *mockOrderServiceForController) OrderService_ApproveReso(ctx context.Context, resoID int) error {
	panic("unimplemented")
}

// OrderService_ConfirmOrdineEntrata implements [orderService].
func (m *mockOrderServiceForController) OrderService_ConfirmOrdineEntrata(ctx context.Context, o *domain.OrdineEntrata) error {
	panic("unimplemented")
}

// OrderService_ConfirmOrdineUscita implements [orderService].
func (m *mockOrderServiceForController) OrderService_ConfirmOrdineUscita(ctx context.Context, o *domain.OrdineUscita) error {
	panic("unimplemented")
}

// OrderService_CreateReso implements [orderService].
func (m *mockOrderServiceForController) OrderService_CreateReso(ctx context.Context, r *domain.Reso) error {
	panic("unimplemented")
}

// OrderService_GeneraDDT implements [orderService].
func (m *mockOrderServiceForController) OrderService_GeneraDDT(ctx context.Context, orderID int, mittente string, destinatario string) (*domain.DDT, error) {
	panic("unimplemented")
}

// OrderService_GeneraFattura implements [orderService].
func (m *mockOrderServiceForController) OrderService_GeneraFattura(ctx context.Context, orderID int) (*domain.Fattura, error) {
	panic("unimplemented")
}

// OrderService_GeneraNotaCredito implements [orderService].
func (m *mockOrderServiceForController) OrderService_GeneraNotaCredito(ctx context.Context, orderID int) (*domain.Fattura, error) {
	panic("unimplemented")
}

// OrderService_GetFatture implements [orderService].
func (m *mockOrderServiceForController) OrderService_GetFatture(ctx context.Context) ([]domain.Fattura, error) {
	panic("unimplemented")
}

// OrderService_GetNoteCredito implements [orderService].
func (m *mockOrderServiceForController) OrderService_GetNoteCredito(ctx context.Context) ([]domain.Fattura, error) {
	panic("unimplemented")
}

// OrderService_GetOrdiniUscita implements [orderService].
func (m *mockOrderServiceForController) OrderService_GetOrdiniUscita(ctx context.Context) ([]domain.OrdineUscita, error) {
	panic("unimplemented")
}

// OrderService_GetOrdiniUscitaByStato implements [orderService].
func (m *mockOrderServiceForController) OrderService_GetOrdiniUscitaByStato(ctx context.Context, stato domain.StatoOrdine) ([]domain.OrdineUscita, error) {
	panic("unimplemented")
}

// OrderService_RifiutaReso implements [orderService].
func (m *mockOrderServiceForController) OrderService_RifiutaReso(ctx context.Context, resoID int) error {
	panic("unimplemented")
}

// OrderService_VerifyReso implements [orderService].
func (m *mockOrderServiceForController) OrderService_VerifyReso(ctx context.Context, resoID int) error {
	panic("unimplemented")
}

func (m *mockOrderServiceForController) ProductService_ConfirmOrdineEntrata(ctx context.Context, o *domain.OrdineEntrata) error {
	args := m.Called(ctx, o)
	return args.Error(0)
}
func (m *mockOrderServiceForController) ProductService_ConfirmOrdineUscita(ctx context.Context, o *domain.OrdineUscita) error {
	args := m.Called(ctx, o)
	return args.Error(0)
}
func (m *mockOrderServiceForController) ProductService_GetOrdiniUscita(ctx context.Context) ([]domain.OrdineUscita, error) {
	args := m.Called(ctx)
	list, _ := args.Get(0).([]domain.OrdineUscita)
	return list, args.Error(1)
}
func (m *mockOrderServiceForController) ProductService_GetOrdiniUscitaByStato(ctx context.Context, stato domain.StatoOrdine) ([]domain.OrdineUscita, error) {
	args := m.Called(ctx, stato)
	list, _ := args.Get(0).([]domain.OrdineUscita)
	return list, args.Error(1)
}
func (m *mockOrderServiceForController) ProductService_GeneraDDT(ctx context.Context, orderID int, mittente, destinatario string) (*domain.DDT, error) {
	args := m.Called(ctx, orderID, mittente, destinatario)
	d, _ := args.Get(0).(*domain.DDT)
	return d, args.Error(1)
}
func (m *mockOrderServiceForController) ProductService_CreateReso(ctx context.Context, r *domain.Reso) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}
func (m *mockOrderServiceForController) ProductService_VerifyReso(ctx context.Context, resoID int) error {
	args := m.Called(ctx, resoID)
	return args.Error(0)
}
func (m *mockOrderServiceForController) ProductService_ApproveReso(ctx context.Context, resoID int) error {
	args := m.Called(ctx, resoID)
	return args.Error(0)
}
func (m *mockOrderServiceForController) ProductService_RifiutaReso(ctx context.Context, resoID int) error {
	args := m.Called(ctx, resoID)
	return args.Error(0)
}
func (m *mockOrderServiceForController) ProductService_GeneraFattura(ctx context.Context, orderID int) (*domain.Fattura, error) {
	args := m.Called(ctx, orderID)
	f, _ := args.Get(0).(*domain.Fattura)
	return f, args.Error(1)
}
func (m *mockOrderServiceForController) ProductService_GeneraNotaCredito(ctx context.Context, orderID int) (*domain.Fattura, error) {
	args := m.Called(ctx, orderID)
	f, _ := args.Get(0).(*domain.Fattura)
	return f, args.Error(1)
}
func (m *mockOrderServiceForController) ProductService_GetFatture(ctx context.Context) ([]domain.Fattura, error) {
	args := m.Called(ctx)
	list, _ := args.Get(0).([]domain.Fattura)
	return list, args.Error(1)
}
func (m *mockOrderServiceForController) ProductService_GetNoteCredito(ctx context.Context) ([]domain.Fattura, error) {
	args := m.Called(ctx)
	list, _ := args.Get(0).([]domain.Fattura)
	return list, args.Error(1)
}

func TestOrderController_RegisterEntryOrder(t *testing.T) {
	svc := new(mockOrderServiceForController)
	svc.On("ProductService_ConfirmOrdineEntrata", mock.Anything, mock.MatchedBy(func(o *domain.OrdineEntrata) bool {
		return o.Numero == "OE-001" && o.Fornitore == "ACME Srl" && len(o.Righe) == 1 && o.Righe[0].ProdottoID == 5
	})).Return(nil)

	c := NewOrderController(svc)
	err := c.OrderController_RegisterEntryOrder(context.Background(), OrderInputDTO{
		Numero:    "OE-001",
		Fornitore: "ACME Srl",
		Righe: []OrderLineInputDTO{
			{ProdottoID: 5, Quantita: 100, PrezzoUnitario: 0.5},
		},
	})

	require.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestOrderController_RegisterExitOrder(t *testing.T) {
	t.Run("mappa correttamente l'input nel domain object", func(t *testing.T) {
		svc := new(mockOrderServiceForController)
		svc.On("ProductService_ConfirmOrdineUscita", mock.Anything, mock.MatchedBy(func(o *domain.OrdineUscita) bool {
			return o.Numero == "OU-001" && o.Cliente == "Cliente Spa" && len(o.Righe) == 1
		})).Return(nil)

		c := NewOrderController(svc)
		err := c.OrderController_RegisterExitOrder(context.Background(), OrderInputDTO{
			Numero:  "OU-001",
			Cliente: "Cliente Spa",
			Righe: []OrderLineInputDTO{
				{ProdottoID: 5, Quantita: 10, PrezzoUnitario: 1.0},
			},
		})

		require.NoError(t, err)
		svc.AssertExpectations(t)
	})

	t.Run("propaga l'errore di giacenza insufficiente", func(t *testing.T) {
		svc := new(mockOrderServiceForController)
		svc.On("ProductService_ConfirmOrdineUscita", mock.Anything, mock.Anything).
			Return(domain.NewGiacenzaInsufficienteError(5, "Vite M6", 3, 10))

		c := NewOrderController(svc)
		err := c.OrderController_RegisterExitOrder(context.Background(), OrderInputDTO{})

		require.Error(t, err)
	})
}

func TestOrderController_GetAll(t *testing.T) {
	svc := new(mockOrderServiceForController)
	ddtID := 7
	svc.On("ProductService_GetOrdiniUscita", mock.Anything).Return([]domain.OrdineUscita{
		{ID: 1, Numero: "OU-001", Cliente: "Cliente A", Stato: domain.StatoEvaso, DDTID: &ddtID},
	}, nil)

	c := NewOrderController(svc)
	out, err := c.OrderController_GetAll(context.Background())

	require.NoError(t, err)
	require.Len(t, out, 1)
	assert.Equal(t, "evaso", out[0].Stato)
	require.NotNil(t, out[0].DDTID)
	assert.Equal(t, 7, *out[0].DDTID)
}

func TestOrderController_GetOrdersByStatus(t *testing.T) {
	svc := new(mockOrderServiceForController)
	svc.On("ProductService_GetOrdiniUscitaByStato", mock.Anything, domain.StatoOrdine("confermato")).
		Return([]domain.OrdineUscita{{ID: 1, Numero: "OU-001", Stato: domain.StatoOrdine("confermato")}}, nil)

	c := NewOrderController(svc)
	out, err := c.OrderController_GetOrdersByStatus(context.Background(), "confermato")

	require.NoError(t, err)
	require.Len(t, out, 1)
	assert.Equal(t, "confermato", out[0].Stato)
}

func TestOrderController_IssueDDT(t *testing.T) {
	t.Run("genera il DDT e lo mappa nel DTO", func(t *testing.T) {
		svc := new(mockOrderServiceForController)
		svc.On("ProductService_GeneraDDT", mock.Anything, 20, "Magazzino", "Cliente Spa").
			Return(&domain.DDT{ID: 5, Numero: "DDT-2026-0001", OrdineUscitaID: 20}, nil)

		c := NewOrderController(svc)
		out, err := c.OrderController_IssueDDT(context.Background(), 20, DDTInputDTO{Mittente: "Magazzino", Destinatario: "Cliente Spa"})

		require.NoError(t, err)
		assert.Equal(t, "DDT-2026-0001", out.Numero)
	})

	t.Run("propaga l'errore se il DDT è già stato generato", func(t *testing.T) {
		svc := new(mockOrderServiceForController)
		svc.On("ProductService_GeneraDDT", mock.Anything, 20, "Magazzino", "Cliente Spa").
			Return(nil, domain.ErrDDTGiaGenerato)

		c := NewOrderController(svc)
		_, err := c.OrderController_IssueDDT(context.Background(), 20, DDTInputDTO{Mittente: "Magazzino", Destinatario: "Cliente Spa"})

		require.ErrorIs(t, err, domain.ErrDDTGiaGenerato)
	})
}

func TestOrderController_CreateReturn(t *testing.T) {
	svc := new(mockOrderServiceForController)
	svc.On("ProductService_CreateReso", mock.Anything, mock.MatchedBy(func(r *domain.Reso) bool {
		return r.OrdineUscitaID == 20 && r.ProdottoID == 5 && r.Quantita == 2 && r.Motivo == "Difettoso"
	})).Return(nil)

	c := NewOrderController(svc)
	err := c.OrderController_CreateReturn(context.Background(), ResoInputDTO{
		OrdineUscitaID: 20,
		ProdottoID:     5,
		Quantita:       2,
		Motivo:         "Difettoso",
		CreatedByID:    1,
	})

	require.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestOrderController_VerifyReturnPhysical(t *testing.T) {
	svc := new(mockOrderServiceForController)
	svc.On("OrderService_VerifyReso", mock.Anything, 3).Return(nil)

	c := NewOrderController(svc)
	err := c.OrderController_VerifyReturnPhysical(context.Background(), 3)

	require.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestOrderController_ApproveReturn(t *testing.T) {
	svc := new(mockOrderServiceForController)
	svc.On("OrderService_ApproveReso", mock.Anything, 3).Return(nil)

	c := NewOrderController(svc)
	err := c.OrderController_ApproveReturn(context.Background(), 3)

	require.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestOrderController_RejectReturn(t *testing.T) {
	svc := new(mockOrderServiceForController)
	svc.On("OrderService_RifiutaReso", mock.Anything, 3).Return(nil)

	c := NewOrderController(svc)
	err := c.OrderController_RejectReturn(context.Background(), 3)

	require.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestOrderController_IssueInvoice(t *testing.T) {
	t.Run("emette la fattura e la mappa nel DTO", func(t *testing.T) {
		svc := new(mockOrderServiceForController)
		svc.On("OrderService_GeneraFattura", mock.Anything, 20).
			Return(&domain.Fattura{ID: 1, Numero: "FT-2026-0001", Tipo: domain.TipoFatturaVendita, Totale: 122.0}, nil)

		c := NewOrderController(svc)
		out, err := c.OrderController_IssueInvoice(context.Background(), 20)

		require.NoError(t, err)
		assert.Equal(t, "FT-2026-0001", out.Numero)
		assert.Equal(t, 122.0, out.Totale)
	})

	t.Run("propaga l'errore se l'ordine non è evaso", func(t *testing.T) {
		svc := new(mockOrderServiceForController)
		svc.On("OrderService_GeneraFattura", mock.Anything, 20).
			Return(nil, errors.New("impossibile emettere fattura per un ordine non ancora evaso"))

		c := NewOrderController(svc)
		_, err := c.OrderController_IssueInvoice(context.Background(), 20)

		require.Error(t, err)
	})
}

func TestOrderController_IssueCreditNote(t *testing.T) {
	svc := new(mockOrderServiceForController)
	svc.On("OrderService_GeneraNotaCredito", mock.Anything, 20).
		Return(&domain.Fattura{ID: 2, Numero: "NC-2026-0001", Tipo: domain.TipoNotaCredito}, nil)

	c := NewOrderController(svc)
	out, err := c.OrderController_IssueCreditNote(context.Background(), 20)

	require.NoError(t, err)
	assert.Equal(t, "NC-2026-0001", out.Numero)
}

func TestOrderController_GetInvoices(t *testing.T) {
	svc := new(mockOrderServiceForController)
	svc.On("OrderService_GetFatture", mock.Anything).Return([]domain.Fattura{
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

// il test di fatturaToDTO è superfluo perché la funzione è privata e viene già testata indirettamente nei test di IssueInvoice e GetInvoices.
