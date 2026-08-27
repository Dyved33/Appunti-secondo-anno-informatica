// Package users fornisce test per la gestione dell'interfaccia amministrativa degli utenti.
package users

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
)

func unwrapModel(t *testing.T, model tea.Model) Model {
	t.Helper()
	switch m := model.(type) {
	case Model:
		return m
	case *Model:
		return *m
	default:
		t.Fatalf("tipo modello inatteso: %T", model)
		return Model{}
	}
}

// TestUsersSearchView testa la visualizzazione iniziale della ricerca utente.
func TestUsersSearchView(t *testing.T) {
	m := New(&api.APIClient{})

	view := m.View().Content
	if !strings.Contains(view, "AMMINISTRAZIONE UTENTI") {
		t.Fatal("la vista deve contenere il titolo della pagina")
	}
	if !strings.Contains(view, "Cerca utente") {
		t.Fatal("la vista deve contenere la barra di ricerca")
	}
	if !strings.Contains(view, "Inserisci uno username") {
		t.Fatal("la vista deve mostrare il messaggio di default quando nessun utente è cercato")
	}
}

// TestUsersFoundAndDisplay testa la ricezione del messaggio userFoundMsg e visualizzazione utente.
func TestUsersFoundAndDisplay(t *testing.T) {
	m := New(&api.APIClient{})

	res, _ := m.Update(userFoundMsg{User: api.User{
		ID:               1,
		Username:         "mario",
		Ruolo:            "magazziniere",
		TentativiFalliti: 0,
		Bloccato:         false,
	}})
	model := unwrapModel(t, res)

	if model.loading {
		t.Fatal("dopo userFoundMsg loading dovrebbe essere false")
	}
	if model.user == nil {
		t.Fatal("user non dovrebbe essere nil")
	}
	if model.user.Username != "mario" {
		t.Fatalf("username atteso mario, ottenuto %s", model.user.Username)
	}

	view := model.View().Content
	if !strings.Contains(view, "mario") {
		t.Fatal("la vista deve contenere i dettagli dell'utente trovato")
	}
	if !strings.Contains(view, "magazziniere") {
		t.Fatal("la vista deve contenere il ruolo dell'utente")
	}
}

// TestUsersSearchInput testa l'invio del comando di ricerca tramite tasto enter.
func TestUsersSearchInput(t *testing.T) {
	m := New(&api.APIClient{})
	m.searchUsername.SetValue("mario")

	res, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	model := unwrapModel(t, res)

	if cmd == nil {
		t.Fatal("premendo enter con testo dovrebbe essere restituito un Cmd di ricerca")
	}
	if !model.loading {
		t.Fatal("durante la ricerca loading dovrebbe essere true")
	}
}

// TestUsersCreateMode testa l'ingresso e l'uscita dalla modalità creazione di un nuovo utente.
func TestUsersCreateMode(t *testing.T) {
	m := New(&api.APIClient{})

	res, _ := m.Update(tea.KeyPressMsg{Code: 'n', Mod: tea.ModCtrl})
	model := unwrapModel(t, res)
	if model.mode != modeCreate {
		t.Fatal("premendo ctrl+n il modello deve entrare in modalità creazione")
	}

	view := model.View().Content
	if !strings.Contains(view, "CREA NUOVO UTENTE") {
		t.Fatal("la vista deve mostrare il form di creazione")
	}

	// Esci premendo esc
	res, _ = model.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	model = unwrapModel(t, res)
	if model.mode != modeSearch {
		t.Fatal("premendo esc il modello deve tornare alla modalità ricerca")
	}
}

// TestUsersEditMode testa l'ingresso e la visualizzazione del form di modifica.
func TestUsersEditMode(t *testing.T) {
	m := New(&api.APIClient{})

	// Carica un utente
	res, _ := m.Update(userFoundMsg{User: api.User{
		ID:       1,
		Username: "mario",
		Ruolo:    "magazziniere",
	}})
	model := unwrapModel(t, res)

	// Entra in modifica
	res, _ = model.Update(tea.KeyPressMsg{Code: 'e', Mod: tea.ModCtrl})
	model = unwrapModel(t, res)
	if model.mode != modeEdit {
		t.Fatal("premendo ctrl+e il modello deve entrare in modalità modifica")
	}

	view := model.View().Content
	if !strings.Contains(view, "MODIFICA UTENTE") {
		t.Fatal("la vista deve mostrare il form di modifica")
	}
	if !strings.Contains(view, "ID:       1") {
		t.Fatal("il form di modifica deve contenere l'ID dell'utente")
	}

	// Annulla con esc
	res, _ = model.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	model = unwrapModel(t, res)
	if model.mode != modeSearch {
		t.Fatal("premendo esc il modello deve tornare alla modalità ricerca")
	}
}

// TestUsersErrorHandling testa la gestione degli errori durante la ricerca o le operazioni.
func TestUsersErrorHandling(t *testing.T) {
	m := New(&api.APIClient{})

	testErr := fmt.Errorf("utente non trovato")
	res, _ := m.Update(usersErrorMsg{Err: testErr})
	model := unwrapModel(t, res)

	if model.err == nil {
		t.Fatal("l'errore non è stato salvato nel modello")
	}
	if model.err.Error() != testErr.Error() {
		t.Fatalf("errore atteso '%v', trovato '%v'", testErr, model.err)
	}
	if model.loading {
		t.Fatal("loading non dovrebbe essere true dopo un errore")
	}

	view := model.View().Content
	if !strings.Contains(view, "Errore") {
		t.Fatal("la vista deve mostrare il messaggio di errore")
	}
}

// TestUsersFormNavigation testa la navigazione tra i campi del form di creazione con tab, down, up, shift+tab.
func TestUsersFormNavigation(t *testing.T) {
	m := New(&api.APIClient{})

	res, _ := m.Update(tea.KeyPressMsg{Code: 'n', Mod: tea.ModCtrl})
	model := unwrapModel(t, res)

	if model.createFocus != 0 {
		t.Fatalf("focus iniziale atteso a 0, trovato %d", model.createFocus)
	}

	// Down passa a 1
	res, _ = model.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	model = unwrapModel(t, res)
	if model.createFocus != 1 {
		t.Fatalf("focus atteso a 1 dopo down, trovato %d", model.createFocus)
	}

	// Tab passa a 2 (Ruolo)
	res, _ = model.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	model = unwrapModel(t, res)
	if model.createFocus != 2 {
		t.Fatalf("focus atteso a 2 dopo tab, trovato %d", model.createFocus)
	}

	// Up torna a 1
	res, _ = model.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	model = unwrapModel(t, res)
	if model.createFocus != 1 {
		t.Fatalf("focus atteso a 1 dopo up, trovato %d", model.createFocus)
	}

	// Shift+Tab torna a 0
	res, _ = model.Update(tea.KeyPressMsg{Code: tea.KeyTab, Mod: tea.ModShift})
	model = unwrapModel(t, res)
	if model.createFocus != 0 {
		t.Fatalf("focus atteso a 0 dopo shift+tab, trovato %d", model.createFocus)
	}
}

// TestUsersRuoloSelection testa la selezione del ruolo nel form di creazione quando focalizzato sul ruolo.
func TestUsersRuoloSelection(t *testing.T) {
	m := New(&api.APIClient{})

	res, _ := m.Update(tea.KeyPressMsg{Code: 'n', Mod: tea.ModCtrl})
	model := unwrapModel(t, res)

	// Spostati sul campo Ruolo (createFocus = 2)
	res, _ = model.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	model = unwrapModel(t, res)
	res, _ = model.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	model = unwrapModel(t, res)
	if model.createFocus != 2 {
		t.Fatalf("focus atteso su ruolo (2), ottenuto %d", model.createFocus)
	}

	initialRuolo := ruoliDisponibili[model.createRuoloIdx]

	res, _ = model.Update(tea.KeyPressMsg{Code: tea.KeyRight})
	model = unwrapModel(t, res)

	nextRuolo := ruoliDisponibili[model.createRuoloIdx]
	if nextRuolo == initialRuolo {
		t.Fatal("il ruolo dovrebbe cambiare navigando a destra")
	}

	res, _ = model.Update(tea.KeyPressMsg{Code: tea.KeyLeft})
	model = unwrapModel(t, res)

	if ruoliDisponibili[model.createRuoloIdx] != initialRuolo {
		t.Fatal("il ruolo dovrebbe tornare al valore iniziale navigando a sinistra")
	}
}

// TestUsersCreateFormValidation testa la validazione del form di creazione.
func TestUsersCreateFormValidation(t *testing.T) {
	m := New(&api.APIClient{})

	res, _ := m.Update(tea.KeyPressMsg{Text: "ctrl+n"})
	model := unwrapModel(t, res)

	res, _ = model.Update(tea.KeyPressMsg{Text: "enter"})
	model = unwrapModel(t, res)

	if model.err == nil {
		t.Fatal("dovrebbe esserci un errore di validazione per campi vuoti")
	}
	if !strings.Contains(model.err.Error(), "obbligatori") {
		t.Fatal("il messaggio di errore dovrebbe menzionare i campi obbligatori")
	}
}

// TestUsersUserCreatedMessage testa il messaggio di successo dopo la creazione di un utente.
func TestUsersUserCreatedMessage(t *testing.T) {
	m := New(&api.APIClient{})

	res, cmd := m.Update(userCreatedMsg{Username: "nuovoutente"})
	model := unwrapModel(t, res)

	if model.mode != modeSearch {
		t.Fatal("dopo la creazione dovrebbe tornare alla modalità ricerca")
	}
	if !strings.Contains(model.success, "nuovoutente") {
		t.Fatal("il messaggio di successo dovrebbe contenere l'username creato")
	}
	if cmd == nil {
		t.Fatal("dovrebbe lanciare il comando di ricerca per il nuovo utente")
	}
}

// TestUsersUserUpdatedMessage testa il messaggio di successo dopo l'aggiornamento di un utente.
func TestUsersUserUpdatedMessage(t *testing.T) {
	m := New(&api.APIClient{})

	res, cmd := m.Update(userUpdatedMsg{Username: "utenteaggiornato"})
	model := unwrapModel(t, res)

	if model.mode != modeSearch {
		t.Fatal("dopo la modifica dovrebbe tornare alla modalità ricerca")
	}
	if !strings.Contains(model.success, "utenteaggiornato") {
		t.Fatal("il messaggio di successo dovrebbe contenere l'username aggiornato")
	}
	if cmd == nil {
		t.Fatal("dovrebbe lanciare il comando di ricerca per l'utente aggiornato")
	}
}

// TestUsersClearSearch testa la pulizia della ricerca con esc.
func TestUsersClearSearch(t *testing.T) {
	m := New(&api.APIClient{})
	m.searchUsername.SetValue("mario")
	m.user = &api.User{ID: 1, Username: "mario"}

	res, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	model := unwrapModel(t, res)

	if model.searchUsername.Value() != "" {
		t.Fatal("il campo di ricerca dovrebbe essere vuoto")
	}
	if model.user != nil {
		t.Fatal("l'utente selezionato dovrebbe essere nil")
	}
}

// TestUsersEditMode_FullFlow testa la modifica completa di un utente con ctrl+e, navigazione e submit.
func TestUsersEditMode_FullFlow(t *testing.T) {
	m := New(&api.APIClient{})
	res, _ := m.Update(userFoundMsg{User: api.User{
		ID:       1,
		Username: "mario",
		Ruolo:    "magazziniere",
	}})
	model := unwrapModel(t, res)

	// Entra in modifica con ctrl+e
	res, _ = model.Update(tea.KeyPressMsg{Code: 'e', Mod: tea.ModCtrl})
	model = unwrapModel(t, res)
	if model.mode != modeEdit {
		t.Fatal("attesa modalità modeEdit dopo ctrl+e")
	}
	if model.editUsername.Value() != "mario" {
		t.Fatalf("atteso username mario, trovato %s", model.editUsername.Value())
	}
	if model.editFocus != 0 {
		t.Fatalf("atteso focus su username (0), trovato %d", model.editFocus)
	}

	// Spostati sul Ruolo con down
	res, _ = model.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	model = unwrapModel(t, res)
	if model.editFocus != 1 {
		t.Fatalf("atteso focus su ruolo (1), trovato %d", model.editFocus)
	}

	// Cambia ruolo con right
	res, _ = model.Update(tea.KeyPressMsg{Code: tea.KeyRight})
	model = unwrapModel(t, res)
	if ruoliDisponibili[model.editRuoloIdx] == "magazziniere" {
		t.Fatal("il ruolo doveva cambiare dopo right")
	}

	// Invia il form con enter
	res, cmd := model.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	model = unwrapModel(t, res)
	if cmd == nil {
		t.Fatal("atteso comando tea.Cmd di aggiornamento utente")
	}
	if !model.loading {
		t.Fatal("atteso loading true durante l'aggiornamento")
	}
}

// TestUsersToggleBlock testa il blocco e sblocco utente con ctrl+b.
func TestUsersToggleBlock(t *testing.T) {
	m := New(&api.APIClient{})
	res, _ := m.Update(userFoundMsg{User: api.User{
		ID:               1,
		Username:         "mario",
		Ruolo:            "magazziniere",
		Bloccato:         false,
		TentativiFalliti: 0,
	}})
	model := unwrapModel(t, res)

	// Premi ctrl+b per bloccare
	res, cmd := model.Update(tea.KeyPressMsg{Code: 'b', Mod: tea.ModCtrl})
	model = unwrapModel(t, res)
	if cmd == nil {
		t.Fatal("atteso comando tea.Cmd dopo ctrl+b")
	}
	if !model.loading {
		t.Fatal("atteso loading true")
	}

	// Ricevi messaggio di utente bloccato
	res, _ = model.Update(userBlockedToggledMsg{User: api.User{
		ID:               1,
		Username:         "mario",
		Ruolo:            "magazziniere",
		Bloccato:         true,
		TentativiFalliti: 3,
	}})
	model = unwrapModel(t, res)
	if !model.user.Bloccato {
		t.Fatal("l'utente dovrebbe risultare bloccato")
	}
	if !strings.Contains(model.success, "bloccato") {
		t.Fatalf("messaggio di successo inatteso: %s", model.success)
	}

	// Premi di nuovo ctrl+b per sbloccare
	res, _ = model.Update(userBlockedToggledMsg{User: api.User{
		ID:               1,
		Username:         "mario",
		Ruolo:            "magazziniere",
		Bloccato:         false,
		TentativiFalliti: 0,
	}})
	model = unwrapModel(t, res)
	if model.user.Bloccato {
		t.Fatal("l'utente dovrebbe risultare sbloccato")
	}
	if !strings.Contains(model.success, "sbloccato") {
		t.Fatalf("messaggio di successo inatteso: %s", model.success)
	}
}

// TestUsersDeleteFlow testa il flusso completo di eliminazione di un utente:
// richiesta conferma (ctrl+d), annullamento (esc), conferma (enter) e pulizia dello stato.
func TestUsersDeleteFlow(t *testing.T) {
	var deletedPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		deletedPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	m := New(api.New(srv.URL))
	m.user = &api.User{ID: 3, Username: "mario", Ruolo: "magazziniere"}

	// ctrl+d richiede la conferma
	res, _ := m.Update(tea.KeyPressMsg{Code: 'd', Mod: tea.ModCtrl})
	model := unwrapModel(t, res)
	if !model.confirmDelete {
		t.Fatal("attesa richiesta di conferma eliminazione dopo ctrl+d")
	}
	view := model.View().Content
	if !strings.Contains(view, "Conferma eliminazione") {
		t.Fatal("la vista deve chiedere conferma dopo ctrl+d")
	}
	if !strings.Contains(view, "mario") {
		t.Fatal("la vista deve mostrare lo username da eliminare")
	}

	// esc annulla la richiesta mantenendo l'utente caricato
	res, _ = model.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	model = unwrapModel(t, res)
	if model.confirmDelete {
		t.Fatal("esc deve annullare la conferma")
	}
	if model.user == nil || model.user.Username != "mario" {
		t.Fatal("l'utente deve restare caricato dopo l'annullamento")
	}

	// enter conferma ed elimina via API
	res, _ = model.Update(tea.KeyPressMsg{Code: 'd', Mod: tea.ModCtrl})
	model = unwrapModel(t, res)
	res, cmd := model.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	model = unwrapModel(t, res)
	if cmd == nil {
		t.Fatal("atteso comando di eliminazione dopo la conferma")
	}
	if model.confirmDelete {
		t.Fatal("la conferma deve risultare disattivata dopo l'invio")
	}

	msg := cmd()
	deleted, ok := msg.(userDeletedMsg)
	if !ok {
		t.Fatalf("atteso userDeletedMsg, ottenuto %T", msg)
	}
	if deleted.Username != "mario" {
		t.Fatalf("username inatteso: %s", deleted.Username)
	}
	if deletedPath != "/api/v1/users/mario" {
		t.Fatalf("chiamata DELETE inattesa: %s", deletedPath)
	}

	// Il messaggio di eliminazione ripulisce lo stato
	res, _ = model.Update(deleted)
	model = unwrapModel(t, res)
	if model.user != nil {
		t.Fatal("l'utente deve essere scaricato dopo l'eliminazione")
	}
	if !strings.Contains(model.success, "eliminato") {
		t.Fatalf("messaggio di successo inatteso: %s", model.success)
	}
	if model.searchUsername.Value() != "" {
		t.Fatal("il campo di ricerca deve essere svuotato dopo l'eliminazione")
	}
}

// TestUsersDeleteErrorKeepsUser verifica che, se l'eliminazione fallisce,
// l'utente resti caricato e l'errore venga mostrato.
func TestUsersDeleteErrorKeepsUser(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "utente non trovato", http.StatusNotFound)
	}))
	defer srv.Close()

	m := New(api.New(srv.URL))
	m.user = &api.User{ID: 3, Username: "ghost"}

	res, _ := m.Update(tea.KeyPressMsg{Code: 'd', Mod: tea.ModCtrl})
	model := unwrapModel(t, res)
	res, cmd := model.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	model = unwrapModel(t, res)

	msg := cmd()
	errMsg, ok := msg.(usersErrorMsg)
	if !ok {
		t.Fatalf("atteso usersErrorMsg, ottenuto %T", msg)
	}
	res, _ = model.Update(errMsg)
	model = unwrapModel(t, res)

	if model.user == nil || model.user.Username != "ghost" {
		t.Fatal("l'utente deve restare caricato se l'eliminazione fallisce")
	}
	if !strings.Contains(model.err.Error(), "utente non trovato") {
		t.Fatalf("errore inatteso: %v", model.err)
	}
}
