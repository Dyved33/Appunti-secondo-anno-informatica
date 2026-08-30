// Package ordini gestisce la schermata TUI relativa alla visualizzazione,
// consultazione e creazione degli ordini di entrata, ordini di uscita e resi.
package ordini

import (
	"charm.land/bubbles/v2/textinput"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
)

// sezioneAttiva rappresenta la scheda/tab corrente selezionata nell'interfaccia.
type sezioneAttiva int

const (
	// ORDINI_ENTRATA indica la sezione degli ordini d'acquisto/fornitore.
	ORDINI_ENTRATA sezioneAttiva = iota
	// ORDINI_USCITA indica la sezione degli ordini di vendita/cliente.
	ORDINI_USCITA
	// RESI indica la sezione della gestione dei resi merci.
	RESI
)

// focusForm identifica il campo di testo del form di creazione su cui risiede il cursore.
type focusForm int

const (
	// FocusControparte seleziona il campo del nome Fornitore/Cliente.
	FocusControparte focusForm = iota
	// FocusNumero seleziona il campo del numero identificativo dell'ordine.
	FocusNumero
	// FocusNote seleziona il campo del testo note opzionale.
	FocusNote
	// FocusRighe seleziona il campo di inserimento degli articoli dell'ordine.
	FocusRighe

	// FocusResoOrdineID seleziona il campo dell'ID dell'ordine di riferimento per il reso.
	FocusResoOrdineID focusForm = iota
	// FocuResoProdottoID seleziona il campo dell'ID del prodotto da rendere.
	FocusResoProdottoID
	// FocusResoQuantita seleziona il campo della quantità del prodotto da rendere.
	FocusResoQuantita
	// FocusResoMotivo seleziona il campo del motivo del reso.
	FocusResoMotivo
)

// Model memorizza lo stato completo del modulo Ordini per l'architettura Bubble Tea.
type Model struct {
	// Client è il client HTTP/API per la comunicazione con il server di backend.
	Client *api.APIClient

	sezione sezioneAttiva

	// Dati caricati dal backend
	ordiniEntrata []api.OrdineEntrata
	ordiniUscita  []api.OrdineUscita
	resi          []api.Reso

	// Cursori e dettagli espansi per ciascuna tabella
	cursoreEntrata int
	cursoreUscita  int
	cursoreResi    int
	offsetEntrata  int
	offsetUscita   int
	offsetResi     int
	espansiEntrata map[int]bool // ID ordine -> true se il dettaglio righe è espanso
	espansiUscita  map[int]bool // ID ordine -> true se il dettaglio righe è espanso
	espansiResi    map[int]bool // ID reso -> true se il dettaglio è espanso

	winW int
	winH int

	// Stato del form di creazione
	creazioneEntrata bool
	creazioneUscita  bool
	creazioneReso    bool

	// Campi di input del form di creazione ordine
	formControparte textinput.Model // Fornitore o Cliente
	formNumero      textinput.Model // Codice/Numero d'ordine
	formNote        textinput.Model // Note accessorie
	formRighe       textinput.Model // Formato: ID:Qtà:Prezzo

	// Campi di input del form di creazione reso
	formResoOrdineID   textinput.Model // ID dell'ordine di riferimento
	formResoProdottoID textinput.Model // ID del prodotto da rendere
	formResoQuantita   textinput.Model // Quantità del prodotto da rendere
	formResoMotivo     textinput.Model // Motivo del reso

	focusCampo focusForm

	// campi per l'eliminazione di un ordine o reso
	confermaEliminazione bool
	idDaEliminare        int

	loading bool
	err     error
}

// rigaAdapter uniforma la struttura delle righe ordine per la resa grafica tabellare.
type rigaAdapter struct {
	ProdottoID     int
	Prodotto       *api.Prodotto
	Quantita       int
	PrezzoUnitario float64
}

// OrdiniLoadedMsg viene emesso in modo asincrono al termine dell'interrogazione del backend
// e trasporta i dati aggiornati di ordini di entrata, uscita e resi.
type OrdiniLoadedMsg struct {
	Entrata []api.OrdineEntrata
	Uscita  []api.OrdineUscita
	Resi    []api.Reso
}

// ErroreMsg notifica un errore riscontrato durante le chiamate di I/O o di validazione.
type ErroreMsg struct {
	Err error
}

// OrdineEntrataCreatoMsg segnala il completamento con successo della creazione di un ordine di entrata.
type OrdineEntrataCreatoMsg struct{}

// OrdineUscitaCreatoMsg segnala il completamento con successo della creazione di un ordine di uscita.
type OrdineUscitaCreatoMsg struct{}

// ResoCreatoMsg segnala il completamento con successo della creazione di un reso.
type ResoCreatoMsg struct{}

// OrdineEntrataEliminatoMsg segnala la corretta eliminazione dell'ordine di entrata.
type OrdineEntrataEliminatoMsg struct{}

// OrdineUscitaEliminatoMsg segnala la corretta eliminazione dell'ordine di uscita.
type OrdineUscitaEliminatoMsg struct{}

// ResoEliminatoMsg segnala la corretta eliminazione del reso.
type ResoEliminatoMsg struct{}

// New inizializza e restituisce un nuovo Model pronto all'uso con i valori di default.
func New(client *api.APIClient) Model {
	return Model{
		Client:  client,
		sezione: ORDINI_ENTRATA,

		espansiEntrata: make(map[int]bool),
		espansiUscita:  make(map[int]bool),
		espansiResi:    make(map[int]bool),
		loading:        true,

		creazioneEntrata: false,
		creazioneUscita:  false,
		formRighe:        textinput.New(),
	}
}
