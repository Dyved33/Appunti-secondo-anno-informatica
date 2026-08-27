package api

import (
	"net/http"
	"time"
)

/* TIPI E COSTANTI */

// StatoOrdine rappresenta lo stato di un ordine (confermato, evaso, annullato).
type StatoOrdine string

// OrdineUscita rappresenta un ordine di uscita di merce dal magazzino.
const (
	// StatoConfermato: ordine registrato ma non ancora evaso.
	StatoConfermato StatoOrdine = "confermato"
	// StatoEvaso: ordine completato, DDT generato (RF-09).
	StatoEvaso StatoOrdine = "evaso"
	// StatoAnnullato: ordine cancellato prima dell'evasione.
	StatoAnnullato StatoOrdine = "annullato"
)

type StatoReso string

const (
	// StatoInAttesa: reso registrato ma non ancora accettato o rifiutato.
	StatoInAttesa StatoReso = "in attesa"
	// StatoAccettato: reso accettato, merce restituita e giacenza aggiornata.
	StatoAccettato StatoReso = "accettato"
	// StatoRifiutato: reso rifiutato, merce non restituita.
	StatoRifiutato StatoReso = "rifiutato"
)

/* APICLIENT */

// APIClient è il client per comunicare con il server API.
type APIClient struct {
	indirizzoServer string
	token           string
	role            string
	httpClient      *http.Client
}

// funzione per inizializzare il client
// prende una stringa che contien l'url del server
// restituisce il puntatore al client
func New(url string) *APIClient {
	return &APIClient{
		indirizzoServer: url,
		httpClient:      &http.Client{},
	}
}

/* DTO SERVER */

// OrderOutputDTO rappresenta la struttura dei dati di un prodotto restituito dal server.
type serverProductDTO struct {
	ID          int     `json:"id"`
	Codice      string  `json:"codice"`
	Barcode     string  `json:"barcode"`
	Nome        string  `json:"nome"`
	Descrizione string  `json:"descrizione"`
	CategoriaID int     `json:"categoria_id"`
	PrezzoAcq   float64 `json:"prezzo_acquisto"`
	PrezzoVend  float64 `json:"prezzo_vendita"`
	SogliaMin   int     `json:"soglia_minima"`
}

// serverCategoryDTO rappresenta una categoria di prodotti nella risposta del server.
type serverCategoryDTO struct {
	ID          int    `json:"id"`
	Nome        string `json:"nome"`
	Descrizione string `json:"descrizione"`
}

// OrderOutputDTO rappresenta la struttura dei dati di una fattura restituita dal server.
type serverFatturaDTO struct {
	ID             int       `json:"id"`
	Numero         string    `json:"numero"`
	Tipo           string    `json:"tipo"`
	OrdineUscitaID int       `json:"ordine_uscita_id"`
	Imponibile     float64   `json:"imponibile"`
	IVA            float64   `json:"iva"`
	Totale         float64   `json:"totale"`
	Data           time.Time `json:"data"`
	PathPDF        string    `json:"path_pdf"`
}

// OrderLineOutputDTO rappresenta una singola riga di ordine nella risposta.
type serverOrderLineDTO struct {
	ID             int     `json:"id"`
	ProdottoID     int     `json:"prodotto_id"`
	Quantita       int     `json:"quantita"`
	PrezzoUnitario float64 `json:"prezzo_unitario"`
}

// OrderOutputDTO rappresenta i dati di un ordine restituito dal server.
type serverOrderDTO struct {
	ID          int                  `json:"id"`
	Numero      string               `json:"numero"`
	Cliente     string               `json:"cliente,omitempty"`
	Fornitore   string               `json:"fornitore,omitempty"`
	Data        time.Time            `json:"data"`
	Stato       string               `json:"stato"`
	Note        string               `json:"note"`
	CreatedByID int                  `json:"created_by_id"`
	DDTID       *int                 `json:"ddt_id,omitempty"`
	Righe       []serverOrderLineDTO `json:"righe,omitempty"`
}

// ResoOutputDTO rappresenta i dati di un reso restituito dal server.
type serverResoDTO struct {
	ID             int       `json:"id"`
	OrdineUscitaID int       `json:"ordine_uscita_id"`
	ProdottoID     int       `json:"prodotto_id"`
	Quantita       int       `json:"quantita"`
	Motivo         string    `json:"motivo"`
	Stato          string    `json:"stato"`
	CreatedByID    int       `json:"created_by_id"`
	Data           time.Time `json:"data"`
}

// serverUserDTO rappresenta i dati utente restituiti dagli endpoint admin users.
type serverUserDTO struct {
	ID               int       `json:"id"`
	Username         string    `json:"username"`
	Ruolo            string    `json:"ruolo"`
	TentativiFalliti int       `json:"tentativi_falliti"`
	Bloccato         bool      `json:"bloccato"`
	CreatedAt        time.Time `json:"created_at"`
}

// serverDDTDTO rappresenta i dati di un DDT (Documento di Trasporto) restituiti dal server.
type serverDDTDTO struct {
	ID             int       `json:"id"`
	Numero         string    `json:"numero"`
	OrdineUscitaID int       `json:"ordine_uscita_id"`
	Data           time.Time `json:"data"`
	Mittente       string    `json:"mittente"`
	Destinatario   string    `json:"destinatario"`
	PathPDF        string    `json:"path_pdf"`
}

/* DTO CLIENT */

// ClientProductDTO rappresenta i dati per la creazione di un prodotto.
type ClientProductDTO struct {
	Codice      string  `json:"codice"`
	Barcode     string  `json:"barcode"`
	Nome        string  `json:"nome"`
	Descrizione string  `json:"descrizione"`
	CategoriaID int     `json:"categoria_id"`
	PrezzoAcq   float64 `json:"prezzo_acquisto"`
	PrezzoVend  float64 `json:"prezzo_vendita"`
	SogliaMin   int     `json:"soglia_minima"`
}

// ClientOrderLineDTO rappresenta una singola riga di ordine nella risposta.
type ClientOrderLineDTO struct {
	ProdottoID     int     `json:"prodotto_id"`
	Quantita       int     `json:"quantita"`
	PrezzoUnitario float64 `json:"prezzo_unitario"`
}

// ClientOrderEntryDTO rappresenta la richiesta di creazione di un ordine di entrata.
// Contiene i dati dell'ordine e l'elenco delle righe (prodotti con quantità e prezzi).
type ClientOrderDTO struct {
	Numero      string               `json:"numero"`
	Cliente     string               `json:"cliente,omitempty"`
	Fornitore   string               `json:"fornitore,omitempty"`
	Stato       string               `json:"stato,omitempty"`
	Note        string               `json:"note"`
	CreatedByID int                  `json:"created_by_id,omitempty"`
	Righe       []ClientOrderLineDTO `json:"righe"`
}

// ClientDDTDTO rappresenta i dati per la generazione di un DDT (Documento di Trasporto).
type ClientDDTDTO struct {
	Mittente     string `json:"mittente"`
	Destinatario string `json:"destinatario"`
}

// ClientResoDTO rappresenta i dati per la creazione di un reso.
type ClientResoDTO struct {
	ID             int       `json:"id"`
	OrdineUscitaID int       `json:"ordine_uscita_id"`
	ProdottoID     int       `json:"prodotto_id"`
	Quantita       int       `json:"quantita"`
	Motivo         string    `json:"motivo"`
	Stato          string    `json:"stato"`
	CreatedByID    int       `json:"created_by_id"`
	Data           time.Time `json:"data"`
}

/* STRUTTURE  CLIENT */

// Prodotto rappresenta un prodotto nel catalogo.
type Prodotto struct {
	ID             int     `json:"id"`
	Codice         string  `json:"codice"`
	Nome           string  `json:"nome"`
	Categoria      string  `json:"categoria"`
	Prezzo         float64 `json:"prezzo"`
	CategoriaID    int     `json:"categoria_id"`
	PrezzoAcquisto float64 `json:"prezzo_acquisto"`
	PrezzoVendita  float64 `json:"prezzo_vendita"`
	Attivo         bool    `json:"attivo"`
}

// Fattura rappresenta una fattura emessa al cliente.
type Fattura struct {
	ID             int       `json:"id"`
	OrdineUscitaID int       `json:"ordine_uscita_id"`
	Numero         string    `json:"numero"`
	Anno           int       `json:"anno"`
	Cliente        string    `json:"cliente"`
	Importo        float64   `json:"importo"`
	NotaCredito    bool      `json:"notaCredito"`
	Data           time.Time `json:"data"`
}

// Ordine rappresenta un ordine di uscita o entrata di merce.
type Ordine struct {
	ID             int     `json:"id"`
	Cliente        string  `json:"cliente"`
	ProdottoID     int     `json:"prodotto_id"`
	Prodotto       string  `json:"prodotto"`
	Quantita       int     `json:"quantita"`
	PrezzoUnitario float64 `json:"prezzoUnitario"`
	Stato          string  `json:"stato"`
}

// KPI rappresenta un indicatore chiave di prestazione (Key Performance Indicator) del magazzino.
type KPI struct {
	Etichetta string `json:"etichetta"`
	Valore    string `json:"valore"`
}

// InventoryItem rappresenta un singolo articolo nel magazzino.
type InventoryItem struct {
	ProdottoID  int     `json:"prodotto_id"`
	Codice      string  `json:"codice"`
	Nome        string  `json:"nome"`
	Categoria   string  `json:"categoria"`
	Quantita    int     `json:"quantita"`
	SogliaMin   int     `json:"soglia_minima"`
	PrezzoVend  float64 `json:"prezzo_vendita"`
	ValoreRiga  float64 `json:"valore_riga"`
	SottoSoglia bool    `json:"sotto_soglia"`
}

// InventoryReport rappresenta il report dell'inventario del magazzino.
type InventoryReport struct {
	GeneratoIl          time.Time       `json:"generato_il"`
	TotaleQuantita      int             `json:"totale_quantita"`
	ValoreTotale        float64         `json:"valore_totale"`
	Articoli            []InventoryItem `json:"articoli"`
	ArticoliSottoSoglia []InventoryItem `json:"articoli_sotto_soglia"`
}

// OrdineEntrata rappresenta un ordine di entrata di merce nel magazzino.
type OrdineEntrata struct {
	ID          int
	Numero      string
	Fornitore   string
	Data        time.Time
	Note        string
	CreatedByID int
	Righe       []RigaOrdineEntrata // non mappata su DB, caricata dal DAO
}

// RigaOrdineEntrata rappresenta una singola riga di un ordine di entrata.
type RigaOrdineEntrata struct {
	ID              int
	OrdineEntrataID int
	ProdottoID      int
	Quantita        int
	PrezzoUnitario  float64
	Prodotto        *Prodotto
}

// rappresenta una spedizione o consegna di merce dal magazzino.
type OrdineUscita struct {
	ID          int
	Numero      string
	Mittente    string
	Cliente     string
	Data        time.Time
	Stato       StatoOrdine
	Note        string
	CreatedByID int
	DDTID       *int
	Righe       []RigaOrdineUscita
}

// RigaOrdineUscita rappresenta una singola riga di un ordine di uscita.
type RigaOrdineUscita struct {
	ID             int
	OrdineUscitaID int
	ProdottoID     int
	Quantita       int
	PrezzoUnitario float64
	Prodotto       *Prodotto
}

// Reso rappresenta la restituzione di merce con ripristino della giacenza.
type Reso struct {
	ID             int
	OrdineUscitaID int
	ProdottoID     int
	Quantita       int
	Motivo         string
	Stato          StatoReso
	Data           time.Time
	CreatedByID    int
}

// User rappresenta un utente amministrabile dalla pagina admin users.
type User struct {
	ID               int
	Username         string
	Ruolo            string
	TentativiFalliti int
	Bloccato         bool
	CreatedAt        time.Time
}

// DDT rappresenta i dati di un DDT (Documento di Trasporto).
type DDT struct {
	ID             int
	Numero         string
	OrdineUscitaID int
	Data           time.Time
	Mittente       string
	Destinatario   string
	PathPDF        string
}
