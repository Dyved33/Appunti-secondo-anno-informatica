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

type Model struct {
	Client *api.APIClient

	resi            []api.Reso
	cursore         int
	offset          int
	dettagliEspansi map[int]bool

	modo    modoVista
	idxStato int

	loading bool
	err     error
}

type ResiLoadedMsg struct {
	Resi []api.Reso
}

type ErroreMsg struct {
	Err error
}

type ResoAggiornatoMsg struct{}

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
