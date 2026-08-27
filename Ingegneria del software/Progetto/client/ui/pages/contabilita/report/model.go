// Package report gestisce lo stato e l'interfaccia utente della sotto-pagina
// Report della Contabilità: le metriche aggregate su fatturato, documenti
// emessi e ordini evasi.
package report

import "github.com/tommasogarofani/ingsw_25-26_group-11/client/api"

// Model è lo stato interno della pagina Report e Marginalità.
//
// Le metriche sono mantenute nel tipo api.KPI restituito dal client: la pagina
// non ridefinisce un proprio tipo di dominio, come già fanno giacenze e resi
// con api.InventoryItem e api.Reso.
type Model struct {
	Client *api.APIClient

	kpi []api.KPI

	loading bool
	err     error
}

// ReportLoadedMsg viene emesso al termine del calcolo/recupero delle metriche di report.
type ReportLoadedMsg struct {
	KPI []api.KPI
}

// ErroreMsg è il messaggio di errore condiviso dalle operazioni async della pagina.
type ErroreMsg struct {
	Err error
}

// New costruisce il modello Report in stato di caricamento, in attesa della
// chiamata API lanciata da Init().
func New(client *api.APIClient) Model {
	return Model{Client: client, loading: true}
}
