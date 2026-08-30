package product

import (
	"context"
	"strings"

	"github.com/tommasogarofani/ingsw_25-26_group-11/server/domain"
)

type productService interface {
	ProductService_Create(ctx context.Context, p *domain.Prodotto) error
	ProductService_Update(ctx context.Context, p *domain.Prodotto) error
	ProductService_Delete(ctx context.Context, id int) error
	ProductService_GetAllProducts(ctx context.Context) ([]domain.Prodotto, error)
	ProductService_GetByID(ctx context.Context, id int) (*domain.Prodotto, error)
	ProductService_GetByBarcode(ctx context.Context, barcode string) (*domain.Prodotto, error)

	ProductService_UpdateListinoPrices(ctx context.Context, productID int, acqPrice, vendPrice float64) error
	ProductService_GetValoreMagazzino(ctx context.Context) (int, float64, error)
	ProductService_GetAllGiacenze(ctx context.Context) (*domain.InventoryReport, error)
	ProductService_GetProdottiSottoSoglia(ctx context.Context) ([]domain.InventoryItem, error)
}

// ProductController espone i metodi di gestione dei prodotti e delle giacenze,
// delegando la logica di business al service iniettato.
type ProductController struct {
	Service productService
}

// NewProductController crea un nuovo ProductController con il service specificato.
func NewProductController(svc productService) *ProductController {
	return &ProductController{Service: svc}
}

// AddProduct gestisce l'aggiunta di un nuovo prodotto al sistema.
func (c *ProductController) ProductController_AddProduct(ctx context.Context, input ProductInputDTO) (ProductOutputDTO, error) {
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

	err := c.Service.ProductService_Create(ctx, p)
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
func (c *ProductController) ProductController_GetProducts(ctx context.Context) ([]ProductOutputDTO, error) {
	list, err := c.Service.ProductService_GetAllProducts(ctx)
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

// GetProductsById recupera un prodotto specifico dal sistema usando il suo ID.
func (c *ProductController) ProductController_GetProductsById(ctx context.Context, id int) ([]ProductOutputDTO, error) {
	p, err := c.Service.ProductService_GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	var dtoList []ProductOutputDTO
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
	return dtoList, nil
}

// GetProductByBarcode recupera un prodotto dal sistema usando il suo barcode.
func (c *ProductController) ProductController_GetProductByBarcode(ctx context.Context, barcode string) (ProductOutputDTO, error) {
	barcode = strings.TrimSpace(barcode)
	p, err := c.Service.ProductService_GetByBarcode(ctx, barcode)
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
func (c *ProductController) ProductController_EditProduct(ctx context.Context, id int, input ProductInputDTO) error {
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
	return c.Service.ProductService_Update(ctx, p)
}

// RemoveProduct rimuove un prodotto dal sistema.
func (c *ProductController) ProductController_RemoveProduct(ctx context.Context, id int) error {
	return c.Service.ProductService_Delete(ctx, id)
}

// UpdatePrices modifica i prezzi di acquisto e vendita di un prodotto.
func (c *ProductController) ProductController_UpdatePrices(ctx context.Context, id int, input PriceUpdateInputDTO) error {
	return c.Service.ProductService_UpdateListinoPrices(ctx, id, input.PrezzoAcq, input.PrezzoVend)
}

// UpdateSellPrice aggiorna solo il prezzo di vendita preservando quello di acquisto.
func (c *ProductController) ProductController_UpdateSellPrice(ctx context.Context, id int, vendPrice float64) error {
	p, err := c.Service.ProductService_GetByID(ctx, id)
	if err != nil {
		return err
	}
	return c.Service.ProductService_UpdateListinoPrices(ctx, id, p.PrezzoAcq, vendPrice)
}

// GetInventoryValue restituisce un riepilogo sintetico del magazzino (quantità totale e valore totale).
func (c *ProductController) ProductController_GetValoreMagazzino(ctx context.Context) (int, float64, error) {
	return c.Service.ProductService_GetValoreMagazzino(ctx)
}

// GetInventoryReport restituisce il report completo dello stato del magazzino.
func (c *ProductController) ProductController_GetInventoryReport(ctx context.Context) (*domain.InventoryReport, error) {
	return c.Service.ProductService_GetAllGiacenze(ctx)
}

// GetLowStockProducts restituisce l'elenco dei prodotti che sono sotto la soglia minima di giacenza.
func (c *ProductController) ProductController_GetProdottiSottoSoglia(ctx context.Context) ([]domain.InventoryItem, error) {
	return c.Service.ProductService_GetProdottiSottoSoglia(ctx)
}
