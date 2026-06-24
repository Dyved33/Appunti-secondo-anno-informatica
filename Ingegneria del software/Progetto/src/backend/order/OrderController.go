package order

import (
	"context"
	"gialosai/backend/domain"
)

// Controller rappresenta il layer di controllo/handler per le operazioni su ordini, DDT, fatture e resi.
// Trasla le richieste in ingresso (DTO) alle chiamate al Service e formatta le risposte.
type Controller struct {
	Service *OrderService
}

// NewController crea una nuova istanza del controller per gli ordini.
func NewController(svc *OrderService) *Controller {
	return &Controller{Service: svc}
}

// RegisterEntryOrder gestisce la registrazione di un ordine di entrata (carico merce).
func (c *Controller) RegisterEntryOrder(ctx context.Context, input OrderInputDTO) error {
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

	return c.Service.ConfirmOrdineEntrata(ctx, o)
}

// RegisterExitOrder gestisce la registrazione di un ordine di uscita (scarico merce).
func (c *Controller) RegisterExitOrder(ctx context.Context, input OrderInputDTO) error {
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

	return c.Service.ConfirmOrdineUscita(ctx, o)
}

// GetOrders recupera tutti gli ordini di uscita.
func (c *Controller) GetOrders(ctx context.Context) ([]OrderOutputDTO, error) {
	orders, err := c.Service.GetOrdiniUscita(ctx)
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
func (c *Controller) GetOrdersByStatus(ctx context.Context, status string) ([]OrderOutputDTO, error) {
	orders, err := c.Service.GetOrdiniUscitaByStato(ctx, domain.StatoOrdine(status))
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
func (c *Controller) IssueDDT(ctx context.Context, orderID int, req DDTInputDTO) (DDTOutputDTO, error) {
	ddt, err := c.Service.GeneraDDT(ctx, orderID, req.Mittente, req.Destinatario)
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
func (c *Controller) CreateReturn(ctx context.Context, req ResoInputDTO) error {
	r := &domain.Reso{
		OrdineUscitaID: req.OrdineUscitaID,
		ProdottoID:     req.ProdottoID,
		Quantita:       req.Quantita,
		Motivo:         req.Motivo,
		CreatedByID:    req.CreatedByID,
	}
	return c.Service.CreateReso(ctx, r)
}

// VerifyReturnPhysical imposta lo stato del reso a in verifica fisica.
func (c *Controller) VerifyReturnPhysical(ctx context.Context, resoID int) error {
	return c.Service.VerifyReso(ctx, resoID)
}

// ApproveReturn approva il reso con ripristino delle giacenze.
func (c *Controller) ApproveReturn(ctx context.Context, resoID int) error {
	return c.Service.ApproveReso(ctx, resoID)
}

// RejectReturn rifiuta il reso.
func (c *Controller) RejectReturn(ctx context.Context, resoID int) error {
	return c.Service.RifiutaReso(ctx, resoID)
}

// IssueInvoice genera una fattura per un ordine evaso.
func (c *Controller) IssueInvoice(ctx context.Context, orderID int) (FatturaOutputDTO, error) {
	f, err := c.Service.GeneraFattura(ctx, orderID)
	if err != nil {
		return FatturaOutputDTO{}, err
	}

	return fatturaToDTO(*f), nil
}

// IssueCreditNote genera una nota di credito per un ordine evaso o stornato.
func (c *Controller) IssueCreditNote(ctx context.Context, orderID int) (FatturaOutputDTO, error) {
	f, err := c.Service.GeneraNotaCredito(ctx, orderID)
	if err != nil {
		return FatturaOutputDTO{}, err
	}

	return fatturaToDTO(*f), nil
}

// GetInvoices recupera tutte le fatture.
func (c *Controller) GetInvoices(ctx context.Context) ([]FatturaOutputDTO, error) {
	fatture, err := c.Service.GetFatture(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]FatturaOutputDTO, 0, len(fatture))
	for _, f := range fatture {
		result = append(result, fatturaToDTO(f))
	}
	return result, nil
}

// GetCreditNotes recupera tutte le note di credito.
func (c *Controller) GetCreditNotes(ctx context.Context) ([]FatturaOutputDTO, error) {
	fatture, err := c.Service.GetNoteCredito(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]FatturaOutputDTO, 0, len(fatture))
	for _, f := range fatture {
		result = append(result, fatturaToDTO(f))
	}
	return result, nil
}

func fatturaToDTO(f domain.Fattura) FatturaOutputDTO {
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
