package Service

import (
	"context"
	"strings"
	"time"

	"gialosai/backend/domain"
	"gialosai/backend/repository"
)

type ProductService struct {
	Repo *Repository.ProductRepository
	Tx   *Repository.TransactionHelper
}

func NewProductService(repo *Repository.ProductRepository, tx *Repository.TransactionHelper) *ProductService {
	return &ProductService{Repo: repo, Tx: tx}
}

func (s *ProductService) CreateProduct(ctx context.Context, p *domain.Prodotto) error {
	p.Barcode = strings.TrimSpace(p.Barcode)
	if !domain.IsValidEAN13(p.Barcode) {
		return domain.ErrBarcodeNonValido
	}

	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		return s.Repo.Create(txCtx, p)
	})
}

func (s *ProductService) UpdateProduct(ctx context.Context, p *domain.Prodotto) error {
	p.Barcode = strings.TrimSpace(p.Barcode)
	if !domain.IsValidEAN13(p.Barcode) {
		return domain.ErrBarcodeNonValido
	}
	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		return s.Repo.Update(txCtx, p)
	})
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int) error {
	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		return s.Repo.Delete(txCtx, id)
	})
}

func (s *ProductService) GetProducts(ctx context.Context) ([]domain.Prodotto, error) {
	return s.Repo.GetAll(ctx)
}

func (s *ProductService) GetProductByID(ctx context.Context, id int) (*domain.Prodotto, error) {
	return s.Repo.GetByID(ctx, id)
}

func (s *ProductService) GetProductByBarcode(ctx context.Context, barcode string) (*domain.Prodotto, error) {
	barcode = strings.TrimSpace(barcode)
	if !domain.IsValidEAN13(barcode) {
		return nil, domain.ErrBarcodeNonValido
	}
	return s.Repo.GetByBarcode(ctx, barcode)
}

// UpdateListinoPrices updates the list prices (RF-14)
func (s *ProductService) UpdateListinoPrices(ctx context.Context, productID int, acqPrice, vendPrice float64) error {
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

// GetInventarioReport costruisce un riepilogo leggibile dello stato del magazzino.
func (s *ProductService) GetInventarioReport(ctx context.Context) (*domain.InventoryReport, error) {
	products, err := s.Repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	giacenze, err := s.Repo.GetAllGiacenze(ctx)
	if err != nil {
		return nil, err
	}

	productsByID := make(map[int]domain.Prodotto, len(products))
	for _, product := range products {
		productsByID[product.ID] = product
	}

	report := &domain.InventoryReport{GeneratoIl: time.Now()}
	for _, giacenza := range giacenze {
		product, ok := productsByID[giacenza.ProdottoID]
		if !ok {
			continue
		}

		categoria := ""
		if product.Categoria != nil {
			categoria = product.Categoria.Nome
		}

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

		report.Articoli = append(report.Articoli, item)
		report.TotaleQuantita += giacenza.Quantita
		report.ValoreTotale += item.ValoreRiga

		if item.SottoSoglia {
			report.ArticoliSottoSoglia = append(report.ArticoliSottoSoglia, item)
		}
	}

	return report, nil
}

// GetProdottiSottoSoglia restituisce solo gli articoli che richiedono attenzione.
func (s *ProductService) GetProdottiSottoSoglia(ctx context.Context) ([]domain.InventoryItem, error) {
	report, err := s.GetInventarioReport(ctx)
	if err != nil {
		return nil, err
	}
	return report.ArticoliSottoSoglia, nil
}

// GetValoreMagazzino restituisce un riepilogo sintetico del magazzino.
func (s *ProductService) GetValoreMagazzino(ctx context.Context) (int, float64, error) {
	report, err := s.GetInventarioReport(ctx)
	if err != nil {
		return 0, 0.0, err
	}

	return report.TotaleQuantita, report.ValoreTotale, nil
}
