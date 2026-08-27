// Package contabilita gestisce la vista principale ed il rendering dei componenti
// TUI per la sezione Contabilità.
package contabilita

import (
	catalogopkg "github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/contabilita/catalogo"
	fatturepkg "github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/contabilita/fatture"
	ordinipkg "github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/contabilita/ordini"
	reportpkg "github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/contabilita/report"
	style "github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/styles"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/tabnav"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// Init inizializza il layout Contabilità avviando la sotto-pagina attiva di default (Ordini).
//
// Come è stato scritto il codice:
// Rispetta l'interfaccia tea.Model delegando l'inizializzazione al metodo Init()
// del sotto-modello ordini, che avvia il caricamento asincrono dei dati.
func (c ContLayout) Init() tea.Cmd {
	return c.ordini.Init()
}

// Update smista i messaggi: prima intercetta i click sulla barra dei tab per
// cambiare sezione, poi inoltra tutto il resto al sotto-modello attualmente attivo.
//
// Come è stato scritto il codice:
//  1. Se il messaggio è un tea.MouseClickMsg su uno dei tab, aggiorna `c.attiva` e richiama
//     Init() sulla nuova sotto-pagina, per garantire dati sempre freschi al cambio tab
//     (stesso comportamento di MagLayout.Update).
//  2. Altrimenti delega l'update al sotto-modello attivo tramite updateSezioneAttiva,
//     con lo stesso pattern di type-assertion visto nel router principale dell'app.
func (c ContLayout) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if click, ok := msg.(tea.MouseClickMsg); ok && click.Button == tea.MouseLeft {
		// Coordinate assolute del click → indice della sezione colpita
		if i := tabnav.HitAt(c.hitboxes(), click.X, click.Y); i >= 0 {
			return c.switchTo(focus(i))
		}
		return c.updateSezioneAttiva(msg)
	}

	return c.updateSezioneAttiva(msg)
}

// switchTo attiva la sezione indicata e ne restituisce il comando di init.
func (c ContLayout) switchTo(dest focus) (tea.Model, tea.Cmd) {
	c.attiva = dest
	switch dest {
	case SezioneOrdini:
		return c, c.ordini.Init()
	case SezioneFatture:
		return c, c.fatture.Init()
	case SezioneCatalogo:
		return c, c.catalogo.Init()
	case SezioneReport:
		return c, c.report.Init()
	}
	return c, nil
}

// hitboxes calcola le aree cliccabili dei tab nelle coordinate schermo correnti.
func (c ContLayout) hitboxes() []tabnav.Rect {
	return tabnav.Rects(c.winW, c.winH, labels[:], int(c.attiva), c.rawPage())
}

// DimUpdate propaga i messaggi di resize a tutte le sotto-pagine di contabilità
// e conserva le dimensioni del terminale per il calcolo dei tab cliccabili.
func (c *ContLayout) DimUpdate(msg tea.WindowSizeMsg) {
	c.winW, c.winH = msg.Width, msg.Height

	updated, _ := c.ordini.Update(msg)
	c.ordini = updated.(ordinipkg.Model)

	updated, _ = c.fatture.Update(msg)
	c.fatture = updated.(fatturepkg.Model)

	updated, _ = c.catalogo.Update(msg)
	c.catalogo = updated.(catalogopkg.Model)

	updated, _ = c.report.Update(msg)
	c.report = updated.(reportpkg.Model)
}

// updateSezioneAttiva inoltra il messaggio solo al sotto-modello corrispondente
// alla sezione attualmente visibile; le altre tre restano "in pausa" e non
// ricevono l'evento, così non consumano messaggi che non le riguardano.
func (c ContLayout) updateSezioneAttiva(msg tea.Msg) (tea.Model, tea.Cmd) {
	var updated tea.Model
	var cmd tea.Cmd

	switch c.attiva {
	case SezioneOrdini:
		updated, cmd = c.ordini.Update(msg)
		c.ordini = updated.(ordinipkg.Model)
	case SezioneFatture:
		updated, cmd = c.fatture.Update(msg)
		c.fatture = updated.(fatturepkg.Model)
	case SezioneCatalogo:
		updated, cmd = c.catalogo.Update(msg)
		c.catalogo = updated.(catalogopkg.Model)
	case SezioneReport:
		updated, cmd = c.report.Update(msg)
		c.report = updated.(reportpkg.Model)
	}

	return c, cmd
}

// View disegna la barra dei tab in alto e sotto il contenuto della sezione attiva.
//
// Come è stato scritto il codice:
//  1. Invoca `c.tabs()` per generare la barra di navigazione con la sezione attiva evidenziata.
//  2. In base alla sezione attiva, invoca la funzione View() del sotto-modello corrispondente.
//  3. Concatena l'header dei tab al contenuto, applica la cornice globale `style.FinestraStile`
//     e propaga il cursore restituito dalla sotto-pagina (usato dal form prezzi di Ordini).
//
// La cornice viene applicata solo qui: le sotto-pagine restituiscono contenuto
// non incorniciato, come avviene per le sotto-pagine di magazzino.
func (c ContLayout) View() tea.View {
	v := tea.NewView(style.FinestraStile.Render(c.rawPage()))
	v.Cursor = c.activeCursor()
	return v
}

// rawPage compone la pagina non incorniciata ("tabs\n\ncorpo") della sola
// sotto-pagina attiva: è la stessa stringa incorniciata da View, quindi le
// hitbox calcolate su di essa coincidono con ciò che l'utente vede a schermo.
func (c ContLayout) rawPage() string {
	var body tea.View
	switch c.attiva {
	case SezioneOrdini:
		body = c.ordini.View()
	case SezioneFatture:
		body = c.fatture.View()
	case SezioneCatalogo:
		body = c.catalogo.View()
	case SezioneReport:
		body = c.report.View()
	}
	return c.tabs() + "\n\n" + body.Content
}

// activeCursor restituisce il cursore della sotto-pagina attualmente visibile.
func (c ContLayout) activeCursor() *tea.Cursor {
	switch c.attiva {
	case SezioneOrdini:
		return c.ordini.View().Cursor
	case SezioneFatture:
		return c.fatture.View().Cursor
	case SezioneCatalogo:
		return c.catalogo.View().Cursor
	case SezioneReport:
		return c.report.View().Cursor
	}
	return nil
}

// tabs costruisce ed applica lo stile alla riga dei tab di navigazione superiore.
//
// Come è stato scritto il codice:
// Scorre la slice delle etichette `labels`. Per la scheda corrispondente al focus attivo
// applica `style.ActiveTab`, per le altre schede applica `style.TabGap`, e restituisce la stringa congiunta.
func (m ContLayout) tabs() string {
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
func (c ContLayout) Hitboxes() []tabnav.Rect { return c.hitboxes() }

// Attiva restituisce l'indice della sezione visibile (usata dai test d'integrazione).
func (c ContLayout) Attiva() int { return int(c.attiva) }

// Size restituisce le dimensioni del terminale con cui il layout calcola le
// aree cliccabili (usata dai test d'integrazione per verificare la propagazione).
func (c ContLayout) Size() (int, int) { return c.winW, c.winH }
