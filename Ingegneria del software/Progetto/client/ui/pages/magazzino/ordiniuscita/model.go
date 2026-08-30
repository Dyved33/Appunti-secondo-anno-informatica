// Package ordiniuscita gestisce lo stato ed i modelli dati per la schermata degli ordini di vendita ed evasione verso i clienti.
package ordiniuscita

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
)

// VistaStato identifica la vista attualmente mostrata dalla pagina Ordini in Uscita.
type VistaStato int

const (
	// MAIN è la vista principale con l'elenco degli ordini in uscita.
	MAIN VistaStato = iota
	// DDT è la vista che mostra il dettaglio del documento di trasporto generato.
	DDT
)

type modoVista int

const (
	modoLista modoVista = iota
	modoModifica
)

var statiOrdine = []string{"confermato", "evaso", "annullato"}

type focusForm int

const (
	focusNumero focusForm = iota
	focusCliente
	focusNote
	focusStato
)

// Model rappresenta lo stato della pagina Ordini in Uscita.
// Mantiene l'elenco degli ordini, il cursore di navigazione, la mappa delle righe
// espanse, il DDT generato, la vista e modalità correnti, i campi del form di modifica
// ed i dati di scorrimento e stato (caricamento/errore).
type Model struct {
	Client *api.APIClient

	ordini          []api.OrdineUscita
	cursore         int
	dettagliEspansi map[int]bool
	DDT             api.DDT

	view     VistaStato
	modo     modoVista
	focusMod focusForm

	formNumero  textinput.Model
	formCliente textinput.Model
	formNote    textinput.Model
	idxStato    int

	winW         int
	winH         int
	scrollOffset int

	loading bool
	err     error
}

// OrdiniUscitaLoadedMsg viene emesso al termine del caricamento degli ordini in uscita dal backend.
type OrdiniUscitaLoadedMsg struct {
	Ordini []api.OrdineUscita
}

// ErroreMsg è il messaggio di errore condiviso dalle operazioni asincrone della pagina Ordini in Uscita.
type ErroreMsg struct {
	Err error
}

// DDTCreatedMsg viene emesso al termine della generazione del documento di trasporto.
type DDTCreatedMsg struct {
	DDT api.DDT
}

// OrdineAggiornatoMsg viene emesso a seguito del corretto aggiornamento di un ordine in uscita.
type OrdineAggiornatoMsg struct{}

// New costruisce ed inizializza un nuovo modello Ordini in Uscita.
// Inizializza la mappa `dettagliEspansi`, un DDT vuoto, `loading = true`, la vista MAIN
// e la modalità lista, memorizzando il riferimento al client API.
func New(client *api.APIClient) Model {
	return Model{
		Client:          client,
		dettagliEspansi: make(map[int]bool),
		loading:         true,
		DDT:             api.DDT{},
		view:            MAIN,
		modo:            modoLista,
	}
}

// preparaFormModifica inizializza i campi del form di modifica con i dati
// dell'ordine sottolineato ed imposta il focus sul campo numero.
func (m *Model) preparaFormModifica() {
	if len(m.ordini) == 0 || m.cursore >= len(m.ordini) {
		return
	}
	o := m.ordini[m.cursore]

	m.formNumero = textinput.New()
	m.formNumero.SetValue(o.Numero)
	m.formNumero.Focus()

	m.formCliente = textinput.New()
	m.formCliente.SetValue(o.Cliente)

	m.formNote = textinput.New()
	m.formNote.SetValue(o.Note)

	m.idxStato = indexStato(string(o.Stato))
	m.focusMod = focusNumero
}

func indexStato(stato string) int {
	for i, s := range statiOrdine {
		if s == stato {
			return i
		}
	}
	return 0
}

func (m *Model) focusCampoForm() tea.Cmd {
	m.formNumero.Blur()
	m.formCliente.Blur()
	m.formNote.Blur()

	switch m.focusMod {
	case focusNumero:
		return m.formNumero.Focus()
	case focusCliente:
		return m.formCliente.Focus()
	case focusNote:
		return m.formNote.Focus()
	default:
		return nil
	}
}

func (m Model) currentFormCursor() *tea.Cursor {
	switch m.focusMod {
	case focusNumero:
		return m.formNumero.Cursor()
	case focusCliente:
		return m.formCliente.Cursor()
	case focusNote:
		return m.formNote.Cursor()
	default:
		return nil
	}
}
