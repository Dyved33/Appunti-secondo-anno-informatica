// Package magazzino gestisce la vista principale e il routing delle sotto-pagine inerenti al magazzino.
package magazzino

import (
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/magazzino/giacenze"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/magazzino/ordinientrata"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/magazzino/ordiniuscita"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/magazzino/resi"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/magazzino/scanner"
)

// focus identifica la sotto-pagina del magazzino attualmente visibile ed attiva.
type focus int

// Costanti enumerate per identificare le sotto-pagine disponibili nel modulo Magazzino.
const (
	GIACENZE focus = iota
	ORDINIENTRATA
	ORDINIUSCITA
	RESI
	SCANNER
)

// labels contiene le etichette usate per il rendering della barra dei tab in View().
// L'indice dell'array corrisponde alla costante enumerata focus.
// I tab si selezionano cliccandoli col mouse (vedi Update).
var labels = [...]string{
	GIACENZE:      "Giacenze",
	ORDINIENTRATA: "Ordini In Entrata",
	ORDINIUSCITA:  "Ordini in Uscita",
	RESI:          "Resi",
	SCANNER:       "Scanner",
}

// MagLayout rappresenta il modello contenitore e router principale per il modulo Magazzino.
//
// Mantiene in memoria lo stato delle cinque sotto-pagine (Giacenze, Ordini in Entrata,
// Ordini in Uscita, Resi, Scanner) ed il riferimento al client API condiviso, consentendo
// il mantenimento del contesto al cambio di tab.
type MagLayout struct {
	Client *api.APIClient
	attiva focus

	// winW e winH conservano le dimensioni correnti del terminale, necessarie
	// per calcolare le aree cliccabili dei tab (vedi Update).
	winW, winH int

	giacenze      giacenze.Model
	ordiniEntrata ordinientrata.Model
	ordiniUscita  ordiniuscita.Model
	resi          resi.Model
	scanner       scanner.Model
}

// New crea ed inizializza una nuova istanza del modello MagLayout.
//
// Come è stato scritto il codice:
// 1. Associa il client API fornito al layout principale.
// 2. Imposta la sotto-pagina attiva iniziale su GIACENZE.
// 3. Inizializza ciascun sotto-modello delegando al rispettivo costruttore New(client).
func New(client *api.APIClient) MagLayout {
	return MagLayout{
		Client: client,
		attiva: GIACENZE,

		giacenze:      giacenze.New(client),
		ordiniEntrata: ordinientrata.New(client),
		ordiniUscita:  ordiniuscita.New(client),
		resi:          resi.New(client),
		scanner:       scanner.New(client),
	}
}
