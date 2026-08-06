package login

import (
	"strings"

	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
	style "github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/styles"

	"charm.land/bubbles/v2/cursor"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

// Init inizializza il componente Login al momento dell'avvio della vista Bubble Tea.
//
// Come è stato scritto il codice:
// La funzione rispetta l'interfaccia tea.Model e restituisce il comando textinput.Blink,
// il quale attiva l'animazione di lampeggio del cursore per il campo di testo correntemente focalizzato.
func (l Login) Init() tea.Cmd {
	return textinput.Blink
}

// Update elabora i messaggi e gli eventi in ingresso aggiornando lo stato del modello Login.
//
// Come è stato scritto il codice:
// 1. Intercetta i messaggi di risposta di sottomissione (SubmitErrorMsg e LoginSubmitMsg):
//   - Su SubmitErrorMsg registra l'errore in `l.err`, imposta `l.loading = false` e `l.success = false`.
//   - Su LoginSubmitMsg resetta l'errore, imposta `l.loading = false` e `l.success = true`.
//
// 2. Intercetta i messaggi di pressione dei tasti (`tea.KeyPressMsg`).
// 3. Se viene premuto "esc", restituisce il modello e il comando `tea.Quit` per chiudere l'applicazione.
// 4. Se vengono premuti i tasti "up" o "down":
//   - Incrementa o decrementa l'indice di focus (`focusIndex`).
//   - Gestisce la navigazione circolare tra i campi (0 e 1).
//   - Invoca l'helper `updateFocus` per impostare `Focus()` sul campo selezionato e `Blur()` sugli altri.
//
// 5. Se viene premuto "enter" (tasto Invio):
//   - Se il cursore è su un campo precedente all'ultimo (`l.focusIndex < len(l.inputs)-1`),
//     sposta il focus al campo successivo (da Username a Password) chiamando `updateFocus(l.focusIndex + 1)`.
//   - Se il cursore si trova sull'ultimo campo (Password) e non c'è una richiesta in corso (`!l.loading`):
//     imposta `l.loading = true`, azzera gli errori precedenti e raccoglie le credenziali immesse (Username e Password).
//     Se è presente un client API (`l.Client`), invoca il comando asincrono `doSubmitCmd`;
//     altrimenti emette un messaggio `LoginSubmitMsg` con i dati raccolti.
//
// 6. Infine, delega l'aggiornamento dello stato interno di ciascun campo di testo al metodo `updateInputs(msg)`.
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

	case tea.KeyPressMsg:
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

// View costruisce ed esegue il rendering dell'interfaccia utente (UI) per la schermata di login.
//
// Come è stato scritto il codice:
//  1. Inserisce l'intestazione stilizzata "AUTENTICAZIONE UTENTE" mediante `style.LoginTitleStile`.
//  2. Scorre i campi di input concatenando le loro viste. Se il cursore è visibile sul campo attivo,
//     calcola la coordinata verticale offset.
//  3. Se l'autenticazione è in corso, aggiunge la stringa di caricamento; se vi è un errore,
//     applica `style.LoginErrorStile`; in caso di successo, applica `style.LoginSuccessStile`.
//  4. Avvolge l'intero blocco testuale nella cornice stilizzata `style.LoginBoxStile` di Lipgloss e restituisce `tea.View`.
func (l Login) View() tea.View {
	var b strings.Builder
	var c *tea.Cursor

	b.WriteString(style.LoginTitleStile.Render("AUTENTICAZIONE UTENTE") + "\n\n")

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
		b.WriteString("\n\n" + style.LoginErrorStile.Render("Errore: "+l.err.Error()))
	} else if l.success {
		b.WriteString("\n\n" + style.LoginSuccessStile.Render("Login effettuato con successo!"))
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
func doSubmitCmd(client *api.ApiClient, username, password string) tea.Cmd {
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
