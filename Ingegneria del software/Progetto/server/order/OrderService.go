package order

import (
	"context"
	"errors"
	"fmt"
	"github.com/tommasogarofani/ingsw_25-26_group-11/server/domain"
	"time"
)

type orderRepository interface {
	OrderRepository_CreateOrdineEntrata(ctx context.Context, o *domain.OrdineEntrata) error
	OrderRepository_GetOrdineEntrataByID(ctx context.Context, id int) (*domain.OrdineEntrata, error)
	OrderRepository_GetOrdiniEntrata(ctx context.Context) ([]domain.OrdineEntrata, error)
	OrderRepository_DeleteOrdineEntrata(ctx context.Context, o *domain.OrdineEntrata) error

	OrderRepository_CreateOrdineUscita(ctx context.Context, o *domain.OrdineUscita) error
	OrderRepository_GetOrdineUscitaByID(ctx context.Context, id int) (*domain.OrdineUscita, error)
	OrderRepository_GetOrdiniUscita(ctx context.Context) ([]domain.OrdineUscita, error)
	OrderRepository_GetOrdiniUscitaByStato(ctx context.Context, stato domain.StatoOrdineUscita) ([]domain.OrdineUscita, error)
	OrderRepository_UpdateOrdineUscita(ctx context.Context, o *domain.OrdineUscita) error
	OrderRepository_DeleteOrdineUscita(ctx context.Context, o *domain.OrdineUscita) error

	OrderRepository_CreateDDT(ctx context.Context, d *domain.DDT) error
	OrderRepository_GetDDTByYear(ctx context.Context, year int) ([]domain.DDT, error)
	OrderRepository_GetDDTByOrdineUscita(ctx context.Context, orderID int) (*domain.DDT, error)

	OrderRepository_CreateReso(ctx context.Context, re *domain.Reso) error
	OrderRepository_GetResoByID(ctx context.Context, id int) (*domain.Reso, string, error)
	OrderRepository_UpdateResoStato(ctx context.Context, id int, stato string) error
	OrderRepository_GetAllResi(ctx context.Context) ([]domain.Reso, error)
	OrderRepository_DeleteReso(ctx context.Context, id int) error

	OrderRepository_CreateFattura(ctx context.Context, f *domain.Fattura) error
	OrderRepository_CountFattureByYearAndType(ctx context.Context, year int, tipo domain.TipoFattura) (int, error)
	OrderRepository_GetAllFatture(ctx context.Context) ([]domain.Fattura, error)
	OrderRepository_DeleteFattura(ctx context.Context, id int) error
}

// productRepoForOrders è il sottoinsieme di product.Repository che serve qui.
type productRepoForOrders interface {
	ProductRepository_GetByID(ctx context.Context, id int) (*domain.Prodotto, error)
	ProductRepository_GetGiacenzaById(ctx context.Context, prodID int) (*domain.Giacenza, error)
	ProductRepository_UpdateGiacenza(ctx context.Context, g *domain.Giacenza) error
}

type transactionRunner interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// OrderService implementa la logica di business per ordini di entrata e uscita,
// DDT, resi e fatture, coordinando il repository degli ordini, quello dei prodotti
// (per la gestione delle giacenze) e le transazioni.
type OrderService struct {
	OrderRepo   orderRepository
	ProductRepo productRepoForOrders
	Tx          transactionRunner
}

// NewOrderService crea un nuovo OrderService con le dipendenze specificate.
func NewOrderService(orderRepo orderRepository, prodRepo productRepoForOrders, tx transactionRunner) *OrderService {
	return &OrderService{OrderRepo: orderRepo, ProductRepo: prodRepo, Tx: tx}
}

// ─────────────────────────────────────────────
// ORDINI ENTRATA
// ─────────────────────────────────────────────

// OrderService_ConfirmOrdineEntrata conferma un ordine di entrata, aggiungendo la data
// corrente, salvando le righe e incrementando la giacenza dei prodotti, il tutto
// all'interno di una transazione.
func (s *OrderService) OrderService_ConfirmOrdineEntrata(ctx context.Context, o *domain.OrdineEntrata) error {
	o.Data = time.Now()

	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		// salva l'ordine di entrata
		err := s.OrderRepo.OrderRepository_CreateOrdineEntrata(txCtx, o)
		if err != nil {
			return err
		}

		// aggiorna la giacenza dei prodotti in entrata
		for _, r := range o.Righe {
			err = s.ProductRepo.ProductRepository_UpdateGiacenza(txCtx, &domain.Giacenza{
				ProdottoID: r.ProdottoID,
				Quantita:   r.Quantita,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// OrderService_GetOrdiniEntrata recupera tutti gli ordini di entrata.
func (s *OrderService) OrderService_GetOrdiniEntrata(ctx context.Context) ([]domain.OrdineEntrata, error) {
	return s.OrderRepo.OrderRepository_GetOrdiniEntrata(ctx)
}

// OrderService_GetOrdiniEntrataById recupera un ordine di entrata in base al suo ID.
func (s *OrderService) OrderService_GetOrdiniEntrataById(ctx context.Context, id int) (*domain.OrdineEntrata, error) {
	return s.OrderRepo.OrderRepository_GetOrdineEntrataByID(ctx, id)
}

// OrderService_DeleteOrdineEntrata elimina un ordine di entrata.
func (s *OrderService) OrderService_DeleteOrdineEntrata(ctx context.Context, o *domain.OrdineEntrata) error {
	return s.OrderRepo.OrderRepository_DeleteOrdineEntrata(ctx, o)
}

// ─────────────────────────────────────────────
// ORDINI USCITA
// ─────────────────────────────────────────────

// OrderService_ConfirmOrdineUscita conferma un ordine di uscita, salvandolo con stato
// iniziale confermato e, in una transazione, verificando che la giacenza sia sufficiente
// per ogni prodotto e decrementandola di conseguenza.
func (s *OrderService) OrderService_ConfirmOrdineUscita(ctx context.Context, o *domain.OrdineUscita) error {
	o.Data = time.Now()
	o.Stato = domain.StatoConfermato // stato iniziale dell'ordine di uscita

	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		// verifica la giacenza dei prodotti in uscita e decrementa la giacenza
		for _, r := range o.Righe {
			g, err := s.ProductRepo.ProductRepository_GetGiacenzaById(txCtx, r.ProdottoID)
			if err != nil {
				return err
			}
			if g.Quantita < r.Quantita {
				p, errP := s.ProductRepo.ProductRepository_GetByID(txCtx, r.ProdottoID)
				name := "Articolo"
				if errP == nil {
					name = p.Nome
				}
				return domain.NewGiacenzaInsufficienteError(r.ProdottoID, name, g.Quantita, r.Quantita)
			}

			// decrementa la giacenza del prodotto
			err = s.ProductRepo.ProductRepository_UpdateGiacenza(txCtx, &domain.Giacenza{
				ProdottoID: r.ProdottoID,
				Quantita:   -r.Quantita,
			})
			if err != nil {
				return err
			}
		}

		// salva l'ordine di uscita
		return s.OrderRepo.OrderRepository_CreateOrdineUscita(txCtx, o)
	})
}

// OrderService_GetOrdiniUscitaNonEvasi recupera gli ordini di uscita in stato confermato
// (cioè non ancora evasi).
func (s *OrderService) OrderService_GetOrdiniUscitaNonEvasi(ctx context.Context) ([]domain.OrdineUscita, error) {
	return s.OrderRepo.OrderRepository_GetOrdiniUscitaByStato(ctx, domain.StatoConfermato)
}

// OrderService_GetOrdiniUscita recupera tutti gli ordini di uscita.
func (s *OrderService) OrderService_GetOrdiniUscita(ctx context.Context) ([]domain.OrdineUscita, error) {
	return s.OrderRepo.OrderRepository_GetOrdiniUscita(ctx)
}

// OrderService_GetOrdiniUscitaById recupera un ordine di uscita in base al suo ID.
func (s *OrderService) OrderService_GetOrdiniUscitaById(ctx context.Context, id int) (*domain.OrdineUscita, error) {
	return s.OrderRepo.OrderRepository_GetOrdineUscitaByID(ctx, id)
}

// OrderService_GetOrdiniUscitaByStato recupera gli ordini di uscita filtrati per stato.
func (s *OrderService) OrderService_GetOrdiniUscitaByStato(ctx context.Context, stato domain.StatoOrdineUscita) ([]domain.OrdineUscita, error) {
	return s.OrderRepo.OrderRepository_GetOrdiniUscitaByStato(ctx, stato)
}

// OrderService_UpdateOrdineUscita aggiorna un ordine di uscita gestendo le transizioni
// di stato relative alle giacenze: annullando un ordine le giacenze vengono ripristinate,
// riattivandolo vengono riscaricate. In assenza di transizione rilevante aggiorna solo
// i dati dell'ordine.
func (s *OrderService) OrderService_UpdateOrdineUscita(ctx context.Context, o *domain.OrdineUscita) error {
	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		// carica lo stato corrente dell'ordine per determinare la transizione di stato
		old, err := s.OrderRepo.OrderRepository_GetOrdineUscitaByID(txCtx, o.ID)
		if err != nil {
			return err
		}

		wasAnnullato := old.Stato == domain.StatoAnnullato
		willBeAnnullato := o.Stato == domain.StatoAnnullato

		switch {
		case !wasAnnullato && willBeAnnullato:
			// l'ordine viene annullato: la merce torna in magazzino, quindi le giacenze
			// precedentemente scaricate vanno ripristinate (incremento).
			for _, r := range o.Righe {
				err = s.ProductRepo.ProductRepository_UpdateGiacenza(txCtx, &domain.Giacenza{
					ProdottoID: r.ProdottoID,
					Quantita:   r.Quantita,
				})
				if err != nil {
					return err
				}
			}

		case wasAnnullato && !willBeAnnullato:
			// l'ordine annullato viene riattivato: le giacenze vanno scaricate di nuovo.
			for _, r := range o.Righe {
				g, err := s.ProductRepo.ProductRepository_GetGiacenzaById(txCtx, r.ProdottoID)
				if err != nil {
					return err
				}
				if g.Quantita < r.Quantita {
					p, errP := s.ProductRepo.ProductRepository_GetByID(txCtx, r.ProdottoID)
					name := "Articolo"
					if errP == nil {
						name = p.Nome
					}
					return domain.NewGiacenzaInsufficienteError(r.ProdottoID, name, g.Quantita, r.Quantita)
				}

				err = s.ProductRepo.ProductRepository_UpdateGiacenza(txCtx, &domain.Giacenza{
					ProdottoID: r.ProdottoID,
					Quantita:   -r.Quantita,
				})
				if err != nil {
					return err
				}
			}
		}

		// negli altri casi (nessuna transizione rilevante con lo stato precedente)
		// la giacenza non viene toccata e l'ordine viene semplicemente aggiornato.

		return s.OrderRepo.OrderRepository_UpdateOrdineUscita(txCtx, o)
	})
}

// OrderService_DeleteOrdineUscita elimina un ordine di uscita; se l'ordine non era
// annullato, ripristina (incrementa) le giacenze dei prodotti precedentemente scaricate.
func (s *OrderService) OrderService_DeleteOrdineUscita(ctx context.Context, o *domain.OrdineUscita) error {
	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		// se l'ordine non era stato annullato, la giacenza era già stata scaricata:
		// eliminarlo ripristina le giacenze (incremento).
		if o.Stato != domain.StatoAnnullato {
			for _, r := range o.Righe {
				err := s.ProductRepo.ProductRepository_UpdateGiacenza(txCtx, &domain.Giacenza{
					ProdottoID: r.ProdottoID,
					Quantita:   r.Quantita,
				})
				if err != nil {
					return err
				}
			}
		}

		return s.OrderRepo.OrderRepository_DeleteOrdineUscita(txCtx, o)
	})
}

// ─────────────────────────────────────────────
// GENERAZIONE DDT
// ─────────────────────────────────────────────

// OrderService_GeneraDDT genera un DDT per un ordine di uscita, assegnandogli un numero
// progressivo per l'anno corrente e portando l'ordine allo stato "evaso" associandogli il
// DDT creato. Se l'ordine ha già un DDT, restituisce quello esistente.
func (s *OrderService) OrderService_GeneraDDT(ctx context.Context, orderID int, mittente, destinatario string) (*domain.DDT, error) {
	var ddt *domain.DDT

	err := s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		o, err := s.OrderRepo.OrderRepository_GetOrdineUscitaByID(txCtx, orderID)
		if err != nil {
			return err
		}

		if o.DDTID != nil {
			ddt, err = s.OrderRepo.OrderRepository_GetDDTByOrdineUscita(txCtx, orderID)
			if err != nil {
				return err
			}

			return nil
		}

		progressive, err := s.OrderRepo.OrderRepository_GetDDTByYear(txCtx, time.Now().Year())
		if err != nil {
			return err
		}

		// crea il DDT
		d := &domain.DDT{
			Numero:         fmt.Sprintf("%s-%d-%04d", domain.PrefissoDDT, time.Now().Year(), len(progressive)+1),
			OrdineUscitaID: o.ID,
			Data:           time.Now(),
			Mittente:       mittente,
			Destinatario:   destinatario,
			PathPDF:        fmt.Sprintf("./output/pdf/ddt_%d.pdf", o.ID),
		}

		err = s.OrderRepo.OrderRepository_CreateDDT(txCtx, d)
		if err != nil {
			return err
		}

		// cambia lo stato dell'ordine a EVASO e associa il DDT
		o.Stato = domain.StatoEvaso
		o.DDTID = &d.ID
		err = s.OrderRepo.OrderRepository_UpdateOrdineUscita(txCtx, o)
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

// OrderService_GetDDTByOrdineUscita recupera il DDT associato a un ordine di uscita.
func (s *OrderService) OrderService_GetDDTByOrdineUscita(ctx context.Context, orderID int) (*domain.DDT, error) {
	return s.OrderRepo.OrderRepository_GetDDTByOrdineUscita(ctx, orderID)
}

// OrderService_GetDDTByYear recupera tutti i DDT emessi nell'anno indicato.
func (s *OrderService) OrderService_GetDDTByYear(ctx context.Context, year int) ([]domain.DDT, error) {
	return s.OrderRepo.OrderRepository_GetDDTByYear(ctx, year)
}

// ─────────────────────────────────────────────
// GESTIONE RESI
// ─────────────────────────────────────────────

// OrderService_CreateReso registra un reso, impostandone la data corrente e verificando
// che l'ordine di uscita associato sia nello stato "evaso"; in caso contrario restituisce
// un errore.
func (s *OrderService) OrderService_CreateReso(ctx context.Context, r *domain.Reso) error {
	r.Data = time.Now()
	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		// verifica che l'ordine di uscita sia evaso prima di permettere il reso
		o, err := s.OrderRepo.OrderRepository_GetOrdineUscitaByID(txCtx, r.OrdineUscitaID)
		if err != nil {
			return err
		}
		if o.Stato != domain.StatoEvaso {
			return errors.New("il reso può essere registrato solo per ordini evasi")
		}

		return s.OrderRepo.OrderRepository_CreateReso(txCtx, r)
	})
}

// OrderService_VerifyReso verifica fisicamente un reso in stato "da_verificare",
// incrementando la giacenza del prodotto reso e portando il reso allo stato "accettato".
func (s *OrderService) OrderService_VerifyReso(ctx context.Context, resoID int) error {
	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		re, stato, err := s.OrderRepo.OrderRepository_GetResoByID(txCtx, resoID)
		if err != nil {
			return err
		}
		if stato != "da_verificare" {
			return errors.New("il reso è già stato elaborato")
		}

		// incrementa la giacenza del prodotto reso
		err = s.ProductRepo.ProductRepository_UpdateGiacenza(txCtx, &domain.Giacenza{
			ProdottoID: re.ProdottoID,
			Quantita:   re.Quantita,
		})
		if err != nil {
			return err
		}

		return s.OrderRepo.OrderRepository_UpdateResoStato(txCtx, resoID, "accettato")
	})
}

// OrderService_ApproveReso approva un reso che è già stato verificato fisicamente
// (stato "accettato"). Restituisce un errore se il reso non è nello stato corretto.
func (s *OrderService) OrderService_ApproveReso(ctx context.Context, resoID int) error {
	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		_, stato, err := s.OrderRepo.OrderRepository_GetResoByID(txCtx, resoID)
		if err != nil {
			return err
		}
		if stato != "accettato" {
			return errors.New("il reso deve essere prima verificato fisicamente")
		}
		return s.OrderRepo.OrderRepository_UpdateResoStato(txCtx, resoID, "accettato")
	})
}

// OrderService_RifiutaReso rifiuta un reso precedentemente accettato, annullando
// l'incremento di giacenza effettuato in verifica e portando il reso allo stato "rifiutato".
func (s *OrderService) OrderService_RifiutaReso(ctx context.Context, resoID int) error {
	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {

		re, stato, err := s.OrderRepo.OrderRepository_GetResoByID(txCtx, resoID)
		if err != nil {
			return err
		}

		if stato != "accettato" {
			return errors.New("reso non verificato o già processato")
		}

		// Rollback della giacenza del prodotto reso
		err = s.ProductRepo.ProductRepository_UpdateGiacenza(txCtx, &domain.Giacenza{
			ProdottoID: re.ProdottoID,
			Quantita:   -re.Quantita,
		})
		if err != nil {
			return err
		}

		return s.OrderRepo.OrderRepository_UpdateResoStato(txCtx, resoID, "rifiutato")
	})
}

// OrderService_GetAllResi recupera tutti i resi convertiti in DTO.
func (s *OrderService) OrderService_GetAllResi(ctx context.Context) ([]ResoDTO, error) {
	resi, err := s.OrderRepo.OrderRepository_GetAllResi(ctx)
	if err != nil {
		return nil, err
	}

	var dtoList []ResoDTO

	for _, r := range resi {
		dtoList = append(dtoList, ResoDTO{
			ID:             r.ID,
			OrdineUscitaID: r.OrdineUscitaID,
			ProdottoID:     r.ProdottoID,
			Quantita:       r.Quantita,
			Motivo:         r.Motivo,
			Stato:          string(r.Stato),
			CreatedByID:    r.CreatedByID,
			Data:           r.Data,
		})
	}

	return dtoList, nil
}

// OrderService_DeleteReso elimina un reso in base al suo ID.
func (s *OrderService) OrderService_DeleteReso(ctx context.Context, resoID int) error {
	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		return s.OrderRepo.OrderRepository_DeleteReso(txCtx, resoID)
	})
}

// ─────────────────────────────────────────────
// GESTIONE FATTURE
// ─────────────────────────────────────────────

// OrderService_GeneraFattura emette una fattura di vendita per un ordine evaso,
// calcolando imponibile, IVA (22%) e totale, e assegnando un numero progressivo per l'anno.
// Restituisce un errore se l'ordine non è ancora evaso.
func (s *OrderService) OrderService_GeneraFattura(ctx context.Context, orderID int) (*domain.Fattura, error) {
	var fat *domain.Fattura

	err := s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		o, err := s.OrderRepo.OrderRepository_GetOrdineUscitaByID(txCtx, orderID)
		if err != nil {
			return err
		}
		if o.Stato != domain.StatoEvaso {
			return errors.New("impossibile emettere fattura per un ordine non ancora evaso")
		}

		progressive, err := s.OrderRepo.OrderRepository_CountFattureByYearAndType(txCtx, time.Now().Year(), domain.TipoFatturaVendita)
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

		err = s.OrderRepo.OrderRepository_CreateFattura(txCtx, f)
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

// OrderService_GeneraNotaCreditoById emette una nota di credito per un ordine evaso,
// calcolando imponibile, IVA (22%) e totale e assegnando un numero progressivo per l'anno.
// Restituisce un errore se l'ordine non è ancora evaso.
func (s *OrderService) OrderService_GeneraNotaCreditoById(ctx context.Context, orderID int) (*domain.Fattura, error) {
	var notaCredito *domain.Fattura

	err := s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		o, err := s.OrderRepo.OrderRepository_GetOrdineUscitaByID(txCtx, orderID)
		if err != nil {
			return err
		}
		if o.Stato != domain.StatoEvaso {
			return errors.New("impossibile emettere nota di credito per un ordine non ancora evaso")
		}

		progressive, err := s.OrderRepo.OrderRepository_CountFattureByYearAndType(txCtx, time.Now().Year(), domain.TipoNotaCredito)
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

		err = s.OrderRepo.OrderRepository_CreateFattura(txCtx, f)
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

// OrderService_GetAllFatture recupera tutte le fatture presenti nel sistema.
func (s *OrderService) OrderService_GetAllFatture(ctx context.Context) ([]domain.Fattura, error) {
	return s.OrderRepo.OrderRepository_GetAllFatture(ctx)
}

// OrderService_GetNoteCredito recupera tutte le note di credito, filtrando le fatture
// di tipo nota di credito.
func (s *OrderService) OrderService_GetNoteCredito(ctx context.Context) ([]domain.Fattura, error) {
	all, err := s.OrderRepo.OrderRepository_GetAllFatture(ctx)
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

// OrderService_DeleteFattura elimina una fattura in base al suo ID.
func (s *OrderService) OrderService_DeleteFattura(ctx context.Context, id int) error {
	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		return s.OrderRepo.OrderRepository_DeleteFattura(txCtx, id)
	})
}
