package ordini

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
	style "github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/styles"
)

// View calcola e restituisce la rappresentazione visiva (layout TUI) dello stato del modello
func (m Model) View() tea.View {
	var sb strings.Builder

	// Titolo
	sb.WriteString(style.TitleStile.Render("GESTIONE ORDINI"))
	sb.WriteString("\n\n")

	if m.loading {
		sb.WriteString("Caricamento ordini in corso...\n")
		return tea.NewView(sb.String())
	}

	if m.err != nil {
		sb.WriteString(style.ErrorStile.Render("Errore: " + m.err.Error()))
		sb.WriteString("\n\n")
		return tea.NewView(sb.String())
	}

	if m.creazioneEntrata {
		v := tea.NewView(m.renderFormCreazioneEntrata())
		v.Cursor = m.currentFormCursor()
		return v
	}
	if m.creazioneUscita {
		v := tea.NewView(m.renderFormCreazioneUscita())
		v.Cursor = m.currentFormCursor()
		return v
	}
	if m.creazioneReso {
		v := tea.NewView(m.renderFormCreazioneReso())
		v.Cursor = m.currentFormCursor()
		return v
	}
	if m.confermaEliminazione {
		var tipoStr string
		switch m.sezione {
		case ORDINI_ENTRATA:
			tipoStr = fmt.Sprintf("l'ordine in entrata #%d", m.idDaEliminare)
		case ORDINI_USCITA:
			tipoStr = fmt.Sprintf("l'ordine in uscita #%d", m.idDaEliminare)
		case RESI:
			tipoStr = fmt.Sprintf("il reso #%d", m.idDaEliminare)
		}

		sb.WriteString(style.TitleStile.Render("CONFERMA ELIMINAZIONE"))
		sb.WriteString("\n\n")
		sb.WriteString(style.ErrorStile.Render(fmt.Sprintf("Sei sicuro di voler eliminare %s?", tipoStr)))
		sb.WriteString("\n\n")
		sb.WriteString(style.FooterStile.Render("[Y/Enter] Conferma  |  [N/Esc] Annulla"))
		return tea.NewView(sb.String())
	}

	// 1. KPI Superiori (In base alla sezione selezionata)
	sb.WriteString(m.renderKPI())
	sb.WriteString("\n\n")

	// 2. Tab Header per selezionare Entrata / Uscita / Resi
	sb.WriteString(m.renderTabHeader())
	sb.WriteString("\n\n")

	// 3. Render Tabella attiva
	switch m.sezione {
	case ORDINI_ENTRATA:
		sb.WriteString(m.renderTabellaEntrata())
	case ORDINI_USCITA:
		sb.WriteString(m.renderTabellaUscita())
	case RESI:
		sb.WriteString(m.renderTabellaResi())
	}

	// 4. Footer con le istruzioni unificate
	sb.WriteString("\n")
	sb.WriteString(style.FooterStile.Render("[tab] Cambia Tabella  |  [↑/↓ o k/j] Naviga  |  [pgup/pgdown o rotellina] Scorri  |  [ctrl+n] Nuovo  |  [ctrl+d] Elimina  |  [enter/space] Dettaglio  |  [r] Ricarica"))

	return tea.NewView(sb.String())
}

// renderKPI calcola e stilizza i riquadri delle metriche chiave (totali, importi, conteggi) per la sezione attiva.
func (m Model) renderKPI() string {
	var confermati, evasi, annullati, totale int
	var inAttesa, accettati, rifiutati int

	switch m.sezione {
	case ORDINI_ENTRATA:
		totale = len(m.ordiniEntrata)
		var totaleSpesa float64
		var totalePezzi int
		fornitoriMap := make(map[string]bool)

		for _, o := range m.ordiniEntrata {
			if o.Fornitore != "" {
				fornitoriMap[o.Fornitore] = true
			}
			for _, r := range o.Righe {
				totalePezzi += r.Quantita
				totaleSpesa += float64(r.Quantita) * r.PrezzoUnitario
			}
		}

		kpiTot := style.BoxStile.Render(fmt.Sprintf("TOTALE ORDINI\n\n%d", totale))
		kpiSpesa := style.BoxStile.Render(fmt.Sprintf("VALORE TOTALE\n\n€ %.2f", totaleSpesa))
		kpiPezzi := style.BoxStile.Render(fmt.Sprintf("TOTALE ARTICOLI\n\n%d", totalePezzi))
		kpiFornitori := style.BoxStile.Render(fmt.Sprintf("FORNITORI UNICI\n\n%d", len(fornitoriMap)))

		return lipgloss.JoinHorizontal(lipgloss.Top, kpiTot, "   ", kpiSpesa, "   ", kpiPezzi, "   ", kpiFornitori)

	case ORDINI_USCITA:
		totale = len(m.ordiniUscita)
		for _, o := range m.ordiniUscita {
			switch o.Stato {
			case api.StatoConfermato:
				confermati++
			case api.StatoEvaso:
				evasi++
			case api.StatoAnnullato:
				annullati++
			}
		}
		kpiTot := style.BoxStile.Render(fmt.Sprintf("TOTALE ORDINI\n\n%d", totale))
		kpiConf := style.BoxStile.Render(fmt.Sprintf("CONFERMATI\n\n%d", confermati))
		kpiEvas := style.BoxStile.Render(fmt.Sprintf("EVASI\n\n%d", evasi))
		kpiAnn := style.BoxStile.Render(fmt.Sprintf("ANNULLATI\n\n%d", annullati))

		return lipgloss.JoinHorizontal(lipgloss.Top, kpiTot, "   ", kpiConf, "   ", kpiEvas, "   ", kpiAnn)

	case RESI:
		totale = len(m.resi)
		for _, r := range m.resi {
			switch r.Stato {
			case api.StatoInAttesa:
				inAttesa++
			case api.StatoAccettato:
				accettati++
			case api.StatoRifiutato:
				rifiutati++
			}
		}
		kpiTot := style.BoxStile.Render(fmt.Sprintf("TOTALE RESI\n\n%d", totale))
		kpiInAttesa := style.BoxStile.Render(fmt.Sprintf("IN ATTESA\n\n%d", inAttesa))
		kpiAcc := style.BoxStile.Render(fmt.Sprintf("ACCETTATI\n\n%d", accettati))
		kpiRif := style.BoxStile.Render(fmt.Sprintf("RIFIUTATI\n\n%d", rifiutati))

		return lipgloss.JoinHorizontal(lipgloss.Top, kpiTot, "   ", kpiInAttesa, "   ", kpiAcc, "   ", kpiRif)

	default:
		return ""
	}
}

// renderTabHeader genera le schede (tabs) cliccabili/selezionabili in cima alla vista.
func (m Model) renderTabHeader() string {
	t1 := " ORDINI IN ENTRATA "
	t2 := " ORDINI IN USCITA "
	t3 := " RESI "

	switch m.sezione {
	case ORDINI_ENTRATA:
		t1 = style.FocusStile.Render("[ " + t1 + " ]")
	case ORDINI_USCITA:
		t2 = style.FocusStile.Render("[ " + t2 + " ]")
	default:
		t3 = style.FocusStile.Render("[ " + t3 + " ]")
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, t1, "   ", t2, "   ", t3)
}

// renderTabellaEntrata formatta in formato tabellare l'elenco degli ordini di entrata attivi.
func (m Model) renderTabellaEntrata() string {
	var sb strings.Builder
	if len(m.ordiniEntrata) == 0 {
		return "Nessun ordine in entrata registrato.\n"
	}

	headerStr := fmt.Sprintf("   %-24s %-30s %-15s %s", "NUMERO", "FORNITORE", "DATA", "RIGHE")
	sb.WriteString(style.FreeHeaderStile.Render(headerStr))
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("─", 100))
	sb.WriteString("\n")

	inizio := m.offsetEntrata
	fine := inizio + m.righeVisibili()
	if fine > len(m.ordiniEntrata) {
		fine = len(m.ordiniEntrata)
	}

	for i, o := range m.ordiniEntrata[inizio:fine] {
		assoluto := inizio + i
		indicatore := "  "
		if assoluto == m.cursoreEntrata {
			indicatore = style.SuccessStile.Render("> ")
		}

		espanso := m.espansiEntrata[o.ID]
		simbolo := " [+] "
		if espanso {
			simbolo = " [-] "
		}

		rigaTesto := fmt.Sprintf("%-20s %-30s %-15s %d art.",
			troncaStringa(o.Numero, 20),
			troncaStringa(o.Fornitore, 30),
			o.Data.Format("02/01/2006"),
			len(o.Righe),
		)

		sb.WriteString(indicatore)
		sb.WriteString(simbolo)
		sb.WriteString(rigaTesto)
		sb.WriteString("\n")
		if espanso {
			sb.WriteString(renderRigheDettaglio(o.Numero, o.Righe, o.Note))
		}
	}

	if len(m.ordiniEntrata) > m.righeVisibili() {
		vis := m.righeVisibili()
		ultimo := m.offsetEntrata + vis
		if ultimo > len(m.ordiniEntrata) {
			ultimo = len(m.ordiniEntrata)
		}
		sb.WriteString(fmt.Sprintf("\nOrdini %d–%d di %d (pgup/pgdown o rotellina per scorrere)", m.offsetEntrata+1, ultimo, len(m.ordiniEntrata)))
	}

	return sb.String()
}

// renderTabellaUscita formatta in formato tabellare l'elenco degli ordini di uscita attivi.
func (m Model) renderTabellaUscita() string {
	var sb strings.Builder
	if len(m.ordiniUscita) == 0 {
		return "Nessun ordine in uscita registrato.\n"
	}

	headerStr := fmt.Sprintf("   %-24s %-30s %-15s %-15s %s", "NUMERO", "CLIENTE", "DATA", "STATO", "RIGHE")
	sb.WriteString(style.FreeHeaderStile.Render(headerStr))
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("─", 100))
	sb.WriteString("\n")

	inizio := m.offsetUscita
	fine := inizio + m.righeVisibili()
	if fine > len(m.ordiniUscita) {
		fine = len(m.ordiniUscita)
	}

	for i, o := range m.ordiniUscita[inizio:fine] {
		assoluto := inizio + i
		indicatore := "  "
		if assoluto == m.cursoreUscita {
			indicatore = style.SuccessStile.Render("> ")
		}

		espanso := m.espansiUscita[o.ID]
		simbolo := " [+] "
		if espanso {
			simbolo = " [-] "
		}

		statoPuro := fmt.Sprintf("%-15s", o.Stato)
		statoBadge := renderingStatoOrdine(o.Stato, statoPuro)

		rigaTesto := fmt.Sprintf("%-20s %-30s %-15s %-15s %d art.",
			troncaStringa(o.Numero, 20),
			troncaStringa(o.Cliente, 30),
			o.Data.Format("02/01/2006"),
			statoBadge,
			len(o.Righe),
		)

		sb.WriteString(indicatore)
		sb.WriteString(simbolo)
		sb.WriteString(rigaTesto)
		sb.WriteString("\n")
		if espanso {
			sb.WriteString(renderRigheDettaglio(o.Numero, o.Righe, o.Note))
		}
	}

	if len(m.ordiniUscita) > m.righeVisibili() {
		vis := m.righeVisibili()
		ultimo := m.offsetUscita + vis
		if ultimo > len(m.ordiniUscita) {
			ultimo = len(m.ordiniUscita)
		}
		sb.WriteString(fmt.Sprintf("\nOrdini %d–%d di %d (pgup/pgdown o rotellina per scorrere)", m.offsetUscita+1, ultimo, len(m.ordiniUscita)))
	}

	return sb.String()
}

// renderRigheDettaglio formatta il sotto-box contenente il dettaglio di ogni articolo e totale ordine.
func renderRigheDettaglio(numeroOrdine string, righeAny any, note string) string {
	var righe []rigaAdapter

	switch v := righeAny.(type) {
	case []api.RigaOrdineEntrata:
		for _, r := range v {
			righe = append(righe, rigaAdapter{
				ProdottoID:     r.ProdottoID,
				Prodotto:       r.Prodotto,
				Quantita:       r.Quantita,
				PrezzoUnitario: r.PrezzoUnitario,
			})
		}
	case []api.RigaOrdineUscita:
		for _, r := range v {
			righe = append(righe, rigaAdapter{
				ProdottoID:     r.ProdottoID,
				Prodotto:       r.Prodotto,
				Quantita:       r.Quantita,
				PrezzoUnitario: r.PrezzoUnitario,
			})
		}
	}

	var sb strings.Builder
	indent := "       "

	sb.WriteString(indent)
	sb.WriteString(style.HeaderStile.Render("┌─ DETTAGLI RIGHE ORDINE " + numeroOrdine))
	sb.WriteString("\n")
	if len(righe) == 0 {
		sb.WriteString(indent)
		sb.WriteString("│  (Nessuna riga presente)\n")
	} else {
		sb.WriteString(indent)
		sb.WriteString(fmt.Sprintf("│  %-10s %-25s %-10s %-12s %s\n", "CODICE", "PRODOTTO", "QUANTITÀ", "PREZZO UNIT.", "TOTALE RIGA"))
		sb.WriteString(indent)
		sb.WriteString("├")
		sb.WriteString(strings.Repeat("─", 65))
		sb.WriteString("\n")

		var totale float64
		for _, r := range righe {
			codice := "N/D"
			nome := fmt.Sprintf("Prodotto #%d", r.ProdottoID)
			if r.Prodotto != nil {
				codice = r.Prodotto.Codice
				if r.Prodotto.Nome != "" {
					nome = r.Prodotto.Nome
				}
			}

			totRiga := float64(r.Quantita) * r.PrezzoUnitario
			totale += totRiga
			sb.WriteString(indent)
			fmt.Fprintf(&sb, "│  %-10s %-25s %-10d € %-10.2f € %.2f\n",
				troncaStringa(codice, 10), troncaStringa(nome, 25), r.Quantita, r.PrezzoUnitario, totRiga)
		}
		sb.WriteString(indent)
		sb.WriteString("├")
		sb.WriteString(strings.Repeat("─", 65))
		sb.WriteString("\n")
		sb.WriteString(indent)
		sb.WriteString(fmt.Sprintf("│  TOTALE COMPLESSIVO ORDINE: € %.2f\n", totale))
	}

	if note != "" {
		sb.WriteString(indent)
		sb.WriteString("│  Note: ")
		sb.WriteString(note)
		sb.WriteString("\n")
	}
	sb.WriteString(indent)
	sb.WriteString("└")
	sb.WriteString(strings.Repeat("─", 65))
	sb.WriteString("\n")
	return sb.String()
}

// renderTabellaResi formatta in formato tabellare l'elenco delle richieste di reso.
func (m Model) renderTabellaResi() string {
	var sb strings.Builder
	if len(m.resi) == 0 {
		return "Nessun reso registrato.\n"
	}

	headerStr := fmt.Sprintf("   %-15s %-20s %-15s %s", "ID RESO", "ID ORDINE USCITA", "DATA", "STATO")
	sb.WriteString(style.FreeHeaderStile.Render(headerStr))
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("─", 100))
	sb.WriteString("\n")

	inizio := m.offsetResi
	fine := inizio + m.righeVisibili()
	if fine > len(m.resi) {
		fine = len(m.resi)
	}

	for i, r := range m.resi[inizio:fine] {
		assoluto := inizio + i
		indicatore := "  "
		if assoluto == m.cursoreResi {
			indicatore = style.SuccessStile.Render("> ")
		}

		espanso := m.espansiResi[r.ID]
		simbolo := " [+] "
		if espanso {
			simbolo = " [-] "
		}

		statoPuro := fmt.Sprintf("%-15s", r.Stato)
		statoBadge := renderingStatoReso(r.Stato, statoPuro)

		rigaTesto := fmt.Sprintf("Reso #%-5d Ordine #%-12d %-15s %s",
			r.ID,
			r.OrdineUscitaID,
			r.Data.Format("02/01/2006"),
			statoBadge,
		)

		sb.WriteString(indicatore)
		sb.WriteString(simbolo)
		sb.WriteString(rigaTesto)
		sb.WriteString("\n")
		if espanso {
			sb.WriteString(renderDettagliReso(r))
		}
	}

	if len(m.resi) > m.righeVisibili() {
		vis := m.righeVisibili()
		ultimo := m.offsetResi + vis
		if ultimo > len(m.resi) {
			ultimo = len(m.resi)
		}
		sb.WriteString(fmt.Sprintf("\nResi %d–%d di %d (pgup/pgdown o rotellina per scorrere)", m.offsetResi+1, ultimo, len(m.resi)))
	}

	return sb.String()
}

// renderDettagliReso formatta il sotto-box contenente le motivazioni del reso.
func renderDettagliReso(reso api.Reso) string {
	var sb strings.Builder
	indent := "       "

	sb.WriteString(indent)
	sb.WriteString(style.HeaderStile.Render(fmt.Sprintf("┌─ DETTAGLI RESO #%d (Ordine Uscita #%d)", reso.ID, reso.OrdineUscitaID)))
	sb.WriteString("\n")
	sb.WriteString(indent)
	sb.WriteString(fmt.Sprintf("│  %-10s %-25s %-10s %s\n", "CODICE", "PRODOTTO", "QUANTITÀ", "MOTIVO RESO"))
	sb.WriteString(indent)
	sb.WriteString("├")
	sb.WriteString(strings.Repeat("─", 65))
	sb.WriteString("\n")

	codice := "N/D"
	nome := fmt.Sprintf("Prodotto #%d", reso.ProdottoID)

	motivo := reso.Motivo
	if motivo == "" {
		motivo = "Nessun motivo specificato"
	}

	sb.WriteString(indent)
	sb.WriteString(fmt.Sprintf("│  %-10s %-25s %-10d %s\n",
		troncaStringa(codice, 10),
		troncaStringa(nome, 25),
		reso.Quantita,
		motivo,
	))

	sb.WriteString(indent)
	sb.WriteString("└")
	sb.WriteString(strings.Repeat("─", 65))
	sb.WriteString("\n")
	return sb.String()
}

// renderFormCreazioneEntrata disegna la vista form per la creazione di un nuovo ordine in entrata.
func (m Model) renderFormCreazioneEntrata() string {
	var sb strings.Builder
	sb.WriteString(style.TitleStile.Render("CREAZIONE ORDINE IN ENTRATA"))
	sb.WriteString("\n\n")

	if m.loading {
		sb.WriteString("Salvataggio ordine in corso...\n")
		return sb.String()
	}

	sb.WriteString(fmt.Sprintf("Fornitore:      %s\n", m.formControparte.View()))
	sb.WriteString(fmt.Sprintf("Numero Ordine:  %s\n", m.formNumero.View()))
	sb.WriteString(fmt.Sprintf("Note:           %s\n", m.formNote.View()))
	sb.WriteString(fmt.Sprintf("Righe Prodotti: %s\n", m.formRighe.View()))
	sb.WriteString("\n")
	sb.WriteString(style.InfoStile.Render("Formato righe: ID:Quantità:Prezzo separate da virgola (es. 1:100:1.50, 2:5:10.00)"))
	sb.WriteString("\n\n")

	if m.err != nil {
		sb.WriteString(style.ErrorStile.Render("Errore: " + m.err.Error()))
		sb.WriteString("\n\n")
	}

	sb.WriteString(style.FooterStile.Render("[tab/↑/↓] Naviga campi  |  [enter] Campo succ./Salva  |  [ctrl+s] Salva Ordine  |  [esc] Annulla"))

	return sb.String()
}

// renderFormCreazioneUscita disegna la vista form per la creazione di un nuovo ordine in uscita.
func (m Model) renderFormCreazioneUscita() string {
	var sb strings.Builder
	sb.WriteString(style.TitleStile.Render("CREAZIONE ORDINE IN USCITA"))
	sb.WriteString("\n\n")

	if m.loading {
		sb.WriteString("Salvataggio ordine in corso...\n")
		return sb.String()
	}

	sb.WriteString(fmt.Sprintf("Cliente:        %s\n", m.formControparte.View()))
	sb.WriteString(fmt.Sprintf("Numero Ordine:  %s\n", m.formNumero.View()))
	sb.WriteString(fmt.Sprintf("Note:           %s\n", m.formNote.View()))
	sb.WriteString(fmt.Sprintf("Righe Prodotti: %s\n", m.formRighe.View()))
	sb.WriteString("\n")
	sb.WriteString(style.InfoStile.Render("Formato righe: ID:Quantità:Prezzo separate da virgola (es. 1:100:1.50, 2:5:10.00)"))
	sb.WriteString("\n\n")

	if m.err != nil {
		sb.WriteString(style.ErrorStile.Render("Errore: " + m.err.Error()))
		sb.WriteString("\n\n")
	}

	sb.WriteString(style.FooterStile.Render("[tab/↑/↓] Naviga campi  |  [enter] Campo succ./Salva  |  [ctrl+s] Salva Ordine  |  [esc] Annulla"))

	return sb.String()
}

// renderFormCreazioneReso disegna la vista form per la creazione di un nuovo reso.
func (m Model) renderFormCreazioneReso() string {
	var sb strings.Builder
	sb.WriteString(style.TitleStile.Render("CREAZIONE NUOVO RESO"))
	sb.WriteString("\n\n")

	if m.loading {
		sb.WriteString("Salvataggio reso in corso...\n")
		return sb.String()
	}

	sb.WriteString(fmt.Sprintf("ID Ordine Uscita: %s\n", m.formResoOrdineID.View()))
	sb.WriteString(fmt.Sprintf("ID Prodotto:      %s\n", m.formResoProdottoID.View()))
	sb.WriteString(fmt.Sprintf("Quantità:         %s\n", m.formResoQuantita.View()))
	sb.WriteString(fmt.Sprintf("Motivo:           %s\n", m.formResoMotivo.View()))
	sb.WriteString("\n\n")

	if m.err != nil {
		sb.WriteString(style.ErrorStile.Render("Errore: " + m.err.Error()))
		sb.WriteString("\n\n")
	}

	sb.WriteString(style.FooterStile.Render("[tab/↑/↓] Naviga campi  |  [enter] Campo succ./Salva  |  [ctrl+s] Salva Reso  |  [esc] Annulla"))

	return sb.String()
}
