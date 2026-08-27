package api

import (
	"fmt"
	"net/http"
	"strconv"
)

func (c *APIClient) GetCatalogo() ([]Prodotto, []string, error) {
	var prodottiServer []serverProductDTO
	if err := c.doJSON(http.MethodGet, "/api/v1/products", nil, &prodottiServer); err != nil {
		return nil, nil, err
	}

	nomeCategoriaByID, categorie, err := c.getCategorie()
	if err != nil {
		return nil, nil, err
	}

	prodotti := make([]Prodotto, 0, len(prodottiServer))
	for _, prodotto := range prodottiServer {
		prodotti = append(prodotti, Prodotto{
			ID:             prodotto.ID,
			Codice:         prodotto.Codice,
			Nome:           prodotto.Nome,
			Categoria:      nomeCategoriaByID[prodotto.CategoriaID],
			Prezzo:         prodotto.PrezzoVend,
			CategoriaID:    prodotto.CategoriaID,
			PrezzoAcquisto: prodotto.PrezzoAcq,
			PrezzoVendita:  prodotto.PrezzoVend,
			Attivo:         true,
		})
	}

	return prodotti, categorie, nil
}

func (c *APIClient) AggiornaPrezzo(prodottoID int, nuovoPrezzo float64) error {
	return c.doJSON(http.MethodPut, "/api/v1/products/"+strconv.Itoa(prodottoID)+"/price", map[string]float64{
		"prezzo_vendita": nuovoPrezzo,
	}, nil)
}

// Ritorna un prodotto dato il suo codice. Se non esiste, ritorna un errore.
func (c *APIClient) GetProductByCode(code string) (Prodotto, error) {
	prodotti, _, err := c.GetCatalogo()
	if err != nil {
		return Prodotto{}, err
	}

	prodottiByCode := make(map[string]Prodotto, len(prodotti))
	for _, p := range prodotti {
		prodottiByCode[p.Codice] = p
	}
	if _, ok := prodottiByCode[code]; !ok {
		return Prodotto{}, fmt.Errorf("prodotto con codice %s non trovato", code)
	}
	return prodottiByCode[code], nil
}

// GetProductByBarcode recupera un prodotto dal backend usando il suo barcode (EAN-13).
// Effettua una richiesta GET a /api/v1/products/barcode/{barcode} e restituisce il prodotto trovato oppure un errore.
func (c *APIClient) GetProductByBarcode(barcode string) (Prodotto, error) {
	var dto serverProductDTO
	if err := c.doJSON(http.MethodGet, "/api/v1/products/barcode/"+barcode, nil, &dto); err != nil {
		return Prodotto{}, err
	}

	nomeCategoriaByID, _, err := c.getCategorie()
	if err != nil {
		return Prodotto{}, err
	}

	return Prodotto{
		ID:             dto.ID,
		Codice:         dto.Codice,
		Nome:           dto.Nome,
		Categoria:      nomeCategoriaByID[dto.CategoriaID],
		CategoriaID:    dto.CategoriaID,
		PrezzoAcquisto: dto.PrezzoAcq,
		PrezzoVendita:  dto.PrezzoVend,
		Attivo:         true,
	}, nil
}

// getCategorie recupera le categorie dal backend e restituisce la mappa ID -> nome
// insieme alla lista dei nomi nell'ordine restituito dal server.
func (c *APIClient) getCategorie() (map[int]string, []string, error) {
	var categorieServer []serverCategoryDTO
	if err := c.doJSON(http.MethodGet, "/api/v1/categories", nil, &categorieServer); err != nil {
		return nil, nil, err
	}

	categorie := make([]string, 0, len(categorieServer))
	nomeCategoriaByID := make(map[int]string, len(categorieServer))
	for _, categoria := range categorieServer {
		nomeCategoriaByID[categoria.ID] = categoria.Nome
		categorie = append(categorie, categoria.Nome)
	}
	return nomeCategoriaByID, categorie, nil
}

// GetReportGiacenze recupera il report dell'inventario dal backend HTTP.
func (c *APIClient) GetReportGiacenze() (*InventoryReport, error) {
	var report InventoryReport

	err := c.doJSON(http.MethodGet, "/api/v1/products/inventory-report", nil, &report)
	if err != nil {
		return nil, fmt.Errorf("errore durante il recupero del report giacenze: %w", err)
	}

	return &report, nil
}

// CreaProdotto invia al server la richiesta di creazione di un nuovo prodotto.
func (c *APIClient) CreaProdotto(input ClientProductDTO) (Prodotto, error) {
	var dto serverProductDTO
	if err := c.doJSON(http.MethodPost, "/api/v1/products/", input, &dto); err != nil {
		return Prodotto{}, err
	}

	nomeCategoriaByID, _, err := c.getCategorie()
	if err != nil {
		return Prodotto{}, err
	}

	return Prodotto{
		ID:             dto.ID,
		Codice:         dto.Codice,
		Nome:           dto.Nome,
		Categoria:      nomeCategoriaByID[dto.CategoriaID],
		CategoriaID:    dto.CategoriaID,
		PrezzoAcquisto: dto.PrezzoAcq,
		PrezzoVendita:  dto.PrezzoVend,
		Attivo:         true,
	}, nil
}

// EliminaProdotto richiede al server di eliminare un prodotto per ID.
func (c *APIClient) EliminaProdotto(id int) error {
	return c.doJSON(http.MethodDelete, "/api/v1/products/"+strconv.Itoa(id), nil, nil)
}

// CreaCategoria invia al server la richiesta di creazione di una nuova categoria.
func (c *APIClient) CreaCategoria(nome, descrizione string) error {
	return c.doJSON(http.MethodPost, "/api/v1/categories/", map[string]string{
		"nome":        nome,
		"descrizione": descrizione,
	}, nil)
}

// EliminaCategoria richiede al server di eliminare una categoria per ID.
func (c *APIClient) EliminaCategoria(id int) error {
	return c.doJSON(http.MethodDelete, "/api/v1/categories/"+strconv.Itoa(id), nil, nil)
}

// GetCategorieIDs recupera le categorie dal backend e restituisce la mappa ID → nome
// e la slice di nomi nell'ordine restituito dal server.
func (c *APIClient) GetCategorieIDs() (map[int]string, []string, error) {
	var categorieServer []serverCategoryDTO
	if err := c.doJSON(http.MethodGet, "/api/v1/categories", nil, &categorieServer); err != nil {
		return nil, nil, err
	}

	categorie := make([]string, 0, len(categorieServer))
	nomeCategoriaByID := make(map[int]string, len(categorieServer))
	for _, cat := range categorieServer {
		nomeCategoriaByID[cat.ID] = cat.Nome
		categorie = append(categorie, cat.Nome)
	}
	return nomeCategoriaByID, categorie, nil
}
