// Package ordiniuscita gestisce la visualizzazione ed interazione per gli ordini in uscita (spedizioni e DDT clienti).
package ordiniuscita

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
	style "github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/styles"
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

// Init avvia l'inizializzazione della pagina richiedendo l'elenco ordini in uscita dal backend.
//
// Come è stato scritto il codice:
// Restituisce la funzione closure `loadOrdiniUscitaCmd()` come comando tea.Cmd per avviare il fetch asincrono.
func (m Model) Init() tea.Cmd {
	return m.loadOrdiniUscitaCmd()
}

// Update gestisce gli eventi utente, la tastiera ed i messaggi di risposta API per gli ordini in uscita.
//
// Come è stato scritto il codice:
// 1. Su `OrdiniUscitaLoadedMsg` salva l'array `ordini`, azzera cursore e scorrimento e disattiva `loading`.
// 2. Su `ErroreMsg` registra il problema e disattiva `loading`.
// 3. Su `tea.WindowSizeMsg` conserva le nuove dimensioni e ricalibra la finestra visibile.
// 4. Su `tea.MouseWheelMsg` (vista lista) scorre nella direzione della rotellina.
// 5. Su `tea.KeyPressMsg`:
//   - "up"/"k" e "down"/"j": muovono il cursore mantenendolo dentro la porzione visibile.
//   - "enter"/"space": commuta l'espansione dei dettagli dell'ordine evidenziato.
//   - "pgup"/"pgdown": scorrono la lista di una pagina; "home"/"end": tornano all'inizio/fine.
//   - "r": imposta `loading = true` e ricarica i dati tramite `loadOrdiniUscitaCmd()`.
//   - "ctrl+d": genera/il DDT dell'ordine selezionato; "esc": torna alla lista dalla vista DDT.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case OrdiniUscitaLoadedMsg:
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

	case DDTCreatedMsg:
		m.loading = false
		m.err = nil
		m.DDT = msg.DDT
		m.view = DDT
		return m, nil

	case OrdineAggiornatoMsg:
		m.loading = false
		m.err = nil
		m.modo = modoLista
		return m, m.loadOrdiniUscitaCmd()

	case tea.WindowSizeMsg:
		m.winW = msg.Width
		m.winH = msg.Height
		if m.view == MAIN {
			m.clampScroll()
		}
		return m, nil

	case tea.MouseWheelMsg:
		if m.view == MAIN && m.modo == modoLista {
			switch msg.Button {
			case tea.MouseWheelUp:
				m.scorri(-passoRotellina)
			case tea.MouseWheelDown:
				m.scorri(passoRotellina)
			}
		}
		return m, nil

	case tea.KeyPressMsg:
		if m.modo == modoModifica {
			return m.handleEditInput(msg)
		}

		switch msg.String() {
		case "up", "k":
			if m.view == MAIN && m.cursore > 0 {
				m.cursore--
			}
			m.mostraCursore()
		case "down", "j":
			if m.view == MAIN && m.cursore < len(m.ordini)-1 {
				m.cursore++
			}
			m.mostraCursore()
		case "enter", "space":
			if m.view == MAIN && len(m.ordini) > 0 {
				id := m.ordini[m.cursore].ID
				m.dettagliEspansi[id] = !m.dettagliEspansi[id]
			}
			m.mostraCursore()
		case "pgup":
			if m.view == MAIN {
				m.scorri(-m.righeVisibili())
			}
		case "pgdown":
			if m.view == MAIN {
				m.scorri(m.righeVisibili())
			}
		case "home":
			if m.view == MAIN {
				m.scrollOffset = 0
			}
		case "end":
			if m.view == MAIN {
				m.scorri(len(m.ordini))
			}
		case "r":
			m.loading = true
			m.scrollOffset = 0
			return m, m.loadOrdiniUscitaCmd()

		case "ctrl+d":
			m.loading = true
			return m, m.showDettagliDDT()

		case "ctrl+e":
			if m.view == MAIN && len(m.ordini) > 0 {
				m.modo = modoModifica
				m.err = nil
				m.preparaFormModifica()
				return m, m.focusCampoForm()
			}

		case "esc":
			m.err = nil
			if m.view == DDT {
				m.view = MAIN
			}
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

		statoPuro := fmt.Sprintf("%-15s", o.Stato)
		statoBadge := renderingStato(o.Stato, statoPuro)

		rigaTesto := fmt.Sprintf("%-20s %-30s %-15s %s %s",
			truncaStringa(o.Numero, 20),
			truncaStringa(o.Cliente, 30),
			dataStr,
			statoBadge,
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
			for _, l := range strings.Split(strings.TrimRight(renderRigheUscita(o), "\n"), "\n") {
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
	return disponibili >= minRigheConKPI
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

// View genera ed impagina l'interfaccia utente stilizzata per gli ordini in uscita.
//
// Come è stato scritto il codice:
//  1. In caso di attesa o errori visualizza lo stato corrente tramite componenti stilizzati.
//  2. Calcola i 4 KPI riassuntivi (Totale, Confermati, Evasi, Con DDT associato).
//  3. Mostra solo la finestra di righe che entra nello schermo (`scrollOffset` + `righeVisibili`),
//     rendendo raggiungibili anche gli ordini oltre la capacità del terminale.
//  4. Se un ordine è espanso (`dettagliEspansi[id] == true`), include la sotto-tabella prodotti.
func (m Model) View() tea.View {
	var sb strings.Builder

	// ── Gestione Caricamento ed Errori ──────────────────────────────────────
	if m.loading {
		sb.WriteString(style.TitleStile.Render("ORDINI IN USCITA"))
		sb.WriteString("\n\n")
		sb.WriteString("Caricamento ordini in uscita in corso...\n")
		return tea.NewView(sb.String())
	}
	if m.err != nil {
		sb.WriteString(style.TitleStile.Render("ORDINI IN USCITA"))
		sb.WriteString("\n\n")
		sb.WriteString(style.ErrorStile.Render("Errore: " + m.err.Error()))
		sb.WriteString("\n\n")
		if m.modo == modoModifica {
			sb.WriteString(style.FooterStile.Render("[tab] Campo successivo  |  [←/→] Stato  |  [enter] Salva  |  [esc] Annulla"))
		} else {
			sb.WriteString(style.FooterStile.Render("[esc] Torna alla lista ordini"))
		}
		return tea.NewView(sb.String())
	}
	if m.view == DDT {
		return m.renderDDT()
	}

	if m.modo == modoModifica {
		sb.WriteString(m.renderEditForm())
		return tea.NewView(sb.String())
	}

	// ── Intestazione fissa ───────────────────────────────────────────────────
	sb.WriteString(m.renderIntestazione(m.mostraKPI()))

	// ── Finestra visibile della tabella ──────────────────────────────────────
	righe := m.listaRighe()

	if len(righe) == 0 {
		sb.WriteString("Nessun ordine in uscita registrato.\n")
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
func (m Model) renderIntestazione(conKPI bool) string {
	var sb strings.Builder

	// Il MarginBottom di TitleStile fornisce già la riga vuota sotto il titolo:
	// qui basta un solo "\n" esplicito per non sprecare righe di terminale.
	sb.WriteString(style.TitleStile.Render("ORDINI IN USCITA (CLIENTI / SPEDIZIONI)"))
	sb.WriteString("\n")

	if conKPI {
		var confermati, evasi, annullati, conDDT int
		for _, o := range m.ordini {
			switch o.Stato {
			case api.StatoConfermato:
				confermati++
			case api.StatoEvaso:
				evasi++
			case api.StatoAnnullato:
				annullati++
			}
			if o.DDTID != nil {
				conDDT++
			}
		}

		kpiTot := style.BoxStile.Render(fmt.Sprintf("TOTALE ORDINI\n\n%d", len(m.ordini)))
		kpiConf := style.BoxStile.Render(fmt.Sprintf("CONFERMATI\n\n%d", confermati))
		kpiEvas := style.BoxStile.Render(fmt.Sprintf("EVASI\n\n%d", evasi))
		kpiDDT := style.BoxStile.Render(fmt.Sprintf("CON DDT\n\n%d", conDDT))

		rigaKPI := lipgloss.JoinHorizontal(lipgloss.Top, kpiTot, "   ", kpiConf, "   ", kpiEvas, "   ", kpiDDT)
		sb.WriteString(rigaKPI)
		sb.WriteString("\n\n")
	}

	headerStr := fmt.Sprintf("   %-24s %-30s %-15s %-15s %s",
		"NUMERO", "CLIENTE", "DATA", "STATO", "RIGHE")
	sb.WriteString(style.FreeHeaderStile.Render(headerStr))
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("─", 100))
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
	scorciatoie := "[↑/↓ o k/j] Naviga | [ctrl+d] Genera DDT | [ctrl+e] Modifica | [enter/space] Espandi/Comprimi Dettaglio | [pgup/pgdn o rotellina] Scorri | [r] Ricarica | [ctrl+l] Logout"

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

// renderDDT compone la vista del documento di trasporto (DDT) mostrandone
// i dati principali (numero, data, mittente, destinatario e percorso PDF).
func (m Model) renderDDT() tea.View {
	var sb strings.Builder

	sb.WriteString(style.TitleStile.Render(fmt.Sprintf("DOCUMENTO DI TRASPORTO (DDT #%d)", m.DDT.ID)))
	sb.WriteString("\n\n")

	// Renderizza le informazioni del DDT memorizzato in m.DDT
	sb.WriteString(fmt.Sprintf("Numero: %s\n", m.DDT.Numero))
	sb.WriteString(fmt.Sprintf("Data:         %s\n", m.DDT.Data.Format("02/01/2006")))
	sb.WriteString(fmt.Sprintf("Mittente:     %s\n", m.DDT.Mittente))
	sb.WriteString(fmt.Sprintf("Destinatario: %s\n", m.DDT.Destinatario))
	sb.WriteString(fmt.Sprintf("PDF: %s\n", m.DDT.PathPDF))
	sb.WriteString("\n")

	// Footer per uscire dalla vista DDT
	sb.WriteString(style.FooterStile.Render("[esc] Torna alla lista ordini"))

	return tea.NewView(sb.String())
}

// renderRigheUscita costruisce la vista dettagliata espansa per un singolo ordine di uscita.
//
// Come è stato scritto il codice:
// Inserisce una sotto-griglia contenente codice prodotto, quantita, prezzo ed importo riga,
// indicando anche l'eventuale identificativo DDT ed eventuali note sulla spedizione.
func renderRigheUscita(o api.OrdineUscita) string {
	var sb strings.Builder
	indent := "       "
	sb.WriteString(indent)
	sb.WriteString(style.HeaderStile.Render("┌─ DETTAGLI RIGHE ORDINE USCITA " + o.Numero))
	sb.WriteString("\n")

	if len(o.Righe) == 0 {
		sb.WriteString(indent)
		sb.WriteString("│  (Nessuna riga presente)\n")
	} else {
		sb.WriteString(indent)
		sb.WriteString(fmt.Sprintf("│  %-10s %-25s %-10s %-12s %s\n",
			"CODICE", "PRODOTTO", "QUANTITÀ", "PREZZO UNIT.", "TOTALE RIGA"))
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
			sb.WriteString(fmt.Sprintf("│  %-10s %-25s %-10d € %-10.2f € %.2f\n",
				truncaStringa(codice, 10),
				truncaStringa(nome, 25),
				r.Quantita,
				r.PrezzoUnitario,
				totRiga,
			))
		}
		sb.WriteString(indent)
		sb.WriteString("├")
		sb.WriteString(strings.Repeat("─", 65))
		sb.WriteString("\n")
		sb.WriteString(indent)
		sb.WriteString(fmt.Sprintf("│  TOTALE ORDINE USCITA: € %.2f\n", totaleOrdine))
	}

	if o.DDTID != nil {
		sb.WriteString(indent)
		sb.WriteString(fmt.Sprintf("│  Riferimento Documento di Trasporto: DDT #%d\n", *o.DDTID))
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

// renderingStato converte lo StatoOrdine nel badge grafico corrispondente.
func renderingStato(s api.StatoOrdine, testoFormattato string) string {
	switch s {
	case api.StatoConfermato:
		return style.InfoStile.Render(testoFormattato)
	case api.StatoEvaso:
		return style.SuccessStile.Render(testoFormattato)
	case api.StatoAnnullato:
		return style.ErrorStile.Render(testoFormattato)
	default:
		return style.InfoStile.Render(testoFormattato)
	}
}

// loadOrdiniUscitaCmd crea un comando tea.Cmd asincrono per recuperare gli ordini in uscita.
//
// Come è stato scritto il codice:
// Invoca `m.Client.GetOrdiniUscita()` e restituisce un `OrdiniUscitaLoadedMsg` o un `ErroreMsg`.
func (m Model) loadOrdiniUscitaCmd() tea.Cmd {
	return func() tea.Msg {
		ordini, err := m.Client.GetOrdiniUscita()
		if err != nil {
			return ErroreMsg{Err: err}
		}
		return OrdiniUscitaLoadedMsg{Ordini: ordini}
	}
}

// truncaStringa tronca la stringa alla lunghezza specificata `maxLen`.
func truncaStringa(s string, maxLen int) string {
	if len(s) > maxLen {
		if maxLen > 3 {
			return s[:maxLen-3] + "..."
		}
		return s[:maxLen]
	}
	return s
}

// showDettagliDDT crea un comando asincrono che recupera il DDT dell'ordine
// selezionato dal backend, emettendo DDTCreatedMsg in caso di successo o ErroreMsg.
func (m Model) showDettagliDDT() tea.Cmd {
	return func() tea.Msg {
		if len(m.ordini) == 0 {
			return ErroreMsg{Err: fmt.Errorf("nessun ordine disponibile per visualizzare il DDT")}
		}
		ordineSelezionato := m.ordini[m.cursore]
		DTO := api.ClientDDTDTO{
			Mittente:     ordineSelezionato.Mittente,
			Destinatario: ordineSelezionato.Cliente,
		}

		DDT, err := m.Client.GetDDT(ordineSelezionato.ID, DTO)
		if err != nil {
			return ErroreMsg{Err: err}
		}
		return DDTCreatedMsg{DDT: DDT}
	}
}

// handleEditInput gestisce gli eventi tastiera in modalità di modifica,
// spostando il focus tra i campi, selezionando lo stato e confermando il salvataggio.
func (m Model) handleEditInput(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.modo = modoLista
		m.err = nil
		return m, nil
	case "tab", "down":
		if m.focusMod == focusStato {
			m.focusMod = focusNumero
		} else {
			m.focusMod++
		}
		return m, m.focusCampoForm()
	case "shift+tab", "up":
		if m.focusMod == focusNumero {
			m.focusMod = focusStato
		} else {
			m.focusMod--
		}
		return m, m.focusCampoForm()
	case "left":
		if m.focusMod == focusStato {
			m.idxStato = (m.idxStato - 1 + len(statiOrdine)) % len(statiOrdine)
			return m, nil
		}
	case "right":
		if m.focusMod == focusStato {
			m.idxStato = (m.idxStato + 1) % len(statiOrdine)
			return m, nil
		}
	case "enter":
		if m.focusMod == focusStato {
			return m.salvaOrdineUscita()
		}
		m.focusMod++
		if m.focusMod > focusStato {
			m.focusMod = focusNumero
		}
		return m, m.focusCampoForm()
	}

	var cmd tea.Cmd
	switch m.focusMod {
	case focusNumero:
		m.formNumero, cmd = m.formNumero.Update(msg)
	case focusCliente:
		m.formCliente, cmd = m.formCliente.Update(msg)
	case focusNote:
		m.formNote, cmd = m.formNote.Update(msg)
	}
	return m, cmd
}

// renderEditForm impagina il form di modifica dell'ordine in uscita,
// evidenziando con lo stile di focus il campo attualmente selezionato.
func (m Model) renderEditForm() string {
	numeroStr := m.formNumero.View()
	clienteStr := m.formCliente.View()
	noteStr := m.formNote.View()

	statoStr := statiOrdine[m.idxStato]
	if m.focusMod == focusStato {
		statoStr = style.FocusStile.Render("< " + statoStr + " >")
	}

	if m.focusMod == focusNumero {
		numeroStr = style.FocusStile.Render(numeroStr)
	}
	if m.focusMod == focusCliente {
		clienteStr = style.FocusStile.Render(clienteStr)
	}
	if m.focusMod == focusNote {
		noteStr = style.FocusStile.Render(noteStr)
	}

	content := strings.Join([]string{
		"MODIFICA ORDINE USCITA",
		"",
		"Numero:  " + numeroStr,
		"Cliente: " + clienteStr,
		"Note:    " + noteStr,
		"Stato:   " + statoStr,
	}, "\n")

	return style.LoginBoxStile.Render(content) + "\n\n" +
		style.FooterStile.Render("[tab] Campo successivo  |  [←/→] Stato  |  [enter] Salva  |  [esc] Annulla")
}

// salvaOrdineUscita raccoglie i valori del form ed invia l'aggiornamento dell'ordine
// selezionato al backend tramite un comando asincrono.
func (m Model) salvaOrdineUscita() (tea.Model, tea.Cmd) {
	if len(m.ordini) == 0 || m.cursore >= len(m.ordini) {
		return m, nil
	}
	o := m.ordini[m.cursore]

	numero := m.formNumero.Value()
	cliente := m.formCliente.Value()
	note := m.formNote.Value()
	stato := statiOrdine[m.idxStato]

	var righe []api.ClientOrderLineDTO
	for _, r := range o.Righe {
		righe = append(righe, api.ClientOrderLineDTO{
			ProdottoID:     r.ProdottoID,
			Quantita:       r.Quantita,
			PrezzoUnitario: r.PrezzoUnitario,
		})
	}

	m.loading = true
	m.err = nil
	return m, m.aggiornaOrdineUscitaCmd(o.ID, api.ClientOrderDTO{
		Numero:  numero,
		Cliente: cliente,
		Note:    note,
		Righe:   righe,
		Stato:   stato,
	})
}

// aggiornaOrdineUscitaCmd crea un comando asincrono che aggiorna l'ordine in uscita,
// emettendo OrdineAggiornatoMsg in caso di successo o ErroreMsg in caso di errore.
func (m Model) aggiornaOrdineUscitaCmd(id int, input api.ClientOrderDTO) tea.Cmd {
	return func() tea.Msg {
		if err := m.Client.AggiornaOrdineUscita(id, input); err != nil {
			return ErroreMsg{Err: err}
		}
		return OrdineAggiornatoMsg{}
	}
}
