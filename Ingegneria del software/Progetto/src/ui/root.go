package ui

import tea "github.com/charmbracelet/bubbletea"

type AppModel struct {
	//contiene il segnaposto con la schermata in cui ci si trova
	statoAttuale sessionState

	//di seguito le schermate possibili
	loginModel     LoginModel
	magazzinoModel MagazzinoModel
}

// Funzione principale pre-avvio interfaccia
func (m AppModel) Init() tea.Cmd {
	return nil
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	//in base alla schermata in cui siamo decide dove eseguire le azioni
	switch m.statoAttuale {

	case loginView:
		//nel caso siamo nella schermata di login
		nuovoLogin, nuovoCmd := m.loginModel.Update(msg)
		m.loginModel = nuovoLogin
		cmd = nuovoCmd
	}

	//nel caso vengano premuti dei tasti speciali
	if KeyMsg, ok := msg.(tea.KeyMsg); ok {
		if KeyMsg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	}

	return m, cmd
}
