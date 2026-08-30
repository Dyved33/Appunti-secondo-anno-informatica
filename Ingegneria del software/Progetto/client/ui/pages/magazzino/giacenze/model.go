// Package giacenze gestisce lo stato e l'interfaccia utente per la consultazione delle giacenze di magazzino e dei report di inventario.
package giacenze

import (
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
)

// Model rappresenta lo stato della pagina Giacenze.
// Mantiene il riferimento al client API per recuperare i report di inventario,
// l'oggetto InventoryReport ricevuto dal backend, l'indice del cursore per la navigazione
// in tabella, il flag per il filtro degli articoli sotto-soglia ed i campi di stato per il caricamento ed errori.
type Model struct {
	Client *api.APIClient

	report            *api.InventoryReport
	filtroSottoSoglia bool
	cursore           int
	offset            int // prima riga visibile nella finestra di scorrimento

	winW int
	winH int

	loading bool
	err     error
}

// GiacenzeLoadedMsg viene emesso al termine del caricamento del report giacenze dal backend.
type GiacenzeLoadedMsg struct {
	Report *api.InventoryReport
}

// ErroreMsg è il messaggio di errore condiviso dalle operazioni asincrone della pagina Giacenze.
type ErroreMsg struct {
	Err error
}

// New costruisce ed inizializza un nuovo modello Giacenze.
//
// Come è stato scritto il codice:
// Riceve il client API ed imposta il flag `loading` a true in modo che la prima vista
// mostri l'indicatore di caricamento in attesa dell'esecuzione della chiamata API in Init().
func New(client *api.APIClient) Model {
	return Model{
		Client:  client,
		loading: true,
	}
}
