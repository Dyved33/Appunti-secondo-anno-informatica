// Package fatture gestisce lo stato e l'interfaccia utente della sotto-pagina
// Fatture della Contabilità: consultazione, creazione ed eliminazione di
// fatture e note di credito, con separazione a tab.
package fatture

import "github.com/tommasogarofani/ingsw_25-26_group-11/client/api"

// sezioneAttiva indica quale tab è attualmente selezionata.
type sezioneAttiva int

const (
	SEZ_FATTURE     sezioneAttiva = iota // tab Fatture
	SEZ_NOTE_CREDITO                     // tab Note di Credito
)

// mode rappresenta le modalità operative della pagina Fatture.
type mode int

const (
	modeView                mode = iota // visualizzazione tabella
	modeCreateFattura                   // selezione ordine per creazione fattura
	modeConfirmDeleteFattura            // conferma eliminazione fattura
	modeCreateNotaCredito               // selezione fattura per creazione nota di credito
)

// Model è lo stato interno della pagina Fatture / Note Credito.
type Model struct {
	Client *api.APIClient

	sezione sezioneAttiva
	mode    mode

	// Dati separati per tab
	fatture     []api.Fattura
	noteCredito []api.Fattura
	cursore     int
	offset      int // pagina corrente nella tabella

	// Dati per creazione fattura
	ordiniEvasi       []api.OrdineUscita // ordini evasi senza fattura
	fatturaCursore    int                // cursore nella lista ordini evasi

	// Dati per eliminazione fattura
	confirmDeleteID   int
	confirmDeleteNome string

	loading bool
	err     error
	success string
}

// FattureLoadedMsg viene emesso al termine del caricamento del registro fatture.
type FattureLoadedMsg struct {
	Fatture     []api.Fattura
	NoteCredito []api.Fattura
}

// ErroreMsg è il messaggio di errore condiviso dalle operazioni async della pagina.
type ErroreMsg struct {
	Err error
}

// fatturaCreataMsg viene emesso dopo la creazione riuscita di una fattura.
type fatturaCreataMsg struct{}

// fatturaEliminataMsg viene emesso dopo l'eliminazione riuscita di una fattura.
type fatturaEliminataMsg struct{}

// ordiniEvasiLoadedMsg viene emesso dopo il caricamento degli ordini evasi.
type ordiniEvasiLoadedMsg struct {
	Ordini []api.OrdineUscita
}

// notaCreditoCreataMsg viene emesso dopo la creazione riuscita di una nota di credito.
type notaCreditoCreataMsg struct{}

// New costruisce il modello Fatture in stato di caricamento, in attesa della
// chiamata API lanciata da Init().
func New(client *api.APIClient) Model {
	return Model{Client: client, loading: true}
}
