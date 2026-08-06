package order

import "time"

// OrderLineInputDTO rappresenta una singola riga (prodotto + quantità) in una richiesta di ordine.
type OrderLineInputDTO struct {
	ProdottoID     int     `json:"prodotto_id"`
	Quantita       int     `json:"quantita"`
	PrezzoUnitario float64 `json:"prezzo_unitario"`
}

// OrderInputDTO rappresenta la richiesta di creazione di un ordine di entrata o uscita.
// Contiene i dati dell'ordine e l'elenco delle righe (prodotti con quantità e prezzi).
type OrderInputDTO struct {
	Numero      string              `json:"numero"`
	Cliente     string              `json:"cliente,omitempty"`
	Fornitore   string              `json:"fornitore,omitempty"`
	Note        string              `json:"note"`
	CreatedByID int                 `json:"created_by_id"`
	Righe       []OrderLineInputDTO `json:"righe"`
}

// OrderLineOutputDTO rappresenta una singola riga di ordine nella risposta.
type OrderLineOutputDTO struct {
	ID             int     `json:"id"`
	ProdottoID     int     `json:"prodotto_id"`
	Quantita       int     `json:"quantita"`
	PrezzoUnitario float64 `json:"prezzo_unitario"`
}

// OrderOutputDTO rappresenta i dati di un ordine nella risposta.
type OrderOutputDTO struct {
	ID          int                  `json:"id"`
	Numero      string               `json:"numero"`
	Cliente     string               `json:"cliente,omitempty"`
	Fornitore   string               `json:"fornitore,omitempty"`
	Data        time.Time            `json:"data"`
	Stato       string               `json:"stato"`
	Note        string               `json:"note"`
	CreatedByID int                  `json:"created_by_id"`
	DDTID       *int                 `json:"ddt_id,omitempty"`
	Righe       []OrderLineOutputDTO `json:"righe,omitempty"`
}

// DDTInputDTO rappresenta i dati per la generazione di un DDT (Documento di Trasporto).
type DDTInputDTO struct {
	Mittente     string `json:"mittente"`
	Destinatario string `json:"destinatario"`
}

// DDTOutputDTO rappresenta i dati di un DDT nella risposta.
type DDTOutputDTO struct {
	ID             int       `json:"id"`
	Numero         string    `json:"numero"`
	OrdineUscitaID int       `json:"ordine_uscita_id"`
	Data           time.Time `json:"data"`
	Mittente       string    `json:"mittente"`
	Destinatario   string    `json:"destinatario"`
	PathPDF        string    `json:"path_pdf"`
}

// FatturaOutputDTO rappresenta i dati di una fattura o nota di credito nella risposta.
type FatturaOutputDTO struct {
	ID             int       `json:"id"`
	Numero         string    `json:"numero"`
	Tipo           string    `json:"tipo"`
	OrdineUscitaID int       `json:"ordine_uscita_id"`
	Imponibile     float64   `json:"imponibile"`
	IVA            float64   `json:"iva"`
	Totale         float64   `json:"totale"`
	Data           time.Time `json:"data"`
	PathPDF        string    `json:"path_pdf"`
}

// ResoInputDTO rappresenta i dati per la registrazione di un reso di merce.
type ResoInputDTO struct {
	OrdineUscitaID int    `json:"ordine_uscita_id"`
	ProdottoID     int    `json:"prodotto_id"`
	Quantita       int    `json:"quantita"`
	Motivo         string `json:"motivo"`
	CreatedByID    int    `json:"created_by_id"`
}
