package api

import (
	"fmt"
	"net/http"
	"strconv"
)

// GetFatture recupera tutte le fatture e note di credito dal server API e le mappa in una slice di Fattura.
func (c *APIClient) GetFatture() ([]Fattura, error) {
	var fattureServer []serverFatturaDTO
	if err := c.doJSON(http.MethodGet, "/api/v1/orders/invoices", nil, &fattureServer); err != nil {
		return nil, err
	}

	var noteCreditoServer []serverFatturaDTO
	if err := c.doJSON(http.MethodGet, "/api/v1/orders/credit-notes", nil, &noteCreditoServer); err != nil {
		return nil, err
	}

	ordini, err := c.GetOrdini()
	if err != nil {
		return nil, err
	}
	clienteByOrdine := make(map[int]string, len(ordini))
	for _, ordine := range ordini {
		clienteByOrdine[ordine.ID] = ordine.Cliente
	}

	fatture := make([]Fattura, 0, len(fattureServer)+len(noteCreditoServer))
	for _, fattura := range fattureServer {
		fatture = append(fatture, Fattura{
			ID:             fattura.ID,
			OrdineUscitaID: fattura.OrdineUscitaID,
			Numero:         fattura.Numero,
			Anno:           fattura.Data.Year(),
			Cliente:        clienteByOrdine[fattura.OrdineUscitaID],
			Importo:        fattura.Totale,
			NotaCredito:    false,
			Data:           fattura.Data,
		})
	}
	for _, fattura := range noteCreditoServer {
		fatture = append(fatture, Fattura{
			ID:             fattura.ID,
			OrdineUscitaID: fattura.OrdineUscitaID,
			Numero:         fattura.Numero,
			Anno:           fattura.Data.Year(),
			Cliente:        clienteByOrdine[fattura.OrdineUscitaID],
			Importo:        fattura.Totale,
			NotaCredito:    true,
			Data:           fattura.Data,
		})
	}

	return fatture, nil
}

// GetOrdini recupera tutti gli ordini dal server API e li mappa in una slice di Ordine.
func (c *APIClient) GetOrdini() ([]Ordine, error) {
	var ordiniServer []serverOrderDTO
	if err := c.doJSON(http.MethodGet, "/api/v1/orders", nil, &ordiniServer); err != nil {
		return nil, err
	}

	prodotti, _, err := c.GetCatalogo()
	if err != nil {
		return nil, err
	}
	prodottiByID := make(map[int]Prodotto, len(prodotti))
	for _, prodotto := range prodotti {
		prodottiByID[prodotto.ID] = prodotto
	}

	ordini := make([]Ordine, 0, len(ordiniServer))
	for _, ordineServer := range ordiniServer {
		ordine := Ordine{
			ID:      ordineServer.ID,
			Cliente: ordineServer.Cliente,
			Stato:   ordineServer.Stato,
		}
		if len(ordineServer.Righe) > 0 {
			riga := ordineServer.Righe[0]
			ordine.ProdottoID = riga.ProdottoID
			ordine.Quantita = riga.Quantita
			ordine.PrezzoUnitario = riga.PrezzoUnitario
			if prodotto, ok := prodottiByID[riga.ProdottoID]; ok {
				ordine.Prodotto = prodotto.Nome
			} else {
				ordine.Prodotto = "Prodotto " + strconv.Itoa(riga.ProdottoID)
			}
		}
		ordini = append(ordini, ordine)
	}

	return ordini, nil
}

// CreaOrdineUscita registra un ordine di uscita e aggiorna le giacenze tramite
// la logica transazionale già presente nel backend.
func (c *APIClient) CreaOrdineUscita(input ClientOrderDTO) error {
	return c.doJSON(http.MethodPost, "/api/v1/orders/crea-uscita", input, nil)
}

// DeleteOrdineUscita elimina un ordine di uscita e aggiorna le giacenze tramite
// la logica transazionale già presente nel backend.
func (c *APIClient) DeleteOrdineUscita(ID int) error {
	path := fmt.Sprintf("/api/v1/orders/delete-uscita/%d", ID)
	return c.doJSON(http.MethodDelete, path, nil, nil)
}

// CreaOrdineEntrata registra un ordine di entrata e aggiorna le giacenze tramite
// la logica transazionale già presente nel backend.
func (c *APIClient) CreaOrdineEntrata(input ClientOrderDTO) error {
	return c.doJSON(http.MethodPost, "/api/v1/orders/crea-entrata", input, nil)
}

// DeleteOrdineEntrata elimina un ordine di entrata e aggiorna le giacenze tramite
// la logica transazionale già presente nel backend.
func (c *APIClient) DeleteOrdineEntrata(ID int) error {
	path := fmt.Sprintf("/api/v1/orders/delete-entrata/%d", ID)
	return c.doJSON(http.MethodDelete, path, nil, nil)
}

// GetReportKPI recupera i dati per il report KPI dal server API e li mappa in una slice di KPI.
func (c *APIClient) GetReportKPI() ([]KPI, error) {
	fatture, err := c.GetFatture()
	if err != nil {
		return nil, err
	}

	ordini, err := c.GetOrdini()
	if err != nil {
		return nil, err
	}

	var fatturatoNetto float64
	var noteCredito int
	for _, fattura := range fatture {
		if fattura.NotaCredito {
			fatturatoNetto -= fattura.Importo
			noteCredito++
			continue
		}
		fatturatoNetto += fattura.Importo
	}

	var ordiniEvasi int
	for _, ordine := range ordini {
		if ordine.Stato == "evaso" {
			ordiniEvasi++
		}
	}

	return []KPI{
		{Etichetta: "Fatturato netto", Valore: fmt.Sprintf("€ %.2f", fatturatoNetto)},
		{Etichetta: "Fatture", Valore: strconv.Itoa(len(fatture) - noteCredito)},
		{Etichetta: "Ordini evasi", Valore: strconv.Itoa(ordiniEvasi)},
		{Etichetta: "Note credito", Valore: strconv.Itoa(noteCredito)},
	}, nil
}

// GetOrdiniEntrata recupera tutti gli ordini di entrata dal server API e li mappa in una slice di OrdineEntrata.
func (c *APIClient) GetOrdiniEntrata() ([]OrdineEntrata, error) {
	var ordiniDTO []serverOrderDTO
	if err := c.doJSON(http.MethodGet, "/api/v1/orders/entrata", nil, &ordiniDTO); err != nil {
		return nil, err
	}

	// Recupero del catalogo prodotti per popolare il puntatore *Prodotto nelle righe
	prodotti, _, err := c.GetCatalogo()
	if err != nil {
		return nil, err
	}
	prodottiByID := make(map[int]Prodotto, len(prodotti))
	for _, p := range prodotti {
		prodottiByID[p.ID] = p
	}

	ordiniEntrata := make([]OrdineEntrata, 0, len(ordiniDTO))
	for _, dto := range ordiniDTO {
		ordine := OrdineEntrata{
			ID:          dto.ID,
			Numero:      dto.Numero,
			Fornitore:   dto.Fornitore,
			Data:        dto.Data,
			Note:        dto.Note,
			CreatedByID: dto.CreatedByID,
			Righe:       make([]RigaOrdineEntrata, 0, len(dto.Righe)),
		}

		// Mappatura di tutte le righe dell'ordine
		for _, rigaDTO := range dto.Righe {
			riga := RigaOrdineEntrata{
				ID:              rigaDTO.ID,
				OrdineEntrataID: dto.ID,
				ProdottoID:      rigaDTO.ProdottoID,
				Quantita:        rigaDTO.Quantita,
				PrezzoUnitario:  rigaDTO.PrezzoUnitario,
			}

			// Assegnazione del riferimento al prodotto se presente nella mappa
			if prod, ok := prodottiByID[rigaDTO.ProdottoID]; ok {
				// Crea una copia locale per evitare che tutte le righe puntino alla stessa variabile di ciclo
				pCopy := prod
				riga.Prodotto = &pCopy
			}

			ordine.Righe = append(ordine.Righe, riga)
		}

		ordiniEntrata = append(ordiniEntrata, ordine)
	}

	return ordiniEntrata, nil
}

// GetOrdiniUscita recupera tutti gli ordini di uscita dal server API e li mappa in una slice di OrdineUscita.
func (c *APIClient) GetOrdiniUscita() ([]OrdineUscita, error) {
	var ordiniDTO []serverOrderDTO
	if err := c.doJSON(http.MethodGet, "/api/v1/orders/uscita", nil, &ordiniDTO); err != nil {
		return nil, err
	}

	// Recupero del catalogo prodotti per popolare il puntatore *Prodotto nelle righe
	prodotti, _, err := c.GetCatalogo()
	if err != nil {
		return nil, err
	}
	prodottiByID := make(map[int]Prodotto, len(prodotti))
	for _, p := range prodotti {
		prodottiByID[p.ID] = p
	}

	ordiniUscita := make([]OrdineUscita, 0, len(ordiniDTO))
	for _, dto := range ordiniDTO {
		ordine := OrdineUscita{
			ID:          dto.ID,
			Numero:      dto.Numero,
			Cliente:     dto.Cliente,
			Data:        dto.Data,
			Stato:       StatoOrdine(dto.Stato),
			Note:        dto.Note,
			CreatedByID: dto.CreatedByID,
			DDTID:       dto.DDTID,
			Righe:       make([]RigaOrdineUscita, 0, len(dto.Righe)),
		}

		// Mappatura di tutte le righe dell'ordine
		for _, rigaDTO := range dto.Righe {
			riga := RigaOrdineUscita{
				ID:             rigaDTO.ID,
				OrdineUscitaID: dto.ID,
				ProdottoID:     rigaDTO.ProdottoID,
				Quantita:       rigaDTO.Quantita,
				PrezzoUnitario: rigaDTO.PrezzoUnitario,
			}

			// Assegnazione del riferimento al prodotto se presente nella mappa
			if prod, ok := prodottiByID[rigaDTO.ProdottoID]; ok {
				// Crea una copia locale per evitare che tutte le righe puntino alla stessa variabile di ciclo
				pCopy := prod
				riga.Prodotto = &pCopy
			}

			ordine.Righe = append(ordine.Righe, riga)
		}

		ordiniUscita = append(ordiniUscita, ordine)
	}

	return ordiniUscita, nil
}

// GetResi recupera tutti i resi dal server API e li mappa in una slice di Reso.
func (c *APIClient) GetResi() ([]Reso, error) {
	var ordiniDTO []serverResoDTO
	if err := c.doJSON(http.MethodGet, "/api/v1/orders/resi", nil, &ordiniDTO); err != nil {
		return nil, err
	}

	resi := make([]Reso, 0, len(ordiniDTO))
	for _, dto := range ordiniDTO {
		ordine := Reso{
			ID:             dto.ID,
			OrdineUscitaID: dto.OrdineUscitaID,
			ProdottoID:     dto.ProdottoID,
			Quantita:       dto.Quantita,
			Motivo:         dto.Motivo,
			Stato:          StatoReso(dto.Stato),
			Data:           dto.Data,
			CreatedByID:    dto.CreatedByID,
		}
		resi = append(resi, ordine)
	}

	return resi, nil
}

// GetDDT richiede al server la generazione di un DDT (Documento di Trasporto)
// per l'ordine di uscita con l'ID specificato, usando i dati mittente e
// destinatario forniti, e restituisce il DDT generato.
func (c *APIClient) GetDDT(ID int, dett ClientDDTDTO) (DDT, error) {
	var ddtDTO serverDDTDTO
	path := fmt.Sprintf("/api/v1/orders/%d/create-ddt", ID)

	if err := c.doJSON(http.MethodPost, path, dett, &ddtDTO); err != nil {
		return DDT{}, err
	}

	return DDT{
		ID:             ddtDTO.ID,
		Numero:         ddtDTO.Numero,
		OrdineUscitaID: ddtDTO.OrdineUscitaID,
		Data:           ddtDTO.Data,
		Mittente:       ddtDTO.Mittente,
		Destinatario:   ddtDTO.Destinatario,
		PathPDF:        ddtDTO.PathPDF,
	}, nil
}

// CreaFattura genera una fattura per l'ordine di uscita con l'ID specificato.
func (c *APIClient) CreaFattura(ordineID int) error {
	path := fmt.Sprintf("/api/v1/orders/create-fattura/%d", ordineID)
	return c.doJSON(http.MethodPost, path, nil, nil)
}

// EliminaFattura elimina una fattura per il suo ID.
func (c *APIClient) EliminaFattura(fatturaID int) error {
	path := fmt.Sprintf("/api/v1/orders/delete-fattura/%d", fatturaID)
	return c.doJSON(http.MethodDelete, path, nil, nil)
}

// CreaNotaCredito genera una nota di credito per l'ordine di uscita con l'ID specificato.
func (c *APIClient) CreaNotaCredito(ordineID int) error {
	path := fmt.Sprintf("/api/v1/orders/create-credit-notes/%d", ordineID)
	return c.doJSON(http.MethodPost, path, nil, nil)
}

// GetOrdiniEvasiLeggero recupera gli ordini di uscita evasi senza risolvere i
// dettagli prodotto (più leggero di GetOrdiniUscita). Serve per la creazione fatture.
func (c *APIClient) GetOrdiniEvasiLeggero() ([]OrdineUscita, error) {
	var ordiniDTO []serverOrderDTO
	if err := c.doJSON(http.MethodGet, "/api/v1/orders/uscita", nil, &ordiniDTO); err != nil {
		return nil, err
	}

	var evasi []OrdineUscita
	for _, dto := range ordiniDTO {
		if dto.Stato != string(StatoEvaso) {
			continue
		}
		evasi = append(evasi, OrdineUscita{
			ID:      dto.ID,
			Numero:  dto.Numero,
			Cliente: dto.Cliente,
			Data:    dto.Data,
			Stato:   StatoOrdine(dto.Stato),
		})
	}
	return evasi, nil
}

// CreaReso invia al server la richiesta di registrazione di un reso.
func (c *APIClient) CreaReso(reso ClientResoDTO) error {
	return c.doJSON(http.MethodPost, "/api/v1/orders/crea-reso", reso, nil)
}

// AggiornaOrdineUscita invia la richiesta di aggiornamento di un ordine di uscita al server.
func (c *APIClient) AggiornaOrdineUscita(id int, input ClientOrderDTO) error {
	path := fmt.Sprintf("/api/v1/orders/update-ordine-uscita/%d", id)
	return c.doJSON(http.MethodPut, path, input, nil)
}

// AggiornaStatoReso aggiorna lo stato di un reso esistente.
func (c *APIClient) AggiornaStatoReso(id int, stato string) error {
	payload := struct {
		ID    int    `json:"id"`
		Stato string `json:"stato"`
	}{ID: id, Stato: stato}
	return c.doJSON(http.MethodPut, "/api/v1/orders/update-resi", payload, nil)
}

// DeleteReso elimina un reso esistente.
func (c *APIClient) DeleteReso(id int) error {
	path := fmt.Sprintf("/api/v1/orders/delete-reso/%d", id)
	return c.doJSON(http.MethodDelete, path, nil, nil)
}
