package ordini

import (
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
	style "github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/styles"
)

// spostaCursore sposta il cursore nella lista degli ordini o dei resi in base alla sezione corrente.
func (m *Model) spostaCursore(delta int) {
	switch m.sezione {
	case ORDINI_ENTRATA:
		nuovo := m.cursoreEntrata + delta
		if nuovo >= 0 && nuovo < len(m.ordiniEntrata) {
			m.cursoreEntrata = nuovo
			m.offsetEntrata = style.PaginaCalcola(m.cursoreEntrata, m.offsetEntrata, len(m.ordiniEntrata))
		}
	case ORDINI_USCITA:
		nuovo := m.cursoreUscita + delta
		if nuovo >= 0 && nuovo < len(m.ordiniUscita) {
			m.cursoreUscita = nuovo
			m.offsetUscita = style.PaginaCalcola(m.cursoreUscita, m.offsetUscita, len(m.ordiniUscita))
		}
	case RESI:
		nuovo := m.cursoreResi + delta
		if nuovo >= 0 && nuovo < len(m.resi) {
			m.cursoreResi = nuovo
			m.offsetResi = style.PaginaCalcola(m.cursoreResi, m.offsetResi, len(m.resi))
		}
	}
}

// toggleEspansione alterna lo stato di espansione dell'ordine o del reso selezionato nella sezione corrente.
func (m *Model) toggleEspansione() {
	switch m.sezione {
	case ORDINI_ENTRATA:
		if len(m.ordiniEntrata) > 0 {
			id := m.ordiniEntrata[m.cursoreEntrata].ID
			m.espansiEntrata[id] = !m.espansiEntrata[id]
		}
	case ORDINI_USCITA:
		if len(m.ordiniUscita) > 0 {
			id := m.ordiniUscita[m.cursoreUscita].ID
			m.espansiUscita[id] = !m.espansiUscita[id]
		}
	case RESI:
		if len(m.resi) > 0 {
			id := m.resi[m.cursoreResi].ID
			m.espansiResi[id] = !m.espansiResi[id]
		}
	}
}

// renderingStatoOrdine restituisce una stringa formattata in base allo stato dell'ordine.
func renderingStatoOrdine(s api.StatoOrdine, testoFormattato string) string {
	switch s {
	case api.StatoConfermato:
		return style.InfoStile.Render(testoFormattato)
	case api.StatoEvaso:
		return style.SuccessStile.Render(testoFormattato)
	case api.StatoAnnullato:
		return style.ErrorStile.Render(testoFormattato)
	default:
		return style.InfoStile.Render(testoFormattato)
	}
}

// renderingStatoReso restituisce una stringa formattata in base allo stato del reso.
func renderingStatoReso(s api.StatoReso, testoFormattato string) string {
	switch s {
	case api.StatoInAttesa:
		return style.InfoStile.Render(testoFormattato)
	case api.StatoAccettato:
		return style.SuccessStile.Render(testoFormattato)
	case api.StatoRifiutato:
		return style.ErrorStile.Render(testoFormattato)
	default:
		return style.InfoStile.Render(testoFormattato)
	}
}

// troncaStringa tronca una stringa se supera la lunghezza massima specificata, aggiungendo "..." alla fine.
func troncaStringa(s string, maxLen int) string {
	if len(s) > maxLen {
		if maxLen > 3 {
			return s[:maxLen-3] + "..."
		}
		return s[:maxLen]
	}
	return s
}
