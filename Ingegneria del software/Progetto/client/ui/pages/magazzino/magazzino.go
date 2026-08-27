// Package magazzino gestisce la vista principale ed il rendering dei componenti TUI per la sezione Magazzino.
package magazzino

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/magazzino/giacenze"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/magazzino/ordinientrata"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/magazzino/ordiniuscita"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/magazzino/resi"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/magazzino/scanner"
	style "github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/styles"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/tabnav"
)

// Init inizializza il layout Magazzino avviando il sotto-modello della sotto-pagina attiva di default (Giacenze).
//
// Come è stato scritto il codice:
// Rispetta l'interfaccia tea.Model delegando la chiamata di inizializzazione iniziale al metodo Init()
// del sotto-modello giacenze per avviare il caricamento asincrono dei dati.
func (m MagLayout) Init() tea.Cmd {
	return m.giacenze.Init()
}

// Update elabora i messaggi ricevuti ed effettua lo smistamento alla sotto-pagina selezionata o gestisce il cambio tab.
//
// Come è stato scritto il codice:
// 1. Se il messaggio è un click del mouse su uno dei tab della barra di navigazione:
//   - Aggiorna il valore di `m.attiva` con la corrispondente costante focus.
//   - Esegue la chiamata `Init()` sulla sotto-pagina selezionata per garantire il refresh dei dati in ingresso.
//
// 2. Se arriva il messaggio speciale "SCANNER" (trigger hardware), attiva la pagina Scanner.
//
// 3. Per qualsiasi altro messaggio, delega l'elaborazione al metodo `UpdateActive(msg)`.
func (m MagLayout) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if click, ok := msg.(tea.MouseClickMsg); ok && click.Button == tea.MouseLeft {
		// Coordinate assolute del click → indice del tab colpito
		if i := tabnav.HitAt(m.hitboxes(), click.X, click.Y); i >= 0 {
			return m.switchTo(focus(i))
		}
		return m.UpdateActive(msg)
	}

	if key, ok := msg.(tea.KeyPressMsg); ok && key.String() == "SCANNER" {
		m.attiva = SCANNER
		return m, m.scanner.Init()
	}

	return m.UpdateActive(msg)
}

// switchTo attiva la sotto-pagina indicata e ne restituisce il comando di init.
func (m MagLayout) switchTo(dest focus) (tea.Model, tea.Cmd) {
	m.attiva = dest
	switch dest {
	case GIACENZE:
		return m, m.giacenze.Init()
	case ORDINIENTRATA:
		return m, m.ordiniEntrata.Init()
	case ORDINIUSCITA:
		return m, m.ordiniUscita.Init()
	case RESI:
		return m, m.resi.Init()
	case SCANNER:
		return m, m.scanner.Init()
	}
	return m, nil
}

// hitboxes calcola le aree cliccabili dei tab nelle coordinate schermo correnti.
func (m MagLayout) hitboxes() []tabnav.Rect {
	return tabnav.Rects(m.winW, m.winH, labels[:], int(m.attiva), m.rawPage())
}

// rawPage compone la pagina non incorniciata ("tabs\n\ncorpo") della sola
// sotto-pagina attiva: è la stessa stringa incorniciata da View, quindi le
// hitbox calcolate su di essa coincidono con ciò che l'utente vede a schermo.
func (m MagLayout) rawPage() string {
	var body tea.View
	switch m.attiva {
	case GIACENZE:
		body = m.giacenze.View()
	case ORDINIENTRATA:
		body = m.ordiniEntrata.View()
	case ORDINIUSCITA:
		body = m.ordiniUscita.View()
	case RESI:
		body = m.resi.View()
	case SCANNER:
		body = m.scanner.View()
	}
	return m.tabs() + "\n\n" + body.Content
}

// DimUpdate propaga i messaggi di resize a tutte le sotto-pagine del magazzino
// e conserva le dimensioni del terminale per il calcolo dei tab cliccabili.
func (m *MagLayout) DimUpdate(msg tea.WindowSizeMsg) {
	m.winW, m.winH = msg.Width, msg.Height

	updated, _ := m.giacenze.Update(msg)
	m.giacenze = updated.(giacenze.Model)

	updated, _ = m.ordiniEntrata.Update(msg)
	m.ordiniEntrata = updated.(ordinientrata.Model)

	updated, _ = m.ordiniUscita.Update(msg)
	m.ordiniUscita = updated.(ordiniuscita.Model)

	updated, _ = m.resi.Update(msg)
	m.resi = updated.(resi.Model)

	updated, _ = m.scanner.Update(msg)
	m.scanner = updated.(scanner.Model)
}

// UpdateActive inoltra l'evento o messaggio esclusivamente alla sotto-pagina attualmente visibile.
//
// Come è stato scritto il codice:
// Utilizza uno switch sul valore di `m.attiva` ed inoltra `msg` al metodo Update() del rispettivo
// sotto-modello, aggiornando l'istanza contenuta in `MagLayout` e restituendo il comando generato.
func (m MagLayout) UpdateActive(msg tea.Msg) (tea.Model, tea.Cmd) {
	var updated tea.Model
	var cmd tea.Cmd

	switch m.attiva {
	case GIACENZE:
		updated, cmd = m.giacenze.Update(msg)
		m.giacenze = updated.(giacenze.Model)
	case ORDINIENTRATA:
		updated, cmd = m.ordiniEntrata.Update(msg)
		m.ordiniEntrata = updated.(ordinientrata.Model)
	case ORDINIUSCITA:
		updated, cmd = m.ordiniUscita.Update(msg)
		m.ordiniUscita = updated.(ordiniuscita.Model)
	case RESI:
		updated, cmd = m.resi.Update(msg)
		m.resi = updated.(resi.Model)
	case SCANNER:
		updated, cmd = m.scanner.Update(msg)
		m.scanner = updated.(scanner.Model)
	}

	return m, cmd
}

// View compone ed impagina l'interfaccia utente complessiva del modulo Magazzino.
//
// Come è stato scritto il codice:
//  1. Invoca `m.tabs()` per generare la stringa formattata contenente la barra di navigazione a schede.
//  2. In base alla sezione attiva `m.attiva`, invoca la funzione `View()` del sotto-modello corrispondente.
//  3. Concatena l'header dei tab al contenuto della vista, applica la cornice globale `style.FinestraStile`
//     e preserva l'eventuale posizione del cursore restituita dalla sotto-pagina attiva.
func (m MagLayout) View() tea.View {
	page := m.rawPage()

	v := tea.NewView(style.FinestraStile.Render(page))
	v.Cursor = m.activeCursor()
	return v
}

// activeCursor restituisce il cursore della sotto-pagina attualmente visibile.
func (m MagLayout) activeCursor() *tea.Cursor {
	switch m.attiva {
	case GIACENZE:
		return m.giacenze.View().Cursor
	case ORDINIENTRATA:
		return m.ordiniEntrata.View().Cursor
	case ORDINIUSCITA:
		return m.ordiniUscita.View().Cursor
	case RESI:
		return m.resi.View().Cursor
	case SCANNER:
		return m.scanner.View().Cursor
	}
	return nil
}

// tabs costruisce ed applica lo stile alla riga dei tab di navigazione superiore.
//
// Come è stato scritto il codice:
// Scorre la slice delle etichette `labels`. Per la scheda corrispondente al focus attivo
// applica `style.ActiveTab`, per le altre schede applica `style.TabGap`, e restituisce la stringa congiunta.
func (m MagLayout) tabs() string {
	parts := make([]string, len(labels))
	for i, label := range labels {
		if focus(i) == m.attiva {
			parts[i] = style.ActiveTab.Render(label)
		} else {
			parts[i] = style.TabGap.Render(label)
		}
	}

	// Unisce tutti gli elementi della slice orizzontalmente con uno spazio di 2 caratteri in mezzo
	return lipgloss.JoinHorizontal(lipgloss.Top, joinWithSpacing(parts, "  ")...)
}

// Funzione di supporto per inserire uno spaziatore tra i tab
func joinWithSpacing(elements []string, spacer string) []string {
	if len(elements) == 0 {
		return nil
	}
	result := []string{elements[0]}
	for _, el := range elements[1:] {
		result = append(result, spacer, el)
	}
	return result
}

// Hitboxes espone le aree cliccabili dei tab (usata dai test d'integrazione dell'app).
func (m MagLayout) Hitboxes() []tabnav.Rect { return m.hitboxes() }

// Attiva restituisce l'indice della sotto-pagina visibile (usata dai test d'integrazione).
func (m MagLayout) Attiva() int { return int(m.attiva) }

// Size restituisce le dimensioni del terminale con cui il layout calcola le
// aree cliccabili (usata dai test d'integrazione per verificare la propagazione).
func (m MagLayout) Size() (int, int) { return m.winW, m.winH }
