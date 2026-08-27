package login

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"

	tea "charm.land/bubbletea/v2"
)

// TestLogin_New e la configurazione iniziale dei campi di input
func TestLogin_New(t *testing.T) {
	l := New()

	if l.Username() != "" {
		t.Errorf("expected empty initial username, got %s", l.Username())
	}
	if l.Password() != "" {
		t.Errorf("expected empty initial password, got %s", l.Password())
	}
}

// TestLogin_Update_Navigation verifica la navigazione del focus tra i campi Username e Password
func TestLogin_Update_Navigation(t *testing.T) {
	l := New()

	// Navigazione verso il basso col tasto "down"
	res, _ := l.Update(tea.KeyPressMsg{Text: "down"})
	l = res.(Login)
	if l.focusIndex != 1 {
		t.Errorf("expected focusIndex 1 after down key, got %d", l.focusIndex)
	}

	// Navigazione col tasto "enter" dal campo 0 al campo 1
	l = New()
	res, _ = l.Update(tea.KeyPressMsg{Text: "enter"})
	l = res.(Login)
	if l.focusIndex != 1 {
		t.Errorf("expected focusIndex 1 after enter on first field, got %d", l.focusIndex)
	}
}

// TestLogin_Submit_Success verifica che l'invio delle credenziali autentichi il client e memorizzi il token JWT
func TestLogin_Submit_Success(t *testing.T) {
	// Server HTTP mock per simulare la risposta di autenticazione del backend
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/login" && r.Method == http.MethodPost {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{
				"token": "mock-jwt-token-12345",
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

	// Simula l'invio premendo Enter sul secondo campo
	res, cmd := l.Update(tea.KeyPressMsg{Text: "enter"})
	l = res.(Login)

	if cmd == nil {
		t.Fatal("expected non-nil cmd on login submit")
	}

	// Esegue il comando restituito per effettuare la chiamata HTTP al server mock
	msg := cmd()

	submitMsg, ok := msg.(LoginSubmitMsg)
	if !ok {
		t.Fatalf("expected LoginSubmitMsg, got %T", msg)
	}

	if submitMsg.Username != l.Username() || submitMsg.Password != l.Password() {
		t.Errorf("credentials mismatch in submit message")
	}

	// Verifica che il client API abbia memorizzato il token e sia ora autenticato per le chiamate successive
	if client.Token() != "mock-jwt-token-12345" {
		t.Errorf("expected client token 'mock-jwt-token-12345', got '%s'", client.Token())
	}
}
