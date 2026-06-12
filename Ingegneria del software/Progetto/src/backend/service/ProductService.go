package Service

import (
	"context"
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
	// Simple validation on barcode length (EAN-13 expects 13 digits)
	if len(p.Barcode) != 13 {
		return domain.ErrBarcodeNonValido
	}

	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		return s.Repo.Create(txCtx, p)
	})
}

func (s *ProductService) UpdateProduct(ctx context.Context, p *domain.Prodotto) error {
	if len(p.Barcode) != 13 {
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
	if len(barcode) != 13 {
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

// GetValoreMagazzino computes total quantity and total monetary value of inventory (UC7)
func (s *ProductService) GetValoreMagazzino(ctx context.Context) (int, float64, error) {
	giacenze, err := s.Repo.GetAllGiacenze(ctx)
	if err != nil {
		return 0, 0.0, err
	}

	var totQty int
	var totVal float64

	for _, g := range giacenze {
		p, err := s.Repo.GetByID(ctx, g.ProdottoID)
		if err == nil {
			totQty += g.Quantita
			totVal += float64(g.Quantita) * p.PrezzoVend
		}
	}
	return totQty, totVal, nil
}
