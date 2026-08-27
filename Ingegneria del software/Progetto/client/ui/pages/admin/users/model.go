// Package users fornisce la gestione dell'interfaccia amministrativa per la gestione degli utenti.
package users

import (
	"charm.land/bubbles/v2/textinput"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
)

// mode rappresenta le modalità operative della pagina di amministrazione utenti.
type mode int

const (
	modeSearch        mode = iota // modeSearch: ricerca di un user per username
	modeCreate                    // modeCreate: creazione di un nuovo utente
	modeEdit                      // modeEdit: modifica di un utente esistente
	modeResetPassword             // modeResetPassword: reset della password di un utente
)

// ruoliDisponibili contiene la lista dei ruoli disponibili nel sistema.
var ruoliDisponibili = []string{"magazziniere", "contabile", "admin"}

// Model rappresenta lo stato dell'interfaccia amministrativa per la gestione degli utenti.
type Model struct {
	Client *api.APIClient // Client: client API per le comunicazioni con il backend
	mode   mode           // mode: modalità operativa corrente (ricerca, creazione, modifica)

	searchUsername textinput.Model // searchUsername: campo input per la ricerca per username
	user           *api.User       // user: utente trovato e caricato dal server
	loading        bool            // loading: indica se è in corso il caricamento dei dati
	err            error           // err: errore eventuale durante le operazioni
	success        string          // success: messaggio di successo dell'ultima operazione
	confirmDelete  bool            // confirmDelete: indica se è attiva la richiesta di conferma per l'eliminazione dell'utente caricato

	// Campi per il form di creazione
	createUsername textinput.Model // createUsername: campo input per lo username nuovo utente
	createPassword textinput.Model // createPassword: campo input per la password del nuovo utente
	createFocus    int             // createFocus: indice del campo focalizzato nel form di creazione (0: Username, 1: Password, 2: Ruolo)
	createRuoloIdx int             // createRuoloIdx: indice del ruolo selezionato nel form di creazione

	// Campi per il form di modifica
	editUsername textinput.Model // editUsername: campo input per lo username nell'edit
	editFocus    int             // editFocus: indice del campo focalizzato nel form di modifica (0: Username, 1: Ruolo)
	editRuoloIdx int             // editRuoloIdx: indice del ruolo selezionato nel form di modifica

	// Campi per il form di reset password
	resetPassword textinput.Model // resetPassword: campo input per la nuova password
	resetConferma textinput.Model // resetConferma: campo input per la conferma della nuova password
	resetFocus    int             // resetFocus: indice del campo focalizzato nel form di reset (0: Password, 1: Conferma)
}

// userFoundMsg rappresenta il messaggio di caricamento riuscito dell'utente cercato.
type userFoundMsg struct {
	User api.User // User: dati dell'utente trovato
}

// userCreatedMsg rappresenta il messaggio di creazione riuscita di un utente.
type userCreatedMsg struct {
	Username string // Username: username dell'utente appena creato
}

// userUpdatedMsg rappresenta il messaggio di modifica riuscita di un utente.
type userUpdatedMsg struct {
	Username string // Username: username dell'utente appena modificato
}

// userBlockedToggledMsg rappresenta il messaggio di cambio stato blocco/sblocco dell'utente.
type userBlockedToggledMsg struct {
	User api.User // User: dati dell'utente aggiornato
}

// userPasswordResetMsg rappresenta il messaggio di reset riuscito della password di un utente.
type userPasswordResetMsg struct {
	Username string // Username: username dell'utente di cui è stata reimpostata la password
}

// userDeletedMsg rappresenta il messaggio di eliminazione riuscita di un utente.
type userDeletedMsg struct {
	Username string // Username: username dell'utente appena eliminato
}

// usersErrorMsg rappresenta il messaggio di errore durante un'operazione.
type usersErrorMsg struct {
	Err error // Err: errore verificatosi durante l'operazione
}

// New crea e inizializza un nuovo modello di gestione utenti con il client API fornito.
// Inizializza i campi input di ricerca, creazione e modifica con i placeholder appropriati.
func New(client *api.APIClient) Model {
	searchUsername := textinput.New()
	searchUsername.Placeholder = "Inserisci username da cercare..."
	searchUsername.Focus()

	createUsername := textinput.New()
	createUsername.Placeholder = "Username"

	createPassword := textinput.New()
	createPassword.Placeholder = "Password"
	createPassword.EchoMode = textinput.EchoPassword

	editUsername := textinput.New()
	editUsername.Placeholder = "Username"

	resetPassword := textinput.New()
	resetPassword.Placeholder = "Nuova password"
	resetPassword.EchoMode = textinput.EchoPassword

	resetConferma := textinput.New()
	resetConferma.Placeholder = "Conferma password"
	resetConferma.EchoMode = textinput.EchoPassword

	return Model{
		Client:         client,
		loading:        false,
		mode:           modeSearch,
		searchUsername: searchUsername,
		createUsername: createUsername,
		createPassword: createPassword,
		createRuoloIdx: 0,
		editUsername:   editUsername,
		editRuoloIdx:   0,
		resetPassword:  resetPassword,
		resetConferma:  resetConferma,
	}
}
