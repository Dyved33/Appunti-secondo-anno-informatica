package order

import (
	"context"
	"github.com/tommasogarofani/ingsw_25-26_group-11/server/domain"
)

type orderService interface {
	OrderService_ConfirmOrdineEntrata(ctx context.Context, o *domain.OrdineEntrata) error
	OrderService_ConfirmOrdineUscita(ctx context.Context, o *domain.OrdineUscita) error
	OrderService_GetOrdiniUscita(ctx context.Context) ([]domain.OrdineUscita, error)
	OrderService_GetOrdiniUscitaByStato(ctx context.Context, stato domain.StatoOrdine) ([]domain.OrdineUscita, error)
	OrderService_GeneraDDT(ctx context.Context, orderID int, mittente, destinatario string) (*domain.DDT, error)
	OrderService_CreateReso(ctx context.Context, r *domain.Reso) error
	OrderService_VerifyReso(ctx context.Context, resoID int) error
	OrderService_ApproveReso(ctx context.Context, resoID int) error
	OrderService_RifiutaReso(ctx context.Context, resoID int) error
	OrderService_GeneraFattura(ctx context.Context, orderID int) (*domain.Fattura, error)
	OrderService_GeneraNotaCredito(ctx context.Context, orderID int) (*domain.Fattura, error)
	OrderService_GetFatture(ctx context.Context) ([]domain.Fattura, error)
	OrderService_GetNoteCredito(ctx context.Context) ([]domain.Fattura, error)
}

type OrderController struct {
	Service orderService
}

func NewOrderController(svc orderService) *OrderController {
	return &OrderController{Service: svc}
}

// RegisterEntryOrder gestisce la registrazione di un ordine di entrata (carico merce).
func (c *OrderController) OrderController_RegisterEntryOrder(ctx context.Context, input OrderInputDTO) error {
	var righe []domain.RigaOrdineEntrata
	for _, r := range input.Righe {
		righe = append(righe, domain.RigaOrdineEntrata{
			ProdottoID:     r.ProdottoID,
			Quantita:       r.Quantita,
			PrezzoUnitario: r.PrezzoUnitario,
		})
	}

	o := &domain.OrdineEntrata{
		Numero:      input.Numero,
		Fornitore:   input.Fornitore,
		Note:        input.Note,
		CreatedByID: input.CreatedByID,
		Righe:       righe,
	}

	return c.Service.OrderService_ConfirmOrdineEntrata(ctx, o)
}

// RegisterExitOrder gestisce la registrazione di un ordine di uscita (scarico merce).
func (c *OrderController) OrderController_RegisterExitOrder(ctx context.Context, input OrderInputDTO) error {
	var righe []domain.RigaOrdineUscita
	for _, r := range input.Righe {
		righe = append(righe, domain.RigaOrdineUscita{
			ProdottoID:     r.ProdottoID,
			Quantita:       r.Quantita,
			PrezzoUnitario: r.PrezzoUnitario,
		})
	}

	o := &domain.OrdineUscita{
		Numero:      input.Numero,
		Cliente:     input.Cliente,
		Note:        input.Note,
		CreatedByID: input.CreatedByID,
		Righe:       righe,
	}

	return c.Service.OrderService_ConfirmOrdineUscita(ctx, o)
}

// GetOrders recupera tutti gli ordini di uscita.
func (c *OrderController) OrderController_GetAll(ctx context.Context) ([]OrderOutputDTO, error) {
	orders, err := c.Service.OrderService_GetOrdiniUscita(ctx)
	if err != nil {
		return nil, err
	}

	var dtoList []OrderOutputDTO
	for _, o := range orders {
		dtoList = append(dtoList, OrderOutputDTO{
			ID:          o.ID,
			Numero:      o.Numero,
			Cliente:     o.Cliente,
			Data:        o.Data,
			Stato:       string(o.Stato),
			Note:        o.Note,
			CreatedByID: o.CreatedByID,
			DDTID:       o.DDTID,
		})
	}
	return dtoList, nil
}

// GetOrdersByStatus recupera gli ordini di uscita filtrati per stato.
func (c *OrderController) OrderController_GetOrdersByStatus(ctx context.Context, status string) ([]OrderOutputDTO, error) {
	orders, err := c.Service.OrderService_GetOrdiniUscitaByStato(ctx, domain.StatoOrdine(status))
	if err != nil {
		return nil, err
	}

	var dtoList []OrderOutputDTO
	for _, o := range orders {
		dtoList = append(dtoList, OrderOutputDTO{
			ID:          o.ID,
			Numero:      o.Numero,
			Cliente:     o.Cliente,
			Data:        o.Data,
			Stato:       string(o.Stato),
			Note:        o.Note,
			CreatedByID: o.CreatedByID,
			DDTID:       o.DDTID,
		})
	}
	return dtoList, nil
}

// IssueDDT genera un Documento di Trasporto per un ordine di uscita.
func (c *OrderController) OrderController_IssueDDT(ctx context.Context, orderID int, req DDTInputDTO) (DDTOutputDTO, error) {
	ddt, err := c.Service.OrderService_GeneraDDT(ctx, orderID, req.Mittente, req.Destinatario)
	if err != nil {
		return DDTOutputDTO{}, err
	}

	return DDTOutputDTO{
		ID:             ddt.ID,
		Numero:         ddt.Numero,
		OrdineUscitaID: ddt.OrdineUscitaID,
		Data:           ddt.Data,
		Mittente:       ddt.Mittente,
		Destinatario:   ddt.Destinatario,
		PathPDF:        ddt.PathPDF,
	}, nil
}

// CreateReturn registra la richiesta di un reso di merce.
func (c *OrderController) OrderController_CreateReturn(ctx context.Context, req ResoInputDTO) error {
	r := &domain.Reso{
		OrdineUscitaID: req.OrdineUscitaID,
		ProdottoID:     req.ProdottoID,
		Quantita:       req.Quantita,
		Motivo:         req.Motivo,
		CreatedByID:    req.CreatedByID,
	}
	return c.Service.OrderService_CreateReso(ctx, r)
}

// VerifyReturnPhysical imposta lo stato del reso a in verifica fisica.
func (c *OrderController) OrderController_VerifyReturnPhysical(ctx context.Context, resoID int) error {
	return c.Service.OrderService_VerifyReso(ctx, resoID)
}

// ApproveReturn approva il reso con ripristino delle giacenze.
func (c *OrderController) OrderController_ApproveReturn(ctx context.Context, resoID int) error {
	return c.Service.OrderService_ApproveReso(ctx, resoID)
}

// RejectReturn rifiuta il reso.
func (c *OrderController) OrderController_RejectReturn(ctx context.Context, resoID int) error {
	return c.Service.OrderService_RifiutaReso(ctx, resoID)
}

// IssueInvoice genera una fattura per un ordine evaso.
func (c *OrderController) OrderController_IssueInvoice(ctx context.Context, orderID int) (FatturaOutputDTO, error) {
	f, err := c.Service.OrderService_GeneraFattura(ctx, orderID)
	if err != nil {
		return FatturaOutputDTO{}, err
	}

	return c.OrderController_FatturaToDTO(*f), nil
}

// IssueCreditNote genera una nota di credito per un ordine evaso o stornato.
func (c *OrderController) OrderController_IssueCreditNote(ctx context.Context, orderID int) (FatturaOutputDTO, error) {
	f, err := c.Service.OrderService_GeneraNotaCredito(ctx, orderID)
	if err != nil {
		return FatturaOutputDTO{}, err
	}

	return c.OrderController_FatturaToDTO(*f), nil
}

// GetInvoices recupera tutte le fatture.
func (c *OrderController) OrderController_GetInvoices(ctx context.Context) ([]FatturaOutputDTO, error) {
	fatture, err := c.Service.OrderService_GetFatture(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]FatturaOutputDTO, 0, len(fatture))
	for _, f := range fatture {
		result = append(result, c.OrderController_FatturaToDTO(f))
	}
	return result, nil
}

// GetCreditNotes recupera tutte le note di credito.
func (c *OrderController) OrderController_GetCreditNotes(ctx context.Context) ([]FatturaOutputDTO, error) {
	fatture, err := c.Service.OrderService_GetNoteCredito(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]FatturaOutputDTO, 0, len(fatture))
	for _, f := range fatture {
		result = append(result, c.OrderController_FatturaToDTO(f))
	}
	return result, nil
}

func (c *OrderController) OrderController_FatturaToDTO(f domain.Fattura) FatturaOutputDTO {
	return FatturaOutputDTO{
		ID:             f.ID,
		Numero:         f.Numero,
		Tipo:           string(f.Tipo),
		OrdineUscitaID: f.OrdineUscitaID,
		Imponibile:     f.Imponibile,
		IVA:            f.IVA,
		Totale:         f.Totale,
		Data:           f.Data,
		PathPDF:        f.PathPDF,
	}
}