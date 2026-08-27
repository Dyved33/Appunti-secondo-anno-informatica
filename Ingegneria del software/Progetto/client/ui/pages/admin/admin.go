package admin

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/admin/dashboard"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/admin/users"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/contabilita/catalogo"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/contabilita/fatture"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/contabilita/ordini"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/contabilita/report"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/magazzino/giacenze"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/magazzino/ordinientrata"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/magazzino/ordiniuscita"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/magazzino/resi"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/magazzino/scanner"
	style "github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/styles"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/tabnav"
)

// Init inizializza il layout Admin avviando il sotto-modello della sotto-pagina attiva di default (dashboard).
//
// Come è stato scritto il codice:
// Rispetta l'interfaccia tea.Model delegando la chiamata di inizializzazione iniziale al metodo Init()
// del sotto-modello giacenze per avviare il caricamento asincrono dei dati.
func (m AdminLayout) Init() tea.Cmd {
	return m.dashboard.Init()
}

// Update elabora i messaggi ricevuti ed effettua lo smistamento alla sotto-pagina selezionata o gestisce il cambio tab.
//
// Come è stato scritto il codice:
// 1. Se il messaggio è un click del mouse su uno dei tab della barra di navigazione:
//   - Aggiorna il valore di `m.attiva` con la corrispondente costante focus.
//   - Esegue la chiamata `Init()` sulla sotto-pagina selezionata per garantire il refresh dei dati in ingresso.
//
// 2. Se il messaggio non rappresenta un click su un tab, delega l'elaborazione al metodo `UpdateActive(msg)`.
func (m AdminLayout) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if click, ok := msg.(tea.MouseClickMsg); ok && click.Button == tea.MouseLeft {
		// Coordinate assolute del click → indice del tab colpito
		if i := tabnav.HitAt(m.hitboxes(), click.X, click.Y); i >= 0 {
			return m.switchTo(focus(i))
		}
		return m.UpdateActive(msg)
	}
	return m.UpdateActive(msg)
}

// switchTo attiva la sotto-pagina indicata e ne restituisce il comando di init.
func (m AdminLayout) switchTo(dest focus) (tea.Model, tea.Cmd) {
	m.attiva = dest
	switch dest {
	case DASHBOARD:
		return m, m.dashboard.Init()
	case USERS:
		return m, m.users.Init()
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
	case CATALOGO:
		return m, m.catalogo.Init()
	case FATTURE:
		return m, m.fatture.Init()
	case ORDINI:
		return m, m.ordini.Init()
	case REPORT:
		return m, m.report.Init()
	}
	return m, nil
}

// hitboxes calcola le aree cliccabili dei tab nelle coordinate schermo correnti.
func (m AdminLayout) hitboxes() []tabnav.Rect {
	return tabnav.Rects(m.winW, m.winH, labels[:], int(m.attiva), m.rawPage())
}

// rawPage compone la pagina non incorniciata ("tabs\n\ncorpo") della sola
// sotto-pagina attiva: è la stessa stringa incorniciata da View, quindi le
// hitbox calcolate su di essa coincidono con ciò che l'utente vede a schermo.
func (m AdminLayout) rawPage() string {
	var body tea.View
	switch m.attiva {
	case DASHBOARD:
		body = m.dashboard.View()
	case USERS:
		body = m.users.View()
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
	case CATALOGO:
		body = m.catalogo.View()
	case FATTURE:
		body = m.fatture.View()
	case ORDINI:
		body = m.ordini.View()
	case REPORT:
		body = m.report.View()
	}
	return m.tabs() + "\n\n" + body.Content
}

// activeCursor restituisce il cursore della sotto-pagina attualmente visibile.
func (m AdminLayout) activeCursor() *tea.Cursor {
	switch m.attiva {
	case DASHBOARD:
		return m.dashboard.View().Cursor
	case USERS:
		return m.users.View().Cursor
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
	case CATALOGO:
		return m.catalogo.View().Cursor
	case FATTURE:
		return m.fatture.View().Cursor
	case ORDINI:
		return m.ordini.View().Cursor
	case REPORT:
		return m.report.View().Cursor
	}
	return nil
}

// DimUpdate propaga i messaggi di resize a tutte le sotto-pagine dell'area admin
// e conserva le dimensioni del terminale per il calcolo dei tab cliccabili.
func (m *AdminLayout) DimUpdate(msg tea.WindowSizeMsg) {
	m.winW, m.winH = msg.Width, msg.Height

	updated, _ := m.dashboard.Update(msg)
	m.dashboard = updated.(dashboard.Model)

	updated, _ = m.users.Update(msg)
	m.users = updated.(users.Model)

	updated, _ = m.giacenze.Update(msg)
	m.giacenze = updated.(giacenze.Model)

	updated, _ = m.ordiniEntrata.Update(msg)
	m.ordiniEntrata = updated.(ordinientrata.Model)

	updated, _ = m.ordiniUscita.Update(msg)
	m.ordiniUscita = updated.(ordiniuscita.Model)

	updated, _ = m.resi.Update(msg)
	m.resi = updated.(resi.Model)

	updated, _ = m.scanner.Update(msg)
	m.scanner = updated.(scanner.Model)

	updated, _ = m.catalogo.Update(msg)
	m.catalogo = updated.(catalogo.Model)

	updated, _ = m.fatture.Update(msg)
	m.fatture = updated.(fatture.Model)

	updated, _ = m.ordini.Update(msg)
	m.ordini = updated.(ordini.Model)

	updated, _ = m.report.Update(msg)
	m.report = updated.(report.Model)
}

// UpdateActive inoltra l'evento o messaggio esclusivamente alla sotto-pagina attualmente visibile.
//
// Come è stato scritto il codice:
// Utilizza uno switch sul valore di `m.attiva` ed inoltra `msg` al metodo Update() del rispettivo
// sotto-modello, aggiornando l'istanza contenuta in `AdminLayout` e restituendo il comando generato.
func (m AdminLayout) UpdateActive(msg tea.Msg) (tea.Model, tea.Cmd) {
	var updated tea.Model
	var cmd tea.Cmd

	switch m.attiva {
	case DASHBOARD:
		updated, cmd = m.dashboard.Update(msg)
		m.dashboard = updated.(dashboard.Model)
	case USERS:
		updated, cmd = m.users.Update(msg)
		m.users = updated.(users.Model)
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
	case CATALOGO:
		updated, cmd = m.catalogo.Update(msg)
		m.catalogo = updated.(catalogo.Model)
	case FATTURE:
		updated, cmd = m.fatture.Update(msg)
		m.fatture = updated.(fatture.Model)
	case ORDINI:
		updated, cmd = m.ordini.Update(msg)
		m.ordini = updated.(ordini.Model)
	case REPORT:
		updated, cmd = m.report.Update(msg)
		m.report = updated.(report.Model)
	}

	return m, cmd
}

// View compone ed impagina l'interfaccia utente complessiva del modulo Admin.
//
// Come è stato scritto il codice:
//  1. Invoca `m.tabs()` per generare la stringa formattata contenente la barra di navigazione a schede.
//  2. In base alla sezione attiva `m.attiva`, invoca la funzione `View()` del sotto-modello corrispondente.
//  3. Concatena l'header dei tab al contenuto della vista, applica la cornice globale `style.FinestraStile`
//     e preserva l'eventuale posizione del cursore restituita dalla sotto-pagina attiva.
func (m AdminLayout) View() tea.View {
	v := tea.NewView(style.FinestraStile.Render(m.rawPage()))
	v.Cursor = m.activeCursor()
	return v
}

// tabs costruisce ed applica lo stile alla riga dei tab di navigazione superiore.
//
// Come è stato scritto il codice:
// Scorre la slice delle etichette `labels`. Per la scheda corrispondente al focus attivo
// applica `style.ActiveTab`, per le altre schede applica `style.TabGap`, e restituisce la stringa congiunta.
func (m AdminLayout) tabs() string {
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
func (m AdminLayout) Hitboxes() []tabnav.Rect { return m.hitboxes() }

// Attiva restituisce l'indice della sotto-pagina visibile (usata dai test d'integrazione).
func (m AdminLayout) Attiva() int { return int(m.attiva) }

// Size restituisce le dimensioni del terminale con cui il layout calcola le
// aree cliccabili (usata dai test d'integrazione per verificare la propagazione).
func (m AdminLayout) Size() (int, int) { return m.winW, m.winH }
