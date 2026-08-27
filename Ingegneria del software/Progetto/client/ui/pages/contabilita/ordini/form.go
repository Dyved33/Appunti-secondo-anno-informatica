package ordini

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
)

// resetForm reimposta lo stato del form di creazione ordine.
// Aggiorna resetForm per pulire anche lo stato del reso
func (m *Model) resetForm() {
	m.creazioneEntrata = false
	m.creazioneUscita = false
	m.creazioneReso = false // NUOVO
	m.err = nil
	m.formControparte = textinput.Model{}
	m.formNumero = textinput.Model{}
	m.formNote = textinput.Model{}
	m.formRighe = textinput.Model{}
	m.focusCampo = FocusControparte
	m.formResoOrdineID = textinput.Model{}
	m.formResoProdottoID = textinput.Model{}
	m.formResoQuantita = textinput.Model{}
	m.formResoMotivo = textinput.Model{}
	m.focusCampo = FocusControparte
}

// avviaCreazione imposta lo stato del form per la creazione di un ordine di entrata o uscita.
func (m *Model) avviaCreazioneOrdine(isEntrata bool) {
	m.creazioneEntrata = isEntrata
	m.creazioneUscita = !isEntrata
	m.err = nil
	m.focusCampo = FocusControparte

	// Inizializzazione dei campi textinput
	m.formControparte = textinput.New()
	if isEntrata {
		m.formControparte.Placeholder = "Nome Fornitore"
	} else {
		m.formControparte.Placeholder = "Nome Cliente"
	}
	m.formControparte.Focus()

	m.formNumero = textinput.New()
	if isEntrata {
		m.formNumero.Placeholder = "es. ORD-ENT-2026-001"
	} else {
		m.formNumero.Placeholder = "es. ORD-USC-2026-001"
	}

	m.formNote = textinput.New()
	m.formNote.Placeholder = "Note opzionali..."

	m.formRighe = textinput.New()
	m.formRighe.Placeholder = "ID:Qtà:Prezzo (es. 1:100:1.50, 2:5:10.00)"
}

// salvaOrdine valida i campi del form e invia la richiesta di creazione ordine al server.
func (m Model) salvaOrdine() (tea.Model, tea.Cmd) {
	controparte := strings.TrimSpace(m.formControparte.Value())
	if controparte == "" {
		if m.creazioneEntrata {
			m.err = errors.New("il nome fornitore è obbligatorio")
		} else {
			m.err = errors.New("il nome cliente è obbligatorio")
		}
		return m, nil
	}

	numero := strings.TrimSpace(m.formNumero.Value())
	if numero == "" {
		m.err = errors.New("il numero d'ordine è obbligatorio")
		return m, nil
	}

	righe, err := parseRigheInput(m.formRighe.Value())
	if err != nil {
		m.err = err
		return m, nil
	}

	m.loading = true
	m.err = nil
	if m.creazioneEntrata {
		return m, m.creaOrdineEntrataCmd(controparte, numero, m.formNote.Value(), righe)
	}
	if m.creazioneUscita {
		return m, m.creaOrdineUscitaCmd(controparte, numero, m.formNote.Value(), righe)
	}

	return m, nil
}

// focusCurrentField imposta il focus sul campo attualmente selezionato nel form.
func (m *Model) focusCurrentField() tea.Cmd {
	m.formControparte.Blur()
	m.formNumero.Blur()
	m.formNote.Blur()
	m.formRighe.Blur()
	m.formResoOrdineID.Blur()
	m.formResoProdottoID.Blur()
	m.formResoQuantita.Blur()
	m.formResoMotivo.Blur()

	switch m.focusCampo {
	case FocusControparte:
		return m.formControparte.Focus()
	case FocusNumero:
		return m.formNumero.Focus()
	case FocusNote:
		return m.formNote.Focus()
	case FocusRighe:
		return m.formRighe.Focus()

	case FocusResoOrdineID:
		return m.formResoOrdineID.Focus()
	case FocusResoProdottoID:
		return m.formResoProdottoID.Focus()
	case FocusResoQuantita:
		return m.formResoQuantita.Focus()
	case FocusResoMotivo:
		return m.formResoMotivo.Focus()

	}
	return nil
}

// currentFormCursor restituisce il cursore del campo attualmente selezionato nel form.
func (m *Model) currentFormCursor() *tea.Cursor {
	switch m.focusCampo {
	case FocusControparte:
		return m.formControparte.Cursor()
	case FocusNumero:
		return m.formNumero.Cursor()
	case FocusNote:
		return m.formNote.Cursor()
	case FocusRighe:
		return m.formRighe.Cursor()
	case FocusResoOrdineID:
		return m.formResoOrdineID.Cursor()
	case FocusResoProdottoID:
		return m.formResoProdottoID.Cursor()
	case FocusResoQuantita:
		return m.formResoQuantita.Cursor()
	case FocusResoMotivo:
		return m.formResoMotivo.Cursor()
	}
	return nil
}

// parseRigheInput analizza l'input dell'utente per le righe dell'ordine e restituisce un array di ClientOrderLineDTO.
func parseRigheInput(raw string) ([]api.ClientOrderLineDTO, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, errors.New("almeno una riga ordine è richiesta")
	}

	parti := strings.Split(raw, ",")
	righe := make([]api.ClientOrderLineDTO, 0, len(parti))

	for _, p := range parti {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		campi := strings.Split(p, ":")
		if len(campi) != 3 {
			return nil, fmt.Errorf("formato riga non valido %q (usare ID:Quantità:Prezzo)", p)
		}

		prodottoID, err := strconv.Atoi(strings.TrimSpace(campi[0]))
		if err != nil || prodottoID <= 0 {
			return nil, fmt.Errorf("ID prodotto non valido in %q", p)
		}

		quantita, err := strconv.Atoi(strings.TrimSpace(campi[1]))
		if err != nil || quantita <= 0 {
			return nil, fmt.Errorf("quantità non valida in %q", p)
		}

		prezzo, err := strconv.ParseFloat(strings.TrimSpace(campi[2]), 64)
		if err != nil || prezzo < 0 {
			return nil, fmt.Errorf("prezzo non valido in %q", p)
		}

		righe = append(righe, api.ClientOrderLineDTO{
			ProdottoID:     prodottoID,
			Quantita:       quantita,
			PrezzoUnitario: prezzo,
		})
	}

	if len(righe) == 0 {
		return nil, errors.New("nessuna riga valida inserita")
	}

	return righe, nil
}

// avviaCreazioneReso imposta lo stato del form per la creazione di un reso.
func (m *Model) avviaCreazioneReso() {
	m.creazioneReso = true
	m.creazioneEntrata = false
	m.creazioneUscita = false
	m.err = nil
	m.focusCampo = FocusResoOrdineID

	m.formResoOrdineID = textinput.New()
	m.formResoOrdineID.Placeholder = "es. 12 (ID dell'ordine di uscita)"
	m.formResoOrdineID.Focus()

	m.formResoProdottoID = textinput.New()
	m.formResoProdottoID.Placeholder = "es. 5 (ID del prodotto da rendere)"

	m.formResoQuantita = textinput.New()
	m.formResoQuantita.Placeholder = "es. 2"

	m.formResoMotivo = textinput.New()
	m.formResoMotivo.Placeholder = "Motivo del reso (es. difettoso)"
}

// salvaReso valida i campi del form e invia la richiesta di creazione reso al server.
func (m Model) salvaReso() (tea.Model, tea.Cmd) {
	ordineID, err := strconv.Atoi(strings.TrimSpace(m.formResoOrdineID.Value()))
	if err != nil || ordineID <= 0 {
		m.err = errors.New("l'ID ordine di uscita deve essere un numero valido > 0")
		return m, nil
	}

	prodottoID, err := strconv.Atoi(strings.TrimSpace(m.formResoProdottoID.Value()))
	if err != nil || prodottoID <= 0 {
		m.err = errors.New("l'ID prodotto deve essere un numero valido > 0")
		return m, nil
	}

	quantita, err := strconv.Atoi(strings.TrimSpace(m.formResoQuantita.Value()))
	if err != nil || quantita <= 0 {
		m.err = errors.New("la quantità deve essere un numero valido > 0")
		return m, nil
	}

	motivo := strings.TrimSpace(m.formResoMotivo.Value())

	m.loading = true
	m.err = nil
	return m, m.creaResoCmd(ordineID, prodottoID, quantita, motivo)
}
