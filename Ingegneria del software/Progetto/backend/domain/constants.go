package domain

// ─────────────────────────────────────────────
// RUOLI — RF-01 (RBAC)
// ─────────────────────────────────────────────

// Ruolo definisce il profilo di accesso di un Utente nel sistema.
type Ruolo string

const (
	RuoloMagazziniere  Ruolo = "magazziniere"  // carico/scarico, DDT, giacenze
	RuoloContabile     Ruolo = "contabile"      // ordini, fatture, listini
	RuoloAmministratore Ruolo = "amministratore" // supervisione globale, reportistica
)

// AllRuoli elenca tutti i ruoli validi — utile per validazione e menu di selezione.
var AllRuoli = []Ruolo{RuoloMagazziniere, RuoloContabile, RuoloAmministratore}

// IsValid verifica che il ruolo sia tra quelli riconosciuti dal sistema.
func (r Ruolo) IsValid() bool {
	for _, v := range AllRuoli {
		if r == v {
			return true
		}
	}
	return false
}

// Label restituisce il nome leggibile del ruolo per la TUI.
func (r Ruolo) Label() string {
	switch r {
	case RuoloMagazziniere:
		return "Magazziniere"
	case RuoloContabile:
		return "Contabile"
	case RuoloAmministratore:
		return "Amministratore"
	default:
		return "Sconosciuto"
	}
}

// ─────────────────────────────────────────────
// STATI ORDINE — RF-09
// ─────────────────────────────────────────────

// StatoOrdine rappresenta il ciclo di vita di un ordine.
type StatoOrdine string

const (
	// StatoConfermato: ordine registrato ma non ancora evaso.
	StatoConfermato StatoOrdine = "confermato"
	// StatoEvaso: ordine completato, DDT generato (RF-09).
	StatoEvaso StatoOrdine = "evaso"
	// StatoAnnullato: ordine cancellato prima dell'evasione.
	StatoAnnullato StatoOrdine = "annullato"
)

// ─────────────────────────────────────────────
// TIPO FATTURA — RF-13
// ─────────────────────────────────────────────

// TipoFattura distingue i documenti fiscali emessi dal contabile.
type TipoFattura string

const (
	TipoFatturaVendita    TipoFattura = "fattura"       // documento fiscale standard
	TipoNotaCredito       TipoFattura = "nota_credito"  // storno/reso — RF-11, RF-13
)

// ─────────────────────────────────────────────
// SOGLIE E LIMITI — RF-02, RF-10, RNF-01, RNF-02
// ─────────────────────────────────────────────

const (
	// MaxTentativiLogin è il numero massimo di tentativi falliti prima
	// del blocco automatico dell'account (RF-02).
	MaxTentativiLogin = 3

	// SogliaRispostaMsBarcode è il tempo massimo di risposta in millisecondi
	// accettabile per il 95% delle scansioni barcode (RNF-01: < 0.5s).
	SogliaRispostaMsBarcode = 500

	// MaxSecondiReportInventario è il tempo massimo CPU in secondi
	// per la generazione del report annuale dell'inventario (RNF-02).
	MaxSecondiReportInventario = 3
)

// ─────────────────────────────────────────────
// PREFISSI NUMERAZIONE DOCUMENTI — RF-08, RF-13
// ─────────────────────────────────────────────

// I prefissi definiscono il formato della numerazione progressiva annuale.
// Esempio: "DDT-2026-0001", "USC-2026-0042", "FAT-2026-0007"
const (
	PrefissoOrdineEntrata = "ENT"
	PrefissoOrdineUscita  = "USC"
	PrefissoDDT           = "DDT"
	PrefissoFattura       = "FAT"
	PrefissoNotaCredito   = "NC"
)

// ─────────────────────────────────────────────
// PERCORSI FILE — RNF-04, RNF-05
// ─────────────────────────────────────────────

const (
	// DBPath è il percorso del file SQLite locale (RNF-05: single-file embedded DB).
	// Tutto rimane in locale, nessun dato esposto all'esterno (RNF-04).
	DBPath = "./data/magazzino.db"

	// OutputPDFDir è la cartella in cui vengono salvati DDT e fatture generati.
	OutputPDFDir = "./output/pdf"

	// OutputCSVDir è la cartella per i report CSV dell'inventario.
	OutputCSVDir = "./output/csv"
)
