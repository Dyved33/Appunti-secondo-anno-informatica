// Package ordinientrata gestisce lo stato ed i modelli dati per la schermata degli ordini di acquisto in entrata dai fornitori.
package ordinientrata

import (
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
)

// Model rappresenta lo stato della pagina Ordini in Entrata.
//
// Oltre ai dati degli ordini mantiene lo stato di scorrimento della lista:
// `winW`/`winH` conservano le dimensioni del terminale (da WindowSizeMsg) e
// `scrollOffset` è l'indice della prima riga di visualizzazione mostrata,
// necessario perché gli ordini oltre la capacità dello schermo resterebbero invisibili.
type Model struct {
	Client *api.APIClient

	ordini          []api.OrdineEntrata
	cursore         int
	dettagliEspansi map[int]bool // ID ordine -> true se i dettagli delle righe sono espansi

	winW         int
	winH         int
	scrollOffset int

	loading bool
	err     error
}

// OrdiniEntrataLoadedMsg viene emesso al termine del caricamento degli ordini in entrata dal backend.
type OrdiniEntrataLoadedMsg struct {
	Ordini []api.OrdineEntrata
}

// ErroreMsg è il messaggio di errore condiviso dalle operazioni asincrone della pagina Ordini in Entrata.
type ErroreMsg struct {
	Err error
}

// New costruisce ed inizializza un nuovo modello OrdiniEntrata.
//
// Come è stato scritto il codice:
// Inizializza la mappa `dettagliEspansi`, imposta il flag `loading = true` per la schermata iniziale,
// e memorizza il riferimento al client API.
func New(client *api.APIClient) Model {
	return Model{
		Client:          client,
		dettagliEspansi: make(map[int]bool),
		loading:         true,
	}
}
