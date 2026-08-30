// Package catalogo contiene le logiche di rendering e di gestione eventi della
// sotto-pagina Catalogo della Contabilità.
package catalogo

import (
	"fmt"
	"strconv"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"

	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
	style "github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/styles"
)

// passoRotellina è il numero di righe di cui la lista scorre ad ogni evento della rotellina del mouse.
const passoRotellina = 3

// overheadChrome è il numero di righe fisse occupate ATTORNO al corpo della pagina
// (logo, intestazione, barra dei tab, bordi della finestra e riga vuota interna).
const overheadChrome = 21

// righeFallback è il numero di righe visibili usato finché non arriva il primo WindowSizeMsg.
const righeFallback = 20

// righeVisibili calcola quante righe di tabella entrano nello spazio residuo dello schermo.
func (m Model) righeVisibili() int {
	if m.winH <= 0 {
		return righeFallback
	}
	// Titolo (4) + etichetta categoria (3) + testata tabella (2) + footer (2).
	disponibili := m.winH - overheadChrome - 11
	if disponibili < 1 {
		disponibili = 1
	}
	return disponibili
}

// clampScroll riporta l'offset nell'intervallo valido rispetto alla lista e allo schermo.
func (m *Model) clampScroll() {
	prodotti := m.prodottiFiltrati()
	maxOffset := len(prodotti) - m.righeVisibili()
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

// scorri avanza la finestra di visualizzazione di `n` righe, mantenendo il cursore visibile.
func (m *Model) scorri(n int) {
	prodotti := m.prodottiFiltrati()
	m.cursore += n
	if m.cursore < 0 {
		m.cursore = 0
	}
	if max := len(prodotti) - 1; m.cursore > max {
		m.cursore = max
	}
	m.mostraCursore()
}

// Init avvia il caricamento di prodotti e categorie.
func (m Model) Init() tea.Cmd {
	return m.loadCatalogoCmd()
}

// Update gestisce la navigazione nella lista filtrata, il cambio di categoria
// con "tab" e le modalità creazione/eliminazione.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case CatalogoLoadedMsg:
		m.loading = false
		m.err = nil
		m.prodotti = msg.Prodotti
		m.categorie = msg.Categorie
		m.categorieID = msg.CategorieIDs
		m.cursore = 0
		m.offset = 0
		return m, nil

	case ErroreMsg:
		m.loading = false
		m.err = msg.Err
		return m, nil

	case prodottoCreatoMsg:
		m.mode = modeView
		m.loading = false
		m.err = nil
		m.success = "Prodotto creato: " + msg.Nome
		m.resetCreateProductForm()
		return m, m.loadCatalogoCmd()

	case prodottoEliminatoMsg:
		m.mode = modeView
		m.loading = false
		m.err = nil
		m.success = "Prodotto eliminato: " + msg.Nome
		return m, m.loadCatalogoCmd()

	case categoriaCreataMsg:
		m.mode = modeView
		m.loading = false
		m.err = nil
		m.success = "Categoria creata: " + msg.Nome
		m.resetCreateCategoryForm()
		return m, m.loadCatalogoCmd()

	case categoriaEliminataMsg:
		m.mode = modeView
		m.loading = false
		m.err = nil
		m.success = "Categoria eliminata: " + msg.Nome
		m.filtroCategoria = nessunFiltro
		return m, m.loadCatalogoCmd()

	case tea.WindowSizeMsg:
		m.winW = msg.Width
		m.winH = msg.Height
		m.clampScroll()
		return m, nil

	case tea.MouseWheelMsg:
		if m.mode == modeView {
			switch msg.Button {
			case tea.MouseWheelUp:
				m.scorri(-passoRotellina)
			case tea.MouseWheelDown:
				m.scorri(passoRotellina)
			}
		}
		return m, nil

	case tea.KeyPressMsg:
		switch m.mode {
		case modeView:
			return m.handleViewInput(msg)
		case modeCreateProduct:
			return m.handleCreateProductInput(msg)
		case modeCreateCategory:
			return m.handleCreateCategoryInput(msg)
		case modeConfirmDeleteProduct, modeConfirmDeleteCategory:
			return m.handleConfirmDeleteInput(msg)
		}
	}

	return m, nil
}

// ── View ──────────────────────────────────────────────────────────────────

// View disegna lo stato attuale della pagina Catalogo: mostra un indicatore di
// caricamento o l'errore se presente, altrimenti rende in base alla modalità
// attiva (tabella prodotti, form creazione prodotto/categoria o conferma
// eliminazione), chiudendo con un eventuale messaggio di successo.
func (m Model) View() tea.View {
	var sb strings.Builder

	sb.WriteString(style.TitleStile.Render("CONTABILITÀ - CATALOGO PRODOTTI"))
	sb.WriteString("\n\n")

	if m.loading {
		sb.WriteString("Caricamento catalogo in corso...\n")
		return tea.NewView(sb.String())
	}
	if m.err != nil {
		sb.WriteString(style.ErrorStile.Render("Errore: " + m.err.Error()))
		sb.WriteString("\n\n")
		sb.WriteString(style.ErrorStile.Render("[esc] Indietro"))

		return tea.NewView(sb.String())
	}

	switch m.mode {
	case modeCreateProduct:
		sb.WriteString(m.renderCreateProductForm())
	case modeCreateCategory:
		sb.WriteString(m.renderCreateCategoryForm())
	case modeConfirmDeleteProduct:
		sb.WriteString(style.ErrorStile.Render(
			fmt.Sprintf("Eliminare il prodotto \"%s\"? L'operazione è definitiva.", m.confirmDeleteNome)))
		sb.WriteString("\n\n")
		sb.WriteString(style.ErrorStile.Render("[enter] Conferma  |  [esc] Annulla"))
	case modeConfirmDeleteCategory:
		sb.WriteString(style.ErrorStile.Render(
			fmt.Sprintf("Eliminare la categoria \"%s\"? L'operazione è definitiva.", m.confirmDeleteNome)))
		sb.WriteString("\n\n")
		sb.WriteString(style.ErrorStile.Render("[enter] Conferma  |  [esc] Annulla"))
	default:
		sb.WriteString(m.renderProductTable())
	}

	// Messaggio di successo
	if m.success != "" {
		sb.WriteString("\n\n")
		sb.WriteString(style.SuccessStile.Render(m.success))
	}

	return tea.NewView(sb.String())
}

// ── Render tabella prodotti (modeView) ───────────────────────────────────

func (m Model) renderProductTable() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Categoria: %s  %s\n\n",
		style.InfoStile.Render(m.etichettaFiltro()),
		style.SubTitleStile.Render("[tab] cambia"),
	))

	prodotti := m.prodottiFiltrati()
	if len(prodotti) == 0 {
		sb.WriteString("Nessun prodotto trovato con i filtri correnti.\n")
	} else {
		headerStr := fmt.Sprintf("   %-10s %-24s %-16s %-12s %s",
			"CODICE", "PRODOTTO", "CATEGORIA", "PREZZO", "STATO")
		sb.WriteString(style.HeaderStile.Render(headerStr))
		sb.WriteString("\n")
		sb.WriteString(strings.Repeat("─", 80))
		sb.WriteString("\n")

		m.clampScroll()
		inizio := m.offset
		fine := inizio + m.righeVisibili()
		if fine > len(prodotti) {
			fine = len(prodotti)
		}

		for i, p := range prodotti[inizio:fine] {
			assoluto := inizio + i
			stato := "attivo"
			if !p.Attivo {
				stato = "non attivo"
			}

			rigaTesto := fmt.Sprintf("%-10s %-24s %-16s %-12s %s",
				troncaStringa(p.Codice, 10),
				troncaStringa(p.Nome, 24),
				troncaStringa(p.Categoria, 16),
				fmt.Sprintf("€ %.2f", p.Prezzo),
				stato,
			)

			if assoluto == m.cursore {
				sb.WriteString(style.SuccessStile.Render("> "))
				sb.WriteString(rigaTesto)
			} else {
				sb.WriteString("  ")
				sb.WriteString(rigaTesto)
			}
			sb.WriteString("\n")
		}
	}

	// Indicatore di scorrimento (se la lista eccede lo schermo)
	if len(prodotti) > m.righeVisibili() {
		vis := m.righeVisibili()
		ultimo := m.offset + vis
		if ultimo > len(prodotti) {
			ultimo = len(prodotti)
		}
		sb.WriteString(fmt.Sprintf("\nProdotti %d–%d di %d  (pgup/pgdown o rotellina per scorrere)", m.offset+1, ultimo, len(prodotti)))
	}

	sb.WriteString("\n")
	sb.WriteString(style.FooterStile.Render(
		"[↑/↓ o k/j] Naviga  |  [pgup/pgdown o rotellina] Scorri  |  [tab] Cambia categoria  |  [ctrl+n] Nuovo prodotto  |  [ctrl+k] Nuova categoria  |  [ctrl+d] Elimina prodotto  |  [ctrl+x] Elimina categoria  |  [r] Ricarica"))

	return sb.String()
}

// ── Render form creazione prodotto ───────────────────────────────────────

func (m Model) renderCreateProductForm() string {
	categoriaLabel := ""
	if len(m.categorie) > 0 {
		nome := m.categorie[m.createCategoriaIdx]
		categoriaLabel = fmt.Sprintf("< %s >", nome)
	} else {
		categoriaLabel = "(nessuna categoria)"
	}

	if m.createFocus == 7 {
		categoriaLabel = style.FocusStile.Render(categoriaLabel)
	}

	fields := []string{
		"CREA NUOVO PRODOTTO",
		"",
		"Codice:      " + m.createCodice.View(),
		"Barcode:     " + m.createBarcode.View(),
		"Nome:        " + m.createNome.View(),
		"Descrizione: " + m.createDescrizione.View(),
		"Prezzo Acq:  " + m.createPrezzoAcq.View(),
		"Prezzo Vend: " + m.createPrezzoVend.View(),
		"Soglia Min:  " + m.createSogliaMin.View(),
		"Categoria:   " + categoriaLabel,
	}

	return style.LoginBoxStile.Render(strings.Join(fields, "\n"))
}

// ── Render form creazione categoria ──────────────────────────────────────

func (m Model) renderCreateCategoryForm() string {
	fields := []string{
		"CREA NUOVA CATEGORIA",
		"",
		"Nome:        " + m.createCatNome.View(),
		"Descrizione: " + m.createCatDescrizione.View(),
	}

	return style.LoginBoxStile.Render(strings.Join(fields, "\n"))
}

// ── Input handlers ────────────────────────────────────────────────────────

func (m Model) handleViewInput(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.err = nil
		m.success = ""
		return m, nil
	case "up", "k":
		if m.cursore > 0 {
			m.cursore--
			m.mostraCursore()
		}
	case "down", "j":
		if m.cursore < len(m.prodottiFiltrati())-1 {
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
		if m.cursore < len(m.prodottiFiltrati())-1 {
			m.cursore += m.righeVisibili()
			if max := len(m.prodottiFiltrati()) - 1; m.cursore > max {
				m.cursore = max
			}
			m.mostraCursore()
		}
	case "tab":
		m.filtroCategoria++
		if m.filtroCategoria >= len(m.categorie) {
			m.filtroCategoria = nessunFiltro
		}
		m.cursore = 0
		m.offset = 0
	case "r":
		m.loading = true
		m.success = ""
		return m, m.loadCatalogoCmd()

	case "ctrl+n":
		m.mode = modeCreateProduct
		m.success = ""
		m.err = nil
		m.resetCreateProductForm()
		return m, m.focusCreateProductField()

	case "ctrl+k":
		m.mode = modeCreateCategory
		m.success = ""
		m.err = nil
		m.resetCreateCategoryForm()
		return m, m.focusCreateCategoryField()

	case "ctrl+d":
		prodotti := m.prodottiFiltrati()
		if len(prodotti) == 0 {
			return m, nil
		}
		p := prodotti[m.cursore]
		m.mode = modeConfirmDeleteProduct
		m.confirmDeleteID = p.ID
		m.confirmDeleteNome = p.Nome
		m.success = ""
		m.err = nil
		return m, nil

	case "ctrl+x":
		if m.filtroCategoria < 0 || m.filtroCategoria >= len(m.categorie) {
			return m, nil
		}
		nomeCat := m.categorie[m.filtroCategoria]
		idCat, ok := m.categorieID[nomeCat]
		if !ok {
			return m, nil
		}
		m.mode = modeConfirmDeleteCategory
		m.confirmDeleteID = idCat
		m.confirmDeleteNome = nomeCat
		m.success = ""
		m.err = nil
		return m, nil
	}

	return m, nil
}

func (m Model) handleCreateProductInput(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	numFields := campiCreateProdotti

	switch msg.String() {
	case "esc":
		m.mode = modeView
		m.err = nil // Reset dell'errore
		m.resetCreateProductForm()
		return m, nil
	case "tab", "down":
		m.createFocus = (m.createFocus + 1) % numFields
		return m, m.focusCreateProductField()
	case "shift+tab", "up":
		m.createFocus = (m.createFocus - 1 + numFields) % numFields
		return m, m.focusCreateProductField()
	case "left":
		if m.createFocus == 7 && len(m.categorie) > 0 {
			m.createCategoriaIdx = (m.createCategoriaIdx - 1 + len(m.categorie)) % len(m.categorie)
			return m, nil
		}
	case "right":
		if m.createFocus == 7 && len(m.categorie) > 0 {
			m.createCategoriaIdx = (m.createCategoriaIdx + 1) % len(m.categorie)
			return m, nil
		}
	case "enter":
		// Validazione campi obbligatori
		if m.createCodice.Value() == "" || m.createNome.Value() == "" {
			m.err = fmt.Errorf("compilare codice e nome")
			m.success = ""
			return m, nil
		}
		if len(m.categorie) == 0 {
			m.err = fmt.Errorf("nessuna categoria disponibile")
			m.success = ""
			return m, nil
		}

		prezzoAcq, _ := strconv.ParseFloat(m.createPrezzoAcq.Value(), 64)
		prezzoVend, _ := strconv.ParseFloat(m.createPrezzoVend.Value(), 64)
		sogliaMin, _ := strconv.Atoi(m.createSogliaMin.Value())

		m.loading = true
		m.err = nil
		m.success = ""
		return m, m.creaProdottoCmd(api.ClientProductDTO{
			Codice:      m.createCodice.Value(),
			Barcode:     m.createBarcode.Value(),
			Nome:        m.createNome.Value(),
			Descrizione: m.createDescrizione.Value(),
			CategoriaID: m.categorieID[m.categorie[m.createCategoriaIdx]],
			PrezzoAcq:   prezzoAcq,
			PrezzoVend:  prezzoVend,
			SogliaMin:   sogliaMin,
		})
	}

	// Delega al campo focalizzato
	var cmd tea.Cmd
	switch m.createFocus {
	case 0:
		m.createCodice, cmd = m.createCodice.Update(msg)
	case 1:
		m.createBarcode, cmd = m.createBarcode.Update(msg)
	case 2:
		m.createNome, cmd = m.createNome.Update(msg)
	case 3:
		m.createDescrizione, cmd = m.createDescrizione.Update(msg)
	case 4:
		m.createPrezzoAcq, cmd = m.createPrezzoAcq.Update(msg)
	case 5:
		m.createPrezzoVend, cmd = m.createPrezzoVend.Update(msg)
	case 6:
		m.createSogliaMin, cmd = m.createSogliaMin.Update(msg)
	}
	return m, cmd
}

func (m Model) handleCreateCategoryInput(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = modeView
		m.resetCreateCategoryForm()
		return m, nil
	case "tab", "down":
		m.createCatFocus = (m.createCatFocus + 1) % 2
		return m, m.focusCreateCategoryField()
	case "shift+tab", "up":
		m.createCatFocus = (m.createCatFocus - 1 + 2) % 2
		return m, m.focusCreateCategoryField()
	case "enter":
		if m.createCatNome.Value() == "" {
			m.err = fmt.Errorf("il nome della categoria è obbligatorio")
			m.success = ""
			return m, nil
		}
		m.loading = true
		m.err = nil
		m.success = ""
		return m, m.creaCategoriaCmd(m.createCatNome.Value(), m.createCatDescrizione.Value())
	}

	var cmd tea.Cmd
	if m.createCatFocus == 0 {
		m.createCatNome, cmd = m.createCatNome.Update(msg)
	} else {
		m.createCatDescrizione, cmd = m.createCatDescrizione.Update(msg)
	}
	return m, cmd
}

func (m Model) handleConfirmDeleteInput(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		m.loading = true
		m.err = nil
		m.success = ""
		if m.mode == modeConfirmDeleteProduct {
			return m, m.eliminaProdottoCmd(m.confirmDeleteID, m.confirmDeleteNome)
		}
		return m, m.eliminaCategoriaCmd(m.confirmDeleteID, m.confirmDeleteNome)
	case "esc":
		m.mode = modeView
		m.confirmDeleteID = 0
		m.confirmDeleteNome = ""
		return m, nil
	}
	return m, nil
}

// ── Async commands ────────────────────────────────────────────────────────

func (m Model) loadCatalogoCmd() tea.Cmd {
	return func() tea.Msg {
		prodotti, categorie, err := m.Client.GetCatalogo()
		if err != nil {
			return ErroreMsg{Err: err}
		}
		idToNome, _, err := m.Client.GetCategorieIDs()
		if err != nil {
			return ErroreMsg{Err: err}
		}
		nomeToID := make(map[string]int, len(idToNome))
		for id, nome := range idToNome {
			nomeToID[nome] = id
		}
		return CatalogoLoadedMsg{Prodotti: prodotti, Categorie: categorie, CategorieIDs: nomeToID}
	}
}

func (m Model) creaProdottoCmd(input api.ClientProductDTO) tea.Cmd {
	return func() tea.Msg {
		_, err := m.Client.CreaProdotto(input)
		if err != nil {
			return ErroreMsg{Err: err}
		}
		return prodottoCreatoMsg{Nome: input.Nome}
	}
}

func (m Model) eliminaProdottoCmd(id int, nome string) tea.Cmd {
	return func() tea.Msg {
		if err := m.Client.EliminaProdotto(id); err != nil {
			return ErroreMsg{Err: err}
		}
		return prodottoEliminatoMsg{Nome: nome}
	}
}

func (m Model) creaCategoriaCmd(nome, descrizione string) tea.Cmd {
	return func() tea.Msg {
		if err := m.Client.CreaCategoria(nome, descrizione); err != nil {
			return ErroreMsg{Err: err}
		}
		return categoriaCreataMsg{Nome: nome}
	}
}

func (m Model) eliminaCategoriaCmd(id int, nome string) tea.Cmd {
	return func() tea.Msg {
		if err := m.Client.EliminaCategoria(id); err != nil {
			return ErroreMsg{Err: err}
		}
		return categoriaEliminataMsg{Nome: nome}
	}
}

// ── Focus helpers ─────────────────────────────────────────────────────────

func (m *Model) focusCreateProductField() tea.Cmd {
	allInputs := []*textinput.Model{
		&m.createCodice, &m.createBarcode, &m.createNome, &m.createDescrizione,
		&m.createPrezzoAcq, &m.createPrezzoVend, &m.createSogliaMin,
	}
	for i, inp := range allInputs {
		if i == m.createFocus {
			(*inp).Focus()
		} else {
			(*inp).Blur()
		}
	}
	return nil
}

func (m *Model) focusCreateCategoryField() tea.Cmd {
	if m.createCatFocus == 0 {
		m.createCatNome.Focus()
		m.createCatDescrizione.Blur()
	} else {
		m.createCatNome.Blur()
		m.createCatDescrizione.Focus()
	}
	return nil
}

func (m *Model) resetCreateProductForm() {
	m.createCodice.SetValue("")
	m.createCodice.Blur()
	m.createBarcode.SetValue("")
	m.createBarcode.Blur()
	m.createNome.SetValue("")
	m.createNome.Blur()
	m.createDescrizione.SetValue("")
	m.createDescrizione.Blur()
	m.createPrezzoAcq.SetValue("")
	m.createPrezzoAcq.Blur()
	m.createPrezzoVend.SetValue("")
	m.createPrezzoVend.Blur()
	m.createSogliaMin.SetValue("")
	m.createSogliaMin.Blur()
	m.createFocus = 0
	m.createCategoriaIdx = 0
}

func (m *Model) resetCreateCategoryForm() {
	m.createCatNome.SetValue("")
	m.createCatNome.Blur()
	m.createCatDescrizione.SetValue("")
	m.createCatDescrizione.Blur()
	m.createCatFocus = 0
}

// ── Utility ───────────────────────────────────────────────────────────────

func (m Model) etichettaFiltro() string {
	if m.filtroCategoria < 0 || m.filtroCategoria >= len(m.categorie) {
		return "Tutte le categorie"
	}
	return m.categorie[m.filtroCategoria]
}

func (m Model) prodottiFiltrati() []api.Prodotto {
	if m.filtroCategoria < 0 || m.filtroCategoria >= len(m.categorie) {
		return m.prodotti
	}
	categoria := m.categorie[m.filtroCategoria]

	var risultato []api.Prodotto
	for _, p := range m.prodotti {
		if p.Categoria == categoria {
			risultato = append(risultato, p)
		}
	}
	return risultato
}

func troncaStringa(s string, maxLen int) string {
	if len(s) > maxLen {
		if maxLen > 3 {
			return s[:maxLen-3] + "..."
		}
		return s[:maxLen]
	}
	return s
}
