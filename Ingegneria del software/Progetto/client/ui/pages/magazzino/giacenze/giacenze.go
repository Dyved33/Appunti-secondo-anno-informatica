// Package giacenze gestisce le logiche di rendering e la gestione eventi per la tabella delle giacenze di magazzino.
package giacenze

import (
	"fmt"
	"strings"

	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
	style "github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/styles"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// Init avvia l'inizializzazione del modello Giacenze ed il caricamento del report dal backend.
//
// Come è stato scritto il codice:
// Restituisce il comando asincrono `loadGiacenzeCmd()` per recuperare i dati dell'inventario tramite API.
func (m Model) Init() tea.Cmd {
	return m.loadGiacenzeCmd()
}

// Update gestisce l'elaborazione dei messaggi e degli eventi di navigazione da tastiera della schermata Giacenze.
//
// Come è stato scritto il codice:
// 1. Intercetta `GiacenzeLoadedMsg` impostando `loading = false`, azzerando l'errore e salvando il report ricevuto.
// 2. Intercetta `ErroreMsg` salvando l'errore ed arrestando l'indicatore di caricamento.
// 3. Gestisce le pressioni dei tasti (`tea.KeyPressMsg`):
//   - "up" / "k": sposta il cursore verso l'alto senza superare l'indice 0.
//   - "down" / "j": sposta il cursore verso il basso senza superare il numero di articoli visibili.
//   - "s": inverte il flag `filtroSottoSoglia` e resetta il cursore a 0.
//   - "r": imposta `loading = true` e lancia nuovamente il comando di ricaricamento `loadGiacenzeCmd()`.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case GiacenzeLoadedMsg:
		m.loading = false
		m.err = nil
		m.report = msg.Report
		m.cursore = 0
		m.offset = 0
		return m, nil

	case ErroreMsg:
		m.loading = false
		m.err = msg.Err
		return m, nil

	case tea.WindowSizeMsg:
		return m, nil

	case tea.KeyPressMsg:
		articoli := m.articoliFiltrati()
		switch msg.String() {
		case "up", "k":
			if m.cursore > 0 {
				m.cursore--
				m.offset = style.PaginaCalcola(m.cursore, m.offset, len(articoli))
			}
		case "down", "j":
			if m.cursore < len(articoli)-1 {
				m.cursore++
				m.offset = style.PaginaCalcola(m.cursore, m.offset, len(articoli))
			}
		case "pgup":
			m.cursore -= style.RighePerPagina
			if m.cursore < 0 {
				m.cursore = 0
			}
			m.offset = style.PaginaCalcola(m.cursore, m.offset, len(articoli))
		case "pgdown":
			m.cursore += style.RighePerPagina
			if max := len(articoli) - 1; m.cursore > max {
				m.cursore = max
			}
			if m.cursore < 0 {
				m.cursore = 0
			}
			m.offset = style.PaginaCalcola(m.cursore, m.offset, len(articoli))
		case "s":
			m.filtroSottoSoglia = !m.filtroSottoSoglia
			m.cursore = 0
			m.offset = 0
		case "r":
			m.loading = true
			return m, m.loadGiacenzeCmd()
		}
	}

	return m, nil
}

// View costruisce ed impagina l'interfaccia utente visibile per la schermata delle Giacenze.
//
// Come è stato scritto il codice:
//  1. Mostra messaggi di stato se in caricamento, errore o in assenza di dati.
//  2. Disegna l'intestazione con 3 box KPI (Totale Pezzi, Valore Totale in €, Articoli Sotto Soglia).
//  3. Mostra l'indicatore dello stato del filtro Sotto-Soglia (`[s]`).
//  4. Formatta la tabella degli articoli con intestazioni di colonna, troncamento dei testi lunghi,
//     evidenziazione della riga su cui si trova il cursore (`>`) e badge per articoli in errore/sotto soglia.
//  5. Inserisce una barra di footer con le scorciatoie di navigazione.
func (m Model) View() tea.View {
	var sb strings.Builder

	// ── Titolo Sezione ──────────────────────────────────────────────────────
	sb.WriteString(style.TitleStile.Render("MAGAZZINO - GIACENZE & INVENTARIO"))
	sb.WriteString("\n\n")

	// ── Gestione Caricamento ed Errori ──────────────────────────────────────
	if m.loading {
		sb.WriteString("Caricamento report giacenze in corso...\n")
		return tea.NewView(sb.String())
	}
	if m.err != nil {
		sb.WriteString(style.ErrorStile.Render("Errore: " + m.err.Error()))
		sb.WriteString("\n")
		return tea.NewView(sb.String())
	}
	if m.report == nil {
		sb.WriteString("Nessun dato di inventario disponibile.\n")
		return tea.NewView(sb.String())
	}

	// ── Blocchi KPI Superiori (Impaginati con Lipgloss) ────────────────────
	sottoSogliaCount := len(m.report.ArticoliSottoSoglia)
	kpiTotal := style.BoxStile.Render(fmt.Sprintf("TOTALE PEZZI\n\n%d", m.report.TotaleQuantita))
	kpiValore := style.BoxStile.Render(fmt.Sprintf("VALORE TOTALE\n\n€ %.2f", m.report.ValoreTotale))
	kpiSottoSoglia := style.BoxStile.Render(fmt.Sprintf("SOTTO SOGLIA\n\n%d", sottoSogliaCount))

	rigaKPI := lipgloss.JoinHorizontal(lipgloss.Top, kpiTotal, "   ", kpiValore, "   ", kpiSottoSoglia)
	sb.WriteString(rigaKPI)
	sb.WriteString("\n\n")

	// ── Indicatore Filtro Sotto Soglia ─────────────────────────────────────
	statoFiltro := "[s] Filtro Sotto-Soglia: OFF"
	if m.filtroSottoSoglia {
		statoFiltro = style.SuccessStile.Render("[s] Filtro Sotto-Soglia: ON")
	}
	sb.WriteString(fmt.Sprintf("%s\n\n", statoFiltro))

	// ── Tabella Giacenze ────────────────────────────────────────────────────
	articoli := m.articoliFiltrati()
	if len(articoli) == 0 {
		sb.WriteString("Nessun articolo trovato con i filtri correnti.\n")
	} else {
		headerStr := fmt.Sprintf("  %-10s %-30s %-25s %-20s %s",
			"CODICE", "PRODOTTO", "CATEGORIA", "GIACENZA", "STATO")
		sb.WriteString(style.FreeHeaderStile.Render(headerStr))
		sb.WriteString("\n")
		sb.WriteString(strings.Repeat("─", 100))
		sb.WriteString("\n")

		inizio := m.offset
		fine := inizio + style.RighePerPagina
		if fine > len(articoli) {
			fine = len(articoli)
		}

		for i, a := range articoli[inizio:fine] {
			assoluto := inizio + i
			giacenzaStr := fmt.Sprintf("%d (min %d)", a.Quantita, a.SogliaMin)

			statoBadge := style.SuccessStile.Render("OK")
			if a.SottoSoglia {
				statoBadge = style.ErrorStile.Render("SOTTO SOGLIA")
			}

			rigaTesto := fmt.Sprintf("%-10s %-30s %-25s %-20s %s",
				troncaStringa(a.Codice, 10),
				troncaStringa(a.Nome, 30),
				troncaStringa(a.Categoria, 25),
				troncaStringa(giacenzaStr, 20),
				statoBadge,
			)

			if assoluto == m.cursore {
				sb.WriteString(style.SuccessStile.Render("> "))
				sb.WriteString(rigaTesto)
				sb.WriteString("\n")
			} else {
				sb.WriteString("  ")
				sb.WriteString(rigaTesto)
				sb.WriteString("\n")
			}
		}

		pagine := (len(articoli) + style.RighePerPagina - 1) / style.RighePerPagina
		paginaCorr := m.offset/style.RighePerPagina + 1
		sb.WriteString(fmt.Sprintf("\nPagina %d/%d", paginaCorr, pagine))
	}

	// ── Footer ───────────────────────────────────────────────────────────────
	sb.WriteString("\n")
	sb.WriteString(style.FooterStile.Render("[↑/↓ o k/j] Naviga  |  [pgup/pgdown] Pagina  |  [s] Filtra Sotto Soglia  |  [r] Ricarica  |  [ctrl+l] Logout"))

	return tea.NewView(sb.String())
}

// articoliFiltrati restituisce la lista degli articoli applicando l'eventuale filtro sotto-soglia.
//
// Come è stato scritto il codice:
// Se `m.filtroSottoSoglia` è vero, restituisce la slice `report.ArticoliSottoSoglia`; altrimenti l'elenco completo `report.Articoli`.
func (m Model) articoliFiltrati() []api.InventoryItem {
	if m.report == nil {
		return nil
	}
	if m.filtroSottoSoglia {
		return m.report.ArticoliSottoSoglia
	}
	return m.report.Articoli
}

// loadGiacenzeCmd invoca l'API in modo asincrono per recuperare il report delle giacenze.
//
// Come è stato scritto il codice:
// Restituisce una closure tea.Cmd che invoca `m.Client.GetReportGiacenze()`. In caso di successo
// emette un `GiacenzeLoadedMsg`, in caso di errore emette un `ErroreMsg`.
func (m Model) loadGiacenzeCmd() tea.Cmd {
	return func() tea.Msg {
		report, err := m.Client.GetReportGiacenze()
		if err != nil {
			return ErroreMsg{Err: err}
		}
		return GiacenzeLoadedMsg{Report: report}
	}
}

// troncaStringa tronca la stringa alla lunghezza massima `maxLen` aggiungendo i puntini di sospensione se necessario.
func troncaStringa(s string, maxLen int) string {
	if len(s) > maxLen {
		if maxLen > 3 {
			return s[:maxLen-3] + "..."
		}
		return s[:maxLen]
	}
	return s
}
