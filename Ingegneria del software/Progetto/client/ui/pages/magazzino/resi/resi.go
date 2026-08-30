// Package resi gestisce il rendering e la logica di navigazione per il registro dei resi e ripristino merce.
package resi

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
	style "github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/styles"
)

// Init avvia l'inizializzazione del componente richiedendo l'elenco dei resi al backend.
func (m Model) Init() tea.Cmd {
	return m.loadResiCmd()
}

// Update elabora i messaggi ricevuti e gestisce la navigazione tra le pratiche di reso.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ResiLoadedMsg:
		m.loading = false
		m.err = nil
		m.resi = msg.Resi
		m.cursore = 0
		m.offset = 0
		return m, nil

	case ErroreMsg:
		m.loading = false
		m.err = msg.Err
		return m, nil

	case ResoAggiornatoMsg:
		m.loading = false
		m.err = nil
		m.modo = modoLista
		return m, m.loadResiCmd()

	case tea.WindowSizeMsg:
		return m, nil

	case tea.KeyPressMsg:
		if m.modo == modoModificaStato {
			return m.handleEditInput(msg)
		}

		switch msg.String() {
		case "up", "k":
			if m.cursore > 0 {
				m.cursore--
				m.offset = style.PaginaCalcola(m.cursore, m.offset, len(m.resi))
			}
		case "down", "j":
			if m.cursore < len(m.resi)-1 {
				m.cursore++
				m.offset = style.PaginaCalcola(m.cursore, m.offset, len(m.resi))
			}
		case "pgup":
			m.cursore -= style.RighePerPagina
			if m.cursore < 0 {
				m.cursore = 0
			}
			m.offset = style.PaginaCalcola(m.cursore, m.offset, len(m.resi))
		case "pgdown":
			m.cursore += style.RighePerPagina
			if max := len(m.resi) - 1; m.cursore > max {
				m.cursore = max
			}
			if m.cursore < 0 {
				m.cursore = 0
			}
			m.offset = style.PaginaCalcola(m.cursore, m.offset, len(m.resi))
		case "enter", "space":
			if len(m.resi) > 0 {
				id := m.resi[m.cursore].ID
				m.dettagliEspansi[id] = !m.dettagliEspansi[id]
			}
		case "ctrl+e":
			if len(m.resi) > 0 {
				m.modo = modoModificaStato
				m.err = nil
				m.idxStato = indexStatoReso(string(m.resi[m.cursore].Stato))
			}
		case "r":
			m.loading = true
			return m, m.loadResiCmd()
		}
	}

	return m, nil
}

// View compone l'interfaccia utente mostrando i KPI, la tabella dei resi ed il dettaglio espanso.
func (m Model) View() tea.View {
	var sb strings.Builder

	// ── Titolo Sezione ──────────────────────────────────────────────────────
	sb.WriteString(style.TitleStile.Render("GESTIONE RESI E RIPRISTINO GIACENZE"))
	sb.WriteString("\n\n")

	// ── Gestione Caricamento ed Errori ──────────────────────────────────────
	if m.loading {
		sb.WriteString("Caricamento registro resi in corso...\n\n")
		return tea.NewView(sb.String())
	}
	if m.err != nil {
		sb.WriteString(style.ErrorStile.Render("Errore: " + m.err.Error()))
		sb.WriteString("\n\n")
		if m.modo == modoModificaStato {
			sb.WriteString(style.FooterStile.Render("[←/→] Seleziona Stato  |  [enter] Conferma  |  [esc] Annulla"))
		} else {
			sb.WriteString(style.FooterStile.Render("[esc] Torna alla lista resi"))
		}
		return tea.NewView(sb.String())
	}

	// ── Modalita modifica stato ─────────────────────────────────────────────
	if m.modo == modoModificaStato {
		sb.WriteString(m.renderEditStatoForm())
		return tea.NewView(sb.String())
	}

	// ── Blocchi KPI Superiori (Impaginati con Lipgloss) ────────────────────
	var totalePezziResi int
	for _, r := range m.resi {
		totalePezziResi += r.Quantita
	}

	kpiTot := style.BoxStile.Render(fmt.Sprintf("PRATICHE RESO\n\n%d", len(m.resi)))
	kpiPezzi := style.BoxStile.Render(fmt.Sprintf("ARTICOLI RESTITUITI\n\n%d", totalePezziResi))

	rigaKPI := lipgloss.JoinHorizontal(lipgloss.Top, kpiTot, "   ", kpiPezzi)
	sb.WriteString(rigaKPI)
	sb.WriteString("\n\n")

	// ── Tabella Resi ────────────────────────────────────────────────────────
	if len(m.resi) == 0 {
		sb.WriteString("Nessuna pratica di reso registrata.\n")
	} else {
		headerStr := fmt.Sprintf("   %-10s %-20s %-15s %-10s %s",
			"ID RESO", "ORDINE ORIGINE", "DATA", "QUANTITÀ", "MOTIVO")
		sb.WriteString(style.FreeHeaderStile.Render(headerStr))
		sb.WriteString("\n")
		sb.WriteString(strings.Repeat("─", 100))
		sb.WriteString("\n")

		inizio := m.offset
		fine := inizio + style.RighePerPagina
		if fine > len(m.resi) {
			fine = len(m.resi)
		}

		for i, r := range m.resi[inizio:fine] {
			assoluto := inizio + i
			dataStr := r.Data.Format("02/01/2006")
			ordUscitaStr := fmt.Sprintf("Ordine #%d", r.OrdineUscitaID)
			qtaStr := fmt.Sprintf("%d pcs", r.Quantita)

			rigaTesto := fmt.Sprintf("%-10d %-20s %-15s %-10s %s",
				r.ID,
				truncaStringa(ordUscitaStr, 20),
				dataStr,
				qtaStr,
				truncaStringa(r.Motivo, 30),
			)

			indicatore := "  "
			if assoluto == m.cursore {
				indicatore = style.SuccessStile.Render("> ")
			}

			espanso := m.dettagliEspansi[r.ID]
			simboloDettaglio := " [+] "
			if espanso {
				simboloDettaglio = " [-] "
			}

			sb.WriteString(indicatore)
			sb.WriteString(simboloDettaglio)
			sb.WriteString(rigaTesto)
			sb.WriteString("\n")

			if espanso {
				sb.WriteString(renderRigheReso(r))
			}
		}

		pagine := (len(m.resi) + style.RighePerPagina - 1) / style.RighePerPagina
		paginaCorr := m.offset/style.RighePerPagina + 1
		sb.WriteString(fmt.Sprintf("\nPagina %d/%d", paginaCorr, pagine))
	}

	// ── Footer ───────────────────────────────────────────────────────────────
	sb.WriteString("\n")
	sb.WriteString(style.FooterStile.Render("[↑/↓ o k/j] Naviga  |  [pgup/pgdown] Pagina  |  [enter/space] Espandi/Comprimi Dettaglio  |  [ctrl+e] Modifica Stato  |  [r] Ricarica  |  [ctrl+l] Logout"))

	return tea.NewView(sb.String())
}

// handleEditInput gestisce gli eventi tastiera in modalita modifica stato.
func (m Model) handleEditInput(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.modo = modoLista
		m.err = nil
		return m, nil
	case "left":
		m.idxStato = (m.idxStato - 1 + len(statiReso)) % len(statiReso)
		return m, nil
	case "right":
		m.idxStato = (m.idxStato + 1) % len(statiReso)
		return m, nil
	case "enter":
		return m.salvaStatoReso()
	}
	return m, nil
}

// renderEditStatoForm renderizza il form di modifica dello stato del reso.
func (m Model) renderEditStatoForm() string {
	if len(m.resi) == 0 || m.cursore >= len(m.resi) {
		return style.FooterStile.Render("[esc] Annulla")
	}
	r := m.resi[m.cursore]

	statoSelezionato := statiReso[m.idxStato]
	statoVisualizzato := style.FocusStile.Render("< " + statoSelezionato + " >")

	content := strings.Join([]string{
		"MODIFICA STATO RESO",
		"",
		fmt.Sprintf("ID Reso:         %d", r.ID),
		fmt.Sprintf("Ordine Origine:  Ordine #%d", r.OrdineUscitaID),
		fmt.Sprintf("Motivo:          %s", r.Motivo),
		"",
		"Nuovo Stato:     " + statoVisualizzato,
	}, "\n")

	return style.LoginBoxStile.Render(content) + "\n\n" +
		style.FooterStile.Render("[←/→] Seleziona Stato  |  [enter] Conferma  |  [esc] Annulla")
}

// salvaStatoReso invia l'aggiornamento dello stato del reso al server.
func (m Model) salvaStatoReso() (tea.Model, tea.Cmd) {
	if len(m.resi) == 0 || m.cursore >= len(m.resi) {
		return m, nil
	}
	r := m.resi[m.cursore]
	stato := statiReso[m.idxStato]

	m.loading = true
	m.err = nil
	return m, m.aggiornaStatoResoCmd(r.ID, stato)
}

// aggiornaStatoResoCmd crea un comando asincrono per aggiornare lo stato del reso.
func (m Model) aggiornaStatoResoCmd(id int, stato string) tea.Cmd {
	return func() tea.Msg {
		if err := m.Client.AggiornaStatoReso(id, stato); err != nil {
			return ErroreMsg{Err: err}
		}
		return ResoAggiornatoMsg{}
	}
}

// renderRigheReso costruisce la vista dettagliata espansa per un singolo reso.
func renderRigheReso(r api.Reso) string {
	var sb strings.Builder
	indent := "       "
	sb.WriteString(indent)
	sb.WriteString(style.HeaderStile.Render(fmt.Sprintf("┌─ DETTAGLI PRATICA RESO #%d", r.ID)))
	sb.WriteString("\n")

	sb.WriteString(indent)
	sb.WriteString(fmt.Sprintf("│  %-18s %-14s %-10s %s\n",
		"ORDINE ORIGINE", "PRODOTTO ID", "QUANTITÀ", "OPERATORE ID"))
	sb.WriteString(indent)
	sb.WriteString("├")
	sb.WriteString(strings.Repeat("─", 65))
	sb.WriteString("\n")

	sb.WriteString(indent)
	sb.WriteString(fmt.Sprintf("│  Ordine #%-11d Prodotto #%-7d %-10d User #%d\n",
		r.OrdineUscitaID,
		r.ProdottoID,
		r.Quantita,
		r.CreatedByID,
	))

	sb.WriteString(indent)
	sb.WriteString("├")
	sb.WriteString(strings.Repeat("─", 65))
	sb.WriteString("\n")
	sb.WriteString(indent)
	sb.WriteString("│  Motivo del reso: ")
	sb.WriteString(r.Motivo)
	sb.WriteString("\n")

	sb.WriteString(indent)
	sb.WriteString("└")
	sb.WriteString(strings.Repeat("─", 65))
	sb.WriteString("\n")

	return sb.String()
}

// loadResiCmd esegue in modo asincrono la chiamata al client API per recuperare l'elenco dei resi.
func (m Model) loadResiCmd() tea.Cmd {
	return func() tea.Msg {
		resi, err := m.Client.GetResi()
		if err != nil {
			return ErroreMsg{Err: err}
		}
		return ResiLoadedMsg{Resi: resi}
	}
}

// truncaStringa tronca la stringa alla lunghezza massima `maxLen` se necessario.
func truncaStringa(s string, maxLen int) string {
	if len(s) > maxLen {
		if maxLen > 3 {
			return s[:maxLen-3] + "..."
		}
		return s[:maxLen]
	}
	return s
}
