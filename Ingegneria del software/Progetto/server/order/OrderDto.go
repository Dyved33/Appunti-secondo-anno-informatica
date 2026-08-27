package order

import "time"

// RigaOrdineEntrataDTO rappresenta una singola riga (prodotto + quantità) in una richiesta di ordine.
type RigaOrdineEntrataDTO struct {
	ProdottoID     int     `json:"prodotto_id"`
	Quantita       int     `json:"quantita"`
	PrezzoUnitario float64 `json:"prezzo_unitario"`
}

// OrdineEntrataDTO rappresenta la richiesta di creazione di un ordine di entrata.
// Contiene i dati dell'ordine e l'elenco delle righe (prodotti con quantità e prezzi).
type OrdineEntrataDTO struct {
	Numero      string                 `json:"numero"`
	Cliente     string                 `json:"cliente,omitempty"`
	Fornitore   string                 `json:"fornitore,omitempty"`
	Note        string                 `json:"note"`
	CreatedByID int                    `json:"created_by_id"`
	Righe       []RigaOrdineEntrataDTO `json:"righe"`
}

// RigaOrdineUscitaDTO rappresenta una singola riga di ordine nella risposta.
type RigaOrdineUscitaDTO struct {
	ID             int     `json:"id"`
	ProdottoID     int     `json:"prodotto_id"`
	Quantita       int     `json:"quantita"`
	PrezzoUnitario float64 `json:"prezzo_unitario"`
}

// OrdineUscitaDTO rappresenta i dati di un ordine di uscita nella risposta.
type OrdineUscitaDTO struct {
	ID          int                   `json:"id"`
	Numero      string                `json:"numero"`
	Cliente     string                `json:"cliente,omitempty"`
	Fornitore   string                `json:"fornitore,omitempty"`
	Data        time.Time             `json:"data"`
	Stato       string                `json:"stato"`
	Note        string                `json:"note"`
	CreatedByID int                   `json:"created_by_id"`
	DDTID       *int                  `json:"ddt_id,omitempty"`
	Righe       []RigaOrdineUscitaDTO `json:"righe,omitempty"`
}

// DDTEntrataDTO rappresenta i dati per la generazione di un DDT (Documento di Trasporto).
type DDTEntrataDTO struct {
	Mittente     string `json:"mittente"`
	Destinatario string `json:"destinatario"`
}

// DDTUscitaDTO rappresenta i dati di un DDT nella risposta.
type DDTUscitaDTO struct {
	ID             int       `json:"id"`
	Numero         string    `json:"numero"`
	OrdineUscitaID int       `json:"ordine_uscita_id"`
	Data           time.Time `json:"data"`
	Mittente       string    `json:"mittente"`
	Destinatario   string    `json:"destinatario"`
	PathPDF        string    `json:"path_pdf"`
}

// FatturaDTO rappresenta i dati di una fattura o nota di credito nella risposta.
type FatturaDTO struct {
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

type ResoDTO struct {
	ID             int       `json:"id"`
	OrdineUscitaID int       `json:"ordine_uscita_id"`
	ProdottoID     int       `json:"prodotto_id"`
	Quantita       int       `json:"quantita"`
	Motivo         string    `json:"motivo"`
	Stato          string    `json:"stato"`
	CreatedByID    int       `json:"created_by_id"`
	Data           time.Time `json:"data"`
}
