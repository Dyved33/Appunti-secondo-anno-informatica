package Handler

import (
	"context"
	"strings"
	"gialosai/backend/DTO"
	"gialosai/backend/domain"
	"gialosai/backend/service"
)

type ProductHandler struct {
	Service *Service.ProductService
}

func NewProductHandler(svc *Service.ProductService) *ProductHandler {
	return &ProductHandler{Service: svc}
}

func (h *ProductHandler) AddProduct(ctx context.Context, input DTO.ProductInputDTO) (DTO.ProductOutputDTO, error) {
	p := &domain.Prodotto{
		Codice:      input.Codice,
		Barcode:     input.Barcode,
		Nome:        input.Nome,
		Descrizione: input.Descrizione,
		CategoriaID: input.CategoriaID,
		PrezzoAcq:   input.PrezzoAcq,
		PrezzoVend:  input.PrezzoVend,
		SogliaMin:   input.SogliaMin,
	}

	err := h.Service.CreateProduct(ctx, p)
	if err != nil {
		return DTO.ProductOutputDTO{}, err
	}

	return DTO.ProductOutputDTO{
		ID:          p.ID,
		Codice:      p.Codice,
		Barcode:     p.Barcode,
		Nome:        p.Nome,
		Descrizione: p.Descrizione,
		CategoriaID: p.CategoriaID,
		PrezzoAcq:   p.PrezzoAcq,
		PrezzoVend:  p.PrezzoVend,
		SogliaMin:   p.SogliaMin,
	}, nil
}

func (h *ProductHandler) GetProducts(ctx context.Context) ([]DTO.ProductOutputDTO, error) {
	list, err := h.Service.GetProducts(ctx)
	if err != nil {
		return nil, err
	}

	var dtoList []DTO.ProductOutputDTO
	for _, p := range list {
		dtoList = append(dtoList, DTO.ProductOutputDTO{
			ID:          p.ID,
			Codice:      p.Codice,
			Barcode:     p.Barcode,
			Nome:        p.Nome,
			Descrizione: p.Descrizione,
			CategoriaID: p.CategoriaID,
			PrezzoAcq:   p.PrezzoAcq,
			PrezzoVend:  p.PrezzoVend,
			SogliaMin:   p.SogliaMin,
		})
	}
	return dtoList, nil
}

func (h *ProductHandler) GetProductByBarcode(ctx context.Context, barcode string) (DTO.ProductOutputDTO, error) {
	barcode = strings.TrimSpace(barcode)
	p, err := h.Service.GetProductByBarcode(ctx, barcode)
	if err != nil {
		return DTO.ProductOutputDTO{}, err
	}

	return DTO.ProductOutputDTO{
		ID:          p.ID,
		Codice:      p.Codice,
		Barcode:     p.Barcode,
		Nome:        p.Nome,
		Descrizione: p.Descrizione,
		CategoriaID: p.CategoriaID,
		PrezzoAcq:   p.PrezzoAcq,
		PrezzoVend:  p.PrezzoVend,
		SogliaMin:   p.SogliaMin,
	}, nil
}

func (h *ProductHandler) EditProduct(ctx context.Context, id int, input DTO.ProductInputDTO) error {
	p := &domain.Prodotto{
		ID:          id,
		Codice:      input.Codice,
		Barcode:     input.Barcode,
		Nome:        input.Nome,
		Descrizione: input.Descrizione,
		CategoriaID: input.CategoriaID,
		PrezzoAcq:   input.PrezzoAcq,
		PrezzoVend:  input.PrezzoVend,
		SogliaMin:   input.SogliaMin,
	}
	return h.Service.UpdateProduct(ctx, p)
}

func (h *ProductHandler) RemoveProduct(ctx context.Context, id int) error {
	return h.Service.DeleteProduct(ctx, id)
}

func (h *ProductHandler) UpdatePrices(ctx context.Context, id int, input DTO.PriceUpdateInputDTO) error {
	return h.Service.UpdateListinoPrices(ctx, id, input.PrezzoAcq, input.PrezzoVend)
}

func (h *ProductHandler) GetInventoryValue(ctx context.Context) (int, float64, error) {
	return h.Service.GetValoreMagazzino(ctx)
}

func (h *ProductHandler) GetInventoryReport(ctx context.Context) (*domain.InventoryReport, error) {
	return h.Service.GetInventarioReport(ctx)
}

func (h *ProductHandler) GetLowStockProducts(ctx context.Context) ([]domain.InventoryItem, error) {
	return h.Service.GetProdottiSottoSoglia(ctx)
}
