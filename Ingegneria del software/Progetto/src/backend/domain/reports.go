package domain

import "time"

// InventoryItem descrive una riga del report inventario in modo esplicito.
type InventoryItem struct {
	ProdottoID  int     `json:"prodotto_id"`
	Codice      string  `json:"codice"`
	Nome        string  `json:"nome"`
	Categoria   string  `json:"categoria"`
	Quantita    int     `json:"quantita"`
	SogliaMin   int     `json:"soglia_minima"`
	PrezzoVend  float64 `json:"prezzo_vendita"`
	ValoreRiga  float64 `json:"valore_riga"`
	SottoSoglia bool    `json:"sotto_soglia"`
}

// InventoryReport riassume lo stato del magazzino con dati semplici da leggere.
type InventoryReport struct {
	GeneratoIl         time.Time       `json:"generato_il"`
	TotaleQuantita      int             `json:"totale_quantita"`
	ValoreTotale        float64         `json:"valore_totale"`
	Articoli            []InventoryItem `json:"articoli"`
	ArticoliSottoSoglia []InventoryItem `json:"articoli_sotto_soglia"`
}