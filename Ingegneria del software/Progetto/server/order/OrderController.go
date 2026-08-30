package order

import (
	"context"

	"github.com/tommasogarofani/ingsw_25-26_group-11/server/domain"
)

type orderService interface {
	OrderService_ConfirmOrdineEntrata(ctx context.Context, o *domain.OrdineEntrata) error
	OrderService_GetOrdiniEntrata(ctx context.Context) ([]domain.OrdineEntrata, error)
	OrderService_GetOrdiniEntrataById(ctx context.Context, id int) (*domain.OrdineEntrata, error)
	OrderService_DeleteOrdineEntrata(ctx context.Context, o *domain.OrdineEntrata) error

	OrderService_ConfirmOrdineUscita(ctx context.Context, o *domain.OrdineUscita) error
	OrderService_GetOrdiniUscitaNonEvasi(ctx context.Context) ([]domain.OrdineUscita, error)
	OrderService_GetOrdiniUscita(ctx context.Context) ([]domain.OrdineUscita, error)
	OrderService_GetOrdiniUscitaByStato(ctx context.Context, stato domain.StatoOrdineUscita) ([]domain.OrdineUscita, error)
	OrderService_GetOrdiniUscitaById(ctx context.Context, id int) (*domain.OrdineUscita, error)
	OrderService_UpdateOrdineUscita(ctx context.Context, o *domain.OrdineUscita) error
	OrderService_DeleteOrdineUscita(ctx context.Context, o *domain.OrdineUscita) error

	OrderService_GetDDTByOrdineUscita(ctx context.Context, orderID int) (*domain.DDT, error)
	OrderService_GeneraDDT(ctx context.Context, orderID int, mittente, destinatario string) (*domain.DDT, error)
	OrderService_GetDDTByYear(ctx context.Context, year int) ([]domain.DDT, error)

	OrderService_CreateReso(ctx context.Context, r *domain.Reso) error
	OrderService_VerifyReso(ctx context.Context, resoID int) error
	OrderService_ApproveReso(ctx context.Context, resoID int) error
	OrderService_RifiutaReso(ctx context.Context, resoID int) error
	OrderService_GetAllResi(ctx context.Context) ([]ResoDTO, error)
	OrderService_DeleteReso(ctx context.Context, resoID int) error

	OrderService_GeneraFattura(ctx context.Context, orderID int) (*domain.Fattura, error)
	OrderService_GeneraNotaCreditoById(ctx context.Context, orderID int) (*domain.Fattura, error)
	OrderService_GetAllFatture(ctx context.Context) ([]domain.Fattura, error)
	OrderService_GetNoteCredito(ctx context.Context) ([]domain.Fattura, error)
	OrderService_DeleteFattura(ctx context.Context, fatturaID int) error
}

// OrderController espone i metodi di gestione di ordini (entrata e uscita),
// DDT, resi e fatture, delegando la logica di business al service iniettato.
type OrderController struct {
	Service orderService
}

// NewOrderController crea un nuovo OrderController con il service specificato.
func NewOrderController(svc orderService) *OrderController {
	return &OrderController{Service: svc}
}

// RegisterEntryOrder gestisce la registrazione di un ordine di entrata (carico merce).
func (c *OrderController) OrderController_RegisterEntryOrder(ctx context.Context, input OrdineEntrataDTO) error {
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
func (c *OrderController) OrderController_RegisterExitOrder(ctx context.Context, input OrdineUscitaDTO) error {
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

// recupera tutti gli ordini di entrata.
func (c *OrderController) OrderController_GetOrdiniEntrata(ctx context.Context) ([]OrdineUscitaDTO, error) {
	orders, err := c.Service.OrderService_GetOrdiniEntrata(ctx)
	if err != nil {
		return nil, err
	}

	var dtoList []OrdineUscitaDTO
	for _, o := range orders {
		detail, err := c.Service.OrderService_GetOrdiniEntrataById(ctx, o.ID)
		if err != nil {
			return nil, err
		}

		dto := OrdineUscitaDTO{
			ID:          o.ID,
			Numero:      o.Numero,
			Fornitore:   o.Fornitore,
			Data:        o.Data,
			Note:        o.Note,
			CreatedByID: o.CreatedByID,
		}
		if detail != nil {
			dto.Righe = make([]RigaOrdineUscitaDTO, 0, len(detail.Righe))
			for _, r := range detail.Righe {
				dto.Righe = append(dto.Righe, RigaOrdineUscitaDTO{
					ID:             r.ID,
					ProdottoID:     r.ProdottoID,
					Quantita:       r.Quantita,
					PrezzoUnitario: r.PrezzoUnitario,
				})
			}
		}
		dtoList = append(dtoList, dto)
	}
	return dtoList, nil
}

// OrderController_GetOrdiniEntrataById recupera un ordine di entrata in base al suo ID.
func (c *OrderController) OrderController_GetOrdiniEntrataById(ctx context.Context, id int) (*OrdineUscitaDTO, error) {
	order, err := c.Service.OrderService_GetOrdiniEntrataById(ctx, id)
	if err != nil {
		return nil, err
	}

	dto := OrdineUscitaDTO{
		ID:          order.ID,
		Numero:      order.Numero,
		Fornitore:   order.Fornitore,
		Data:        order.Data,
		Note:        order.Note,
		CreatedByID: order.CreatedByID,
	}
	return &dto, nil
}

// OrderController_DeleteOrdineEntrata elimina un ordine di entrata in base al suo ID.
func (c *OrderController) OrderController_DeleteOrdineEntrata(ctx context.Context, id int) error {
	order, err := c.Service.OrderService_GetOrdiniEntrataById(ctx, id)
	if err != nil {
		return err
	}

	return c.Service.OrderService_DeleteOrdineEntrata(ctx, order)
}

//============================================================================================

// GetOrders recupera tutti gli ordini di uscita.
func (c *OrderController) OrderController_GetOrdiniUscita(ctx context.Context) ([]OrdineUscitaDTO, error) {
	orders, err := c.Service.OrderService_GetOrdiniUscita(ctx)
	if err != nil {
		return nil, err
	}

	var dtoList []OrdineUscitaDTO
	for _, o := range orders {
		detail, err := c.Service.OrderService_GetOrdiniUscitaById(ctx, o.ID)
		if err != nil {
			return nil, err
		}

		dto := OrdineUscitaDTO{
			ID:          o.ID,
			Numero:      o.Numero,
			Cliente:     o.Cliente,
			Data:        o.Data,
			Stato:       string(o.Stato),
			Note:        o.Note,
			CreatedByID: o.CreatedByID,
			DDTID:       o.DDTID,
		}
		if detail != nil {
			dto.Righe = make([]RigaOrdineUscitaDTO, 0, len(detail.Righe))
			for _, r := range detail.Righe {
				dto.Righe = append(dto.Righe, RigaOrdineUscitaDTO{
					ID:             r.ID,
					ProdottoID:     r.ProdottoID,
					Quantita:       r.Quantita,
					PrezzoUnitario: r.PrezzoUnitario,
				})
			}
		}
		dtoList = append(dtoList, dto)
	}
	return dtoList, nil
}

// GetOrdersByStatus recupera gli ordini di uscita filtrati per stato.
func (c *OrderController) OrderController_GetOrdersByStatus(ctx context.Context, status string) ([]OrdineUscitaDTO, error) {
	orders, err := c.Service.OrderService_GetOrdiniUscitaByStato(ctx, domain.StatoOrdineUscita(status))
	if err != nil {
		return nil, err
	}

	var dtoList []OrdineUscitaDTO
	for _, o := range orders {
		dtoList = append(dtoList, OrdineUscitaDTO{
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

// OrderController_UpdateOrdineUscita aggiorna i dati e lo stato di un ordine di uscita,
// inclusa la sostituzione delle relative righe.
func (c *OrderController) OrderController_UpdateOrdineUscita(ctx context.Context, id int, input OrdineUscitaDTO) error {
	order, err := c.Service.OrderService_GetOrdiniUscitaById(ctx, id)
	if err != nil {
		return err
	}

	order.Numero = input.Numero
	order.Cliente = input.Cliente
	order.Note = input.Note
	order.Stato = domain.StatoOrdineUscita(input.Stato)

	// Aggiorna le righe dell'ordine
	var righe []domain.RigaOrdineUscita
	for _, r := range input.Righe {
		righe = append(righe, domain.RigaOrdineUscita{
			ID:             r.ID,
			ProdottoID:     r.ProdottoID,
			Quantita:       r.Quantita,
			PrezzoUnitario: r.PrezzoUnitario,
		})
	}
	order.Righe = righe

	return c.Service.OrderService_UpdateOrdineUscita(ctx, order)
}

// OrderController_DeleteOrdineUscita elimina un ordine di uscita in base al suo ID.
func (c *OrderController) OrderController_DeleteOrdineUscita(ctx context.Context, id int) error {
	order, err := c.Service.OrderService_GetOrdiniUscitaById(ctx, id)
	if err != nil {
		return err
	}

	return c.Service.OrderService_DeleteOrdineUscita(ctx, order)
}

//============================================================================================

// GeneraDDT genera un Documento di Trasporto per un ordine di uscita.
func (c *OrderController) OrderController_GeneraDDT(ctx context.Context, orderID int, req DDTEntrataDTO) (DDTUscitaDTO, error) {
	ddt, err := c.Service.OrderService_GeneraDDT(ctx, orderID, req.Mittente, req.Destinatario)
	if err != nil {
		return DDTUscitaDTO{}, err
	}

	return DDTUscitaDTO{
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
func (c *OrderController) OrderController_CreateReturn(ctx context.Context, req ResoDTO) error {
	r := &domain.Reso{
		OrdineUscitaID: req.OrdineUscitaID,
		ProdottoID:     req.ProdottoID,
		Quantita:       req.Quantita,
		Stato:          domain.StatoReso(req.Stato),
		Motivo:         req.Motivo,
		CreatedByID:    req.CreatedByID,
	}
	return c.Service.OrderService_CreateReso(ctx, r)
}

// OrderController_UpdateStatoReso aggiorna lo stato di un reso in base al nuovo stato
// accettato ("da_verificare", "approvato", "rifiutato"). Se lo stato non è tra questi,
// non esegue alcuna operazione.
func (c *OrderController) OrderController_UpdateStatoReso(ctx context.Context, resoID int, newStato string) error {
	switch newStato {
	case "da_verificare":
		return c.Service.OrderService_VerifyReso(ctx, resoID)
	case "approvato":
		return c.Service.OrderService_ApproveReso(ctx, resoID)
	case "rifiutato":
		return c.Service.OrderService_RifiutaReso(ctx, resoID)
	default:
		return nil // Stato non valido, non fare nulla
	}
}

// GetReturns recupera tutti i resi.
func (c *OrderController) OrderController_GetReturns(ctx context.Context) ([]ResoDTO, error) {
	return c.Service.OrderService_GetAllResi(ctx)
}

// OrderController_DeleteReso elimina un reso in base al suo ID.
func (c *OrderController) OrderController_DeleteReso(ctx context.Context, resoID int) error {
	return c.Service.OrderService_DeleteReso(ctx, resoID)
}

// GeneraFattura genera una fattura per un ordine evaso.
func (c *OrderController) OrderController_GeneraFattura(ctx context.Context, orderID int) (FatturaDTO, error) {
	f, err := c.Service.OrderService_GeneraFattura(ctx, orderID)
	if err != nil {
		return FatturaDTO{}, err
	}

	return c.OrderController_FatturaToDTO(*f), nil
}

// OrderController_DeleteFattura elimina una fattura in base al suo ID.
func (c *OrderController) OrderController_DeleteFattura(ctx context.Context, fatturaID int) error {
	return c.Service.OrderService_DeleteFattura(ctx, fatturaID)
}

// GeneraNotaCredito genera una nota di credito per un ordine evaso o stornato.
func (c *OrderController) OrderController_GeneraNotaCredito(ctx context.Context, orderID int) (FatturaDTO, error) {
	f, err := c.Service.OrderService_GeneraNotaCreditoById(ctx, orderID)
	if err != nil {
		return FatturaDTO{}, err
	}

	return c.OrderController_FatturaToDTO(*f), nil
}

// GetInvoices recupera tutte le fatture.
func (c *OrderController) OrderController_GetInvoices(ctx context.Context) ([]FatturaDTO, error) {
	fatture, err := c.Service.OrderService_GetAllFatture(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]FatturaDTO, 0, len(fatture))
	for _, f := range fatture {
		result = append(result, c.OrderController_FatturaToDTO(f))
	}
	return result, nil
}

// GetCreditNotes recupera tutte le note di credito.
func (c *OrderController) OrderController_GetCreditNotes(ctx context.Context) ([]FatturaDTO, error) {
	fatture, err := c.Service.OrderService_GetNoteCredito(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]FatturaDTO, 0, len(fatture))
	for _, f := range fatture {
		result = append(result, c.OrderController_FatturaToDTO(f))
	}
	return result, nil
}

// OrderController_FatturaToDTO converte una Fattura del dominio nel corrispondente DTO.
func (c *OrderController) OrderController_FatturaToDTO(f domain.Fattura) FatturaDTO {
	return FatturaDTO{
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
