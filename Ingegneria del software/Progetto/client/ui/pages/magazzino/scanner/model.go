// Package scanner gestisce lo stato ed i modelli dati per la schermata di scansione dei codici a barre e ricerca prodotti.
package scanner

import (
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
)

// Model rappresenta lo stato della pagina Scanner per i codici a barre e prodotti.
// Mantiene il riferimento al client API, il testo immesso nel buffer `input`, il pointer
// al `prodottoTrovato`, lo `storicoScansioni` delle ultime scansioni (max 10), ed i flag di caricamento ed errore.
type Model struct {
	Client *api.APIClient

	input            string
	prodottoTrovato  *api.Prodotto
	storicoScansioni []api.Prodotto

	loading bool
	err     error
}

// ProdottoScansionatoMsg viene emesso al termine della ricerca del prodotto per codice dal backend.
type ProdottoScansionatoMsg struct {
	Prodotto api.Prodotto
}

// ErroreMsg è il messaggio di errore condiviso dalle operazioni asincrone della pagina Scanner.
type ErroreMsg struct {
	Err error
}

// New costruisce ed inizializza un nuovo modello Scanner.
//
// Come è stato scritto il codice:
// Memorizza il client API ed inizializza lo `storicoScansioni` come slice vuota con allocazione iniziale 0.
func New(client *api.APIClient) Model {
	return Model{
		Client:           client,
		storicoScansioni: make([]api.Prodotto, 0),
	}
}
