package ordini

import (
	tea "charm.land/bubbletea/v2"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
)

// loadOrdiniCmd costruisce un tea.Cmd asincrono per il recupero in parallelo/sequenza
// degli ordini di entrata, ordini di uscita e resi tramite l'API backend.
func (m Model) loadOrdiniCmd() tea.Cmd {
	return func() tea.Msg {
		entrata, err := m.Client.GetOrdiniEntrata()
		if err != nil {
			return ErroreMsg{Err: err}
		}
		uscita, err := m.Client.GetOrdiniUscita()
		if err != nil {
			return ErroreMsg{Err: err}
		}
		resi, err := m.Client.GetResi()
		if err != nil {
			return ErroreMsg{Err: err}
		}
		return OrdiniLoadedMsg{Entrata: entrata, Uscita: uscita, Resi: resi}
	}
}

// creaOrdineEntrataCmd invia l'ordine di acquisto al backend per la persistenza.
func (m Model) creaOrdineEntrataCmd(fornitore, numero, note string, righe []api.ClientOrderLineDTO) tea.Cmd {
	return func() tea.Msg {
		dto := api.ClientOrderDTO{
			Numero:    numero,
			Fornitore: fornitore,
			Note:      note,
			Righe:     righe,
		}
		err := m.Client.CreaOrdineEntrata(dto)
		if err != nil {
			return ErroreMsg{Err: err}
		}
		return OrdineEntrataCreatoMsg{}
	}
}

// creaOrdineUscitaCmd invia l'ordine di vendita al backend per la persistenza.
func (m Model) creaOrdineUscitaCmd(cliente, numero, note string, righe []api.ClientOrderLineDTO) tea.Cmd {
	return func() tea.Msg {
		dto := api.ClientOrderDTO{
			Numero:  numero,
			Cliente: cliente,
			Note:    note,
			Righe:   righe,
		}
		err := m.Client.CreaOrdineUscita(dto)
		if err != nil {
			return ErroreMsg{Err: err}
		}
		return OrdineUscitaCreatoMsg{}
	}
}

// creaResoCmd invia la richiesta di reso al backend per la persistenza.
func (m Model) creaResoCmd(ordineID, prodottoID, quantita int, motivo string) tea.Cmd {
	return func() tea.Msg {
		dto := api.ClientResoDTO{
			OrdineUscitaID: ordineID,
			ProdottoID:     prodottoID,
			Quantita:       quantita,
			Motivo:         motivo,
			Stato:          "In attesa",
		}
		err := m.Client.CreaReso(dto)
		if err != nil {
			return ErroreMsg{Err: err}
		}
		return ResoCreatoMsg{}
	}
}

// eliminaOrdineEntrata invia la richiesta di eliminazione dell'ordine di entrata al backend.
func (m Model) deleteOrdineEntrataCmd(ID int) tea.Cmd {
	return func() tea.Msg {
		err := m.Client.DeleteOrdineEntrata(ID)
		if err != nil {
			return ErroreMsg{Err: err}
		}
		return OrdineEntrataEliminatoMsg{}
	}
}

// eliminaOrdineUscita invia la richiesta di eliminazione dell'ordine di uscita al backend.
func (m Model) deleteOrdineUscitaCmd(ID int) tea.Cmd {
	return func() tea.Msg {
		err := m.Client.DeleteOrdineUscita(ID)
		if err != nil {
			return ErroreMsg{Err: err}
		}
		return OrdineUscitaEliminatoMsg{}
	}
}

// eliminaReso invia la richiesta di eliminazione del reso al backend.
func (m Model) deleteResoCmd(ID int) tea.Cmd {
	return func() tea.Msg {
		err := m.Client.DeleteReso(ID)
		if err != nil {
			return ErroreMsg{Err: err}
		}
		return ResoEliminatoMsg{}
	}
}
