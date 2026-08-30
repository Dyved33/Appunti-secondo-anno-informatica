package ordini

import (
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
	style "github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/styles"
)

// passoRotellina è il numero di righe di cui la lista scorre ad ogni evento della rotellina del mouse.
const passoRotellina = 3

// overheadChrome è il numero di righe fisse occupate ATTORNO al corpo della pagina
// (logo, intestazione, barra dei tab, bordi della finestra e riga vuota interna).
const overheadChrome = 23

// righeFallback è il numero di righe visibili usato finché non arriva il primo WindowSizeMsg.
const righeFallback = 14

// righeVisibili calcola quante righe di tabella entrano nello spazio residuo dello schermo.
//
// Il corpo della pagina Ordini ha sempre: intestazione (titolo + KPI + barra tab,
// misurata in 10 righe), testata della tabella (2 righe) e footer (2 righe).
func (m Model) righeVisibili() int {
	if m.winH <= 0 {
		return righeFallback
	}
	disponibili := m.winH - overheadChrome - 14
	if disponibili < 1 {
		disponibili = 1
	}
	return disponibili
}

// clampScroll riporta gli offset di tutte le sezioni nei rispettivi intervalli validi.
func (m *Model) clampScroll() {
	m.clampSection(len(m.ordiniEntrata), &m.offsetEntrata)
	m.clampSection(len(m.ordiniUscita), &m.offsetUscita)
	m.clampSection(len(m.resi), &m.offsetResi)
}

// clampSection riporta il singolo offset nell'intervallo valido [0, n-ruote].
func (m Model) clampSection(n int, offset *int) {
	maxOffset := n - m.righeVisibili()
	if maxOffset < 0 {
		maxOffset = 0
	}
	if *offset > maxOffset {
		*offset = maxOffset
	}
	if *offset < 0 {
		*offset = 0
	}
}

// mostraCursore aggiorna l'offset della sezione corrente così che il cursore resti visibile.
func (m *Model) mostraCursore() {
	vis := m.righeVisibili()
	switch m.sezione {
	case ORDINI_ENTRATA:
		if m.cursoreEntrata < m.offsetEntrata {
			m.offsetEntrata = m.cursoreEntrata
		} else if m.cursoreEntrata >= m.offsetEntrata+vis {
			m.offsetEntrata = m.cursoreEntrata - vis + 1
		}
	case ORDINI_USCITA:
		if m.cursoreUscita < m.offsetUscita {
			m.offsetUscita = m.cursoreUscita
		} else if m.cursoreUscita >= m.offsetUscita+vis {
			m.offsetUscita = m.cursoreUscita - vis + 1
		}
	case RESI:
		if m.cursoreResi < m.offsetResi {
			m.offsetResi = m.cursoreResi
		} else if m.cursoreResi >= m.offsetResi+vis {
			m.offsetResi = m.cursoreResi - vis + 1
		}
	}
	m.clampScroll()
}

// spostaCursore sposta il cursore nella lista degli ordini o dei resi in base alla sezione corrente.
func (m *Model) spostaCursore(delta int) {
	switch m.sezione {
	case ORDINI_ENTRATA:
		nuovo := m.cursoreEntrata + delta
		if nuovo >= 0 && nuovo < len(m.ordiniEntrata) {
			m.cursoreEntrata = nuovo
			m.mostraCursore()
		} else if nuovo > len(m.ordiniEntrata)-1 {
			m.cursoreEntrata = len(m.ordiniEntrata) - 1
			m.mostraCursore()
		}
	case ORDINI_USCITA:
		nuovo := m.cursoreUscita + delta
		if nuovo >= 0 && nuovo < len(m.ordiniUscita) {
			m.cursoreUscita = nuovo
			m.mostraCursore()
		} else if nuovo > len(m.ordiniUscita)-1 {
			m.cursoreUscita = len(m.ordiniUscita) - 1
			m.mostraCursore()
		}
	case RESI:
		nuovo := m.cursoreResi + delta
		if nuovo >= 0 && nuovo < len(m.resi) {
			m.cursoreResi = nuovo
			m.mostraCursore()
		} else if nuovo > len(m.resi)-1 {
			m.cursoreResi = len(m.resi) - 1
			m.mostraCursore()
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
