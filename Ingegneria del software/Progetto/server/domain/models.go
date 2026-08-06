package domain

import "time"

// ─────────────────────────────────────────────
// UTENTE — RF-01, RF-02, RF-03
// ─────────────────────────────────────────────

// Utente rappresenta un utente del sistema con ruolo RBAC.
type Utente struct {
	ID                int       `db:"id"`
	Username          string    `db:"username"`
	PasswordHash      string    `db:"password_hash"`        // bcrypt hash (x/crypto/bcrypt)
	Ruolo             Ruolo     `db:"ruolo"`                // Magazziniere | Contabile | Amministratore
	DomandaSicurezza  string    `db:"domanda_sicurezza"`    // RF-03
	RispostaHash      string    `db:"risposta_hash"`        // bcrypt hash della risposta
	TentativiFalliti  int       `db:"tentativi_falliti"`    // RF-02: blocco a 3
	Bloccato          bool      `db:"bloccato"`             // RF-02
	CreatedAt         time.Time `db:"created_at"`
}

// ─────────────────────────────────────────────
// CATALOGO — RF-04
// ─────────────────────────────────────────────

// Categoria raggruppa i prodotti per tipo merceologico.
type Categoria struct {
	ID          int    `db:"id"`
	Nome        string `db:"nome"`
	Descrizione string `db:"descrizione"`
}

// Prodotto rappresenta l'anagrafica di un articolo a magazzino.
type Prodotto struct {
	ID          int     `db:"id"`
	Codice      string  `db:"codice"`       // codice interno
	Barcode     string  `db:"barcode"`      // EAN-13 — RF-15, US-02
	Nome        string  `db:"nome"`
	Descrizione string  `db:"descrizione"`
	CategoriaID int     `db:"categoria_id"`
	PrezzoAcq   float64 `db:"prezzo_acquisto"` // RF-14
	PrezzoVend  float64 `db:"prezzo_vendita"`  // RF-14
	SogliaMin   int     `db:"soglia_minima"`   // RF-10: notifica sotto soglia
	Categoria   *Categoria // join opzionale, non mappato su DB
}

// ─────────────────────────────────────────────
// GIACENZA — RF-05, RF-06, RF-07, RF-10
// ─────────────────────────────────────────────

// Giacenza traccia la quantità fisica di un prodotto a magazzino.
type Giacenza struct {
	ID         int       `db:"id"`
	ProdottoID int       `db:"prodotto_id"`
	Quantita   int       `db:"quantita"`
	UpdatedAt  time.Time `db:"updated_at"`
	Prodotto   *Prodotto // join opzionale
}

// SottoSoglia restituisce true se la giacenza è sotto la soglia minima (RF-10).
func (g *Giacenza) SottoSoglia() bool {
	if g.Prodotto == nil {
		return false
	}
	return g.Quantita < g.Prodotto.SogliaMin
}

// ─────────────────────────────────────────────
// ORDINI DI ENTRATA — RF-06
// ─────────────────────────────────────────────

// OrdineEntrata rappresenta un carico di merce in arrivo al magazzino.
type OrdineEntrata struct {
	ID          int                `db:"id"`
	Numero      string             `db:"numero"`      // es. "ENT-2026-0042"
	Fornitore   string             `db:"fornitore"`
	Data        time.Time          `db:"data"`
	Stato       StatoOrdine        `db:"stato"`
	Note        string             `db:"note"`
	CreatedByID int                `db:"created_by"`
	Righe       []RigaOrdineEntrata // non mappata su DB, caricata dal DAO
}

// RigaOrdineEntrata è una singola riga (prodotto + quantità) di un OrdineEntrata.
type RigaOrdineEntrata struct {
	ID              int     `db:"id"`
	OrdineEntrataID int     `db:"ordine_entrata_id"`
	ProdottoID      int     `db:"prodotto_id"`
	Quantita        int     `db:"quantita"`
	PrezzoUnitario  float64 `db:"prezzo_unitario"`
	Prodotto        *Prodotto
}

// ─────────────────────────────────────────────
// ORDINI DI USCITA — RF-07, RF-08, RF-09
// ─────────────────────────────────────────────

// OrdineUscita rappresenta una spedizione o consegna di merce dal magazzino.
type OrdineUscita struct {
	ID          int               `db:"id"`
	Numero      string            `db:"numero"`    // es. "USC-2026-0017"
	Cliente     string            `db:"cliente"`
	Data        time.Time         `db:"data"`
	Stato       StatoOrdine       `db:"stato"`     // RF-09: Confermato → Evaso
	Note        string            `db:"note"`
	CreatedByID int               `db:"created_by"`
	DDTID       *int              `db:"ddt_id"`    // FK al DDT generato (RF-08)
	Righe       []RigaOrdineUscita
}

// RigaOrdineUscita è una singola riga (prodotto + quantità) di un OrdineUscita.
type RigaOrdineUscita struct {
	ID             int     `db:"id"`
	OrdineUscitaID int     `db:"ordine_uscita_id"`
	ProdottoID     int     `db:"prodotto_id"`
	Quantita       int     `db:"quantita"`     // TC-01: validato prima del ROLLBACK
	PrezzoUnitario float64 `db:"prezzo_unitario"`
	Prodotto       *Prodotto
}

// ─────────────────────────────────────────────
// DDT — RF-08, RF-09
// ─────────────────────────────────────────────

// DDT (Documento di Trasporto) generato alla conferma di un OrdineUscita.
type DDT struct {
	ID             int       `db:"id"`
	Numero         string    `db:"numero"`       // es. "DDT-2026-0017"
	OrdineUscitaID int       `db:"ordine_uscita_id"`
	Data           time.Time `db:"data"`
	Mittente       string    `db:"mittente"`
	Destinatario   string    `db:"destinatario"`
	PathPDF        string    `db:"path_pdf"`     // percorso file PDF generato
	OrdineUscita   *OrdineUscita
}

// ─────────────────────────────────────────────
// FATTURE E NOTE DI CREDITO — RF-13
// ─────────────────────────────────────────────

// Fattura rappresenta un documento fiscale emesso per un ordine di uscita evaso.
type Fattura struct {
	ID             int         `db:"id"`
	Numero         string      `db:"numero"`          // numerazione progressiva annuale — RF-13
	Anno           int         `db:"anno"`
	Tipo           TipoFattura `db:"tipo"`            // Fattura | NotaCredito
	OrdineUscitaID int         `db:"ordine_uscita_id"`
	Imponibile     float64     `db:"imponibile"`
	IVA            float64     `db:"iva"`
	Totale         float64     `db:"totale"`
	Data           time.Time   `db:"data"`
	PathPDF        string      `db:"path_pdf"`
}

// ─────────────────────────────────────────────
// RESI — RF-11
// ─────────────────────────────────────────────

// Reso rappresenta la restituzione di merce con ripristino della giacenza.
type Reso struct {
	ID             int       `db:"id"`
	OrdineUscitaID int       `db:"ordine_uscita_id"`
	ProdottoID     int       `db:"prodotto_id"`
	Quantita       int       `db:"quantita"`
	Motivo         string    `db:"motivo"`
	Data           time.Time `db:"data"`
	CreatedByID    int       `db:"created_by"`
}

// ─────────────────────────────────────────────
// SESSIONE — Stateful JWT / Whitelist
// ─────────────────────────────────────────────

type Sessione struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	TokenHash string    `db:"token_hash"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}