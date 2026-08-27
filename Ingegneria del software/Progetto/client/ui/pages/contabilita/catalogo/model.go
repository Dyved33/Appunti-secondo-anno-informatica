// Package catalogo gestisce lo stato e l'interfaccia utente della sotto-pagina
// Catalogo della Contabilità: consultazione, creazione ed eliminazione di
// prodotti e categorie.
package catalogo

import (
	"charm.land/bubbles/v2/textinput"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
)

// mode rappresenta le modalità operative della pagina Catalogo.
type mode int

const (
	modeView                  mode = iota // visualizzazione lista prodotti
	modeCreateProduct                    // form creazione prodotto
	modeCreateCategory                   // form creazione categoria
	modeConfirmDeleteProduct             // conferma eliminazione prodotto
	modeConfirmDeleteCategory            // conferma eliminazione categoria
)

// nessunFiltro è il valore di filtroCategoria che indica "tutte le categorie".
const nessunFiltro = -1

// campiCreateProdotti indica il numero di campi testuali nel form creazione prodotto.
// 0=Codice, 1=Barcode, 2=Nome, 3=Descrizione, 4=PrezzoAcquisto, 5=PrezzoVendita, 6=SogliaMin, 7=Categoria (selettore)
const campiCreateProdotti = 8

// Model è lo stato interno della pagina Catalogo (Prodotti e Categorie).
type Model struct {
	Client *api.APIClient

	mode mode

	// Dati caricati dal backend
	prodotti  []api.Prodotto
	categorie []string
	// categorieID mappa il nome della categoria al suo ID server.
	categorieID map[string]int

	filtroCategoria int // indice in categorie; nessunFiltro = mostra tutto
	cursore         int
	offset          int // pagina corrente nella tabella

	loading bool
	err     error
	success string

	// ── Form creazione prodotto ──────────────────────────────────────────
	createCodice      textinput.Model
	createBarcode     textinput.Model
	createNome        textinput.Model
	createDescrizione textinput.Model
	createPrezzoAcq   textinput.Model
	createPrezzoVend  textinput.Model
	createSogliaMin   textinput.Model
	createFocus       int // 0..7 (7 = selettore categoria)
	createCategoriaIdx int

	// ── Form creazione categoria ─────────────────────────────────────────
	createCatNome        textinput.Model
	createCatDescrizione textinput.Model
	createCatFocus       int // 0=Nome, 1=Descrizione

	// ── Eliminazione ────────────────────────────────────────────────────
	confirmDeleteID   int
	confirmDeleteNome string
}

// ── Messaggi asincroni ────────────────────────────────────────────────────

// CatalogoLoadedMsg viene emesso al termine del caricamento di prodotti e categorie.
type CatalogoLoadedMsg struct {
	Prodotti     []api.Prodotto
	Categorie    []string
	CategorieIDs map[string]int
}

// ErroreMsg è il messaggio di errore condiviso dalle operazioni async della pagina.
type ErroreMsg struct {
	Err error
}

// prodottoCreatoMsg viene emesso dopo la creazione riuscita di un prodotto.
type prodottoCreatoMsg struct {
	Nome string
}

// prodottoEliminatoMsg viene emesso dopo l'eliminazione riuscita di un prodotto.
type prodottoEliminatoMsg struct {
	Nome string
}

// categoriaCreataMsg viene emesso dopo la creazione riuscita di una categoria.
type categoriaCreataMsg struct {
	Nome string
}

// categoriaEliminataMsg viene emesso dopo l'eliminazione riuscita di una categoria.
type categoriaEliminataMsg struct {
	Nome string
}

// New costruisce il modello Catalogo con nessun filtro categoria attivo e in
// stato di caricamento, in attesa della chiamata API lanciata da Init().
func New(client *api.APIClient) Model {
	// Campi form prodotto
	createCodice := textinput.New()
	createCodice.Placeholder = "es. PROD-001"

	createBarcode := textinput.New()
	createBarcode.Placeholder = "EAN-13 (opzionale)"

	createNome := textinput.New()
	createNome.Placeholder = "Nome prodotto"

	createDescrizione := textinput.New()
	createDescrizione.Placeholder = "Descrizione (opzionale)"

	createPrezzoAcq := textinput.New()
	createPrezzoAcq.Placeholder = "0.00"

	createPrezzoVend := textinput.New()
	createPrezzoVend.Placeholder = "0.00"

	createSogliaMin := textinput.New()
	createSogliaMin.Placeholder = "0"

	// Campi form categoria
	createCatNome := textinput.New()
	createCatNome.Placeholder = "Nome categoria"

	createCatDescrizione := textinput.New()
	createCatDescrizione.Placeholder = "Descrizione (opzionale)"

	return Model{
		Client:              client,
		loading:             true,
		filtroCategoria:     nessunFiltro,
		mode:                modeView,
		createCodice:        createCodice,
		createBarcode:       createBarcode,
		createNome:          createNome,
		createDescrizione:   createDescrizione,
		createPrezzoAcq:     createPrezzoAcq,
		createPrezzoVend:    createPrezzoVend,
		createSogliaMin:     createSogliaMin,
		createCategoriaIdx:  0,
		createCatNome:       createCatNome,
		createCatDescrizione: createCatDescrizione,
		categorieID:         make(map[string]int),
	}
}
