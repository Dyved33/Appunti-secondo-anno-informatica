// Package scanner gestisce la logica di digitazione, ricerca asincrona e rendering dello scanner barcode.
package scanner

import (
	"fmt"
	"strings"

	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
	style "github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/styles"

	tea "charm.land/bubbletea/v2"
)

// Init inizializza il sotto-modello Scanner.
//
// Come è stato scritto il codice:
// Non richiede chiamate di rete in fase di avvio iniziale, pertanto restituisce `nil`.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update gestisce la digitazione del codice a barre, la ricerca via API e l'aggiornamento dello storico scansioni.
//
// Come è stato scritto il codice:
//  1. Su `ProdottoScansionatoMsg`: azzera `loading` ed `err`, memorizza in `prodottoTrovato` il puntatore al prodotto riscontrato,
//     lo inserisce in testa a `storicoScansioni` (mantenendo massimo 10 elementi) e svuota il buffer `input`.
//  2. Su `ErroreMsg`: imposta `loading = false`, memorizza l'errore e resetta `prodottoTrovato` a nil.
//  3. Su `tea.KeyPressMsg`:
//     - "enter": se `input` non è vuoto e non è in caricamento, imposta `loading = true` e lancia `scansionaCmd`.
//     - "backspace": rimuove l'ultimo carattere dal buffer `input`.
//     - "esc": resetta il buffer `input`, il prodotto ed eventuali errori.
//     - caratteri stampabili: appende il carattere al buffer `input`.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ProdottoScansionatoMsg:
		m.loading = false
		m.err = nil
		p := msg.Prodotto
		m.prodottoTrovato = &p

		// Aggiunge allo storico in testa evitando duplicati consecutivi
		m.storicoScansioni = append([]api.Prodotto{p}, m.storicoScansioni...)
		if len(m.storicoScansioni) > 10 {
			m.storicoScansioni = m.storicoScansioni[:10]
		}
		m.input = ""
		return m, nil

	case ErroreMsg:
		m.loading = false
		m.err = msg.Err
		m.prodottoTrovato = nil
		return m, nil

	case tea.WindowSizeMsg:
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			if strings.TrimSpace(m.input) != "" && !m.loading {
				m.loading = true
				m.err = nil
				return m, m.scansionaCmd(strings.TrimSpace(m.input))
			}
		case "backspace":
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
		case "esc":
			m.input = ""
			m.prodottoTrovato = nil
			m.err = nil
		default:
			// Se viene digitato un singolo carattere stampabile
			keyStr := msg.String()
			if len(keyStr) == 1 && keyStr >= " " && keyStr <= "~" {
				m.input += keyStr
			}
		}
	}

	return m, nil
}

// View disegna la casella di testo dell'input barcode, la scheda dettagli del prodotto identificato ed il registro storico.
//
// Come è stato scritto il codice:
// 1. Renderizza l'header "SCANNER BARCODE & RICERCA PRODOTTO".
// 2. Mostra il prompt dell'input con il cursore a blocco (`█`).
// 3. Se `loading` è vero o `err` non nil, mostra lo stato relativo.
// 4. Se `prodottoTrovato != nil`, genera una card dettagliata con codice, nome, categoria, prezzi e stato attivo.
// 5. Se `storicoScansioni` contiene elementi, mostra l'elenco delle ultime 10 scansioni effettuate.
func (m Model) View() tea.View {
	var sb strings.Builder

	// ── Titolo Sezione ──────────────────────────────────────────────────────
	sb.WriteString(style.TitleStile.Render("SCANNER BARCODE & RICERCA PRODOTTO"))
	sb.WriteString("\n\n")

	// ── Prompt Input Scansione ─────────────────────────────────────────────
	promptText := fmt.Sprintf("Codice a Barre / Prodotto: > %s█", m.input)
	boxInput := style.LoginBoxStile.Render(promptText)
	sb.WriteString(boxInput)
	sb.WriteString("\n\n")

	// ── Stato di Caricamento ed Errori ──────────────────────────────────────
	if m.loading {
		sb.WriteString(style.InfoStile.Render("Ricerca prodotto in corso per codice..."))
		sb.WriteString("\n\n")
	} else if m.err != nil {
		sb.WriteString(style.ErrorStile.Render("Errore: " + m.err.Error()))
		sb.WriteString("\n\n")
	}

	// ── Scheda Dettaglio Prodotto Identificato ──────────────────────────────
	if m.prodottoTrovato != nil {
		p := m.prodottoTrovato
		statoStr := style.SuccessStile.Render("ATTIVO")
		if !p.Attivo {
			statoStr = style.ErrorStile.Render("DISATTIVATO")
		}

		cardContent := fmt.Sprintf(
			"PRODOTTO IDENTIFICATO\n\n"+
				"Codice:          %s\n"+
				"Nome:            %s\n"+
				"Categoria:       %s\n"+
				"Prezzo Acquisto: € %.2f\n"+
				"Prezzo Vendita:  € %.2f\n"+
				"Stato:           %s",
			p.Codice, p.Nome, p.Categoria, p.PrezzoAcquisto, p.PrezzoVendita, statoStr,
		)
		sb.WriteString(style.LoginBoxStile.Render(cardContent))
		sb.WriteString("\n\n")
	}

	// ── Registro Ultime Scansioni Effettuate ───────────────────────────────
	if len(m.storicoScansioni) > 0 {
		sb.WriteString(style.SubTitleStile.Render("Ultime Scansioni Effettuate:"))
		sb.WriteString("\n\n")

		headerStr := fmt.Sprintf("   %-10s %-24s %-16s %s",
			"CODICE", "PRODOTTO", "CATEGORIA", "PREZZO VENDITA")
		sb.WriteString(style.HeaderStile.Render(headerStr))
		sb.WriteString("\n")
		sb.WriteString(strings.Repeat("─", 70))
		sb.WriteString("\n")

		for _, s := range m.storicoScansioni {
			rigaTesto := fmt.Sprintf("   %-10s %-24s %-16s € %.2f",
				troncaStringa(s.Codice, 10),
				troncaStringa(s.Nome, 24),
				troncaStringa(s.Categoria, 16),
				s.PrezzoVendita,
			)
			sb.WriteString(rigaTesto)
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	// ── Footer ───────────────────────────────────────────────────────────────
	sb.WriteString("\n\n")
	sb.WriteString(style.FooterStile.Render("[Digita Codice]  |  [enter] Scansiona / Cerca  |  [esc] Pulisci Field  |  [ctrl+l] Logout"))

	return tea.NewView(sb.String())
}

// scansionaCmd esegue in modo asincrono la ricerca del prodotto tramite il client API.
//
// Come è stato scritto il codice:
// Restituisce una closure tea.Cmd che invoca prima `m.Client.GetProductByBarcode(codice)` (ricerca per barcode EAN-13);
// in caso di esito negativo riprova con `m.Client.GetProductByCode(codice)` (ricerca per codice prodotto).
// Emette `ProdottoScansionatoMsg` oppure `ErroreMsg`.
func (m Model) scansionaCmd(codice string) tea.Cmd {
	return func() tea.Msg {
		prodotto, err := m.Client.GetProductByBarcode(codice)
		if err != nil {
			prodotto, err = m.Client.GetProductByCode(codice)
			if err != nil {
				return ErroreMsg{Err: err}
			}
		}
		return ProdottoScansionatoMsg{Prodotto: prodotto}
	}
}

// troncaStringa tronca la stringa alla lunghezza massima maxLen.
func troncaStringa(s string, maxLen int) string {
	if len(s) > maxLen {
		if maxLen > 3 {
			return s[:maxLen-3] + "..."
		}
		return s[:maxLen]
	}
	return s
}
