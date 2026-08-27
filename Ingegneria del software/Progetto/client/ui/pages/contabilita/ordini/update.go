package ordini

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	style "github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/styles"
)

// Init avvia il comando iniziale di caricamento dati dal backend.
func (m Model) Init() tea.Cmd {
	return m.loadOrdiniCmd()
}

// Update risponde alle interazioni dell'utente (pressioni di tasti) e ai messaggi asincroni.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case OrdiniLoadedMsg:
		m.loading = false
		m.err = nil
		m.ordiniEntrata = msg.Entrata
		m.ordiniUscita = msg.Uscita
		m.resi = msg.Resi
		m.cursoreEntrata = 0
		m.cursoreUscita = 0
		m.cursoreResi = 0
		m.offsetEntrata = 0
		m.offsetUscita = 0
		m.offsetResi = 0
		return m, nil

	case ErroreMsg:
		m.loading = false
		m.err = msg.Err
		return m, nil

	case tea.WindowSizeMsg:
		return m, nil

	case OrdineEntrataCreatoMsg, OrdineUscitaCreatoMsg, ResoCreatoMsg:
		m.resetForm()
		m.loading = true
		return m, m.loadOrdiniCmd()

	case OrdineEntrataEliminatoMsg, OrdineUscitaEliminatoMsg, ResoEliminatoMsg:
		m.loading = true
		return m, m.loadOrdiniCmd()

	case tea.KeyPressMsg:
		if m.creazioneEntrata || m.creazioneUscita {
			return m.updateFormCreazione(msg)
		}
		if m.creazioneReso {
			return m.updateFormCreazioneReso(msg)
		}
		if m.confermaEliminazione {
			switch strings.ToLower(msg.String()) {
			case "y", "enter":
				m.loading = true
				m.confermaEliminazione = false
				switch m.sezione {
				case ORDINI_ENTRATA:
					return m, m.deleteOrdineEntrataCmd(m.idDaEliminare)
				case ORDINI_USCITA:
					return m, m.deleteOrdineUscitaCmd(m.idDaEliminare)
				case RESI:
					return m, m.deleteResoCmd(m.idDaEliminare)
				}
			case "n", "esc":
				m.confermaEliminazione = false
				return m, nil
			}
			return m, nil
		}
		switch msg.String() {
		case "tab":
			switch m.sezione {
			case ORDINI_ENTRATA:
				m.sezione = ORDINI_USCITA
			case ORDINI_USCITA:
				m.sezione = RESI
			default:
				m.sezione = ORDINI_ENTRATA
			}

		case "up":
			m.spostaCursore(-1)

		case "down":
			m.spostaCursore(1)

		case "pgup":
			m.spostaCursore(-style.RighePerPagina)

		case "pgdown":
			m.spostaCursore(style.RighePerPagina)

		case "enter", "space":
			m.toggleEspansione()

		case "r":
			m.loading = true
			m.err = nil
			return m, m.loadOrdiniCmd()

		case "ctrl+n":
			switch m.sezione {
			case ORDINI_ENTRATA:
				m.avviaCreazioneOrdine(true)
				return m, m.focusCurrentField()
			case ORDINI_USCITA:
				m.avviaCreazioneOrdine(false)
				return m, m.focusCurrentField()
			case RESI:
				m.avviaCreazioneReso()
				return m, m.focusCurrentField()
			}

		case "ctrl+d":
			switch m.sezione {
			case ORDINI_ENTRATA:
				return m.avviaEliminazioneOrdineEntrata()
			case ORDINI_USCITA:
				return m.avviaEliminazioneOrdineUscita()
			case RESI:
				return m.avviaEliminazioneReso()
			}
		}
	}

	return m, nil
}

// updateFormCreazione gestisce la navigazione e l'input di testo all'interno del form di creazione.
func (m Model) updateFormCreazione(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.resetForm()
		return m, nil

	case "tab", "down":
		m.focusCampo = (m.focusCampo + 1) % 4
		return m, m.focusCurrentField()

	case "shift+tab", "up":
		m.focusCampo = (m.focusCampo - 1 + 4) % 4
		return m, m.focusCurrentField()

	case "ctrl+s":
		return m.salvaOrdine()

	case "enter":
		if m.focusCampo == FocusRighe {
			return m.salvaOrdine()
		}
		m.focusCampo = (m.focusCampo + 1) % 4
		return m, m.focusCurrentField()
	}

	var cmd tea.Cmd
	switch m.focusCampo {
	case FocusControparte:
		m.formControparte, cmd = m.formControparte.Update(msg)
	case FocusNumero:
		m.formNumero, cmd = m.formNumero.Update(msg)
	case FocusNote:
		m.formNote, cmd = m.formNote.Update(msg)
	case FocusRighe:
		m.formRighe, cmd = m.formRighe.Update(msg)
	}

	return m, cmd
}

// updateFormCreazioneReso gestisce la navigazione e l'input di testo nel form dei resi.
func (m Model) updateFormCreazioneReso(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.resetForm()
		return m, nil

	case "tab", "down":
		if m.focusCampo == FocusResoMotivo {
			m.focusCampo = FocusResoOrdineID
		} else {
			m.focusCampo++
		}
		return m, m.focusCurrentField()

	case "shift+tab", "up":
		if m.focusCampo == FocusResoOrdineID {
			m.focusCampo = FocusResoMotivo
		} else {
			m.focusCampo--
		}
		return m, m.focusCurrentField()

	case "ctrl+s":
		return m.salvaReso()

	case "enter":
		if m.focusCampo == FocusResoMotivo {
			return m.salvaReso()
		}
		m.focusCampo++
		return m, m.focusCurrentField()
	}

	var cmd tea.Cmd
	switch m.focusCampo {
	case FocusResoOrdineID:
		m.formResoOrdineID, cmd = m.formResoOrdineID.Update(msg)
	case FocusResoProdottoID:
		m.formResoProdottoID, cmd = m.formResoProdottoID.Update(msg)
	case FocusResoQuantita:
		m.formResoQuantita, cmd = m.formResoQuantita.Update(msg)
	case FocusResoMotivo:
		m.formResoMotivo, cmd = m.formResoMotivo.Update(msg)
	}

	return m, cmd
}

func (m Model) avviaEliminazioneOrdineEntrata() (tea.Model, tea.Cmd) {
	if len(m.ordiniEntrata) == 0 || m.cursoreEntrata < 0 || m.cursoreEntrata >= len(m.ordiniEntrata) {
		return m, nil
	}
	m.confermaEliminazione = true
	m.sezione = ORDINI_ENTRATA
	m.idDaEliminare = m.ordiniEntrata[m.cursoreEntrata].ID
	return m, nil
}

func (m Model) avviaEliminazioneOrdineUscita() (tea.Model, tea.Cmd) {
	if len(m.ordiniUscita) == 0 || m.cursoreUscita < 0 || m.cursoreUscita >= len(m.ordiniUscita) {
		return m, nil
	}
	m.confermaEliminazione = true
	m.sezione = ORDINI_USCITA
	m.idDaEliminare = m.ordiniUscita[m.cursoreUscita].ID
	return m, nil
}

func (m Model) avviaEliminazioneReso() (tea.Model, tea.Cmd) {
	if len(m.resi) == 0 || m.cursoreResi < 0 || m.cursoreResi >= len(m.resi) {
		return m, nil
	}
	m.confermaEliminazione = true
	m.sezione = RESI
	m.idDaEliminare = m.resi[m.cursoreResi].ID
	return m, nil
}
