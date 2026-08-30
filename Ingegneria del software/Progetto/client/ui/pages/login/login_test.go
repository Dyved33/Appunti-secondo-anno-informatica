package login

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"charm.land/bubbles/v2/cursor"
	tea "charm.land/bubbletea/v2"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
)

// TestLogin_New verifica il costruttore New, l'allocazione dei campi input,
// i placeholder e i valori iniziali di username e password.
func TestLogin_New(t *testing.T) {
	t.Parallel()

	l := New()

	if l.Username() != "" {
		t.Errorf("atteso username iniziale vuoto, ottenuto %s", l.Username())
	}
	if l.Password() != "" {
		t.Errorf("atteso password iniziale vuota, ottenuto %s", l.Password())
	}
	if l.IsAuthenticated() {
		t.Error("l'utente non dovrebbe risultare autenticato all'avvio")
	}
	if l.Error() != nil {
		t.Errorf("nessun errore atteso all'avvio, ottenuto: %v", l.Error())
	}
	if len(l.inputs) != 2 {
		t.Fatalf("attesi 2 campi input, ottenuti %d", len(l.inputs))
	}
	if l.inputs[0].Placeholder != "Username" {
		t.Errorf("placeholder atteso 'Username', ottenuto '%s'", l.inputs[0].Placeholder)
	}
	if l.inputs[1].Placeholder != "Password" {
		t.Errorf("placeholder atteso 'Password', ottenuto '%s'", l.inputs[1].Placeholder)
	}
}

// TestLogin_GettersEdgeCases verifica i metodi Username() e Password() quando la slice
// inputs è vuota o parziale.
func TestLogin_GettersEdgeCases(t *testing.T) {
	t.Parallel()

	emptyLogin := Login{}
	if emptyLogin.Username() != "" {
		t.Errorf("atteso username vuoto su Login senza inputs, ottenuto %q", emptyLogin.Username())
	}
	if emptyLogin.Password() != "" {
		t.Errorf("atteso password vuota su Login senza inputs, ottenuto %q", emptyLogin.Password())
	}

	full := New()
	oneInputLogin := Login{inputs: full.inputs[0:1]}
	if oneInputLogin.Password() != "" {
		t.Errorf("atteso password vuota su Login con 1 solo input, ottenuto %q", oneInputLogin.Password())
	}
}

// TestLogin_Init verifica che Init restituisca il comando di lampeggio del cursore.
func TestLogin_Init(t *testing.T) {
	t.Parallel()

	l := New()
	cmd := l.Init()
	if cmd == nil {
		t.Fatal("Init() dovrebbe restituire un tea.Cmd non nil (textinput.Blink)")
	}
}

// TestLogin_Update_SubmitMessages verifica la corretta gestione dei messaggi
// di successo (LoginSubmitMsg) e di errore (SubmitErrorMsg).
func TestLogin_Update_SubmitMessages(t *testing.T) {
	t.Parallel()

	l := New()
	l.loading = true

	// Test SubmitErrorMsg
	errFake := errors.New("credenziali errate")
	res, cmd := l.Update(SubmitErrorMsg{Err: errFake})
	l = res.(Login)

	if cmd != nil {
		t.Error("cmd atteso nil su SubmitErrorMsg")
	}
	if l.loading {
		t.Error("loading dovrebbe essere false dopo SubmitErrorMsg")
	}
	if l.success {
		t.Error("success dovrebbe essere false dopo SubmitErrorMsg")
	}
	if l.Error() == nil || l.Error().Error() != "credenziali errate" {
		t.Errorf("errore inatteso memorizzato: %v", l.Error())
	}

	// Test LoginSubmitMsg
	res, cmd = l.Update(LoginSubmitMsg{Username: "admin", Password: "pwd"})
	l = res.(Login)

	if cmd != nil {
		t.Error("cmd atteso nil su LoginSubmitMsg")
	}
	if l.loading {
		t.Error("loading dovrebbe essere false dopo LoginSubmitMsg")
	}
	if !l.IsAuthenticated() {
		t.Error("IsAuthenticated() dovrebbe restituire true dopo LoginSubmitMsg")
	}
	if l.Error() != nil {
		t.Errorf("errore dovrebbe essere nil dopo login con successo, ottenuto: %v", l.Error())
	}
}

// TestLogin_Update_Navigation verifica la navigazione ciclica del focus con up, down ed enter.
func TestLogin_Update_Navigation(t *testing.T) {
	t.Parallel()

	l := New()

	// Navigazione verso il basso col tasto "down"
	res, _ := l.Update(tea.KeyPressMsg{Text: "down"})
	l = res.(Login)
	if l.focusIndex != 1 {
		t.Errorf("atteso focusIndex 1 dopo tasto down, ottenuto %d", l.focusIndex)
	}

	// Wrapping circolare verso il basso (da 1 a 0)
	res, _ = l.Update(tea.KeyPressMsg{Text: "down"})
	l = res.(Login)
	if l.focusIndex != 0 {
		t.Errorf("atteso wrap focusIndex a 0 dopo down su ultimo campo, ottenuto %d", l.focusIndex)
	}

	// Wrapping circolare verso l'alto (da 0 a 1)
	res, _ = l.Update(tea.KeyPressMsg{Text: "up"})
	l = res.(Login)
	if l.focusIndex != 1 {
		t.Errorf("atteso wrap focusIndex a 1 dopo up su primo campo, ottenuto %d", l.focusIndex)
	}

	// Navigazione con enter dal campo 0 al campo 1
	l = New()
	res, _ = l.Update(tea.KeyPressMsg{Text: "enter"})
	l = res.(Login)
	if l.focusIndex != 1 {
		t.Errorf("atteso focusIndex 1 dopo enter sul primo campo, ottenuto %d", l.focusIndex)
	}
}

// TestLogin_Update_WindowSizeAndEsc verifica WindowSizeMsg e il tasto esc (tea.Quit).
func TestLogin_Update_WindowSizeAndEsc(t *testing.T) {
	t.Parallel()

	l := New()

	// WindowSizeMsg non modifica lo stato
	res, cmd := l.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
	l = res.(Login)
	if cmd != nil {
		t.Error("cmd atteso nil su WindowSizeMsg")
	}

	// Tasto esc restituisce tea.Quit
	res, cmd = l.Update(tea.KeyPressMsg{Text: "esc"})
	if cmd == nil {
		t.Fatal("atteso tea.Quit su esc, ottenuto cmd nil")
	}
}

// TestLogin_Update_TypingInputs verifica che la digitazione di caratteri
// popoli i campi Username e Password.
func TestLogin_Update_TypingInputs(t *testing.T) {
	t.Parallel()

	l := New()

	// Digita 'm', 'a', 'r', 'i', 'o' su campo 0
	for _, ch := range "mario" {
		res, _ := l.Update(tea.KeyPressMsg{Text: string(ch)})
		l = res.(Login)
	}
	if l.Username() != "mario" {
		t.Errorf("username atteso 'mario', ottenuto '%s'", l.Username())
	}

	// Passa a password
	res, _ := l.Update(tea.KeyPressMsg{Text: "enter"})
	l = res.(Login)

	// Digita 'p', 'a', 's', 's' su campo 1
	for _, ch := range "pass" {
		res, _ := l.Update(tea.KeyPressMsg{Text: string(ch)})
		l = res.(Login)
	}
	if l.Password() != "pass" {
		t.Errorf("password attesa 'pass', ottenuta '%s'", l.Password())
	}
}

// TestLogin_Submit_Success verifica il flusso completo di sottomissione con backend mock.
func TestLogin_Submit_Success(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/login" && r.Method == http.MethodPost {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"token": "mock-jwt-token-12345",
				"role":  "admin",
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	client := api.New(ts.URL)
	l := New()
	l.Client = client

	// Sposta il focus al secondo campo (Password)
	res, _ := l.Update(tea.KeyPressMsg{Text: "enter"})
	l = res.(Login)

	// Simula invio premendo Enter sul secondo campo
	res, cmd := l.Update(tea.KeyPressMsg{Text: "enter"})
	l = res.(Login)

	if !l.loading {
		t.Error("loading dovrebbe essere true durante la sottomissione")
	}
	if cmd == nil {
		t.Fatal("atteso tea.Cmd non-nil su submit")
	}

	msg := cmd()
	submitMsg, ok := msg.(LoginSubmitMsg)
	if !ok {
		t.Fatalf("atteso LoginSubmitMsg, ottenuto %T", msg)
	}

	if submitMsg.Username != l.Username() || submitMsg.Password != l.Password() {
		t.Errorf("credenziali non corrispondenti nel messaggio")
	}
	if client.Token() != "mock-jwt-token-12345" {
		t.Errorf("token salvato non valido: %s", client.Token())
	}

	// Invio durante il loading viene ignorato
	res, ignoreCmd := l.Update(tea.KeyPressMsg{Text: "enter"})
	if ignoreCmd != nil {
		t.Error("enter durante loading dovrebbe restituire cmd nil")
	}
	_ = res
}

// TestLogin_Submit_WithoutClient verifica che la sottomissione senza client API
// emetta direttamente il messaggio LoginSubmitMsg.
func TestLogin_Submit_WithoutClient(t *testing.T) {
	t.Parallel()

	l := New()
	l.Client = nil

	// Sposta a password e invia
	res, _ := l.Update(tea.KeyPressMsg{Text: "enter"})
	l = res.(Login)

	res, cmd := l.Update(tea.KeyPressMsg{Text: "enter"})
	if cmd == nil {
		t.Fatal("atteso tea.Cmd non-nil")
	}

	msg := cmd()
	if _, ok := msg.(LoginSubmitMsg); !ok {
		t.Fatalf("atteso LoginSubmitMsg, ottenuto %T", msg)
	}
}

// TestLogin_Submit_BackendError verifica la gestione di un errore HTTP restituito dal backend.
func TestLogin_Submit_BackendError(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("credenziali errate"))
	}))
	defer ts.Close()

	client := api.New(ts.URL)
	l := New()
	l.Client = client

	res, _ := l.Update(tea.KeyPressMsg{Text: "enter"})
	l = res.(Login)

	res, cmd := l.Update(tea.KeyPressMsg{Text: "enter"})
	l = res.(Login)

	msg := cmd()
	errMsg, ok := msg.(SubmitErrorMsg)
	if !ok {
		t.Fatalf("atteso SubmitErrorMsg, ottenuto %T", msg)
	}
	if !strings.Contains(errMsg.Err.Error(), "credenziali errate") {
		t.Errorf("errore inatteso: %v", errMsg.Err)
	}
}

// TestLogin_BlockedAccount verifica il blocco dell'account quando l'errore contiene "bloccato".
func TestLogin_BlockedAccount(t *testing.T) {
	t.Parallel()

	l := New()
	res, _ := l.Update(SubmitErrorMsg{Err: errors.New("account utente bloccato")})
	l = res.(Login)

	if !l.isBlocked() {
		t.Error("isBlocked() atteso true")
	}

	// Con account bloccato, i tasti ordinari sono ignorati
	res, cmd := l.Update(tea.KeyPressMsg{Text: "a"})
	if cmd != nil {
		t.Error("tasto durante blocco dovrebbe restituire cmd nil")
	}

	// Con account bloccato, esc produce tea.Quit
	res, cmd = l.Update(tea.KeyPressMsg{Text: "esc"})
	if cmd == nil {
		t.Fatal("esc con account bloccato deve restituire tea.Quit")
	}

	// Verifica la vista con account bloccato
	v := l.View()
	if !strings.Contains(v.Content, "Account bloccato") {
		t.Errorf("vista bloccata attesa con messaggio 'Account bloccato', ottenuta:\n%s", v.Content)
	}
	if !strings.Contains(v.Content, "[esc] Chiudi applicazione") {
		t.Errorf("vista bloccata deve mostrare footer di chiusura, ottenuta:\n%s", v.Content)
	}
}

// TestLogin_View_VariStates verifica il rendering di View() in tutti gli stati:
// normale, in caricamento, con errore generico e con successo.
func TestLogin_View_VariStates(t *testing.T) {
	t.Parallel()

	// 1. Stato normale
	l := New()
	v := l.View()
	if !strings.Contains(v.Content, "AUTENTICAZIONE UTENTE") {
		t.Error("View normale non contiene il titolo atteso")
	}
	if !strings.Contains(v.Content, "[↑/↓] Naviga") {
		t.Error("View normale non contiene il footer di navigazione")
	}

	// Con cursore nascosto
	l.cursorMode = cursor.CursorHide
	vHidden := l.View()
	if vHidden.Cursor != nil {
		t.Error("con cursorMode CursorHide il cursore deve essere nil")
	}

	// 2. Stato in caricamento
	l = New()
	l.loading = true
	vLoading := l.View()
	if !strings.Contains(vLoading.Content, "Autenticazione in corso...") {
		t.Errorf("View di caricamento non contiene il testo atteso:\n%s", vLoading.Content)
	}

	// 3. Stato con errore generico
	l = New()
	l.err = errors.New("password non valida")
	vErr := l.View()
	if !strings.Contains(vErr.Content, "Autenticazione Fallita") {
		t.Errorf("View di errore non contiene 'Autenticazione Fallita':\n%s", vErr.Content)
	}

	// 4. Stato con successo
	l = New()
	l.success = true
	vSuccess := l.View()
	if !strings.Contains(vSuccess.Content, "Login effettuato con successo!") {
		t.Errorf("View di successo non contiene il testo atteso:\n%s", vSuccess.Content)
	}
}
