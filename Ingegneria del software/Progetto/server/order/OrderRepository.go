package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/tommasogarofani/ingsw_25-26_group-11/server/domain"
)

type OrderRepository struct {
	DB *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

// ─────────────────────────────────────────────
// ORDINI DI ENTRATA
// ─────────────────────────────────────────────

func (r *OrderRepository) OrderRepository_CreateOrdineEntrata(ctx context.Context, o *domain.OrdineEntrata) error {
	exec := domain.GetExecutor(ctx, r.DB)
	query := `INSERT INTO OrdiniEntrata (numero, fornitore, data, stato, note, created_by)
	          VALUES (?, ?, ?, ?, ?, ?)`
	res, err := exec.ExecContext(ctx, query, o.Numero, o.Fornitore, o.Data, string(o.Stato), o.Note, o.CreatedByID)
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

func (r *OrderRepository) OrderRepository_GetOrdineEntrataByID(ctx context.Context, id int) (*domain.OrdineEntrata, error) {
	exec := domain.GetExecutor(ctx, r.DB)
	query := `SELECT id, numero, fornitore, data, stato, note, created_by FROM OrdiniEntrata WHERE id = ?`
	var o domain.OrdineEntrata
	var statoStr string
	err := exec.QueryRowContext(ctx, query, id).Scan(&o.ID, &o.Numero, &o.Fornitore, &o.Data, &statoStr, &o.Note, &o.CreatedByID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("ordine non trovato")
		}
		return nil, err
	}
	o.Stato = domain.StatoOrdine(statoStr)

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

// ─────────────────────────────────────────────
// ORDINI DI USCITA
// ─────────────────────────────────────────────

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

func (r *OrderRepository) OrderRepository_GetOrdineUscitaByID(ctx context.Context, id int) (*domain.OrdineUscita, error) {
	exec := domain.GetExecutor(ctx, r.DB)
	query := `SELECT id, numero, cliente, data, stato, note, created_by, ddt_id FROM OrdiniUscita WHERE id = ?`
	var o domain.OrdineUscita
	var statoStr string
	var ddtID sql.NullInt64
	err := exec.QueryRowContext(ctx, query, id).Scan(&o.ID, &o.Numero, &o.Cliente, &o.Data, &statoStr, &o.Note, &o.CreatedByID, &ddtID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("ordine non trovato")
		}
		return nil, err
	}
	o.Stato = domain.StatoOrdine(statoStr)
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
		var statoStr string
		var ddtID sql.NullInt64
		err := rows.Scan(&o.ID, &o.Numero, &o.Cliente, &o.Data, &statoStr, &o.Note, &o.CreatedByID, &ddtID)
		if err != nil {
			return nil, err
		}
		o.Stato = domain.StatoOrdine(statoStr)
		if ddtID.Valid {
			val := int(ddtID.Int64)
			o.DDTID = &val
		}
		list = append(list, o)
	}
	return list, nil
}

func (r *OrderRepository) OrderRepository_GetOrdiniUscitaByStato(ctx context.Context, stato domain.StatoOrdine) ([]domain.OrdineUscita, error) {
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
		var statoStr string
		var ddtID sql.NullInt64
		err := rows.Scan(&o.ID, &o.Numero, &o.Cliente, &o.Data, &statoStr, &o.Note, &o.CreatedByID, &ddtID)
		if err != nil {
			return nil, err
		}
		o.Stato = domain.StatoOrdine(statoStr)
		if ddtID.Valid {
			val := int(ddtID.Int64)
			o.DDTID = &val
		}
		list = append(list, o)
	}
	return list, nil
}

func (r *OrderRepository) OrderRepository_UpdateOrdineUscita(ctx context.Context, o *domain.OrdineUscita) error {
	exec := domain.GetExecutor(ctx, r.DB)
	query := `UPDATE OrdiniUscita SET stato = ?, ddt_id = ? WHERE id = ?`
	_, err := exec.ExecContext(ctx, query, string(o.Stato), o.DDTID, o.ID)
	return err
}

// ─────────────────────────────────────────────
// DDT
// ─────────────────────────────────────────────

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

func (r *OrderRepository) OrderRepository_CountDDTByYear(ctx context.Context, year int) (int, error) {
	exec := domain.GetExecutor(ctx, r.DB)
	pattern := fmt.Sprintf("DDT-%d-%%", year)
	var count int
	if err := exec.QueryRowContext(ctx, `SELECT COUNT(*) FROM DDT WHERE numero LIKE ?`, pattern).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *OrderRepository) OrderRepository_GetDDTByOrdineUscita(ctx context.Context, orderID int) (*domain.DDT, error) {
	exec := domain.GetExecutor(ctx, r.DB)
	query := `SELECT id, numero, ordine_uscita_id, data, mittente, destinatario, path_pdf FROM DDT WHERE ordine_uscita_id = ?`
	var d domain.DDT
	err := exec.QueryRowContext(ctx, query, orderID).Scan(&d.ID, &d.Numero, &d.OrdineUscitaID, &d.Data, &d.Mittente, &d.Destinatario, &d.PathPDF)
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

func (r *OrderRepository) OrderRepository_CountFattureByYearAndType(ctx context.Context, year int, tipo domain.TipoFattura) (int, error) {
	exec := domain.GetExecutor(ctx, r.DB)
	var count int
	if err := exec.QueryRowContext(ctx, `SELECT COUNT(*) FROM Fatture WHERE anno = ? AND tipo = ?`, year, string(tipo)).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

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

func (r *OrderRepository) OrderRepository_DeleteFattura(ctx context.Context, id int) error {
	exec := domain.GetExecutor(ctx, r.DB)
	_, err := exec.ExecContext(ctx, `DELETE FROM Fatture WHERE id = ?`, id)
	return err
}

// ─────────────────────────────────────────────
// RESI
// ─────────────────────────────────────────────

func (r *OrderRepository) OrderRepository_CreateReso(ctx context.Context, re *domain.Reso, stato string) error {
	exec := domain.GetExecutor(ctx, r.DB)
	query := `INSERT INTO Resi (ordine_uscita_id, prodotto_id, quantita, motivo, stato, data, created_by)
	          VALUES (?, ?, ?, ?, ?, ?, ?)`
	res, err := exec.ExecContext(ctx, query, re.OrdineUscitaID, re.ProdottoID, re.Quantita, re.Motivo, stato, re.Data, re.CreatedByID)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		re.ID = int(id)
	}
	return nil
}

func (r *OrderRepository) OrderRepository_GetResoByID(ctx context.Context, id int) (*domain.Reso, string, error) {
	exec := domain.GetExecutor(ctx, r.DB)
	query := `SELECT id, ordine_uscita_id, prodotto_id, quantita, motivo, stato, data, created_by FROM Resi WHERE id = ?`
	var re domain.Reso
	var stato string
	err := exec.QueryRowContext(ctx, query, id).Scan(&re.ID, &re.OrdineUscitaID, &re.ProdottoID, &re.Quantita, &re.Motivo, &stato, &re.Data, &re.CreatedByID)
	if err != nil {
		return nil, "", err
	}
	return &re, stato, nil
}

func (r *OrderRepository) OrderRepository_UpdateResoStato(ctx context.Context, id int, stato string) error {
	exec := domain.GetExecutor(ctx, r.DB)
	query := `UPDATE Resi SET stato = ? WHERE id = ?`
	_, err := exec.ExecContext(ctx, query, stato, id)
	return err
}

func (r *OrderRepository) OrderRepository_GetAllResi(ctx context.Context) ([]map[string]interface{}, error) {
	exec := domain.GetExecutor(ctx, r.DB)
	query := `SELECT id, ordine_uscita_id, prodotto_id, quantita, motivo, stato, data, created_by FROM Resi`
	rows, err := exec.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []map[string]interface{}
	for rows.Next() {
		var id, outID, prodID, qty, createdBy int
		var motivo, stato, data string
		err := rows.Scan(&id, &outID, &prodID, &qty, &motivo, &stato, &data, &createdBy)
		if err != nil {
			return nil, err
		}
		item := map[string]interface{}{
			"id":        id,
			"order_id":  outID,
			"prod_id":   prodID,
			"qty":       qty,
			"motivo":    motivo,
			"stato":     stato,
			"data":      data,
			"created":   createdBy,
		}
		list = append(list, item)
	}
	return list, nil
}
