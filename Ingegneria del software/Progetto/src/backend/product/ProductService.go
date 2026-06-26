package product

import (
	"context"
	"strings"
	"time"
	"gialosai/backend/domain"
)

type Service struct {
	Repo *Repository
	Tx   *domain.TransactionHelper
}

func NewService(repo *Repository, tx *domain.TransactionHelper) *Service {
	return &Service{Repo: repo, Tx: tx}
}

// CreateProduct crea un nuovo prodotto nel sistema
func (s *Service) CreateProduct(ctx context.Context, p *domain.Prodotto) error {
	// Trimma e valida il barcode
	p.Barcode = strings.TrimSpace(p.Barcode)
	if !domain.IsValidEAN13(p.Barcode) {
		return domain.ErrBarcodeNonValido
	}

	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		return s.Repo.Create(txCtx, p)
	})
}

// UpdateProduct modifica un prodotto esistente nel sistema
func (s *Service) UpdateProduct(ctx context.Context, p *domain.Prodotto) error {
	// Trimma e valida il barcode
	p.Barcode = strings.TrimSpace(p.Barcode)
	if !domain.IsValidEAN13(p.Barcode) {
		return domain.ErrBarcodeNonValido
	}
	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		return s.Repo.Update(txCtx, p)
	})
}

// DeleteProduct rimuove un prodotto dal sistema. Nota: La giacenza associata viene eliminata automaticamente dal repository
func (s *Service) DeleteProduct(ctx context.Context, id int) error {
	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		return s.Repo.Delete(txCtx, id)
	})
}

// GetProducts recupera tutti i prodotti disponibili nel sistema.
func (s *Service) GetProducts(ctx context.Context) ([]domain.Prodotto, error) {
	return s.Repo.GetAll(ctx)
}

// GetProductByID recupera un singolo prodotto dal sistema usando il suo ID.
func (s *Service) GetProductByID(ctx context.Context, id int) (*domain.Prodotto, error) {
	return s.Repo.GetByID(ctx, id)
}

// GetProductByBarcode recupera un prodotto dal sistema usando il suo barcode (EAN-13). Valida il barcode prima di effettuare la ricerca.
func (s *Service) GetProductByBarcode(ctx context.Context, barcode string) (*domain.Prodotto, error) {
	barcode = strings.TrimSpace(barcode)
	if !domain.IsValidEAN13(barcode) {
		return nil, domain.ErrBarcodeNonValido
	}
	return s.Repo.GetByBarcode(ctx, barcode)
}

// UpdateListinoPrices aggiorna i prezzi di acquisto e vendita di un prodotto
func (s *Service) UpdateListinoPrices(ctx context.Context, productID int, acqPrice, vendPrice float64) error {
	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		p, err := s.Repo.GetByID(txCtx, productID)
		if err != nil {
			return err
		}
		p.PrezzoAcq = acqPrice
		p.PrezzoVend = vendPrice
		return s.Repo.Update(txCtx, p)
	})
}

// GetInventarioReport costruisce un riepilogo completo dello stato del magazzino
func (s *Service) GetInventarioReport(ctx context.Context) (*domain.InventoryReport, error) {
	// Recupera tutti i prodotti
	products, err := s.Repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// Recupera tutte le giacenze
	giacenze, err := s.Repo.GetAllGiacenze(ctx)
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
func (s *Service) GetProdottiSottoSoglia(ctx context.Context) ([]domain.InventoryItem, error) {
	report, err := s.GetInventarioReport(ctx)
	if err != nil {
		return nil, err
	}
	return report.ArticoliSottoSoglia, nil
}

// GetValoreMagazzino restituisce un riepilogo sintetico del magazzino: quantità totale e valore totale.
func (s *Service) GetValoreMagazzino(ctx context.Context) (int, float64, error) {
	report, err := s.GetInventarioReport(ctx)
	if err != nil {
		return 0, 0.0, err
	}

	return report.TotaleQuantita, report.ValoreTotale, nil
}
