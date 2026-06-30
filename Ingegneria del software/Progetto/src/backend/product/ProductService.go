package product

import (
	"context"
	"github.com/tommasogarofani/ingsw_25-26_group-11/src/backend/domain"
	"strings"
	"time"
)

type productRepository interface {
	ProductRepository_GetAll(ctx context.Context) ([]domain.Prodotto, error)
	ProductRepository_GetByID(ctx context.Context, id int) (*domain.Prodotto, error)
	ProductRepository_GetByBarcode(ctx context.Context, barcode string) (*domain.Prodotto, error)
	ProductRepository_Create(ctx context.Context, p *domain.Prodotto) error
	ProductRepository_Update(ctx context.Context, p *domain.Prodotto) error
	ProductRepository_Delete(ctx context.Context, id int) error
	ProductRepository_GetAllGiacenze(ctx context.Context) ([]domain.Giacenza, error)
}

type transactionRunner interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type ProductService struct {
	Repo productRepository
	Tx   transactionRunner
}

// ProductService_GetProducts implements [productService].
func (s *ProductService) ProductService_GetProducts(ctx context.Context) ([]domain.Prodotto, error) {
	panic("unimplemented")
}

func NewProductService(repo productRepository, tx transactionRunner) *ProductService {
	return &ProductService{Repo: repo, Tx: tx}
}

// CreateProduct crea un nuovo prodotto nel sistema
func (s *ProductService) ProductService_Create(ctx context.Context, p *domain.Prodotto) error {
	// Trimma e valida il barcode
	p.Barcode = strings.TrimSpace(p.Barcode)
	if !domain.IsValidEAN13(p.Barcode) {
		return domain.ErrBarcodeNonValido
	}

	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		return s.Repo.ProductRepository_Create(txCtx, p)
	})
}

// UpdateProduct modifica un prodotto esistente nel sistema
func (s *ProductService) ProductService_Update(ctx context.Context, p *domain.Prodotto) error {
	// Trimma e valida il barcode
	p.Barcode = strings.TrimSpace(p.Barcode)
	if !domain.IsValidEAN13(p.Barcode) {
		return domain.ErrBarcodeNonValido
	}
	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		return s.Repo.ProductRepository_Update(txCtx, p)
	})
}

// DeleteProduct rimuove un prodotto dal sistema. Nota: La giacenza associata viene eliminata automaticamente dal repository
func (s *ProductService) ProductService_Delete(ctx context.Context, id int) error {
	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		return s.Repo.ProductRepository_Delete(txCtx, id)
	})
}

// GetProducts recupera tutti i prodotti disponibili nel sistema.
func (s *ProductService) ProductService_GetAll(ctx context.Context) ([]domain.Prodotto, error) {
	return s.Repo.ProductRepository_GetAll(ctx)
}

// GetProductByID recupera un singolo prodotto dal sistema usando il suo ID.
func (s *ProductService) ProductService_GetByID(ctx context.Context, id int) (*domain.Prodotto, error) {
	return s.Repo.ProductRepository_GetByID(ctx, id)
}

// GetProductByBarcode recupera un prodotto dal sistema usando il suo barcode (EAN-13). Valida il barcode prima di effettuare la ricerca.
func (s *ProductService) ProductService_GetByBarcode(ctx context.Context, barcode string) (*domain.Prodotto, error) {
	barcode = strings.TrimSpace(barcode)
	if !domain.IsValidEAN13(barcode) {
		return nil, domain.ErrBarcodeNonValido
	}
	return s.Repo.ProductRepository_GetByBarcode(ctx, barcode)
}

// UpdateListinoPrices aggiorna i prezzi di acquisto e vendita di un prodotto
func (s *ProductService) ProductService_UpdateListinoPrices(ctx context.Context, productID int, acqPrice, vendPrice float64) error {
	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		p, err := s.Repo.ProductRepository_GetByID(txCtx, productID)
		if err != nil {
			return err
		}
		p.PrezzoAcq = acqPrice
		p.PrezzoVend = vendPrice
		return s.Repo.ProductRepository_Update(txCtx, p)
	})
}

// GetInventarioReport costruisce un riepilogo completo dello stato del magazzino
func (s *ProductService) ProductService_GetInventarioReport(ctx context.Context) (*domain.InventoryReport, error) {
	// Recupera tutti i prodotti
	products, err := s.Repo.ProductRepository_GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// Recupera tutte le giacenze
	giacenze, err := s.Repo.ProductRepository_GetAllGiacenze(ctx)
	if err != nil {
		return nil, err
	}

	// Crea una mappa dei prodotti per accesso rapido
	productsByID := make(map[int]domain.Prodotto, len(products))
	for _, product := range products {
		productsByID[product.ID] = product
	}

	// Costruisce il report
	report := &domain.InventoryReport{GeneratoIl: time.Now()}
	for _, giacenza := range giacenze {
		product, ok := productsByID[giacenza.ProdottoID]
		if !ok {
			continue
		}

		// Estrae il nome della categoria se presente
		categoria := ""
		if product.Categoria != nil {
			categoria = product.Categoria.Nome
		}

		// Calcola i valori dell'articolo
		item := domain.InventoryItem{
			ProdottoID:  giacenza.ProdottoID,
			Codice:      product.Codice,
			Nome:        product.Nome,
			Categoria:   categoria,
			Quantita:    giacenza.Quantita,
			SogliaMin:   product.SogliaMin,
			PrezzoVend:  product.PrezzoVend,
			ValoreRiga:  float64(giacenza.Quantita) * product.PrezzoVend,
			SottoSoglia: giacenza.Quantita < product.SogliaMin,
		}

		// Aggiunge l'articolo al report
		report.Articoli = append(report.Articoli, item)
		report.TotaleQuantita += giacenza.Quantita
		report.ValoreTotale += item.ValoreRiga

		// Se l'articolo è sotto soglia, lo aggiunge anche alla lista apposita
		if item.SottoSoglia {
			report.ArticoliSottoSoglia = append(report.ArticoliSottoSoglia, item)
		}
	}

	return report, nil
}

// GetProdottiSottoSoglia restituisce solo gli articoli che richiedono attenzione perché sotto soglia minima.
func (s *ProductService) ProductService_GetProdottiSottoSoglia(ctx context.Context) ([]domain.InventoryItem, error) {
	report, err := s.ProductService_GetInventarioReport(ctx)
	if err != nil {
		return nil, err
	}
	return report.ArticoliSottoSoglia, nil
}

// GetValoreMagazzino restituisce un riepilogo sintetico del magazzino: quantità totale e valore totale.
func (s *ProductService) ProductService_GetValoreMagazzino(ctx context.Context) (int, float64, error) {
	report, err := s.ProductService_GetInventarioReport(ctx)
	if err != nil {
		return 0, 0.0, err
	}

	return report.TotaleQuantita, report.ValoreTotale, nil
}
