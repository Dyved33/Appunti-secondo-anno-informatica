// Package users fornisce la gestione dell'interfaccia amministrativa per la gestione degli utenti.
// Supporta operazioni di ricerca per username, creazione, modifica, reset password, blocco/sblocco ed eliminazione degli utenti tramite un'interfaccia CLI interattiva.
package users

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
	style "github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/styles"
)

// Init inizializza il modello focalizzando il campo di ricerca per username.
func (m Model) Init() tea.Cmd {
	return m.searchUsername.Focus()
}

// Update elabora gli eventi in input e aggiorna lo stato del modello.
// Gestisce i messaggi di ricerca utente, creazione, modifica, errori e input da tastiera.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case userFoundMsg:
		u := msg.User
		m.user = &u
		m.loading = false
		m.err = nil
		return m, nil

	case userCreatedMsg:
		m.mode = modeSearch
		m.loading = false
		m.success = "Utente creato: " + msg.Username
		m.err = nil
		m.resetCreateForm()
		m.searchUsername.SetValue(msg.Username)
		m.searchUsername.Focus()
		return m, m.searchUserCmd(msg.Username)

	case userUpdatedMsg:
		m.mode = modeSearch
		m.loading = false
		m.success = "Utente aggiornato: " + msg.Username
		m.err = nil
		m.searchUsername.SetValue(msg.Username)
		m.searchUsername.Focus()
		return m, m.searchUserCmd(msg.Username)

	case userBlockedToggledMsg:
		u := msg.User
		m.user = &u
		m.loading = false
		m.err = nil
		if u.Bloccato {
			m.success = "Utente " + u.Username + " bloccato."
		} else {
			m.success = "Utente " + u.Username + " sbloccato con successo."
		}
		return m, nil

	case userPasswordResetMsg:
		m.mode = modeSearch
		m.loading = false
		m.err = nil
		m.success = "Password reimpostata per l'utente " + msg.Username + "."
		m.searchUsername.Focus()
		return m, nil

	case userDeletedMsg:
		m.mode = modeSearch
		m.loading = false
		m.confirmDelete = false
		m.user = nil
		m.err = nil
		m.success = "Utente eliminato: " + msg.Username
		m.searchUsername.SetValue("")
		m.searchUsername.Focus()
		return m, nil

	case usersErrorMsg:
		m.loading = false
		m.err = msg.Err
		m.success = ""
		return m, nil

	case tea.WindowSizeMsg:
		return m, nil

	case tea.KeyPressMsg:
		if m.mode == modeCreate {
			return m.handleCreateInput(msg)
		}
		if m.mode == modeEdit {
			return m.handleEditInput(msg)
		}
		if m.mode == modeResetPassword {
			return m.handleResetPasswordInput(msg)
		}
		return m.handleSearchInput(msg)
	}

	return m, nil
}

// View renderizza l'interfaccia visuale della pagina amministrazione utenti.
// Visualizza la barra di ricerca, i dettagli dell'utente trovato, i form di creazione/modifica e le istruzioni di aiuto.
func (m Model) View() tea.View {
	var sb strings.Builder

	// ── Titolo Sezione ──────────────────────────────────────────────────────
	sb.WriteString(style.TitleStile.Render("AMMINISTRAZIONE UTENTI"))
	sb.WriteString("\n\n")

	// ── Modalità Operative ──────────────────────────────────────────────────
	switch m.mode {
	case modeCreate:
		sb.WriteString(m.renderCreateForm())
		sb.WriteString("\n\n")
		sb.WriteString(style.FooterStile.Render("[tab] Campo successivo  |  [←/→] Ruolo  |  [enter] Crea  |  [esc] Annulla  |  [ctrl+l] Logout"))

	case modeEdit:
		sb.WriteString(m.renderEditForm())
		sb.WriteString("\n\n")
		sb.WriteString(style.FooterStile.Render("[tab] Campo successivo  |  [←/→] Ruolo  |  [enter] Salva  |  [esc] Annulla  |  [ctrl+l] Logout"))

	case modeResetPassword:
		sb.WriteString(m.renderResetForm())
		sb.WriteString("\n\n")
		sb.WriteString(style.FooterStile.Render("[tab] Campo successivo  |  [enter] Conferma  |  [esc] Annulla  |  [ctrl+l] Logout"))

	default:
		sb.WriteString("Cerca utente: ")
		sb.WriteString(m.searchUsername.View())
		sb.WriteString("\n\n")
		sb.WriteString(m.renderUser())
		sb.WriteString("\n\n")
		if m.user != nil && m.confirmDelete {
			sb.WriteString(style.ErrorStile.Render(fmt.Sprintf("Eliminazione di \"%s\": l'operazione è definitiva.", m.user.Username)))
			sb.WriteString("\n\n")
			sb.WriteString(style.ErrorStile.Render("[enter] Conferma eliminazione  |  [esc] Annulla"))
		} else if m.user != nil {
			btnBlock := "[ctrl+b] Blocca utente"
			if m.user.Bloccato {
				btnBlock = "[ctrl+b] Sblocca utente"
			}
			sb.WriteString(style.FooterStile.Render("[enter] Cerca  |  [ctrl+n] Crea  |  [ctrl+e] Modifica  |  " + btnBlock + "  |  [ctrl+d] Elimina utente  |  [ctrl+p] Resetta password  |  [esc] Pulisci  |  [ctrl+l] Logout"))
		} else {
			sb.WriteString(style.FooterStile.Render("[enter] Cerca  |  [ctrl+n] Crea utente  |  [esc] Pulisci  |  [ctrl+l] Logout"))
		}
	}

	// ── Messaggio di Successo ───────────────────────────────────────────────
	if m.success != "" {
		sb.WriteString("\n\n")
		sb.WriteString(style.SuccessStile.Render(m.success))
	}

	// ── Costruzione della Vista e Cursore ───────────────────────────────────
	view := tea.NewView(style.FinestraInternaStile.Render(sb.String()))
	switch m.mode {
	case modeCreate:
		if cursor := m.currentCreateCursor(); cursor != nil {
			view.Cursor = cursor
		}
	case modeEdit:
		if cursor := m.currentEditCursor(); cursor != nil {
			view.Cursor = cursor
		}
	case modeResetPassword:
		if cursor := m.currentResetCursor(); cursor != nil {
			view.Cursor = cursor
		}
	case modeSearch:
		if cursor := m.searchUsername.Cursor(); cursor != nil {
			view.Cursor = cursor
		}
	}

	return view
}

// renderUser genera la rappresentazione dei dati dell'utente trovato.
// Mostra ID, username, ruolo, stato (attivo/bloccato) e numero di tentativi falliti.
func (m Model) renderUser() string {

	if m.err != nil {
		return style.ErrorStile.Render("Errore: Nessun utente trovato.")
	}
	if m.user == nil {
		return style.InfoStile.Render("Inserisci uno username e premi [enter] per cercare.")
	}

	stato := "Attivo"
	if m.user.Bloccato {
		stato = "Bloccato"
	}

	var sb strings.Builder
	header := fmt.Sprintf("   %-4s %-20s %-16s %-12s %-8s", "ID", "USERNAME", "RUOLO", "STATO", "TENT.")
	sb.WriteString(style.HeaderStile.Render(header))
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("─", 70))
	sb.WriteString("\n")

	riga := fmt.Sprintf("   %-4d %-20s %-16s %-12s %-8d", m.user.ID, tronca(m.user.Username, 20), m.user.Ruolo, stato, m.user.TentativiFalliti)
	sb.WriteString(style.SuccessStile.Render("> "))
	sb.WriteString(riga)
	sb.WriteString("\n")

	return sb.String()
}

// renderCreateForm genera il form per la creazione di un nuovo utente.
// Include campi per username, password e ruolo.
func (m Model) renderCreateForm() string {
	ruoloStr := ruoliDisponibili[m.createRuoloIdx]
	if m.createFocus == 2 {
		ruoloStr = style.FocusStile.Render("< " + ruoloStr + " >")
	}

	return style.LoginBoxStile.Render(strings.Join([]string{
		"CREA NUOVO UTENTE",
		"",
		"Username: " + m.createUsername.View(),
		"Password: " + m.createPassword.View(),
		"Ruolo:    " + ruoloStr,
	}, "\n"))
}

// renderEditForm genera il form per la modifica dell'utente trovato.
// Visualizza ID, username e ruolo dell'utente.
func (m Model) renderEditForm() string {
	if m.user == nil {
		return style.LoginBoxStile.Render("Nessun utente da modificare")
	}

	ruoloStr := ruoliDisponibili[m.editRuoloIdx]
	if m.editFocus == 1 {
		ruoloStr = style.FocusStile.Render("< " + ruoloStr + " >")
	}

	return style.LoginBoxStile.Render(strings.Join([]string{
		"MODIFICA UTENTE",
		"",
		fmt.Sprintf("ID:       %d", m.user.ID),
		"Username: " + m.editUsername.View(),
		"Ruolo:    " + ruoloStr,
	}, "\n"))
}

// renderResetForm genera il form per il reset della password dell'utente trovato.
// Include campi per la nuova password e la relativa conferma.
func (m Model) renderResetForm() string {
	if m.user == nil {
		return style.LoginBoxStile.Render("Nessun utente selezionato")
	}

	return style.LoginBoxStile.Render(strings.Join([]string{
		"RESET PASSWORD",
		"",
		fmt.Sprintf("Utente:   %s", m.user.Username),
		"Password: " + m.resetPassword.View(),
		"Conferma: " + m.resetConferma.View(),
	}, "\n"))
}

// handleSearchInput gestisce gli input quando in modalità ricerca per username.
func (m Model) handleSearchInput(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	if m.confirmDelete {
		switch msg.String() {
		case "enter":
			username := m.user.Username
			m.loading = true
			m.err = nil
			m.success = ""
			m.confirmDelete = false
			return m, m.deleteUserCmd(username)
		case "esc":
			m.confirmDelete = false
			m.err = nil
			return m, nil
		default:
			// In attesa di conferma vengono ignorati gli altri tasti
			return m, nil
		}
	}

	switch msg.String() {
	case "ctrl+n":
		m.mode = modeCreate
		m.success = ""
		m.err = nil
		m.resetCreateForm()
		return m, m.focusCreateField()

	case "ctrl+e":
		if m.user == nil {
			return m, nil
		}
		m.mode = modeEdit
		m.success = ""
		m.err = nil
		m.prepareEditForm()
		return m, m.focusEditField()

	case "ctrl+b":
		if m.user == nil {
			return m, nil
		}
		m.loading = true
		m.err = nil
		m.success = ""
		newBlocked := !m.user.Bloccato
		return m, m.toggleBlockUserCmd(m.user.Username, newBlocked)

	case "ctrl+p":
		if m.user == nil {
			return m, nil
		}
		m.mode = modeResetPassword
		m.success = ""
		m.err = nil
		m.prepareResetForm()
		return m, m.focusResetField()

	case "ctrl+d":
		if m.user == nil {
			return m, nil
		}
		m.confirmDelete = true
		m.err = nil
		m.success = ""
		return m, nil

	case "ctrl+r":
		if m.user != nil {
			m.loading = true
			m.err = nil
			m.success = ""
			return m, m.searchUserCmd(m.user.Username)
		}

	case "esc":
		m.searchUsername.SetValue("")
		m.user = nil
		m.err = nil
		m.success = ""
		return m, nil

	case "enter":
		val := strings.TrimSpace(m.searchUsername.Value())
		if val == "" {
			return m, nil
		}
		m.loading = true
		m.err = nil
		m.success = ""
		return m, m.searchUserCmd(val)
	}

	var cmd tea.Cmd
	m.searchUsername, cmd = m.searchUsername.Update(msg)
	return m, cmd
}

// handleResetPasswordInput gestisce gli input da tastiera quando in modalità reset della password.
// Supporta navigazione tra i campi (tab), validazione di coerenza delle password e invio del form.
func (m Model) handleResetPasswordInput(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = modeSearch
		m.resetResetForm()
		m.searchUsername.Focus()
		return m, nil
	case "tab", "down":
		m.resetFocus = (m.resetFocus + 1) % 2
		return m, m.focusResetField()
	case "shift+tab", "up":
		m.resetFocus = (m.resetFocus - 1 + 2) % 2
		return m, m.focusResetField()
	case "enter":
		if m.user == nil {
			return m, nil
		}
		if m.resetPassword.Value() == "" || m.resetConferma.Value() == "" {
			m.err = fmt.Errorf("compilare tutti i campi obbligatori")
			m.success = ""
			return m, nil
		}
		if m.resetPassword.Value() != m.resetConferma.Value() {
			m.err = fmt.Errorf("le password non coincidono")
			m.success = ""
			return m, nil
		}
		m.loading = true
		m.err = nil
		return m, m.resetPasswordCmd(m.user.Username, m.resetPassword.Value())
	}

	var cmd tea.Cmd
	switch m.resetFocus {
	case 0:
		m.resetPassword, cmd = m.resetPassword.Update(msg)
	case 1:
		m.resetConferma, cmd = m.resetConferma.Update(msg)
	}
	return m, cmd
}

// handleCreateInput gestisce gli input da tastiera quando in modalità creazione di un nuovo utente.
// Supporta navigazione tra i campi (tab, frecce su/giù), selezione del ruolo (frecce sx/dx), validazione e invio del form.
func (m Model) handleCreateInput(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = modeSearch
		m.resetCreateForm()
		m.searchUsername.Focus()
		return m, nil
	case "tab", "down":
		m.createFocus = (m.createFocus + 1) % 3
		return m, m.focusCreateField()
	case "shift+tab", "up":
		m.createFocus = (m.createFocus - 1 + 3) % 3
		return m, m.focusCreateField()
	case "left":
		if m.createFocus == 2 {
			m.createRuoloIdx = (m.createRuoloIdx - 1 + len(ruoliDisponibili)) % len(ruoliDisponibili)
			return m, nil
		}
	case "right":
		if m.createFocus == 2 {
			m.createRuoloIdx = (m.createRuoloIdx + 1) % len(ruoliDisponibili)
			return m, nil
		}
	case "enter":
		if m.createUsername.Value() == "" || m.createPassword.Value() == "" {
			m.err = fmt.Errorf("compilare tutti i campi obbligatori")
			m.success = ""
			return m, nil
		}
		m.loading = true
		m.err = nil
		return m, m.createUserCmd(api.UserCreateInput{
			Username: m.createUsername.Value(),
			Password: m.createPassword.Value(),
			Ruolo:    ruoliDisponibili[m.createRuoloIdx],
		})
	}

	var cmd tea.Cmd
	switch m.createFocus {
	case 0:
		m.createUsername, cmd = m.createUsername.Update(msg)
	case 1:
		m.createPassword, cmd = m.createPassword.Update(msg)
	}
	return m, cmd
}

// handleEditInput gestisce gli input da tastiera quando in modalità modifica di un utente.
// Supporta navigazione tra i campi (tab, frecce su/giù), selezione del ruolo (frecce sx/dx), validazione e invio del form.
func (m Model) handleEditInput(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = modeSearch
		m.searchUsername.Focus()
		return m, nil
	case "tab", "down":
		m.editFocus = (m.editFocus + 1) % 2
		return m, m.focusEditField()
	case "shift+tab", "up":
		m.editFocus = (m.editFocus - 1 + 2) % 2
		return m, m.focusEditField()
	case "left":
		if m.editFocus == 1 {
			m.editRuoloIdx = (m.editRuoloIdx - 1 + len(ruoliDisponibili)) % len(ruoliDisponibili)
			return m, nil
		}
	case "right":
		if m.editFocus == 1 {
			m.editRuoloIdx = (m.editRuoloIdx + 1) % len(ruoliDisponibili)
			return m, nil
		}
	case "enter":
		if m.user == nil {
			return m, nil
		}
		if m.editUsername.Value() == "" {
			m.err = fmt.Errorf("username non valido")
			m.success = ""
			return m, nil
		}
		m.loading = true
		m.err = nil
		return m, m.updateUserCmd(m.user.Username, api.UserUpdateInput{
			Username: m.editUsername.Value(),
			Ruolo:    ruoliDisponibili[m.editRuoloIdx],
		})
	}

	var cmd tea.Cmd
	if m.editFocus == 0 {
		m.editUsername, cmd = m.editUsername.Update(msg)
	}
	return m, cmd
}

// searchUserCmd crea un comando asincrono per la ricerca di un utente per username.
func (m Model) searchUserCmd(username string) tea.Cmd {
	return func() tea.Msg {
		user, err := m.Client.GetUsersByUsername(username)
		if err != nil {
			return usersErrorMsg{Err: err}
		}
		return userFoundMsg{User: user}
	}
}

// createUserCmd crea un comando asincrono per la creazione di un nuovo utente.
func (m Model) createUserCmd(input api.UserCreateInput) tea.Cmd {
	return func() tea.Msg {
		created, err := m.Client.CreateUser(input)
		if err != nil {
			return usersErrorMsg{Err: err}
		}
		return userCreatedMsg{Username: created.Username}
	}
}

// updateUserCmd crea un comando asincrono per l'aggiornamento di un utente esistente.
func (m Model) updateUserCmd(username string, input api.UserUpdateInput) tea.Cmd {
	return func() tea.Msg {
		updated, err := m.Client.UpdateUser(username, input)
		if err != nil {
			return usersErrorMsg{Err: err}
		}
		return userUpdatedMsg{Username: updated.Username}
	}
}

// resetPasswordCmd crea un comando asincrono per reimpostare la password di un utente.
func (m Model) resetPasswordCmd(username, newPassword string) tea.Cmd {
	return func() tea.Msg {
		_, err := m.Client.UpdateUser(username, api.UserUpdateInput{Password: newPassword})
		if err != nil {
			return usersErrorMsg{Err: err}
		}
		return userPasswordResetMsg{Username: username}
	}
}

// toggleBlockUserCmd crea un comando asincrono per bloccare o sbloccare un utente.
func (m Model) toggleBlockUserCmd(username string, newBlocked bool) tea.Cmd {
	return func() tea.Msg {
		updated, err := m.Client.ToggleUserBlock(username, newBlocked)
		if err != nil {
			return usersErrorMsg{Err: err}
		}
		return userBlockedToggledMsg{User: updated}
	}
}

// deleteUserCmd crea un comando asincrono per eliminare definitivamente un utente.
func (m Model) deleteUserCmd(username string) tea.Cmd {
	return func() tea.Msg {
		if err := m.Client.DeleteUser(username); err != nil {
			return usersErrorMsg{Err: err}
		}
		return userDeletedMsg{Username: username}
	}
}

// resetCreateForm ripristina il form di creazione agli stati predefiniti e riporta il focus al primo campo.
func (m *Model) resetCreateForm() {
	m.createUsername.SetValue("")
	m.createPassword.SetValue("")
	m.createRuoloIdx = 0
	m.createFocus = 0
	m.createUsername.Focus()
	m.createPassword.Blur()
}

// prepareEditForm inizializza il form di modifica con i dati dell'utente attualmente visualizzato.
func (m *Model) prepareEditForm() {
	if m.user == nil {
		return
	}
	m.editUsername.SetValue(m.user.Username)
	m.editRuoloIdx = indexRuolo(m.user.Ruolo)
	m.editFocus = 0
	m.editUsername.Focus()
}

// prepareResetForm inizializza il form di reset password agli stati predefiniti.
func (m *Model) prepareResetForm() {
	m.resetResetForm()
	m.resetFocus = 0
}

// resetResetForm ripristina i campi del form di reset password e riporta il focus al primo campo.
func (m *Model) resetResetForm() {
	m.resetPassword.SetValue("")
	m.resetConferma.SetValue("")
	m.resetFocus = 0
}

// focusCreateField sposta il focus al campo specificato nel form di creazione.
func (m *Model) focusCreateField() tea.Cmd {
	m.createUsername.Blur()
	m.createPassword.Blur()

	switch m.createFocus {
	case 0:
		return m.createUsername.Focus()
	case 1:
		return m.createPassword.Focus()
	default:
		return nil
	}
}

// focusEditField sposta il focus al campo username nel form di modifica.
func (m *Model) focusEditField() tea.Cmd {
	m.editUsername.Blur()
	if m.editFocus == 0 {
		return m.editUsername.Focus()
	}
	return nil
}

// focusResetField sposta il focus al campo specificato nel form di reset password.
func (m *Model) focusResetField() tea.Cmd {
	m.resetPassword.Blur()
	m.resetConferma.Blur()

	switch m.resetFocus {
	case 0:
		return m.resetPassword.Focus()
	case 1:
		return m.resetConferma.Focus()
	default:
		return nil
	}
}

// currentCreateCursor restituisce il cursore del campo attualmente focalizzato nel form di creazione.
func (m *Model) currentCreateCursor() *tea.Cursor {
	switch m.createFocus {
	case 0:
		return m.createUsername.Cursor()
	case 1:
		return m.createPassword.Cursor()
	default:
		return nil
	}
}

// currentEditCursor restituisce il cursore del campo attualmente focalizzato nel form di modifica.
func (m *Model) currentEditCursor() *tea.Cursor {
	if m.editFocus == 0 {
		return m.editUsername.Cursor()
	}
	return nil
}

// currentResetCursor restituisce il cursore del campo attualmente focalizzato nel form di reset password.
func (m *Model) currentResetCursor() *tea.Cursor {
	switch m.resetFocus {
	case 0:
		return m.resetPassword.Cursor()
	case 1:
		return m.resetConferma.Cursor()
	default:
		return nil
	}
}

// indexRuolo restituisce l'indice della posizione del ruolo nella lista dei ruoli disponibili.
func indexRuolo(ruolo string) int {
	for i, r := range ruoliDisponibili {
		if r == ruolo {
			return i
		}
	}
	return 0
}

// tronca limita la lunghezza di una stringa, aggiungendo "..." se necessario.
func tronca(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
