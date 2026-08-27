// Package report contiene le logiche di rendering e di gestione eventi della
// sotto-pagina Report della Contabilità.
package report

import (
	"fmt"
	"strings"

	style "github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/styles"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// Init avvia il caricamento delle metriche di report (fatturato, documenti, ordini evasi).
func (m Model) Init() tea.Cmd {
	return m.loadReportCmd()
}

// Update gestisce solo il refresh manuale ("r"): il report è una vista di sola
// lettura, senza navigazione interna né form di modifica, quindi non serve
// nessun cursore né delega a widget Bubbles.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ReportLoadedMsg:
		m.loading = false
		m.err = nil
		m.kpi = msg.KPI
		return m, nil

	case ErroreMsg:
		m.loading = false
		m.err = msg.Err
		return m, nil

	case tea.WindowSizeMsg:
		return m, nil

	case tea.KeyPressMsg:
		if msg.String() == "r" {
			m.loading = true
			return m, m.loadReportCmd()
		}
	}

	return m, nil
}

// View disegna le metriche come box affiancati, riusando lo stesso LoginBoxStile
// già impiegato per i KPI di giacenze e resi.
//
// La cornice esterna non viene applicata qui: è ContLayout.View() ad avvolgere
// la sotto-pagina attiva in style.FinestraStile, come avviene in magazzino.
func (m Model) View() tea.View {
	var sb strings.Builder

	// ── Titolo ──────────────────────────────────────────────────────────────
	sb.WriteString(style.TitleStile.Render("CONTABILITÀ - REPORT E MARGINALITÀ"))
	sb.WriteString("\n\n")

	if m.loading {
		sb.WriteString("Calcolo metriche in corso...\n")
		return tea.NewView(sb.String())
	}

	if m.err != nil {
		sb.WriteString(style.ErrorStile.Render("Errore: " + m.err.Error()))
		sb.WriteString("\n")
		return tea.NewView(sb.String())
	}

	if len(m.kpi) == 0 {
		sb.WriteString("Nessuna metrica disponibile.\n")
	} else {
		// ── Griglia di Box KPI ────────────────────────────────────────────────
		var boxes []string
		for _, k := range m.kpi {
			contenuto := fmt.Sprintf("%s\n\n%s", strings.ToUpper(k.Etichetta), style.SuccessStile.Render(k.Valore))
			boxes = append(boxes, style.BoxStile.Render(contenuto))
		}

		// Disponiamo i box a coppie di 2 per riga usando JoinHorizontal
		for i := 0; i < len(boxes); i += 2 {
			if i+1 < len(boxes) {
				// Riga con due box affiancati da 3 spazi
				riga := lipgloss.JoinHorizontal(lipgloss.Top, boxes[i], "   ", boxes[i+1])
				sb.WriteString(riga)
			} else {
				// Singolo box rimanente sull'ultima riga
				sb.WriteString(boxes[i])
			}
			sb.WriteString("\n\n")
		}
	}

	// ── Footer ───────────────────────────────────────────────────────────────
	sb.WriteString(style.FooterStile.Render("[r] Ricalcola  |  [ctrl+l] Logout"))

	return tea.NewView(sb.String())
}

// loadReportCmd recupera dal backend le metriche aggregate per la dashboard di report.
func (m Model) loadReportCmd() tea.Cmd {
	return func() tea.Msg {
		kpi, err := m.Client.GetReportKPI()
		if err != nil {
			return ErroreMsg{Err: err}
		}
		return ReportLoadedMsg{KPI: kpi}
	}
}
