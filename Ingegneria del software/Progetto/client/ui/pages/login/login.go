package login

import (
	"strings"

	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
	style "github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/styles"

	"charm.land/bubbles/v2/cursor"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

// Init inizializza il componente Login all'avvio della vista, restituendo il
// comando che attiva il lampeggio del cursore sul campo focalizzato.
func (l Login) Init() tea.Cmd {
	return textinput.Blink
}

// Update elabora i messaggi ed eventi in ingresso aggiornando lo stato del
// modello Login: gestisce i risultati della sottomissione (errore o successo),
// la pressione dei tasti per navigare tra i campi o inviare il form, ed
// eventualmente chiude l'applicazione con il tasto "esc".
func (l Login) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case SubmitErrorMsg:
		l.err = msg.Err
		l.loading = false
		l.success = false
		return l, nil

	case LoginSubmitMsg:
		l.err = nil
		l.loading = false
		l.success = true
		return l, nil

	case tea.WindowSizeMsg:
		return l, nil

	case tea.KeyPressMsg:
		if l.isBlocked() {
			if msg.String() == "esc" {
				return l, tea.Quit
			}
			return l, nil
		}

		switch msg.String() {
		case "esc":
			return l, tea.Quit

		case "up", "down":
			if msg.String() == "up" {
				l.focusIndex--
			} else {
				l.focusIndex++
			}

			// Gestione del wrapping circolare dell'indice di focus
			if l.focusIndex >= len(l.inputs) {
				l.focusIndex = 0
			} else if l.focusIndex < 0 {
				l.focusIndex = len(l.inputs) - 1
			}

			cmd := l.updateFocus(l.focusIndex)
			return l, cmd

		case "enter":
			// Se il focus si trova su un campo precedente all'ultimo (es. Username), avanziamo al campo successivo (Password).
			if l.focusIndex < len(l.inputs)-1 {
				cmd := l.updateFocus(l.focusIndex + 1)
				return l, cmd
			}

			// Evita sottomissioni duplicate durante il caricamento
			if l.loading {
				return l, nil
			}

			l.loading = true
			l.err = nil

			// Se il focus è sull'ultimo campo (Password), effettuiamo l'invio delle credenziali
			username := l.Username()
			password := l.Password()

			if l.Client != nil {
				return l, doSubmitCmd(l.Client, username, password)
			}

			return l, func() tea.Msg {
				return LoginSubmitMsg{
					Username: username,
					Password: password,
				}
			}
		}
	}

	cmd := l.updateInputs(msg)

	return l, cmd
}

// isBlocked verifica se l'account è stato bloccato per troppi tentativi falliti.
func (l Login) isBlocked() bool {
	return l.err != nil && strings.Contains(strings.ToLower(l.err.Error()), "bloccato")
}

// View costruisce il rendering dell'interfaccia utente della schermata di login:
// mostra titolo, campi di input, il feedback di caricamento/errore/successo e la
// barra dei comandi, il tutto racchiuso nella cornice stilizzata del login.
func (l Login) View() tea.View {
	var b strings.Builder
	var c *tea.Cursor

	b.WriteString(style.TitleStile.Render("AUTENTICAZIONE UTENTE"))
	b.WriteString("\n\n")

	for i, in := range l.inputs {
		b.WriteString(l.inputs[i].View())
		if i < len(l.inputs)-1 {
			b.WriteRune('\n')
		}
		if l.cursorMode != cursor.CursorHide && in.Focused() {
			c = in.Cursor()
			if c != nil {
				c.Y += i + 2
			}
		}
	}

	if l.loading {
		b.WriteString("\n\nAutenticazione in corso...")
	} else if l.err != nil {
		b.WriteString("\n\n")
		if l.isBlocked() {
			b.WriteString(style.ErrorStile.Render("Account bloccato. Contattare l'admin."))
		} else {
			b.WriteString(style.ErrorStile.Render("Autenticazione Fallita"))
		}
	} else if l.success {
		b.WriteString("\n\n")
		b.WriteString(style.SuccessStile.Render("Login effettuato con successo!"))
	}

	b.WriteString("\n\n")
	if l.isBlocked() {
		b.WriteString(style.FooterStile.Render("[esc] Chiudi applicazione"))
	} else {
		b.WriteString(style.FooterStile.Render("[↑/↓] Naviga  |  [enter] Accedi  |  [esc] Esci"))
	}

	boxFinale := style.LoginBoxStile.Render(b.String())

	v := tea.NewView(boxFinale)
	v.Cursor = c
	return v
}

// updateFocus è una funzione di supporto per aggiornare lo stato di focus dei campi di testo.
//
// Come è stato scritto il codice:
// Aggiorna `focusIndex` al valore `targetIndex`. Per ciascun campo della slice `inputs`, applica `Focus()`
// al campo che corrisponde all'indice selezionato e `Blur()` su tutti gli altri, restituendo i relativi comandi raggruppati con `tea.Batch`.
func (l *Login) updateFocus(targetIndex int) tea.Cmd {
	l.focusIndex = targetIndex
	cmds := make([]tea.Cmd, len(l.inputs))
	for i := range l.inputs {
		if i == l.focusIndex {
			cmds[i] = l.inputs[i].Focus()
		} else {
			l.inputs[i].Blur()
		}
	}
	return tea.Batch(cmds...)
}

// updateInputs è una funzione di supporto interna per l'aggiornamento bidirezionale dello stato di tutti i campi di input.
//
// Come è stato scritto il codice:
// 1. Inizializza una slice di comandi `tea.Cmd` di lunghezza pari a quella del vettore di input (`l.inputs`).
// 2. Itera su ciascun modulo di input invocando il relativo metodo `Update(msg)`, aggiornando l'istanza del campo e memorizzando il comando ritornato.
// 3. Raggruppa tutti i comandi risultanti tramite `tea.Batch(cmds...)` e li restituisce alla catena di esecuzione di Bubble Tea.
func (l *Login) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(l.inputs))

	for i := range l.inputs {
		l.inputs[i], cmds[i] = l.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

// doSubmitCmd crea un comando tea.Cmd per eseguire in modo asincrono la chiamata di autenticazione API.
//
// Come è stato scritto il codice:
// Restituisce una closure che esegue `client.Login(username, password)`. In caso di errore restituisce un `SubmitErrorMsg`;
// in caso di successo restituisce un `LoginSubmitMsg`.
func doSubmitCmd(client *api.APIClient, username, password string) tea.Cmd {
	return func() tea.Msg {
		err := client.Login(username, password)
		if err != nil {
			return SubmitErrorMsg{Err: err}
		}
		return LoginSubmitMsg{
			Username: username,
			Password: password,
		}
	}
}
