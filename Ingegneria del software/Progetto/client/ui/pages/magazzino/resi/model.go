// Package resi gestisce lo stato ed i modelli dati per la schermata di gestione dei resi merce e ripristino giacenze.
package resi

import (
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
)

type modoVista int

const (
	modoLista modoVista = iota
	modoModificaStato
)

var statiReso = []string{"da_verificare", "accettato", "rifiutato"}

// Model rappresenta lo stato della pagina di gestione dei Resi.
// Mantiene l'elenco delle pratiche di reso, il cursore di navigazione con relativo
// offset di scorrimento, la mappa delle pratiche con dettagli espansi, la modalità
// corrente (lista o modifica stato) ed i campi di stato per caricamento ed errori.
type Model struct {
	Client *api.APIClient

	resi            []api.Reso
	cursore         int
	offset          int
	dettagliEspansi map[int]bool

	modo     modoVista
	idxStato int

	loading bool
	err     error
}

// ResiLoadedMsg viene emesso al termine del caricamento dell'elenco dei resi dal backend.
type ResiLoadedMsg struct {
	Resi []api.Reso
}

// ErroreMsg è il messaggio di errore condiviso dalle operazioni asincrone della pagina Resi.
type ErroreMsg struct {
	Err error
}

// ResoAggiornatoMsg viene emesso a seguito del corretto aggiornamento dello stato di un reso.
type ResoAggiornatoMsg struct{}

// New costruisce ed inizializza un nuovo modello Resi.
// Inizializza la mappa `dettagliEspansi`, imposta `loading = true` per la schermata iniziale
// e `modo = modoLista`, memorizzando il riferimento al client API.
func New(client *api.APIClient) Model {
	return Model{
		Client:          client,
		dettagliEspansi: make(map[int]bool),
		loading:         true,
		modo:            modoLista,
	}
}

func indexStatoReso(stato string) int {
	for i, s := range statiReso {
		if s == stato {
			return i
		}
	}
	return 0
}
