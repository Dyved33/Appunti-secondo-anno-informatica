package login

import (
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"

	"charm.land/bubbles/v2/cursor"
	"charm.land/bubbles/v2/textinput"
)

// utente rappresenta la struttura dati interna per la gestione delle credenziali (username e password).
type utente struct {
	username string
	password string
}

// Login rappresenta il modello di stato della pagina di login nell'architettura TUI Bubble Tea.
// Mantiene l'indice del campo attualmente focalizzato (focusIndex), la slice dei moduli di input di testo (inputs),
// la modalità del cursore (cursorMode), il client API (Client), ed i campi di stato (err, success, loading)
// per tracciare il risultato dell'autenticazione ed il feedback grafico.
type Login struct {
	focusIndex int
	inputs     []textinput.Model
	cursorMode cursor.Mode
	Client     *api.ApiClient
	err        error
	success    bool
	loading    bool
}

// LoginSubmitMsg è il messaggio emesso quando l'utente invia il form di login premendo invio sull'ultimo campo.
type LoginSubmitMsg struct {
	Username string
	Password string
}

// SubmitErrorMsg rappresenta un messaggio di errore ritornato in seguito a un tentativo di autenticazione fallito.
type SubmitErrorMsg struct {
	Err error
}

// New inizializza e restituisce un nuovo modello Login con i campi di testo preconfigurati.
//
// Come è stato scritto il codice:
// 1. Alloca una slice di 2 componenti textinput.Model per gestire l'input di Username (indice 0) e Password (indice 1).
// 2. Cicla sugli elementi impostando un limite di 12 caratteri per ciascun campo.
// 3. Tramite uno switch sull'indice:
//   - Indice 0 (Username): viene impostato il testo placeholder "Username" e riceve il focus iniziale tramite Focus().
//   - Indice 1 (Password): viene impostato il placeholder "Password" e la modalità EchoPassword con il carattere speciale '🞽'.
//
// 4. Assegna ciascuna istanza configurata al relativo indice della slice e restituisce la struttura Login.
func New() Login {

	l := Login{
		inputs: make([]textinput.Model, 2),
	}

	var t textinput.Model

	for i := range l.inputs {
		t = textinput.New()
		t.CharLimit = 12

		// Impostazione del focus e delle proprietà specifiche per ciascun campo
		switch i {
		case 0:
			t.Placeholder = "Username"
			t.Focus()
		case 1:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '🞽'
		}

		l.inputs[i] = t

	}

	return l
}

// Username restituisce il testo correntemente inserito nel campo Username (indice 0).
//
// Come è stato scritto il codice:
// Verifica che la slice degli input contenga almeno un elemento ed ne estrae il valore corrente.
func (l Login) Username() string {
	if len(l.inputs) > 0 {
		return l.inputs[0].Value()
	}
	return ""
}

// Password restituisce il testo correntemente inserito nel campo Password (indice 1).
//
// Come è stato scritto il codice:
// Verifica che la slice degli input contenga almeno due elementi ed ne estrae il valore corrente.
func (l Login) Password() string {
	if len(l.inputs) > 1 {
		return l.inputs[1].Value()
	}
	return ""
}

// IsAuthenticated restituisce true se il login è stato completato con successo.
//
// Come è stato scritto il codice:
// Controlla il flag bool interno success.
func (l Login) IsAuthenticated() bool {
	return l.success
}

// Error restituisce l'eventuale errore riscontrato durante la fase di login.
//
// Come è stato scritto il codice:
// Restituisce l'interfaccia error memorizzata nella struttura Login.
func (l Login) Error() error {
	return l.err
}
