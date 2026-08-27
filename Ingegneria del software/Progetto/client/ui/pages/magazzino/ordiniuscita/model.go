// Package ordiniuscita gestisce lo stato ed i modelli dati per la schermata degli ordini di vendita ed evasione verso i clienti.
package ordiniuscita

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
)

type VistaStato int

const (
	MAIN VistaStato = iota
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

	winW        int
	winH        int
	scrollOffset int

	loading bool
	err     error
}

type OrdiniUscitaLoadedMsg struct {
	Ordini []api.OrdineUscita
}

type ErroreMsg struct {
	Err error
}

type DDTCreatedMsg struct {
	DDT api.DDT
}

type OrdineAggiornatoMsg struct{}

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
