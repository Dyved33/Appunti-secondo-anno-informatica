package Handler

import (
	"context"
	"gialosai/backend/DTO"
	"gialosai/backend/domain"
	"gialosai/backend/service"
)

type OrderHandler struct {
	Service *Service.OrderService
}

func NewOrderHandler(svc *Service.OrderService) *OrderHandler {
	return &OrderHandler{Service: svc}
}

func (h *OrderHandler) RegisterEntryOrder(ctx context.Context, input DTO.OrderInputDTO) error {
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

	return h.Service.ConfirmOrdineEntrata(ctx, o)
}

func (h *OrderHandler) RegisterExitOrder(ctx context.Context, input DTO.OrderInputDTO) error {
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

	return h.Service.ConfirmOrdineUscita(ctx, o)
}

func (h *OrderHandler) GetOrders(ctx context.Context) ([]DTO.OrderOutputDTO, error) {
	orders, err := h.Service.GetOrdiniUscita(ctx)
	if err != nil {
		return nil, err
	}

	var dtoList []DTO.OrderOutputDTO
	for _, o := range orders {
		dtoList = append(dtoList, DTO.OrderOutputDTO{
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

func (h *OrderHandler) GetOrdersByStatus(ctx context.Context, status string) ([]DTO.OrderOutputDTO, error) {
	orders, err := h.Service.GetOrdiniUscitaByStato(ctx, domain.StatoOrdine(status))
	if err != nil {
		return nil, err
	}

	var dtoList []DTO.OrderOutputDTO
	for _, o := range orders {
		dtoList = append(dtoList, DTO.OrderOutputDTO{
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

func (h *OrderHandler) IssueDDT(ctx context.Context, orderID int, req DTO.DDTInputDTO) (DTO.DDTOutputDTO, error) {
	ddt, err := h.Service.GeneraDDT(ctx, orderID, req.Mittente, req.Destinatario)
	if err != nil {
		return DTO.DDTOutputDTO{}, err
	}

	return DTO.DDTOutputDTO{
		ID:             ddt.ID,
		Numero:         ddt.Numero,
		OrdineUscitaID: ddt.OrdineUscitaID,
		Data:           ddt.Data,
		Mittente:       ddt.Mittente,
		Destinatario:   ddt.Destinatario,
		PathPDF:        ddt.PathPDF,
	}, nil
}

func (h *OrderHandler) CreateReturn(ctx context.Context, req DTO.ResoInputDTO) error {
	r := &domain.Reso{
		OrdineUscitaID: req.OrdineUscitaID,
		ProdottoID:     req.ProdottoID,
		Quantita:       req.Quantita,
		Motivo:         req.Motivo,
		CreatedByID:    req.CreatedByID,
	}
	return h.Service.CreateReso(ctx, r)
}

func (h *OrderHandler) VerifyReturnPhysical(ctx context.Context, resoID int) error {
	return h.Service.VerifyReso(ctx, resoID)
}

func (h *OrderHandler) ApproveReturn(ctx context.Context, resoID int) error {
	return h.Service.ApproveReso(ctx, resoID)
}

func (h *OrderHandler) RejectReturn(ctx context.Context, resoID int) error {
	return h.Service.RifiutaReso(ctx, resoID)
}

func (h *OrderHandler) IssueInvoice(ctx context.Context, orderID int) (DTO.OrderOutputDTO, error) {
	f, err := h.Service.GeneraFattura(ctx, orderID)
	if err != nil {
		return DTO.OrderOutputDTO{}, err
	}

	// Fetch order details
	o, err := h.Service.GetDettaglioOrdineUscita(ctx, f.OrdineUscitaID)
	if err != nil {
		return DTO.OrderOutputDTO{}, err
	}

	return DTO.OrderOutputDTO{
		ID:          o.ID,
		Numero:      o.Numero,
		Cliente:     o.Cliente,
		Data:        o.Data,
		Stato:       string(o.Stato),
		Note:        o.Note,
		CreatedByID: o.CreatedByID,
		DDTID:       o.DDTID,
	}, nil
}
