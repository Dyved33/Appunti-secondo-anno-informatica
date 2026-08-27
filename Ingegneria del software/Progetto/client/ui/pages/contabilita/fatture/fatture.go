// Package fatture contiene le logiche di rendering e di gestione eventi della
// sotto-pagina Fatture della Contabilità.
package fatture

import (
	"fmt"
	"strings"

	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
	style "github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/styles"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// Init avvia il caricamento del registro fatture.
func (m Model) Init() tea.Cmd {
	return m.loadFattureCmd()
}

// Update gestisce il caricamento asincrono, il cambio tab, la navigazione
// e le modalità creazione/eliminazione.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case FattureLoadedMsg:
		m.loading = false
		m.err = nil
		m.fatture = msg.Fatture
		m.noteCredito = msg.NoteCredito
		m.cursore = 0
		m.offset = 0
		return m, nil

	case ErroreMsg:
		m.loading = false
		m.err = msg.Err
		return m, nil

	case fatturaCreataMsg:
		m.mode = modeView
		m.loading = false
		m.err = nil
		m.success = "Fattura creata con successo."
		return m, m.loadFattureCmd()

	case fatturaEliminataMsg:
		m.mode = modeView
		m.loading = false
		m.err = nil
		m.success = "Fattura eliminata."
		return m, m.loadFattureCmd()

	case ordiniEvasiLoadedMsg:
		m.loading = false
		m.err = nil
		m.ordiniEvasi = msg.Ordini
		m.fatturaCursore = 0
		return m, nil

	case notaCreditoCreataMsg:
		m.mode = modeView
		m.loading = false
		m.err = nil
		m.sezione = SEZ_NOTE_CREDITO
		m.success = "Nota di credito creata con successo."
		return m, m.loadFattureCmd()

	case tea.WindowSizeMsg:
		return m, nil

	case tea.KeyPressMsg:
		switch m.mode {
		case modeCreateFattura:
			return m.handleCreateFatturaInput(msg)
		case modeConfirmDeleteFattura:
			return m.handleConfirmDeleteInput(msg)
		case modeCreateNotaCredito:
			return m.handleCreateNotaCreditoInput(msg)
		default:
			return m.handleViewInput(msg)
		}
	}

	return m, nil
}

// ── View ──────────────────────────────────────────────────────────────────

func (m Model) View() tea.View {
	var sb strings.Builder

	sb.WriteString(style.TitleStile.Render("CONTABILITÀ - FATTURE E NOTE DI CREDITO"))
	sb.WriteString("\n\n")

	if m.loading {
		sb.WriteString("Caricamento in corso...\n")
		return tea.NewView(sb.String())
	}
	if m.err != nil {
		sb.WriteString(style.ErrorStile.Render("Errore: " + m.err.Error()))
		sb.WriteString("\n")
		return tea.NewView(sb.String())
	}

	switch m.mode {
	case modeCreateFattura:
		sb.WriteString(m.renderCreateFattura())
	case modeCreateNotaCredito:
		sb.WriteString(m.renderCreateNotaCredito())
	case modeConfirmDeleteFattura:
		sb.WriteString(style.ErrorStile.Render(
			fmt.Sprintf("Eliminare la fattura \"%s\"? L'operazione è definitiva.", m.confirmDeleteNome)))
		sb.WriteString("\n\n")
		sb.WriteString(style.ErrorStile.Render("[enter] Conferma  |  [esc] Annulla"))
	default:
		sb.WriteString(m.renderMainView())
	}

	// Messaggio di successo
	if m.success != "" {
		sb.WriteString("\n\n")
		sb.WriteString(style.SuccessStile.Render(m.success))
	}

	return tea.NewView(sb.String())
}

// ── Render principale (tab + tabella) ─────────────────────────────────────

func (m Model) renderMainView() string {
	var sb strings.Builder

	// KPI compatti su una riga
	netto := m.fatturatoNetto()
	kpi := fmt.Sprintf("  Fatturato netto: € %.2f  |  Fatture: %d  |  Note di credito: %d",
		netto, len(m.fatture), len(m.noteCredito))
	sb.WriteString(style.InfoStile.Render(kpi))
	sb.WriteString("\n\n")

	// Tab header
	sb.WriteString(m.renderTabHeader())
	sb.WriteString("\n")

	// Tabella corrente
	lista := m.listaCorrente()
	if len(lista) == 0 {
		if m.sezione == SEZ_FATTURE {
			sb.WriteString("Nessuna fattura registrata.\n")
		} else {
			sb.WriteString("Nessuna nota di credito registrata.\n")
		}
	} else {
		headerStr := fmt.Sprintf("   %-12s %-8s %-24s %-14s %s",
			"NUMERO", "ANNO", "CLIENTE", "IMPORTO", "TIPO")
		sb.WriteString(style.HeaderStile.Render(headerStr))
		sb.WriteString("\n")
		sb.WriteString(strings.Repeat("─", 80))
		sb.WriteString("\n")

		inizio := m.offset
		fine := inizio + style.RighePerPagina
		if fine > len(lista) {
			fine = len(lista)
		}

		for i, f := range lista[inizio:fine] {
			assoluto := inizio + i
			tipo := "fattura"
			if f.NotaCredito {
				tipo = "nota di credito"
			}

			rigaTesto := fmt.Sprintf("%-12s %-8d %-24s %-14s %s",
				f.Numero,
				f.Anno,
				troncaStringa(f.Cliente, 24),
				fmt.Sprintf("€ %.2f", f.Importo),
				tipo,
			)

			if assoluto == m.cursore {
				sb.WriteString(style.SuccessStile.Render("> "))
				sb.WriteString(rigaTesto)
			} else {
				sb.WriteString("  ")
				sb.WriteString(rigaTesto)
			}
			sb.WriteString("\n")
		}

		if len(lista) > style.RighePerPagina {
			pagine := (len(lista) + style.RighePerPagina - 1) / style.RighePerPagina
			paginaCorr := m.offset/style.RighePerPagina + 1
			sb.WriteString(fmt.Sprintf("\nPagina %d/%d", paginaCorr, pagine))
		}
	}

	// Footer compatto su riga unica
	sb.WriteString("\n")
	if m.sezione == SEZ_FATTURE {
		sb.WriteString(style.FooterStile.Render(
			"[↑/↓] Naviga  [pgup/pgdown] Pagina  [tab] Tab  [ctrl+n] Nuova fattura  [ctrl+k] Nota credito  [ctrl+d] Elimina  [r] Ricarica"))
	} else {
		sb.WriteString(style.FooterStile.Render(
			"[↑/↓] Naviga  [pgup/pgdown] Pagina  [tab] Tab  [r] Ricarica"))
	}

	return sb.String()
}

// ── Render creazione fattura ──────────────────────────────────────────────

func (m Model) renderCreateFattura() string {
	return m.renderSelezioneOrdine("SELEZIONA ORDINE DA FATTURARE", "Crea fattura")
}

func (m Model) renderCreateNotaCredito() string {
	return m.renderSelezioneOrdine("SELEZIONA ORDINE PER NOTA DI CREDITO", "Crea nota di credito")
}

func (m Model) renderSelezioneOrdine(titolo, azione string) string {
	var sb strings.Builder

	sb.WriteString(style.SubTitleStile.Render(titolo))
	sb.WriteString("\n\n")

	if len(m.ordiniEvasi) == 0 {
		sb.WriteString("Nessun ordine evaso disponibile.\n")
	} else {
		headerStr := fmt.Sprintf("   %-10s %-24s %-20s %s",
			"ID", "NUMERO", "CLIENTE", "DATA")
		sb.WriteString(style.HeaderStile.Render(headerStr))
		sb.WriteString("\n")
		sb.WriteString(strings.Repeat("─", 70))
		sb.WriteString("\n")

		for i, o := range m.ordiniEvasi {
			rigaTesto := fmt.Sprintf("%-10d %-24s %-20s %s",
				o.ID,
				troncaStringa(o.Numero, 24),
				troncaStringa(o.Cliente, 20),
				o.Data.Format("02/01/2006"),
			)

			if i == m.fatturaCursore {
				sb.WriteString(style.SuccessStile.Render("> "))
				sb.WriteString(rigaTesto)
			} else {
				sb.WriteString("  ")
				sb.WriteString(rigaTesto)
			}
			sb.WriteString("\n")
		}
	}

	sb.WriteString("\n")
	sb.WriteString(style.FooterStile.Render(fmt.Sprintf("[↑/↓] Naviga  |  [enter] %s  |  [esc] Annulla", azione)))

	return sb.String()
}

// ── Tab header ────────────────────────────────────────────────────────────

func (m Model) renderTabHeader() string {
	t1 := " FATTURE "
	t2 := " NOTE DI CREDITO "

	switch m.sezione {
	case SEZ_FATTURE:
		t1 = style.FocusStile.Render("[ " + t1 + " ]")
	case SEZ_NOTE_CREDITO:
		t2 = style.FocusStile.Render("[ " + t2 + " ]")
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, t1, "   ", t2)
}

// ── Input handlers ────────────────────────────────────────────────────────

func (m Model) handleViewInput(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	lista := m.listaCorrente()
	switch msg.String() {
	case "tab":
		if m.sezione == SEZ_FATTURE {
			m.sezione = SEZ_NOTE_CREDITO
		} else {
			m.sezione = SEZ_FATTURE
		}
		m.cursore = 0
		m.offset = 0

	case "up", "k":
		if m.cursore > 0 {
			m.cursore--
			m.offset = style.PaginaCalcola(m.cursore, m.offset, len(lista))
		}
	case "down", "j":
		if m.cursore < len(lista)-1 {
			m.cursore++
			m.offset = style.PaginaCalcola(m.cursore, m.offset, len(lista))
		}
	case "pgup":
		m.cursore -= style.RighePerPagina
		if m.cursore < 0 {
			m.cursore = 0
		}
		m.offset = style.PaginaCalcola(m.cursore, m.offset, len(lista))
	case "pgdown":
		m.cursore += style.RighePerPagina
		if max := len(lista) - 1; m.cursore > max {
			m.cursore = max
		}
		if m.cursore < 0 {
			m.cursore = 0
		}
		m.offset = style.PaginaCalcola(m.cursore, m.offset, len(lista))
	case "r":
		m.loading = true
		m.success = ""
		return m, m.loadFattureCmd()

	case "ctrl+n":
		if m.sezione != SEZ_FATTURE {
			return m, nil
		}
		m.mode = modeCreateFattura
		m.success = ""
		m.err = nil
		return m, m.loadOrdiniEvasiCmd()

	case "ctrl+k":
		m.mode = modeCreateNotaCredito
		m.success = ""
		m.err = nil
		return m, m.loadOrdiniEvasiCmd()

	case "ctrl+d":
		if m.sezione != SEZ_FATTURE {
			return m, nil
		}
		if len(m.fatture) == 0 {
			return m, nil
		}
		f := m.fatture[m.cursore]
		m.mode = modeConfirmDeleteFattura
		m.confirmDeleteID = f.ID
		m.confirmDeleteNome = f.Numero
		m.success = ""
		m.err = nil
		return m, nil
	}

	return m, nil
}

func (m Model) handleCreateFatturaInput(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = modeView
		m.fatturaCursore = 0
		return m, nil
	case "up", "k":
		if m.fatturaCursore > 0 {
			m.fatturaCursore--
		}
	case "down", "j":
		if m.fatturaCursore < len(m.ordiniEvasi)-1 {
			m.fatturaCursore++
		}
	case "enter":
		if len(m.ordiniEvasi) == 0 {
			return m, nil
		}
		ordine := m.ordiniEvasi[m.fatturaCursore]
		m.loading = true
		m.err = nil
		m.success = ""
		return m, m.creaFatturaCmd(ordine.ID)
	}

	return m, nil
}

func (m Model) handleConfirmDeleteInput(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		m.loading = true
		m.err = nil
		m.success = ""
		return m, m.eliminaFatturaCmd(m.confirmDeleteID)
	case "esc":
		m.mode = modeView
		m.confirmDeleteID = 0
		m.confirmDeleteNome = ""
		return m, nil
	}
	return m, nil
}

func (m Model) handleCreateNotaCreditoInput(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = modeView
		m.fatturaCursore = 0
		return m, nil
	case "up", "k":
		if m.fatturaCursore > 0 {
			m.fatturaCursore--
		}
	case "down", "j":
		if m.fatturaCursore < len(m.ordiniEvasi)-1 {
			m.fatturaCursore++
		}
	case "enter":
		if len(m.ordiniEvasi) == 0 {
			return m, nil
		}
		ordine := m.ordiniEvasi[m.fatturaCursore]
		m.loading = true
		m.err = nil
		m.success = ""
		return m, m.creaNotaCreditoCmd(ordine.ID)
	}
	return m, nil
}

// ── Dati ──────────────────────────────────────────────────────────────────

func (m Model) listaCorrente() []api.Fattura {
	if m.sezione == SEZ_FATTURE {
		return m.fatture
	}
	return m.noteCredito
}

func (m Model) fatturatoNetto() float64 {
	var netto float64
	for _, f := range m.fatture {
		netto += f.Importo
	}
	for _, nc := range m.noteCredito {
		netto -= nc.Importo
	}
	return netto
}

// ── Async commands ────────────────────────────────────────────────────────

func (m Model) loadFattureCmd() tea.Cmd {
	return func() tea.Msg {
		fatture, err := m.Client.GetFatture()
		if err != nil {
			return ErroreMsg{Err: err}
		}

		var soleFatture []api.Fattura
		var soleNC []api.Fattura
		for _, f := range fatture {
			if f.NotaCredito {
				soleNC = append(soleNC, f)
			} else {
				soleFatture = append(soleFatture, f)
			}
		}

		return FattureLoadedMsg{Fatture: soleFatture, NoteCredito: soleNC}
	}
}

func (m Model) loadOrdiniEvasiCmd() tea.Cmd {
	return func() tea.Msg {
		ordini, err := m.Client.GetOrdiniEvasiLeggero()
		if err != nil {
			return ErroreMsg{Err: err}
		}
		return ordiniEvasiLoadedMsg{Ordini: ordini}
	}
}

func (m Model) creaFatturaCmd(ordineID int) tea.Cmd {
	return func() tea.Msg {
		if err := m.Client.CreaFattura(ordineID); err != nil {
			return ErroreMsg{Err: err}
		}
		return fatturaCreataMsg{}
	}
}

func (m Model) eliminaFatturaCmd(id int) tea.Cmd {
	return func() tea.Msg {
		if err := m.Client.EliminaFattura(id); err != nil {
			return ErroreMsg{Err: err}
		}
		return fatturaEliminataMsg{}
	}
}

func (m Model) creaNotaCreditoCmd(ordineID int) tea.Cmd {
	return func() tea.Msg {
		if err := m.Client.CreaNotaCredito(ordineID); err != nil {
			return ErroreMsg{Err: err}
		}
		return notaCreditoCreataMsg{}
	}
}

func troncaStringa(s string, maxLen int) string {
	if len(s) > maxLen {
		if maxLen > 3 {
			return s[:maxLen-3] + "..."
		}
		return s[:maxLen]
	}
	return s
}
