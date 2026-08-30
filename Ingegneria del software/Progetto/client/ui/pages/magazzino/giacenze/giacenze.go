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

// passoRotellina è il numero di righe di cui la lista scorre ad ogni evento della rotellina del mouse.
const passoRotellina = 3

// overheadChrome è il numero di righe fisse occupate ATTORNO al corpo della pagina
// (logo, intestazione, barra dei tab, bordi della finestra e riga vuota interna).
// È stato misurato sul layout reale composto da tabnav.PlacePage: sommandolo alle
// dimensioni del terminale si ricava quante righe di lista possono essere mostrate.
const overheadChrome = 21

// righeFallback è il numero di righe visibili usato finché non arriva il primo WindowSizeMsg.
const righeFallback = 14

// minRigheConKPI è il numero minimo di righe di lista richiesto per mostrare i blocchi KPI:
// nei terminali troppo bassi i KPI vengono nascosti per lasciare spazio alla tabella e al footer.
const minRigheConKPI = 3

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
// 3. Intercetta `tea.WindowSizeMsg` conservando le nuove dimensioni e ricalibrando la finestra visibile.
// 4. Intercetta `tea.MouseWheelMsg` scorrendo la lista nella direzione della rotellina.
// 5. Gestisce le pressioni dei tasti (`tea.KeyPressMsg`):
//   - "up" / "k": sposta il cursore verso l'alto senza superare l'indice 0.
//   - "down" / "j": sposta il cursore verso il basso senza superare il numero di articoli visibili.
//   - "pgup" / "pgdown": scorrono la lista di una pagina; "home" / "end": tornano all'inizio/fine.
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
		m.winW = msg.Width
		m.winH = msg.Height
		m.clampScroll()
		return m, nil

	case tea.MouseWheelMsg:
		switch msg.Button {
		case tea.MouseWheelUp:
			m.scorri(-passoRotellina)
		case tea.MouseWheelDown:
			m.scorri(passoRotellina)
		}
		return m, nil

	case tea.KeyPressMsg:
		articoli := m.articoliFiltrati()
		switch msg.String() {
		case "up", "k":
			if m.cursore > 0 {
				m.cursore--
				m.mostraCursore()
			}
		case "down", "j":
			if m.cursore < len(articoli)-1 {
				m.cursore++
				m.mostraCursore()
			}
		case "pgup":
			if m.cursore > 0 {
				m.cursore -= m.righeVisibili()
				if m.cursore < 0 {
					m.cursore = 0
				}
				m.mostraCursore()
			}
		case "pgdown":
			if m.cursore < len(articoli)-1 {
				m.cursore += m.righeVisibili()
				if max := len(articoli) - 1; m.cursore > max {
					m.cursore = max
				}
				m.mostraCursore()
			}
		case "home":
			m.cursore = 0
			m.mostraCursore()
		case "end":
			if len(articoli) > 0 {
				m.cursore = len(articoli) - 1
			}
			m.mostraCursore()
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

// ── Scorrimento della lista ──────────────────────────────────────────────────

// righeVisibili calcola quante righe di tabella entrano nello spazio residuo dello schermo.
//
// Come è stato scritto il codice:
// Al termine `winH` vengono sottratte le righe fisse esterne al corpo (overheadChrome)
// e quelle dell'intestazione e del footer della pagina, misurate sul rendering effettivo:
// così il risultato resta corretto anche se il layout cambierà in futuro.
func (m Model) righeVisibili() int {
	if m.winH <= 0 || m.report == nil {
		return righeFallback
	}
	disponibili := m.winH - overheadChrome -
		lipgloss.Height(m.renderIntestazione(m.mostraKPI())) -
		lipgloss.Height(m.renderFooter(""))
	if disponibili < 1 {
		disponibili = 1
	}
	return disponibili
}

// mostraKPI indica se i blocchi KPI entrano nello schermo lasciando almeno
// minRigheConKPI righe di lista: nei terminali bassi vengono nascosti per
// privilegiare la visibilità della tabella e del footer.
func (m Model) mostraKPI() bool {
	if m.winH <= 0 || m.report == nil {
		return true
	}
	disponibili := m.winH - overheadChrome -
		lipgloss.Height(m.renderIntestazione(true)) -
		lipgloss.Height(m.renderFooter(""))
	return disponibili >= minRigheConKPI
}

// scorri avanza la finestra di visualizzazione di `n` righe, mantenendo il cursore visibile.
func (m *Model) scorri(n int) {
	m.cursore += n
	articoli := m.articoliFiltrati()
	if m.cursore < 0 {
		m.cursore = 0
	}
	if max := len(articoli) - 1; m.cursore > max {
		m.cursore = max
	}
	m.mostraCursore()
}

// clampScroll riporta l'offset nell'intervallo valido rispetto alla lista e allo schermo.
func (m *Model) clampScroll() {
	articoli := m.articoliFiltrati()
	maxOffset := len(articoli) - m.righeVisibili()
	if maxOffset < 0 {
		maxOffset = 0
	}
	if m.offset > maxOffset {
		m.offset = maxOffset
	}
	if m.offset < 0 {
		m.offset = 0
	}
}

// mostraCursore aggiorna lo scorrimento così che la riga puntata dal cursore resti visibile.
func (m *Model) mostraCursore() {
	vis := m.righeVisibili()
	if m.cursore < m.offset {
		m.offset = m.cursore
	} else if m.cursore >= m.offset+vis {
		m.offset = m.cursore - vis + 1
	}
	m.clampScroll()
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

	// ── Gestione Caricamento ed Errori ──────────────────────────────────────
	if m.loading {
		sb.WriteString(style.TitleStile.Render("MAGAZZINO - GIACENZE & INVENTARIO"))
		sb.WriteString("\n\n")
		sb.WriteString("Caricamento report giacenze in corso...\n")
		return tea.NewView(sb.String())
	}
	if m.err != nil {
		sb.WriteString(style.TitleStile.Render("MAGAZZINO - GIACENZE & INVENTARIO"))
		sb.WriteString("\n\n")
		sb.WriteString(style.ErrorStile.Render("Errore: " + m.err.Error()))
		sb.WriteString("\n")
		return tea.NewView(sb.String())
	}
	if m.report == nil {
		sb.WriteString(style.TitleStile.Render("MAGAZZINO - GIACENZE & INVENTARIO"))
		sb.WriteString("\n\n")
		sb.WriteString("Nessun dato di inventario disponibile.\n")
		return tea.NewView(sb.String())
	}

	// ── Intestazione fissa ───────────────────────────────────────────────────
	sb.WriteString(m.renderIntestazione(m.mostraKPI()))

	// ── Finestra visibile della tabella ──────────────────────────────────────
	articoli := m.articoliFiltrati()

	if len(articoli) == 0 {
		sb.WriteString("Nessun articolo trovato con i filtri correnti.\n")
		sb.WriteString("\n\n")
		sb.WriteString(m.renderFooter(""))
		return tea.NewView(sb.String())
	}

	m.clampScroll()
	inizio := m.offset
	fine := inizio + m.righeVisibili()
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

	// ── Footer con indicatore di scorrimento ────────────────────────────────
	sb.WriteString(m.renderFooter(m.indicatoreArticoli()))

	return tea.NewView(sb.String())
}

// renderIntestazione compone la parte fissa superiore della pagina: titolo,
// blocchi KPI aggregati (se `conKPI`, vedi mostraKPI), indicatore del filtro
// sotto-soglia e testata della tabella.
func (m Model) renderIntestazione(conKPI bool) string {
	var sb strings.Builder

	// Il MarginBottom di TitleStile fornisce già la riga vuota sotto il titolo:
	// qui basta un solo "\n" esplicito per non sprecare righe di terminale.
	sb.WriteString(style.TitleStile.Render("MAGAZZINO - GIACENZE & INVENTARIO"))
	sb.WriteString("\n")

	// ── Blocchi KPI Superiori ───────────────────────────────────────────────
	if conKPI {
		sottoSogliaCount := len(m.report.ArticoliSottoSoglia)
		kpiTotal := style.BoxStile.Render(fmt.Sprintf("TOTALE PEZZI\n\n%d", m.report.TotaleQuantita))
		kpiValore := style.BoxStile.Render(fmt.Sprintf("VALORE TOTALE\n\n€ %.2f", m.report.ValoreTotale))
		kpiSottoSoglia := style.BoxStile.Render(fmt.Sprintf("SOTTO SOGLIA\n\n%d", sottoSogliaCount))

		rigaKPI := lipgloss.JoinHorizontal(lipgloss.Top, kpiTotal, "   ", kpiValore, "   ", kpiSottoSoglia)
		sb.WriteString(rigaKPI)
		sb.WriteString("\n")
	}

	// ── Indicatore Filtro Sotto Soglia ─────────────────────────────────────
	statoFiltro := "[s] Filtro Sotto-Soglia: OFF"
	if m.filtroSottoSoglia {
		statoFiltro = style.SuccessStile.Render("[s] Filtro Sotto-Soglia: ON")
	}
	sb.WriteString(fmt.Sprintf("%s", statoFiltro))
	sb.WriteString("\n")

	// ── Testata Tabella ─────────────────────────────────────────────────────
	headerStr := fmt.Sprintf("  %-10s %-30s %-25s %-20s %s",
		"CODICE", "PRODOTTO", "CATEGORIA", "GIACENZA", "STATO")
	sb.WriteString(style.FreeHeaderStile.Render(headerStr))
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("─", 100))
	sb.WriteString("\n")

	return sb.String()
}

// renderFooter compone le istruzioni di tastiera e, se la lista eccede lo schermo,
// l'indicatore delle righe correntemente visibili. Il blocco ha SEMPRE due righe
// (scorciatoie + indicatore anche vuoto): così il calcolo delle righe visibili è
// stabile e l'ultima riga di lista non viene mai spinta fuori schermo.
func (m Model) renderFooter(indicatore string) string {
	scorciatoie := "[↑/↓ o k/j] Naviga  |  [pgup/pgdown o rotellina] Scorri  |  [s] Filtra Sotto Soglia  |  [r] Ricarica  |  [ctrl+l] Logout"
	return style.FooterStile.Render(scorciatoie + "\n" + indicatore)
}

// indicatoreArticoli descrive, se la lista eccede lo schermo, le righe correntemente visibili.
func (m Model) indicatoreArticoli() string {
	articoli := m.articoliFiltrati()
	vis := m.righeVisibili()
	if len(articoli) <= vis {
		return ""
	}
	primo := m.offset + 1
	ultimo := m.offset + vis
	if ultimo > len(articoli) {
		ultimo = len(articoli)
	}
	return fmt.Sprintf("Articoli %d–%d di %d  (pgup/pgdown o rotellina per scorrere)", primo, ultimo, len(articoli))
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
