package order

import (
	"context"
	"errors"
	"fmt"
	"gialosai/backend/domain"
	"gialosai/backend/product"
	"time"
)

type OrderService struct {
	OrderRepo   *OrderRepository
	ProductRepo *product.Repository
	Tx          *db.TransactionHelper
}

func NewOrderService(orderRepo *OrderRepository, prodRepo *product.Repository, tx *db.TransactionHelper) *OrderService {
	return &OrderService{OrderRepo: orderRepo, ProductRepo: prodRepo, Tx: tx}
}

// ─────────────────────────────────────────────
// ORDINI ENTRATA
// ─────────────────────────────────────────────

func (s *OrderService) ConfirmOrdineEntrata(ctx context.Context, o *domain.OrdineEntrata) error {
	o.Data = time.Now()
	o.Stato = domain.StatoConfermato

	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		// salva l'ordine di entrata
		err := s.OrderRepo.CreateOrdineEntrata(txCtx, o)
		if err != nil {
			return err
		}

		// aggiorna la giacenza dei prodotti in entrata
		for _, r := range o.Righe {
			err = s.ProductRepo.UpdateGiacenza(txCtx, r.ProdottoID, r.Quantita)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// ─────────────────────────────────────────────
// ORDINI USCITA
// ─────────────────────────────────────────────

func (s *OrderService) ConfirmOrdineUscita(ctx context.Context, o *domain.OrdineUscita) error {
	o.Data = time.Now()
	o.Stato = domain.StatoConfermato // stato iniziale dell'ordine di uscita

	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		// verifica la giacenza dei prodotti in uscita e decrementa la giacenza
		for _, r := range o.Righe {
			g, err := s.ProductRepo.GetGiacenza(txCtx, r.ProdottoID)
			if err != nil {
				return err
			}
			if g.Quantita < r.Quantita {
				p, errP := s.ProductRepo.GetByID(txCtx, r.ProdottoID)
				name := "Articolo"
				if errP == nil {
					name = p.Nome
				}
				return domain.NewGiacenzaInsufficienteError(r.ProdottoID, name, g.Quantita, r.Quantita)
			}

			// decrementa la giacenza del prodotto
			err = s.ProductRepo.UpdateGiacenza(txCtx, r.ProdottoID, -r.Quantita)
			if err != nil {
				return err
			}
		}

		// salva l'ordine di uscita
		return s.OrderRepo.CreateOrdineUscita(txCtx, o)
	})
}

func (s *OrderService) GetOrdiniUscitaNonEvasi(ctx context.Context) ([]domain.OrdineUscita, error) {
	return s.OrderRepo.GetOrdiniUscitaByStato(ctx, domain.StatoConfermato)
}

func (s *OrderService) GetOrdiniUscita(ctx context.Context) ([]domain.OrdineUscita, error) {
	return s.OrderRepo.GetOrdiniUscita(ctx)
}

func (s *OrderService) GetOrdiniUscitaByStato(ctx context.Context, stato domain.StatoOrdine) ([]domain.OrdineUscita, error) {
	return s.OrderRepo.GetOrdiniUscitaByStato(ctx, stato)
}

func (s *OrderService) GetDettaglioOrdineUscita(ctx context.Context, id int) (*domain.OrdineUscita, error) {
	return s.OrderRepo.GetOrdineUscitaByID(ctx, id)
}

// ─────────────────────────────────────────────
// GENERAZIONE DDT
// ─────────────────────────────────────────────

func (s *OrderService) GeneraDDT(ctx context.Context, orderID int, mittente, destinatario string) (*domain.DDT, error) {
	var ddt *domain.DDT

	err := s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		o, err := s.OrderRepo.GetOrdineUscitaByID(txCtx, orderID)
		if err != nil {
			return err
		}

		if o.Stato == domain.StatoEvaso {
			return domain.ErrDDTGiaGenerato
		}

		progressive, err := s.OrderRepo.CountDDTByYear(txCtx, time.Now().Year())
		if err != nil {
			return err
		}

		// crea il DDT
		d := &domain.DDT{
			Numero:         fmt.Sprintf("%s-%d-%04d", domain.PrefissoDDT, time.Now().Year(), progressive+1),
			OrdineUscitaID: o.ID,
			Data:           time.Now(),
			Mittente:       mittente,
			Destinatario:   destinatario,
			PathPDF:        fmt.Sprintf("./output/pdf/ddt_%d.pdf", o.ID),
		}

		err = s.OrderRepo.CreateDDT(txCtx, d)
		if err != nil {
			return err
		}

		// cambia lo stato dell'ordine a EVASO e associa il DDT
		o.Stato = domain.StatoEvaso
		o.DDTID = &d.ID
		err = s.OrderRepo.UpdateOrdineUscita(txCtx, o)
		if err != nil {
			return err
		}

		ddt = d
		return nil
	})

	if err != nil {
		return nil, err
	}
	return ddt, nil
}

// ─────────────────────────────────────────────
// GESTIONE RESI
// ─────────────────────────────────────────────

func (s *OrderService) CreateReso(ctx context.Context, r *domain.Reso) error {
	r.Data = time.Now()
	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		// verifica che l'ordine di uscita sia evaso prima di permettere il reso
		o, err := s.OrderRepo.GetOrdineUscitaByID(txCtx, r.OrdineUscitaID)
		if err != nil {
			return err
		}
		if o.Stato != domain.StatoEvaso {
			return errors.New("il reso può essere registrato solo per ordini evasi")
		}

		return s.OrderRepo.CreateReso(txCtx, r, "da_verificare")
	})
}

func (s *OrderService) VerifyReso(ctx context.Context, resoID int) error {
	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		re, stato, err := s.OrderRepo.GetResoByID(txCtx, resoID)
		if err != nil {
			return err
		}
		if stato != "da_verificare" {
			return errors.New("il reso è già stato elaborato")
		}

		// incrementa la giacenza del prodotto reso
		err = s.ProductRepo.UpdateGiacenza(txCtx, re.ProdottoID, re.Quantita)
		if err != nil {
			return err
		}

		return s.OrderRepo.UpdateResoStato(txCtx, resoID, "verificato")
	})
}

func (s *OrderService) ApproveReso(ctx context.Context, resoID int) error {
	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		_, stato, err := s.OrderRepo.GetResoByID(txCtx, resoID)
		if err != nil {
			return err
		}
		if stato != "verificato" {
			return errors.New("il reso deve essere prima verificato fisicamente")
		}
		return s.OrderRepo.UpdateResoStato(txCtx, resoID, "approvato")
	})
}

func (s *OrderService) RifiutaReso(ctx context.Context, resoID int) error {
	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		re, stato, err := s.OrderRepo.GetResoByID(txCtx, resoID)
		if err != nil {
			return err
		}
		if stato != "verificato" {
			return errors.New("reso non verificato o già processato")
		}

		// Rollback della giacenza del prodotto reso
		err = s.ProductRepo.UpdateGiacenza(txCtx, re.ProdottoID, -re.Quantita)
		if err != nil {
			return err
		}

		return s.OrderRepo.UpdateResoStato(txCtx, resoID, "rifiutato")
	})
}

func (s *OrderService) GetResi(ctx context.Context) ([]map[string]interface{}, error) {
	return s.OrderRepo.GetAllResi(ctx)
}

// ─────────────────────────────────────────────
// GESTIONE FATTURE
// ─────────────────────────────────────────────

func (s *OrderService) GeneraFattura(ctx context.Context, orderID int) (*domain.Fattura, error) {
	var fat *domain.Fattura

	err := s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		o, err := s.OrderRepo.GetOrdineUscitaByID(txCtx, orderID)
		if err != nil {
			return err
		}
		if o.Stato != domain.StatoEvaso {
			return errors.New("impossibile emettere fattura per un ordine non ancora evaso")
		}

		progressive, err := s.OrderRepo.CountFattureByYearAndType(txCtx, time.Now().Year(), domain.TipoFatturaVendita)
		if err != nil {
			return err
		}

		// calcola imponibile, IVA e totale
		var imponibile float64
		for _, r := range o.Righe {
			imponibile += float64(r.Quantita) * r.PrezzoUnitario
		}
		iva := imponibile * 0.22 // 22% IVA standard
		totale := imponibile + iva

		f := &domain.Fattura{
			Numero:         fmt.Sprintf("%s-%d-%04d", domain.PrefissoFattura, time.Now().Year(), progressive+1),
			Anno:           time.Now().Year(),
			Tipo:           domain.TipoFatturaVendita,
			OrdineUscitaID: o.ID,
			Imponibile:     imponibile,
			IVA:            iva,
			Totale:         totale,
			Data:           time.Now(),
			PathPDF:        fmt.Sprintf("./output/pdf/fattura_%d.pdf", o.ID),
		}

		err = s.OrderRepo.CreateFattura(txCtx, f)
		if err != nil {
			return err
		}
		fat = f
		return nil
	})

	if err != nil {
		return nil, err
	}
	return fat, nil
}

func (s *OrderService) GeneraNotaCredito(ctx context.Context, orderID int) (*domain.Fattura, error) {
	var notaCredito *domain.Fattura

	err := s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		o, err := s.OrderRepo.GetOrdineUscitaByID(txCtx, orderID)
		if err != nil {
			return err
		}
		if o.Stato != domain.StatoEvaso {
			return errors.New("impossibile emettere nota di credito per un ordine non ancora evaso")
		}

		progressive, err := s.OrderRepo.CountFattureByYearAndType(txCtx, time.Now().Year(), domain.TipoNotaCredito)
		if err != nil {
			return err
		}

		var imponibile float64
		for _, r := range o.Righe {
			imponibile += float64(r.Quantita) * r.PrezzoUnitario
		}
		iva := imponibile * 0.22
		totale := imponibile + iva

		f := &domain.Fattura{
			Numero:         fmt.Sprintf("%s-%d-%04d", domain.PrefissoNotaCredito, time.Now().Year(), progressive+1),
			Anno:           time.Now().Year(),
			Tipo:           domain.TipoNotaCredito,
			OrdineUscitaID: o.ID,
			Imponibile:     imponibile,
			IVA:            iva,
			Totale:         totale,
			Data:           time.Now(),
			PathPDF:        fmt.Sprintf("./output/pdf/nota_credito_%d.pdf", o.ID),
		}

		err = s.OrderRepo.CreateFattura(txCtx, f)
		if err != nil {
			return err
		}
		notaCredito = f
		return nil
	})

	if err != nil {
		return nil, err
	}
	return notaCredito, nil
}

func (s *OrderService) GetFatture(ctx context.Context) ([]domain.Fattura, error) {
	return s.OrderRepo.GetAllFatture(ctx)
}

func (s *OrderService) GetNoteCredito(ctx context.Context) ([]domain.Fattura, error) {
	all, err := s.OrderRepo.GetAllFatture(ctx)
	if err != nil {
		return nil, err
	}

	var noteCredito []domain.Fattura
	for _, fattura := range all {
		if fattura.Tipo == domain.TipoNotaCredito {
			noteCredito = append(noteCredito, fattura)
		}
	}
	return noteCredito, nil
}

func (s *OrderService) DeleteFattura(ctx context.Context, id int) error {
	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		return s.OrderRepo.DeleteFattura(txCtx, id)
	})
}
