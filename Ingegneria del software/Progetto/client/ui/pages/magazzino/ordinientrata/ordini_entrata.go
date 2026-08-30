// Package ordinientrata gestisce le interazioni e la renderizzazione della schermata Ordini in Entrata.
package ordinientrata

import (
	"fmt"
	"strings"

	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
	style "github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/styles"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// overheadChrome è il numero di righe fisse occupate ATTORNO al corpo della pagina
// (logo, intestazione, barra dei tab, bordi della finestra e riga vuota interna).
// È stato misurato sul layout reale composto da tabnav.PlacePage: sommandolo alle
// dimensioni del terminale si ricava quante righe di lista possono essere mostrate.
const overheadChrome = 21

// passoRotellina è il numero di righe di cui la lista scorre per ogni evento della rotellina.
const passoRotellina = 3

// righeFallback è il numero di righe visibili usato finché non arriva il primo WindowSizeMsg.
const righeFallback = 12

// minRigheConKPI è il numero minimo di righe di lista richiesto per mostrare i blocchi KPI:
// nei terminali troppo bassi i KPI vengono nascosti per lasciare spazio agli ordini.
const minRigheConKPI = 3

// lineaLista rappresenta una singola riga renderizzabile della lista ordini.
//
// `ordine` è l'indice dell'ordine di riferimento nell'array `ordini`; per le righe
// di dettaglio di un ordine espanso vale lo stesso indice con `inizioOrd = false`,
// mentre le righe principali hanno `inizioOrd = true`.
type lineaLista struct {
	testo     string
	ordine    int
	inizioOrd bool
}

// Init avvia l'inizializzazione della vista Ordini in Entrata richiedendo la lista ordini al backend.
//
// Come è stato scritto il codice:
// Restituisce il comando asincrono `loadOrdiniEntrataCmd()` che recupera la lista degli ordini fornitori dal client API.
func (m Model) Init() tea.Cmd {
	return m.loadOrdiniEntrataCmd()
}

// Update gestisce l'aggiornamento dello stato in seguito alla ricezione di messaggi o eventi di navigazione.
//
// Come è stato scritto il codice:
//  1. Su `OrdiniEntrataLoadedMsg` aggiorna la lista `ordini`, resetta cursore e scorrimento ed azzera `loading` ed `err`.
//  2. Su `ErroreMsg` registra l'errore e disattiva lo stato di caricamento.
//  3. Su `tea.WindowSizeMsg` conserva le nuove dimensioni e ricalibra la finestra visibile.
//  4. Su `tea.MouseWheelMsg` scorre la lista nella direzione della rotellina.
//  5. Su `tea.KeyPressMsg`:
//     - "up"/"k" e "down"/"j": muovono il cursore mantenendolo dentro la porzione visibile.
//     - "enter"/"space": commute l'espansione dei dettagli dell'ordine evidenziato.
//     - "pgup"/"pgdown": scorrono la lista di una pagina; "home"/"end": tornano all'inizio/fine.
//     - "r": riattiva `loading = true` e richiama `loadOrdiniEntrataCmd()`.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case OrdiniEntrataLoadedMsg:
		m.loading = false
		m.err = nil
		m.ordini = msg.Ordini
		m.cursore = 0
		m.scrollOffset = 0
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
		switch msg.String() {
		case "up", "k":
			if m.cursore > 0 {
				m.cursore--
			}
			m.mostraCursore()
		case "down", "j":
			if m.cursore < len(m.ordini)-1 {
				m.cursore++
			}
			m.mostraCursore()
		case "enter", "space":
			if len(m.ordini) > 0 {
				id := m.ordini[m.cursore].ID
				m.dettagliEspansi[id] = !m.dettagliEspansi[id]
			}
			m.mostraCursore()
		case "pgup":
			m.scorri(-m.righeVisibili())
		case "pgdown":
			m.scorri(m.righeVisibili())
		case "home":
			m.scrollOffset = 0
		case "end":
			m.scorri(len(m.ordini))
		case "r":
			m.loading = true
			return m, m.loadOrdiniEntrataCmd()
		}
	}

	return m, nil
}

// ── Scorrimento della lista ──────────────────────────────────────────────────

// listaRighe compone TUTTE le righe visualizzabili: una riga per ciascun ordine
// più, per gli ordini espansi, le righe di dettaglio prodotto.
func (m Model) listaRighe() []lineaLista {
	righe := make([]lineaLista, 0, len(m.ordini))

	for i, o := range m.ordini {
		dataStr := o.Data.Format("02/01/2006")
		numRighe := fmt.Sprintf("%d art.", len(o.Righe))

		rigaTesto := fmt.Sprintf("%-24s %-34s %-15s %s",
			troncaStringa(o.Numero, 20),
			troncaStringa(o.Fornitore, 30),
			dataStr,
			numRighe,
		)

		prefisso := "  "
		if i == m.cursore {
			prefisso = style.SuccessStile.Render("> ")
		}

		simboloDettaglio := " [+] "
		if m.dettagliEspansi[o.ID] {
			simboloDettaglio = " [-] "
		}

		righe = append(righe, lineaLista{
			testo:     prefisso + simboloDettaglio + rigaTesto,
			ordine:    i,
			inizioOrd: true,
		})

		if m.dettagliEspansi[o.ID] {
			for _, l := range strings.Split(strings.TrimRight(renderRigheEntrata(o), "\n"), "\n") {
				righe = append(righe, lineaLista{testo: l, ordine: i, inizioOrd: false})
			}
		}
	}

	return righe
}

// mostraKPI indica se i blocchi KPI entrano nello schermo lasciando almeno
// minRigheConKPI righe di lista: nei terminali bassi vengono nascosti per
// privilegiare la visibilità degli ordini.
func (m Model) mostraKPI() bool {
	if m.winH <= 0 {
		return true
	}
	disponibili := m.winH - overheadChrome -
		lipgloss.Height(m.renderIntestazione(true)) -
		lipgloss.Height(m.renderFooter(-1, -1))
	if disponibili >= minRigheConKPI {
		return true
	}
	return false
}

// righeVisibili calcola quante righe di lista entrano nello spazio residuo dello schermo.
//
// Come è stato scritto il codice:
// Al termine `winH` vengono sottratte le righe fisse esterne al corpo (overheadChrome)
// e quelle dell'intestazione e del footer della pagina, misurate sul rendering effettivo:
// così il risultato resta corretto anche se il layout cambierà in futuro.
func (m Model) righeVisibili() int {
	if m.winH <= 0 {
		return righeFallback
	}
	disponibili := m.winH - overheadChrome -
		lipgloss.Height(m.renderIntestazione(m.mostraKPI())) -
		lipgloss.Height(m.renderFooter(-1, -1))
	if disponibili < 1 {
		disponibili = 1
	}
	return disponibili
}

// scorri avanza la finestra di visualizzazione di `n` righe, fermandosi agli estremi.
func (m *Model) scorri(n int) {
	m.scrollOffset += n
	m.clampScroll()
}

// clampScroll riporta scrollOffset nell'intervallo valido rispetto a lista e schermo.
func (m *Model) clampScroll() {
	maxOffset := len(m.listaRighe()) - m.righeVisibili()
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

// mostraCursore aggiorna lo scorrimento così che l'ordine puntato dal cursore resti visibile.
func (m *Model) mostraCursore() {
	righe := m.listaRighe()
	inizio := -1
	for i, l := range righe {
		if l.inizioOrd && l.ordine == m.cursore {
			inizio = i
			break
		}
	}
	if inizio < 0 {
		return
	}
	vis := m.righeVisibili()
	if inizio < m.scrollOffset {
		m.scrollOffset = inizio
	} else if inizio >= m.scrollOffset+vis {
		m.scrollOffset = inizio - vis + 1
	}
	m.clampScroll()
}

// ── Rendering ────────────────────────────────────────────────────────────────

// View costruisce il layout grafico stilizzato per la lista e per i dettagli degli ordini in entrata.
//
// Come è stato scritto il codice:
//  1. Mostra messaggi di caricamento o errore se il modello non è pronto.
//  2. Renderizza l'intestazione fissa (titolo, KPI e testata tabella) tramite `renderIntestazione`.
//  3. Mostra solo la finestra di righe che entra nello schermo (`scrollOffset` + `righeVisibili`),
//     rendendo raggiungibile anche gli ordini oltre la capacità del terminale.
//  4. Aggiunge a piè di pagina le istruzioni di tastiera e l'indicatore "Ordini X–Y di Z".
func (m Model) View() tea.View {
	var sb strings.Builder

	// ── Gestione Caricamento ed Errori ──────────────────────────────────────
	if m.loading {
		sb.WriteString(style.TitleStile.Render("ORDINI IN ENTRATA"))
		sb.WriteString("\n\n")
		sb.WriteString("Caricamento ordini in entrata in corso...\n")
		return tea.NewView(sb.String())
	}
	if m.err != nil {
		sb.WriteString(style.TitleStile.Render("ORDINI IN ENTRATA"))
		sb.WriteString("\n\n")
		sb.WriteString(style.ErrorStile.Render("Errore: " + m.err.Error()))
		sb.WriteString("\n")
		return tea.NewView(sb.String())
	}

	// ── Intestazione fissa ───────────────────────────────────────────────────
	sb.WriteString(m.renderIntestazione(m.mostraKPI()))

	// ── Finestra visibile della tabella ──────────────────────────────────────
	righe := m.listaRighe()

	if len(righe) == 0 {
		sb.WriteString("Nessun ordine in entrata registrato.\n")
		sb.WriteString("\n\n")
		sb.WriteString(m.renderFooter(-1, -1))
		return tea.NewView(sb.String())
	}

	m.clampScroll()
	inizio := m.scrollOffset
	fine := inizio + m.righeVisibili()
	if fine > len(righe) {
		fine = len(righe)
	}

	for _, l := range righe[inizio:fine] {
		sb.WriteString(l.testo)
		sb.WriteString("\n")
	}

	// ── Footer con indicatore di scorrimento ────────────────────────────────
	primo, ultimo := estremiVisibili(righe, inizio, fine)
	sb.WriteString(m.renderFooter(primo, ultimo))

	return tea.NewView(sb.String())
}

// renderIntestazione compone la parte fissa superiore della pagina: titolo,
// blocchi KPI aggregati (se `conKPI`, vedi mostraKPI) e testata della tabella ordini.
//
// Come è stato scritto il codice:
// Gli ordini in entrata non hanno più uno stato nel backend, quindi i KPI mostrano
// dati reali aggregati: numero ordini, fornitori distinti, articoli totali e valore complessivo.
func (m Model) renderIntestazione(conKPI bool) string {
	var sb strings.Builder

	// Il MarginBottom di TitleStile fornisce già la riga vuota sotto il titolo:
	// qui basta un solo "\n" esplicito per non sprecare righe di terminale.
	sb.WriteString(style.TitleStile.Render("ORDINI IN ENTRATA"))
	sb.WriteString("\n")

	if conKPI {
		fornitori, articoli, valore := statisticheEntrata(m.ordini)

		kpiTot := style.BoxStile.Render(fmt.Sprintf("TOTALE ORDINI\n\n%d", len(m.ordini)))
		kpiForn := style.BoxStile.Render(fmt.Sprintf("FORNITORI\n\n%d", fornitori))
		kpiArticoli := style.BoxStile.Render(fmt.Sprintf("ARTICOLI\n\n%d", articoli))
		kpiValore := style.BoxStile.Render(fmt.Sprintf("VALORE TOTALE\n\n€ %.2f", valore))

		rigaKPI := lipgloss.JoinHorizontal(lipgloss.Top, kpiTot, "   ", kpiForn, "   ", kpiArticoli, "   ", kpiValore)
		sb.WriteString(rigaKPI)
		sb.WriteString("\n\n")
	}

	headerStr := fmt.Sprintf("   %-24s %-34s %-15s %s",
		"NUMERO", "FORNITORE", "DATA", "RIGHE")
	sb.WriteString(style.FreeHeaderStile.Render(headerStr))
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("─", 90))
	sb.WriteString("\n")

	return sb.String()
}

// renderFooter compone le istruzioni di tastiera e, se la lista eccede lo schermo,
// l'indicatore degli ordini correntmente visibili ("Ordini X–Y di Z").
//
// `primo` e `ultimo` sono indici 0-based degli ordini visibili; `primo < 0`
// significa "nessun indicatore". Il blocco ha SEMPRE due righe (scorciatoie +
// indicatore anche vuoto): così il calcolo delle righe visibili è stabile e
// l'ultima riga di lista non viene mai spinta fuori schermo dall'indicatore stesso.
func (m Model) renderFooter(primo, ultimo int) string {
	scorciatoie := "[↑/↓ o k/j] Naviga  |  [enter/space] Espandi/Comprimi Dettaglio  |  [pgup/pgdn o rotellina] Scorri  |  [r] Ricarica  |  [ctrl+l] Logout"

	indicatore := ""
	if primo >= 0 && len(m.ordini) > 0 && (primo > 0 || ultimo < len(m.ordini)-1) {
		indicatore = fmt.Sprintf("Ordini %d–%d di %d  (pgup/pgdn o rotellina per vederli tutti)", primo+1, ultimo+1, len(m.ordini))
	}

	return style.FooterStile.Render(scorciatoie + "\n" + indicatore)
}

// estremiVisibili individua gli indici (0-based) del primo e dell'ultimo ordine
// la cui riga principale cade nella finestra [inizio, fine).
//
// Come è stato scritto il codice:
// Scorre la finestra cercando le righe marcate `inizioOrd`; se nessuna riga principale
// è visibile (finestra aperta in mezzo ai dettagli di un ordine espanso) usa l'ordine
// di riferimento della prima riga visibile.
func estremiVisibili(righe []lineaLista, inizio, fine int) (primo, ultimo int) {
	primo, ultimo = -1, -1
	for _, l := range righe[inizio:fine] {
		if l.inizioOrd {
			if primo == -1 {
				primo = l.ordine
			}
			ultimo = l.ordine
		}
	}
	if primo == -1 && fine > inizio {
		primo = righe[inizio].ordine
		ultimo = righe[fine-1].ordine
	}
	return primo, ultimo
}

// statisticheEntrata calcola fornitori distinti, quantità totale di articoli e
// valore complessivo (quantità × prezzo unitario) sugli ordini in entrata.
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

// renderRigheEntrata genera ed impagina la tabella subordinata per le singole righe di un ordine in entrata.
//
// Come è stato scritto il codice:
// Scorre la slice `o.Righe`, formatta il codice e nome prodotto, calcola il prezzo totale per riga
// ed il totale complessivo dell'ordine, includendo eventuali note a margine.
func renderRigheEntrata(o api.OrdineEntrata) string {
	var sb strings.Builder
	indent := "       "

	sb.WriteString(indent)
	sb.WriteString(style.HeaderStile.Render("┌─ DETTAGLI RIGHE ORDINE " + o.Numero))
	sb.WriteString("\n")

	if len(o.Righe) == 0 {
		sb.WriteString(indent)
		sb.WriteString("│  (Nessuna riga presente)\n")
	} else {
		sb.WriteString(indent)
		fmt.Fprintf(&sb, "│  %-10s %-25s %-10s %-12s %s\n",
			"CODICE", "PRODOTTO", "QUANTITÀ", "PREZZO UNIT.", "TOTALE RIGA")
		sb.WriteString(indent)
		sb.WriteString("├")
		sb.WriteString(strings.Repeat("─", 65))
		sb.WriteString("\n")

		var totaleOrdine float64
		for _, r := range o.Righe {
			codice := "N/D"
			nome := fmt.Sprintf("Prodotto #%d", r.ProdottoID)
			if r.Prodotto != nil {
				codice = r.Prodotto.Codice
				if r.Prodotto.Nome != "" {
					nome = r.Prodotto.Nome
				}
			}

			totRiga := float64(r.Quantita) * r.PrezzoUnitario
			totaleOrdine += totRiga

			sb.WriteString(indent)
			fmt.Fprintf(&sb, "│  %-10s %-25s %-10d € %-10.2f € %.2f\n",
				troncaStringa(codice, 10),
				troncaStringa(nome, 25),
				r.Quantita,
				r.PrezzoUnitario,
				totRiga)
		}
		sb.WriteString(indent)
		sb.WriteString("├")
		sb.WriteString(strings.Repeat("─", 65))
		sb.WriteString("\n")
		sb.WriteString(indent)
		fmt.Fprintf(&sb, "│  TOTALE COMPLESSIVO ORDINE: € %.2f\n", totaleOrdine)
	}

	if o.Note != "" {
		sb.WriteString(indent)
		sb.WriteString("│  Note: ")
		sb.WriteString(o.Note)
		sb.WriteString("\n")
	}
	sb.WriteString(indent)
	sb.WriteString("└")
	sb.WriteString(strings.Repeat("─", 65))
	sb.WriteString("\n")

	return sb.String()
}

// loadOrdiniEntrataCmd effettua in modo asincrono la chiamata al client API per recuperare la lista degli ordini in entrata.
//
// Come è stato scritto il codice:
// Restituisce una closure tea.Cmd che invoca `m.Client.GetOrdiniEntrata()`. Restituisce `OrdiniEntrataLoadedMsg` in caso di successo
// oppure `ErroreMsg` in caso di errore.
func (m Model) loadOrdiniEntrataCmd() tea.Cmd {
	return func() tea.Msg {
		ordini, err := m.Client.GetOrdiniEntrata()
		if err != nil {
			return ErroreMsg{Err: err}
		}
		return OrdiniEntrataLoadedMsg{Ordini: ordini}
	}
}

// troncaStringa tronca una stringa a `maxLen` caratteri aggiungendo i puntini di sospensione se la lunghezza viene superata.
func troncaStringa(s string, maxLen int) string {
	if len(s) > maxLen {
		if maxLen > 3 {
			return s[:maxLen-3] + "..."
		}
		return s[:maxLen]
	}
	return s
}
