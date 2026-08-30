package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/tommasogarofani/ingsw_25-26_group-11/server/domain"
)

// OrderRepository fornisce l'accesso persistente (SQL) alle entità degli ordini:
// ordini di entrata e uscita, relativi DDT, fatture e resi.
type OrderRepository struct {
	DB *sql.DB
}

// NewOrderRepository crea un nuovo OrderRepository legato al database specificato.
func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

// ─────────────────────────────────────────────
// ORDINI DI ENTRATA
// ─────────────────────────────────────────────

// OrderRepository_CreateOrdineEntrata inserisce un ordine di entrata e le sue righe
// nel database, assegnando gli ID auto-generati all'ordine e alle righe.
func (r *OrderRepository) OrderRepository_CreateOrdineEntrata(ctx context.Context, o *domain.OrdineEntrata) error {
	exec := domain.GetExecutor(ctx, r.DB)
	query := `INSERT INTO OrdiniEntrata (numero, fornitore, data, note, created_by)
	          VALUES (?, ?, ?, ?, ?)`
	res, err := exec.ExecContext(ctx, query, o.Numero, o.Fornitore, o.Data, o.Note, o.CreatedByID)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		o.ID = int(id)
	}

	for i := range o.Righe {
		o.Righe[i].OrdineEntrataID = o.ID
		queryR := `INSERT INTO RigheOrdineEntrata (ordine_entrata_id, prodotto_id, quantita, prezzo_unitario)
		           VALUES (?, ?, ?, ?)`
		resR, err := exec.ExecContext(ctx, queryR, o.Righe[i].OrdineEntrataID, o.Righe[i].ProdottoID, o.Righe[i].Quantita, o.Righe[i].PrezzoUnitario)
		if err != nil {
			return err
		}
		idR, err := resR.LastInsertId()
		if err == nil {
			o.Righe[i].ID = int(idR)
		}
	}
	return nil
}

// OrderRepository_GetOrdineEntrataByID recupera un ordine di entrata e le sue righe
// in base all'ID. Restituisce un errore "ordine non trovato" se l'ordine non esiste.
func (r *OrderRepository) OrderRepository_GetOrdineEntrataByID(ctx context.Context, id int) (*domain.OrdineEntrata, error) {
	exec := domain.GetExecutor(ctx, r.DB)
	query := `SELECT id, numero, fornitore, data, note, created_by FROM OrdiniEntrata WHERE id = ?`
	var o domain.OrdineEntrata
	err := exec.QueryRowContext(ctx, query, id).Scan(&o.ID, &o.Numero, &o.Fornitore, &o.Data, &o.Note, &o.CreatedByID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("ordine non trovato")
		}
		return nil, err
	}

	// Load lines
	queryR := `SELECT id, ordine_entrata_id, prodotto_id, quantita, prezzo_unitario FROM RigheOrdineEntrata WHERE ordine_entrata_id = ?`
	rows, err := exec.QueryContext(ctx, queryR, o.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r domain.RigaOrdineEntrata
		err := rows.Scan(&r.ID, &r.OrdineEntrataID, &r.ProdottoID, &r.Quantita, &r.PrezzoUnitario)
		if err != nil {
			return nil, err
		}
		o.Righe = append(o.Righe, r)
	}
	return &o, nil
}

// OrderRepository_GetOrdiniEntrata recupera tutti gli ordini di entrata presenti nel database.
func (r *OrderRepository) OrderRepository_GetOrdiniEntrata(ctx context.Context) ([]domain.OrdineEntrata, error) {
	exec := domain.GetExecutor(ctx, r.DB)
	query := `SELECT id, numero, fornitore, data, note, created_by FROM OrdiniEntrata`
	rows, err := exec.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.OrdineEntrata
	for rows.Next() {
		var o domain.OrdineEntrata
		var ddtID sql.NullInt64
		err := rows.Scan(&o.ID, &o.Numero, &o.Fornitore, &o.Data, &o.Note, &o.CreatedByID)
		if err != nil {
			return nil, err
		}

		if ddtID.Valid {
			val := int(ddtID.Int64)
			o.DDTID = val
		}
		list = append(list, o)
	}
	return list, nil
}

// OrderRepository_DeleteOrdineEntrata rimuove un ordine di entrata dal database in base al suo ID.
func (r *OrderRepository) OrderRepository_DeleteOrdineEntrata(ctx context.Context, o *domain.OrdineEntrata) error {
	exec := domain.GetExecutor(ctx, r.DB)
	query := `DELETE FROM OrdiniEntrata WHERE id = ?`
	_, err := exec.ExecContext(ctx, query, o.ID)
	return err
}

// ─────────────────────────────────────────────
// ORDINI DI USCITA
// ─────────────────────────────────────────────

// OrderRepository_CreateOrdineUscita inserisce un ordine di uscita e le sue righe
// nel database, assegnando gli ID auto-generati all'ordine e alle righe.
func (r *OrderRepository) OrderRepository_CreateOrdineUscita(ctx context.Context, o *domain.OrdineUscita) error {
	exec := domain.GetExecutor(ctx, r.DB)
	query := `INSERT INTO OrdiniUscita (numero, cliente, data, stato, note, created_by, ddt_id)
	          VALUES (?, ?, ?, ?, ?, ?, ?)`
	res, err := exec.ExecContext(ctx, query, o.Numero, o.Cliente, o.Data, string(o.Stato), o.Note, o.CreatedByID, o.DDTID)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		o.ID = int(id)
	}

	for i := range o.Righe {
		o.Righe[i].OrdineUscitaID = o.ID
		queryR := `INSERT INTO RigheOrdineUscita (ordine_uscita_id, prodotto_id, quantita, prezzo_unitario)
		           VALUES (?, ?, ?, ?)`
		resR, err := exec.ExecContext(ctx, queryR, o.Righe[i].OrdineUscitaID, o.Righe[i].ProdottoID, o.Righe[i].Quantita, o.Righe[i].PrezzoUnitario)
		if err != nil {
			return err
		}
		idR, err := resR.LastInsertId()
		if err == nil {
			o.Righe[i].ID = int(idR)
		}
	}
	return nil
}

// OrderRepository_GetOrdineUscitaByID recupera un ordine di uscita e le sue righe
// in base all'ID. Restituisce un errore "ordine non trovato" se l'ordine non esiste.
func (r *OrderRepository) OrderRepository_GetOrdineUscitaByID(ctx context.Context, id int) (*domain.OrdineUscita, error) {
	exec := domain.GetExecutor(ctx, r.DB)
	query := `SELECT id, numero, cliente, data, stato, note, created_by, ddt_id FROM OrdiniUscita WHERE id = ?`
	var o domain.OrdineUscita
	var ddtID sql.NullInt64
	err := exec.QueryRowContext(ctx, query, id).Scan(&o.ID, &o.Numero, &o.Cliente, &o.Data, &o.Stato, &o.Note, &o.CreatedByID, &ddtID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("ordine non trovato")
		}
		return nil, err
	}
	if ddtID.Valid {
		val := int(ddtID.Int64)
		o.DDTID = &val
	}

	// Load lines
	queryR := `SELECT id, ordine_uscita_id, prodotto_id, quantita, prezzo_unitario FROM RigheOrdineUscita WHERE ordine_uscita_id = ?`
	rows, err := exec.QueryContext(ctx, queryR, o.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r domain.RigaOrdineUscita
		err := rows.Scan(&r.ID, &r.OrdineUscitaID, &r.ProdottoID, &r.Quantita, &r.PrezzoUnitario)
		if err != nil {
			return nil, err
		}
		o.Righe = append(o.Righe, r)
	}
	return &o, nil
}

// OrderRepository_GetOrdiniUscita recupera tutti gli ordini di uscita presenti nel database.
func (r *OrderRepository) OrderRepository_GetOrdiniUscita(ctx context.Context) ([]domain.OrdineUscita, error) {
	exec := domain.GetExecutor(ctx, r.DB)
	query := `SELECT id, numero, cliente, data, stato, note, created_by, ddt_id FROM OrdiniUscita`
	rows, err := exec.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.OrdineUscita
	for rows.Next() {
		var o domain.OrdineUscita
		var ddtID sql.NullInt64
		err := rows.Scan(&o.ID, &o.Numero, &o.Cliente, &o.Data, &o.Stato, &o.Note, &o.CreatedByID, &ddtID)
		if err != nil {
			return nil, err
		}
		if ddtID.Valid {
			val := int(ddtID.Int64)
			o.DDTID = &val
		}
		list = append(list, o)
	}
	return list, nil
}

// OrderRepository_GetOrdiniUscitaByStato recupera gli ordini di uscita filtrati per stato.
func (r *OrderRepository) OrderRepository_GetOrdiniUscitaByStato(ctx context.Context, stato domain.StatoOrdineUscita) ([]domain.OrdineUscita, error) {
	exec := domain.GetExecutor(ctx, r.DB)
	query := `SELECT id, numero, cliente, data, stato, note, created_by, ddt_id FROM OrdiniUscita WHERE stato = ?`
	rows, err := exec.QueryContext(ctx, query, string(stato))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.OrdineUscita
	for rows.Next() {
		var o domain.OrdineUscita
		var ddtID sql.NullInt64
		err := rows.Scan(&o.ID, &o.Numero, &o.Cliente, &o.Data, &o.Stato, &o.Note, &o.CreatedByID, &ddtID)
		if err != nil {
			return nil, err
		}
		if ddtID.Valid {
			val := int(ddtID.Int64)
			o.DDTID = &val
		}
		list = append(list, o)
	}
	return list, nil
}

// OrderRepository_UpdateOrdineUscita aggiorna lo stato e l'ID del DDT associato a un ordine di uscita.
func (r *OrderRepository) OrderRepository_UpdateOrdineUscita(ctx context.Context, o *domain.OrdineUscita) error {
	exec := domain.GetExecutor(ctx, r.DB)
	query := `UPDATE OrdiniUscita SET stato = ?, ddt_id = ? WHERE id = ?`
	_, err := exec.ExecContext(ctx, query, string(o.Stato), o.DDTID, o.ID)
	return err
}

// OrderRepository_DeleteOrdineUscita rimuove un ordine di uscita dal database in base al suo ID.
func (r *OrderRepository) OrderRepository_DeleteOrdineUscita(ctx context.Context, o *domain.OrdineUscita) error {
	exec := domain.GetExecutor(ctx, r.DB)
	query := `DELETE FROM OrdiniUscita WHERE id = ?`
	_, err := exec.ExecContext(ctx, query, o.ID)
	return err
}

// ─────────────────────────────────────────────
// DDT
// ─────────────────────────────────────────────

// OrderRepository_CreateDDT inserisce un nuovo DDT nel database, assegnando l'ID
// auto-generato al DDT passato.
func (r *OrderRepository) OrderRepository_CreateDDT(ctx context.Context, d *domain.DDT) error {
	exec := domain.GetExecutor(ctx, r.DB)
	query := `INSERT INTO DDT (numero, ordine_uscita_id, data, mittente, destinatario, path_pdf)
	          VALUES (?, ?, ?, ?, ?, ?)`
	res, err := exec.ExecContext(ctx, query, d.Numero, d.OrdineUscitaID, d.Data, d.Mittente, d.Destinatario, d.PathPDF)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		d.ID = int(id)
	}
	return nil
}

// OrderRepository_GetDDTByYear recupera tutti i DDT il cui numero progressivo
// appartiene all'anno indicato.
func (r *OrderRepository) OrderRepository_GetDDTByYear(ctx context.Context, year int) ([]domain.DDT, error) {
	exec := domain.GetExecutor(ctx, r.DB)

	pattern := fmt.Sprintf("DDT-%d-%%", year)

	rows, err := exec.QueryContext(
		ctx,
		`SELECT id, numero, ordine_uscita_id, data, mittente, destinatario, path_pdf
		 FROM DDT
		 WHERE numero LIKE ?`,
		pattern,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ddtList []domain.DDT

	for rows.Next() {
		var ddt domain.DDT

		err := rows.Scan(
			&ddt.ID,
			&ddt.Numero,
			&ddt.OrdineUscitaID,
			&ddt.Data,
			&ddt.Mittente,
			&ddt.Destinatario,
			&ddt.PathPDF,
		)
		if err != nil {
			return nil, err
		}

		ddtList = append(ddtList, ddt)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ddtList, nil
}

// OrderRepository_GetDDTByOrdineUscita recupera il DDT associato a un ordine di uscita.
// Restituisce (nil, nil) se l'ordine non ha un DDT associato.
func (r *OrderRepository) OrderRepository_GetDDTByOrdineUscita(ctx context.Context, orderID int) (*domain.DDT, error) {
	exec := domain.GetExecutor(ctx, r.DB)
	query := `SELECT id, numero, ordine_uscita_id, data, mittente, destinatario, path_pdf
	          FROM DDT WHERE ordine_uscita_id = ?`
	var d domain.DDT
	err := exec.QueryRowContext(ctx, query, orderID).Scan(
		&d.ID,
		&d.Numero,
		&d.OrdineUscitaID,
		&d.Data,
		&d.Mittente,
		&d.Destinatario,
		&d.PathPDF,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &d, nil
}

// ─────────────────────────────────────────────
// FATTURE
// ─────────────────────────────────────────────

// OrderRepository_CreateFattura inserisce una nuova fattura nel database, assegnando
// l'ID auto-generato alla fattura passata.
func (r *OrderRepository) OrderRepository_CreateFattura(ctx context.Context, f *domain.Fattura) error {
	exec := domain.GetExecutor(ctx, r.DB)
	query := `INSERT INTO Fatture (numero, anno, tipo, ordine_uscita_id, imponibile, iva, totale, data, path_pdf)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	res, err := exec.ExecContext(ctx, query, f.Numero, f.Anno, string(f.Tipo), f.OrdineUscitaID, f.Imponibile, f.IVA, f.Totale, f.Data, f.PathPDF)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		f.ID = int(id)
	}
	return nil
}

// OrderRepository_CountFattureByYearAndType conta le fatture di un determinato anno e tipo,
// utile per calcolare la numerazione progressiva.
func (r *OrderRepository) OrderRepository_CountFattureByYearAndType(ctx context.Context, year int, tipo domain.TipoFattura) (int, error) {
	exec := domain.GetExecutor(ctx, r.DB)
	var count int
	if err := exec.QueryRowContext(ctx, `SELECT COUNT(*) FROM Fatture WHERE anno = ? AND tipo = ?`, year, string(tipo)).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// OrderRepository_GetAllFatture recupera tutte le fatture presenti nel database.
func (r *OrderRepository) OrderRepository_GetAllFatture(ctx context.Context) ([]domain.Fattura, error) {
	exec := domain.GetExecutor(ctx, r.DB)
	query := `SELECT id, numero, anno, tipo, ordine_uscita_id, imponibile, iva, totale, data, path_pdf FROM Fatture`
	rows, err := exec.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Fattura
	for rows.Next() {
		var f domain.Fattura
		var tipoStr string
		err := rows.Scan(&f.ID, &f.Numero, &f.Anno, &tipoStr, &f.OrdineUscitaID, &f.Imponibile, &f.IVA, &f.Totale, &f.Data, &f.PathPDF)
		if err != nil {
			return nil, err
		}
		f.Tipo = domain.TipoFattura(tipoStr)
		list = append(list, f)
	}
	return list, nil
}

// OrderRepository_DeleteFattura rimuove una fattura dal database in base al suo ID.
func (r *OrderRepository) OrderRepository_DeleteFattura(ctx context.Context, id int) error {
	exec := domain.GetExecutor(ctx, r.DB)
	_, err := exec.ExecContext(ctx, `DELETE FROM Fatture WHERE id = ?`, id)
	return err
}

// ─────────────────────────────────────────────
// RESI
// ─────────────────────────────────────────────

// OrderRepository_CreateReso inserisce un nuovo reso nel database, assegnando l'ID
// auto-generato al reso passato.
func (r *OrderRepository) OrderRepository_CreateReso(ctx context.Context, re *domain.Reso) error {
	exec := domain.GetExecutor(ctx, r.DB)
	query := `INSERT INTO Resi (ordine_uscita_id, prodotto_id, quantita, motivo, stato, data, created_by)
	          VALUES (?, ?, ?, ?, ?, ?, ?)`
	res, err := exec.ExecContext(ctx, query, re.OrdineUscitaID, re.ProdottoID, re.Quantita, re.Motivo, re.Stato, re.Data, re.CreatedByID)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		re.ID = int(id)
	}
	return nil
}

// OrderRepository_GetResoByID recupera un reso in base al suo ID e ne restituisce anche
// lo stato corrente come stringa.
func (r *OrderRepository) OrderRepository_GetResoByID(ctx context.Context, id int) (*domain.Reso, string, error) {
	exec := domain.GetExecutor(ctx, r.DB)

	query := `SELECT id, ordine_uscita_id, prodotto_id, quantita, motivo, stato, data, created_by
              FROM Resi WHERE id = ?`

	var re domain.Reso

	err := exec.QueryRowContext(ctx, query, id).Scan(
		&re.ID,
		&re.OrdineUscitaID,
		&re.ProdottoID,
		&re.Quantita,
		&re.Motivo,
		&re.Stato,
		&re.Data,
		&re.CreatedByID,
	)

	if err != nil {
		return nil, "", err
	}

	//fmt.Printf("DEBUG RESO: %+v, stato=%s\n", re, re.Stato)

	return &re, string(re.Stato), nil
}

// OrderRepository_UpdateResoStato aggiorna lo stato di un reso nel database.
func (r *OrderRepository) OrderRepository_UpdateResoStato(ctx context.Context, id int, stato string) error {
	exec := domain.GetExecutor(ctx, r.DB)
	query := `UPDATE Resi SET stato = ? WHERE id = ?`
	_, err := exec.ExecContext(ctx, query, stato, id)
	return err
}

// OrderRepository_GetAllResi recupera tutti i resi presenti nel database.
func (r *OrderRepository) OrderRepository_GetAllResi(ctx context.Context) ([]domain.Reso, error) {
	exec := domain.GetExecutor(ctx, r.DB)

	query := `
		SELECT id, ordine_uscita_id, prodotto_id, quantita,
		       motivo, stato, data, created_by
		FROM Resi
	`

	rows, err := exec.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Reso

	for rows.Next() {
		var reso domain.Reso

		err := rows.Scan(
			&reso.ID,
			&reso.OrdineUscitaID,
			&reso.ProdottoID,
			&reso.Quantita,
			&reso.Motivo,
			&reso.Stato,
			&reso.Data,
			&reso.CreatedByID,
		)
		if err != nil {
			return nil, err
		}

		list = append(list, reso)
	}

	return list, nil
}

// OrderRepository_DeleteReso rimuove un reso dal database in base al suo ID.
func (r *OrderRepository) OrderRepository_DeleteReso(ctx context.Context, id int) error {
	exec := domain.GetExecutor(ctx, r.DB)
	query := `DELETE FROM Resi WHERE id = ?`
	_, err := exec.ExecContext(ctx, query, id)
	return err
}
