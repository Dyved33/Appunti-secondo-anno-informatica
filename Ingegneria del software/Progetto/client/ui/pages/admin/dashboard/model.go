// Package dashboard gestisce lo stato ed i modelli dati per la schermata di panoramica amministrativa.
// La dashboard aggrega in un'unica vista i KPI provenienti da quattro sorgenti dati distinte:
// report giacenze, report KPI contabilità, ordini in uscita e ordini in entrata.
package dashboard

import (
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
)

// Model rappresenta lo stato della pagina Dashboard admin.
//
// I campi prefissati con "kpi" raccolgono i dati aggregati ricevuti dalle API;
// i campi "loading*" tracciano quali richieste sono ancora in caricamento;
// `cursoreScorteCritiche` permette la navigazione nella lista degli articoli sotto-soglia.
type Model struct {
	client *api.APIClient

	// Dati magazzino/giacenze
	reportGiacenze  *api.InventoryReport
	loadingGiacenze bool
	erroreGiacenze  error

	// Dati KPI contabilità (fatturato, fatture, ordini evasi, note credito)
	kpiContabilita []api.KPI
	loadingKPI     bool
	erroreKPI      error

	// Dati ordini in uscita
	ordiniUscita        []api.OrdineUscita
	loadingOrdiniUscita bool
	erroreOrdiniUscita  error

	// Dati ordini in entrata
	ordiniEntrata        []api.OrdineEntrata
	loadingOrdiniEntrata bool
	erroreOrdiniEntrata  error

	// Navigazione nella lista delle scorte critiche
	cursoreScorteCritiche int

	// Stato di scorrimento della pagina: `winW`/`winH` conservano le dimensioni
	// del terminale (da WindowSizeMsg) e `scrollOffset` è l'indice della prima
	// riga di contenuto mostrata quando la dashboard eccede lo schermo.
	winW         int
	winH         int
	scrollOffset int
}

// New costruisce ed inizializza un nuovo modello Dashboard.
//
// Come è stato scritto il codice:
// Imposta tutti i flag di caricamento a true, in modo che la prima View()
// mostri indicatori di caricamento per ciascun blocco KPI.
// I dati arriveranno tramite i messaggi emessi dai comandi lanciati in Init().
func New(client *api.APIClient) Model {
	return Model{
		client:               client,
		loadingGiacenze:      true,
		loadingKPI:           true,
		loadingOrdiniUscita:  true,
		loadingOrdiniEntrata: true,
	}
}

// — Messaggi interni —

// GiacenzeCaricateMsg viene emesso al termine del caricamento del report giacenze.
type GiacenzeCaricateMsg struct {
	Report *api.InventoryReport
}

// KPICaricatiMsg viene emesso al termine del caricamento dei KPI contabilità.
type KPICaricatiMsg struct {
	KPI []api.KPI
}

// OrdiniUscitaCaricatiMsg viene emesso al termine del caricamento degli ordini in uscita.
type OrdiniUscitaCaricatiMsg struct {
	Ordini []api.OrdineUscita
}

// OrdiniEntrataCaricatiMsg viene emesso al termine del caricamento degli ordini in entrata.
type OrdiniEntrataCaricatiMsg struct {
	Ordini []api.OrdineEntrata
}

// ErroreGiacenzeMsg viene emesso in caso di errore nel caricamento delle giacenze.
type ErroreGiacenzeMsg struct{ Err error }

// ErroreKPIMsg viene emesso in caso di errore nel caricamento dei KPI contabilità.
type ErroreKPIMsg struct{ Err error }

// ErroreOrdiniUscitaMsg viene emesso in caso di errore nel caricamento degli ordini in uscita.
type ErroreOrdiniUscitaMsg struct{ Err error }

// ErroreOrdiniEntrataMsg viene emesso in caso di errore nel caricamento degli ordini in entrata.
type ErroreOrdiniEntrataMsg struct{ Err error }
