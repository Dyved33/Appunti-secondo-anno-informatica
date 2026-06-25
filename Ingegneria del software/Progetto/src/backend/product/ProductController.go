package product

import (
	"context"
	"strings"
	"gialosai/backend/domain"
)

type Controller struct {
	Service *Service
}


func NewController(svc *Service) *Controller {
	return &Controller{Service: svc}
}

// AddProduct gestisce l'aggiunta di un nuovo prodotto al sistema.
func (c *Controller) AddProduct(ctx context.Context, input ProductInputDTO) (ProductOutputDTO, error) {
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

	err := c.Service.CreateProduct(ctx, p)
	if err != nil {
		return ProductOutputDTO{}, err
	}

	return ProductOutputDTO{
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

// GetProducts recupera la lista di tutti i prodotti disponibili nel sistema.
func (c *Controller) GetProducts(ctx context.Context) ([]ProductOutputDTO, error) {
	list, err := c.Service.GetProducts(ctx)
	if err != nil {
		return nil, err
	}

	var dtoList []ProductOutputDTO
	for _, p := range list {
		dtoList = append(dtoList, ProductOutputDTO{
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

// GetProductByBarcode recupera un prodotto dal sistema usando il suo barcode.
func (c *Controller) GetProductByBarcode(ctx context.Context, barcode string) (ProductOutputDTO, error) {
	barcode = strings.TrimSpace(barcode)
	p, err := c.Service.GetProductByBarcode(ctx, barcode)
	if err != nil {
		return ProductOutputDTO{}, err
	}

	return ProductOutputDTO{
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

// EditProduct modifica un prodotto esistente nel sistema.
func (c *Controller) EditProduct(ctx context.Context, id int, input ProductInputDTO) error {
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
	return c.Service.UpdateProduct(ctx, p)
}

// RemoveProduct rimuove un prodotto dal sistema.
func (c *Controller) RemoveProduct(ctx context.Context, id int) error {
	return c.Service.DeleteProduct(ctx, id)
}

// UpdatePrices modifica i prezzi di acquisto e vendita di un prodotto.
func (c *Controller) UpdatePrices(ctx context.Context, id int, input PriceUpdateInputDTO) error {
	return c.Service.UpdateListinoPrices(ctx, id, input.PrezzoAcq, input.PrezzoVend)
}

// GetInventoryValue restituisce un riepilogo sintetico del magazzino (quantità totale e valore totale).
func (c *Controller) GetInventoryValue(ctx context.Context) (int, float64, error) {
	return c.Service.GetValoreMagazzino(ctx)
}

// GetInventoryReport restituisce il report completo dello stato del magazzino.
func (c *Controller) GetInventoryReport(ctx context.Context) (*domain.InventoryReport, error) {
	return c.Service.GetInventarioReport(ctx)
}

// GetLowStockProducts restituisce l'elenco dei prodotti che sono sotto la soglia minima di giacenza.
func (c *Controller) GetLowStockProducts(ctx context.Context) ([]domain.InventoryItem, error) {
	return c.Service.GetProdottiSottoSoglia(ctx)
}
