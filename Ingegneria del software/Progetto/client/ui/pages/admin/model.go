// Package admin implementa il layout contenitore e router del modulo amministrativo.
// Raggruppa le sotto-pagine (dashboard, utenti, magazzino, contabilità) in una
// barra di tab navigabile e gestisce lo smistamento degli eventi alla pagina attiva.
package admin

import (
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/admin/dashboard"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/admin/users"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/contabilita/catalogo"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/contabilita/fatture"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/contabilita/ordini"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/contabilita/report"
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
	DASHBOARD focus = iota
	USERS
	GIACENZE
	ORDINIENTRATA
	ORDINIUSCITA
	RESI
	SCANNER
	CATALOGO
	FATTURE
	ORDINI
	REPORT
)

// labels contiene le etichette usate per il rendering della barra dei tab in View().
// L'indice dell'array corrisponde alla costante enumerata focus.
// I tab si selezionano cliccandoli col mouse (vedi Update).
var labels = [...]string{
	DASHBOARD:     "Dashboard",
	USERS:         "Utenti",
	GIACENZE:      "Giacenze",
	ORDINIENTRATA: "Ordini In Entrata",
	ORDINIUSCITA:  "Ordini in Uscita",
	RESI:          "Resi",
	SCANNER:       "Scanner",
	CATALOGO:      "Catalogo",
	FATTURE:       "Fatture",
	ORDINI:        "Ordini",
	REPORT:        "Report",
}

// AdminLayout rappresenta il modello contenitore e router principale per il modulo Admin.
//
// Mantiene in memoria lo stato delle 12 sotto-pagine
// ed il riferimento al client API condiviso, consentendo
// il mantenimento del contesto al cambio di tab.
type AdminLayout struct {
	Client *api.APIClient
	attiva focus

	// winW e winH conservano le dimensioni correnti del terminale, necessarie
	// per calcolare le aree cliccabili dei tab (vedi Update).
	winW, winH int

	dashboard     dashboard.Model
	users         users.Model
	giacenze      giacenze.Model
	ordiniEntrata ordinientrata.Model
	ordiniUscita  ordiniuscita.Model
	resi          resi.Model
	scanner       scanner.Model
	catalogo      catalogo.Model
	fatture       fatture.Model
	ordini        ordini.Model
	report        report.Model
}

// New crea ed inizializza una nuova istanza del modello AdminLayout.
//
// Come è stato scritto il codice:
// 1. Associa il client API fornito al layout principale.
// 2. Imposta la sotto-pagina attiva iniziale su DASHBOARD.
// 3. Inizializza ciascun sotto-modello delegando al rispettivo costruttore New(client).
func New(client *api.APIClient) AdminLayout {
	return AdminLayout{
		Client: client,
		attiva: DASHBOARD,

		dashboard:     dashboard.New(client),
		users:         users.New(client),
		giacenze:      giacenze.New(client),
		ordiniEntrata: ordinientrata.New(client),
		ordiniUscita:  ordiniuscita.New(client),
		resi:          resi.New(client),
		scanner:       scanner.New(client),
		catalogo:      catalogo.New(client),
		fatture:       fatture.New(client),
		ordini:        ordini.New(client),
		report:        report.New(client),
	}
}
