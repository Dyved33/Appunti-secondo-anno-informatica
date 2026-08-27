// Package contabilita gestisce la vista principale e il routing delle sotto-pagine
// inerenti alla contabilità.
package contabilita

import (
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/contabilita/catalogo"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/contabilita/fatture"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/contabilita/ordini"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/contabilita/report"
)

// sezione identifica quale sotto-pagina della Contabilità è attualmente visibile.
type focus int

// Costanti enumerate per identificare le sotto-pagine disponibili nel modulo Contabilità.
const (
	SezioneOrdini focus = iota
	SezioneFatture
	SezioneCatalogo
	SezioneReport
)

// sezioneLabels contiene le etichette usate per disegnare la barra dei tab in View().
// L'indice dell'array corrisponde al valore della costante sezione.
// I tab si selezionano cliccandoli col mouse (vedi Update).
var labels = [...]string{
	SezioneOrdini:   "Ordini & Listini",
	SezioneFatture:  "Fatture",
	SezioneCatalogo: "Catalogo",
	SezioneReport:   "Report",
}

// ContLayout è il modello contenitore e router principale della sezione Contabilità.
//
// Mantiene vive in memoria tutte e quattro le istanze dei sotto-modelli: così
// quando l'utente cambia tab e poi torna indietro, lo stato della pagina
// precedente (cursore, filtri, dati già caricati) non va perso.
type ContLayout struct {
	Client *api.APIClient
	attiva focus

	// winW e winH conservano le dimensioni correnti del terminale, necessarie
	// per calcolare le aree cliccabili dei tab (vedi Update).
	winW, winH int

	ordini   ordini.Model
	fatture  fatture.Model
	catalogo catalogo.Model
	report   report.Model
}

// New costruisce il layout di Contabilità e inizializza i quattro sotto-modelli.
//
// Come è stato scritto il codice:
//  1. Associa il client API fornito al layout principale.
//  2. Imposta la sotto-pagina attiva iniziale su SezioneOrdini.
//  3. Inizializza ciascun sotto-modello delegando al rispettivo costruttore New(client),
//     che riceve lo stesso *api.APIClient condiviso: ogni sotto-pagina effettua le
//     proprie chiamate API in autonomia senza passare dal genitore.
func New(client *api.APIClient) ContLayout {
	return ContLayout{
		Client: client,
		attiva: SezioneOrdini,

		ordini:   ordini.New(client),
		fatture:  fatture.New(client),
		catalogo: catalogo.New(client),
		report:   report.New(client),
	}
}
