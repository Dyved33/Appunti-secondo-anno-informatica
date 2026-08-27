// Package dashboard gestisce la logica di caricamento e il rendering della panoramica amministrativa.
package dashboard

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
	style "github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/styles"
)

// ── Costanti di layout per lo scorrimento della pagina ──────────────────────

// overheadChrome è il numero di righe fisse occupate ATTORNO al corpo della
// pagina (logo ASCII, intestazione, barra dei tab, cornice della finestra).
// È stato misurato sul layout reale composto da tabnav.PlacePage: sottraendolo
// dall'altezza del terminale si ottiene lo spazio disponibile per il contenuto.
const overheadChrome = 21

// passoRotellina è il numero di righe di cui la pagina scorre per ogni evento della rotellina.
const passoRotellina = 3

// righeFallback è il numero di righe visibili usato finché non arriva il primo WindowSizeMsg.
const righeFallback = 14

// ── Scorrimento della pagina ─────────────────────────────────────────────────

// contenutoRighe restituisce TUTTE le righe del corpo dashboard (blocchi KPI,
// scorte critiche e footer), già renderizzate con gli stili.
//
// Come è stato scritto il codice:
// Ogni riga prodotta da lipgloss porta con sé i propri codici ANSI, quindi le
// righe possono essere affettate in sicurezza per realizzare la finestra visibile.
func (m Model) contenutoRighe() []string {
	return strings.Split(m.corpo(), "\n")
}

// righeVisibili calcola quante righe di contenuto entrano nello spazio residuo dello schermo.
func (m Model) righeVisibili() int {
	if m.winH <= 0 {
		return righeFallback
	}
	disponibili := m.winH - overheadChrome
	if disponibili < 1 {
		disponibili = 1
	}
	return disponibili
}

// finestraEffettiva è il numero di righe effettivamente mostrate: quando il
// contenuto eccede lo schermo viene riservata una riga per l'indicatore di scorrimento.
func (m Model) finestraEffettiva() int {
	visibili := m.righeVisibili()
	if len(m.contenutoRighe()) > visibili {
		visibili--
		if visibili < 1 {
			visibili = 1
		}
	}
	return visibili
}

// scorri avanza la finestra di visualizzazione di `n` righe, fermandosi agli estremi.
func (m *Model) scorri(n int) {
	m.scrollOffset += n
	m.clampScroll()
}

// clampScroll riporta scrollOffset nell'intervallo valido rispetto a contenuto e schermo.
func (m *Model) clampScroll() {
	maxOffset := len(m.contenutoRighe()) - m.finestraEffettiva()
	if maxOffset < 0 {
		maxOffset = 0
	}
	if m.scrollOffset > maxOffset {
		m.scrollOffset = maxOffset
	}
	if m.scrollOffset < 0 {
		m.scrollOffset = 0
	}
}

// Init avvia il caricamento parallelo di tutte e quattro le sorgenti dati della dashboard.
//
// Come è stato scritto il codice:
// Utilizza tea.Batch per lanciare contemporaneamente quattro comandi asincroni indipendenti
// (giacenze, KPI contabilità, ordini in uscita, ordini in entrata). I rispettivi messaggi
// di completamento aggiorneranno il modello via Update() non appena ciascuna chiamata termina.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.loadGiacenzeCmd(),
		m.loadKPICmd(),
		m.loadOrdiniUscitaCmd(),
		m.loadOrdiniEntrataCmd(),
	)
}

// Update elabora i messaggi di risposta dalle API e gli eventi di navigazione da tastiera.
//
// Come è stato scritto il codice:
//  1. Intercetta i quattro messaggi di caricamento completato (*CaricatoMsg) aggiornando il campo
//     corrispondente e abbassando il relativo flag loading.
//  2. Intercetta i quattro messaggi di errore (*ErroreMsg) salvando l'errore e abbassando loading.
//  3. Su tea.KeyPressMsg gestisce la navigazione "up"/"k" e "down"/"j" nella lista delle scorte
//     critiche, e "r" per ricaricare tutti i blocchi in parallelo.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// — Caricamenti completati —

	case GiacenzeCaricateMsg:
		m.reportGiacenze = msg.Report
		m.loadingGiacenze = false
		m.erroreGiacenze = nil
		m.cursoreScorteCritiche = 0
		return m, nil

	case KPICaricatiMsg:
		m.kpiContabilita = msg.KPI
		m.loadingKPI = false
		m.erroreKPI = nil
		return m, nil

	case OrdiniUscitaCaricatiMsg:
		m.ordiniUscita = msg.Ordini
		m.loadingOrdiniUscita = false
		m.erroreOrdiniUscita = nil
		return m, nil

	case OrdiniEntrataCaricatiMsg:
		m.ordiniEntrata = msg.Ordini
		m.loadingOrdiniEntrata = false
		m.erroreOrdiniEntrata = nil
		return m, nil

	// — Errori —

	case ErroreGiacenzeMsg:
		m.erroreGiacenze = msg.Err
		m.loadingGiacenze = false
		return m, nil

	case ErroreKPIMsg:
		m.erroreKPI = msg.Err
		m.loadingKPI = false
		return m, nil

	case ErroreOrdiniUscitaMsg:
		m.erroreOrdiniUscita = msg.Err
		m.loadingOrdiniUscita = false
		return m, nil

	case ErroreOrdiniEntrataMsg:
		m.erroreOrdiniEntrata = msg.Err
		m.loadingOrdiniEntrata = false
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

	// — Tastiera —

	case tea.KeyPressMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursoreScorteCritiche > 0 {
				m.cursoreScorteCritiche--
			}
		case "down", "j":
			if m.reportGiacenze != nil &&
				m.cursoreScorteCritiche < len(m.reportGiacenze.ArticoliSottoSoglia)-1 {
				m.cursoreScorteCritiche++
			}
		case "pgup":
			m.scorri(-m.righeVisibili())
		case "pgdown":
			m.scorri(m.righeVisibili())
		case "home":
			m.scrollOffset = 0
		case "end":
			m.scorri(len(m.contenutoRighe()))
		case "r":
			// Ricarica tutto in parallelo reimpostando i flag di loading
			m.loadingGiacenze = true
			m.loadingKPI = true
			m.loadingOrdiniUscita = true
			m.loadingOrdiniEntrata = true
			return m, tea.Batch(
				m.loadGiacenzeCmd(),
				m.loadKPICmd(),
				m.loadOrdiniUscitaCmd(),
				m.loadOrdiniEntrataCmd(),
			)
		}
	}

	return m, nil
}

// View genera ed impagina l'interfaccia utente della dashboard, applicando lo
// scorrimento quando il contenuto eccede l'altezza del terminale.
//
// Come è stato scritto il codice:
//  1. Il corpo completo viene costruito da corpo() (blocchi KPI su due righe +
//     lista scorte critiche + footer).
//  2. Se il corpo entra nello schermo viene restituito integralmente.
//  3. Altrimenti viene mostrata solo la finestra di righe puntata da scrollOffset,
//     sostituendo il footer con un indicatore "Righe X–Y di Z": così anche nei
//     terminali bassi tutta la dashboard resta raggiungibile con rotellina o PgSu/PgGiù.
func (m Model) View() tea.View {
	righe := m.contenutoRighe()

	if len(righe) <= m.righeVisibili() {
		return tea.NewView(strings.Join(righe, "\n"))
	}

	visibili := m.finestraEffettiva()
	m.clampScroll()
	inizio := m.scrollOffset
	fine := inizio + visibili
	if fine > len(righe) {
		fine = len(righe)
	}

	var sb strings.Builder
	sb.WriteString(strings.Join(righe[inizio:fine], "\n"))
	sb.WriteString("\n")
	sb.WriteString(style.FooterStile.Render(fmt.Sprintf(
		"Righe %d–%d di %d — rotellina o PgSu/PgGiù per scorrere",
		inizio+1, fine, len(righe),
	)))

	return tea.NewView(sb.String())
}

// corpo costruisce il contenuto COMPLETO della dashboard: titolo, i quattro
// blocchi KPI affiancati su due righe, la lista delle scorte critiche ed il footer.
//
// Come è stato scritto il codice:
//  1. Per ciascuna delle quattro sezioni, se il blocco è ancora in caricamento mostra un indicatore
//     testuale; se vi è un errore lo mostra stilizzato; altrimenti renderizza i dati.
//  2. Le quattro sezioni KPI vengono affiancate su due colonne usando i box di stile.
//  3. La lista scorte critiche è navigabile con ↑/↓ e mostra codice, nome, quantità e soglia minima.
func (m Model) corpo() string {
	var sb strings.Builder

	// ── Titolo ──────────────────────────────────────────────────────────────
	sb.WriteString(style.TitleStile.Render("DASHBOARD ADMIN"))
	sb.WriteString("\n\n")

	// ── Riga 1: Magazzino | Contabilità ─────────────────────────────────────
	bloccoGiacenze := m.renderBloccoGiacenze()
	bloccoKPI := m.renderBloccoKPI()

	// Unisce orizzontalmente i due blocchi con 3 spazi di separazione
	riga1 := lipgloss.JoinHorizontal(lipgloss.Top, bloccoGiacenze, "   ", bloccoKPI)
	sb.WriteString(riga1)
	sb.WriteString("\n\n")

	// ── Riga 2: Ordini Uscita | Ordini Entrata ───────────────────────────────
	bloccoUscita := m.renderBloccoOrdiniUscita()
	bloccoEntrata := m.renderBloccoOrdiniEntrata()

	riga2 := lipgloss.JoinHorizontal(lipgloss.Top, bloccoUscita, "   ", bloccoEntrata)
	sb.WriteString(riga2)
	sb.WriteString("\n\n")

	// ── Scorte Critiche ──────────────────────────────────────────────────────
	sb.WriteString(m.renderScorteCritiche())

	// ── Footer ───────────────────────────────────────────────────────────────
	sb.WriteString("\n")
	sb.WriteString(style.FooterStile.Render(
		"[↑/↓ o k/j] Naviga Scorte Critiche  |  [r] Ricarica  |  [ctrl+l] Logout",
	))

	return sb.String()
}

// ── Blocchi di rendering ─────────────────────────────────────────────────────

// righePerBlocco è il numero ESATTO di righe di contenuto di ciascun blocco KPI.
//
// BoxStile fissa solo un'altezza minima: se due blocchi affiancati hanno
// contenuti con numero di righe diverso, i bordi inferiori risultano disallineati
// ("bordi sfalzati"). Passando sempre da renderBlocco, ogni riquadro ha la
// stessa geometria in qualunque stato (caricamento, errore, dati).
const righePerBlocco = 6

// renderBlocco normalizza le righe di contenuto a righePerBlocco e le incornicia.
//
// Come è stato scritto il codice:
// Aggiunge righe vuote se il contenuto è più corto e tronca se più lungo,
// così il box reso da style.BoxStile ha sempre la medesima altezza dei box
// affiancati sulla stessa riga della dashboard.
func renderBlocco(righe []string) string {
	if len(righe) > righePerBlocco {
		righe = righe[:righePerBlocco]
	}
	for len(righe) < righePerBlocco {
		righe = append(righe, "")
	}
	return style.BoxStile.Render(strings.Join(righe, "\n"))
}

// renderBloccoGiacenze costruisce il box KPI relativo al magazzino e alle giacenze.
func (m Model) renderBloccoGiacenze() string {
	titolo := "MAGAZZINO"
	switch {
	case m.loadingGiacenze:
		return renderBlocco([]string{titolo, "", "Caricamento..."})
	case m.erroreGiacenze != nil:
		return renderBlocco([]string{titolo, "",
			style.ErrorStile.Render("Errore: " + troncaStringa(m.erroreGiacenze.Error(), 28))})
	case m.reportGiacenze == nil:
		return renderBlocco([]string{titolo, "", "Nessun dato disponibile."})
	}

	r := m.reportGiacenze
	sottoSoglia := len(r.ArticoliSottoSoglia)
	stato := style.SuccessStile.Render("OK")
	if sottoSoglia > 0 {
		stato = style.ErrorStile.Render(fmt.Sprintf("! %d critici", sottoSoglia))
	}

	return renderBlocco([]string{
		titolo,
		"",
		fmt.Sprintf("Totale Pezzi    %d", r.TotaleQuantita),
		fmt.Sprintf("Valore Stock    € %.2f", r.ValoreTotale),
		fmt.Sprintf("Scorte Critiche %s", stato),
	})
}

// renderBloccoKPI costruisce il box dei KPI contabilità (fatturato, fatture, ordini, note credito).
func (m Model) renderBloccoKPI() string {
	titolo := "CONTABILITÀ"
	switch {
	case m.loadingKPI:
		return renderBlocco([]string{titolo, "", "Caricamento..."})
	case m.erroreKPI != nil:
		return renderBlocco([]string{titolo, "",
			style.ErrorStile.Render("Errore: " + troncaStringa(m.erroreKPI.Error(), 28))})
	case len(m.kpiContabilita) == 0:
		return renderBlocco([]string{titolo, "", "Nessun dato disponibile."})
	}

	righe := []string{titolo, ""}
	for _, kpi := range m.kpiContabilita {
		righe = append(righe, fmt.Sprintf("%-14s %s", kpi.Etichetta, kpi.Valore))
	}
	return renderBlocco(righe)
}

// renderBloccoOrdiniUscita costruisce il box riepilogativo degli ordini in uscita per stato.
func (m Model) renderBloccoOrdiniUscita() string {
	titolo := "ORDINI IN USCITA"
	switch {
	case m.loadingOrdiniUscita:
		return renderBlocco([]string{titolo, "", "Caricamento..."})
	case m.erroreOrdiniUscita != nil:
		return renderBlocco([]string{titolo, "",
			style.ErrorStile.Render("Errore: " + troncaStringa(m.erroreOrdiniUscita.Error(), 28))})
	}

	var confermati, evasi, annullati int
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

	return renderBlocco([]string{
		titolo,
		"",
		fmt.Sprintf("Totale      %d", len(m.ordiniUscita)),
		fmt.Sprintf("Confermati  %d", confermati),
		fmt.Sprintf("Evasi       %s", style.SuccessStile.Render(fmt.Sprintf("%d", evasi))),
		fmt.Sprintf("Annullati   %d", annullati),
	})
}

// renderBloccoOrdiniEntrata costruisce il box riepilogativo degli ordini in entrata.
//
// Come è stato scritto il codice:
// Gli ordini in entrata non hanno più uno stato nel backend: il box mostra quindi
// dati reali aggregati (totale ordini, fornitori distinti, articoli caricati e valore complessivo).
func (m Model) renderBloccoOrdiniEntrata() string {
	titolo := "ORDINI IN ENTRATA"
	switch {
	case m.loadingOrdiniEntrata:
		return renderBlocco([]string{titolo, "", "Caricamento..."})
	case m.erroreOrdiniEntrata != nil:
		return renderBlocco([]string{titolo, "",
			style.ErrorStile.Render("Errore: " + troncaStringa(m.erroreOrdiniEntrata.Error(), 28))})
	}

	fornitori, articoli, valore := statisticheEntrata(m.ordiniEntrata)

	return renderBlocco([]string{
		titolo,
		"",
		fmt.Sprintf("Totale     %d", len(m.ordiniEntrata)),
		fmt.Sprintf("Fornitori  %d", fornitori),
		fmt.Sprintf("Articoli   %d", articoli),
		fmt.Sprintf("Valore     € %.2f", valore),
	})
}

// statisticheEntrata calcola fornitori distinti, quantità totale e valore complessivo degli ordini in entrata.
func statisticheEntrata(ordini []api.OrdineEntrata) (fornitori, articoli int, valore float64) {
	insFornitori := make(map[string]bool)
	for _, o := range ordini {
		if o.Fornitore != "" {
			insFornitori[o.Fornitore] = true
		}
		for _, r := range o.Righe {
			articoli += r.Quantita
			valore += float64(r.Quantita) * r.PrezzoUnitario
		}
	}
	return len(insFornitori), articoli, valore
}

// renderScorteCritiche costruisce la tabella navigabile degli articoli sotto soglia minima.
//
// Come è stato scritto il codice:
// Se il report giacenze non è ancora disponibile mostra un messaggio di attesa.
// Se non ci sono articoli sotto soglia mostra un messaggio positivo.
// Altrimenti stampa l'header e una riga per ogni articolo, evidenziando quella
// puntata dal cursore con il simbolo "> " stilizzato.
func (m Model) renderScorteCritiche() string {
	var sb strings.Builder

	sb.WriteString(style.TitleStile.Render("SCORTE CRITICHE"))
	sb.WriteString("\n")

	if m.loadingGiacenze {
		sb.WriteString("  Caricamento...\n")
		return sb.String()
	}
	if m.erroreGiacenze != nil {
		sb.WriteString(style.ErrorStile.Render("  Errore: " + m.erroreGiacenze.Error()))
		sb.WriteString("\n")
		return sb.String()
	}
	if m.reportGiacenze == nil || len(m.reportGiacenze.ArticoliSottoSoglia) == 0 {
		sb.WriteString(style.SuccessStile.Render("  Nessuna scorta critica — inventario nella norma."))
		sb.WriteString("\n")
		return sb.String()
	}

	articoli := m.reportGiacenze.ArticoliSottoSoglia

	// Header tabella
	headerStr := fmt.Sprintf("   %-10s %-24s %-12s %-10s %s",
		"CODICE", "PRODOTTO", "GIACENZA", "SOGLIA MIN", "CATEGORIA")
	sb.WriteString(style.FreeHeaderStile.Render(headerStr))
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("─", 75))
	sb.WriteString("\n")

	for i, a := range articoli {
		giacenzaStr := fmt.Sprintf("%d pz", a.Quantita)
		sogliaStr := fmt.Sprintf("%d pz", a.SogliaMin)

		riga := fmt.Sprintf("%-10s %-24s %-12s %-10s %s",
			troncaStringa(a.Codice, 10),
			troncaStringa(a.Nome, 24),
			giacenzaStr,
			sogliaStr,
			troncaStringa(a.Categoria, 15),
		)

		if i == m.cursoreScorteCritiche {
			sb.WriteString(style.ErrorStile.Render("> "))
			sb.WriteString(riga)
		} else {
			sb.WriteString("  ")
			sb.WriteString(riga)
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// ── Comandi API ──────────────────────────────────────────────────────────────

// loadGiacenzeCmd recupera in modo asincrono il report inventario dal backend.
func (m Model) loadGiacenzeCmd() tea.Cmd {
	return func() tea.Msg {
		report, err := m.client.GetReportGiacenze()
		if err != nil {
			return ErroreGiacenzeMsg{Err: err}
		}
		return GiacenzeCaricateMsg{Report: report}
	}
}

// loadKPICmd recupera in modo asincrono i KPI di contabilità dal backend.
func (m Model) loadKPICmd() tea.Cmd {
	return func() tea.Msg {
		kpi, err := m.client.GetReportKPI()
		if err != nil {
			return ErroreKPIMsg{Err: err}
		}
		return KPICaricatiMsg{KPI: kpi}
	}
}

// loadOrdiniUscitaCmd recupera in modo asincrono gli ordini in uscita dal backend.
func (m Model) loadOrdiniUscitaCmd() tea.Cmd {
	return func() tea.Msg {
		ordini, err := m.client.GetOrdiniUscita()
		if err != nil {
			return ErroreOrdiniUscitaMsg{Err: err}
		}
		return OrdiniUscitaCaricatiMsg{Ordini: ordini}
	}
}

// loadOrdiniEntrataCmd recupera in modo asincrono gli ordini in entrata dal backend.
func (m Model) loadOrdiniEntrataCmd() tea.Cmd {
	return func() tea.Msg {
		ordini, err := m.client.GetOrdiniEntrata()
		if err != nil {
			return ErroreOrdiniEntrataMsg{Err: err}
		}
		return OrdiniEntrataCaricatiMsg{Ordini: ordini}
	}
}

// ── Utility ──────────────────────────────────────────────────────────────────

// troncaStringa tronca la stringa a maxLen caratteri aggiungendo "..." se necessario.
func troncaStringa(s string, maxLen int) string {
	if len(s) > maxLen {
		if maxLen > 3 {
			return s[:maxLen-3] + "..."
		}
		return s[:maxLen]
	}
	return s
}
