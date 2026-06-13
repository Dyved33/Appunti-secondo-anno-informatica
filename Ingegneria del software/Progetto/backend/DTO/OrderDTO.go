package DTO

import "time"

type OrderLineInputDTO struct {
	ProdottoID     int     `json:"prodotto_id"`
	Quantita       int     `json:"quantita"`
	PrezzoUnitario float64 `json:"prezzo_unitario"`
}

type OrderInputDTO struct {
	Numero      string              `json:"numero"`
	Cliente     string              `json:"cliente,omitempty"`
	Fornitore   string              `json:"fornitore,omitempty"`
	Note        string              `json:"note"`
	CreatedByID int                 `json:"created_by_id"`
	Righe       []OrderLineInputDTO `json:"righe"`
}

type OrderLineOutputDTO struct {
	ID             int     `json:"id"`
	ProdottoID     int     `json:"prodotto_id"`
	Quantita       int     `json:"quantita"`
	PrezzoUnitario float64 `json:"prezzo_unitario"`
}

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

type DDTInputDTO struct {
	Mittente     string `json:"mittente"`
	Destinatario string `json:"destinatario"`
}

type DDTOutputDTO struct {
	ID             int       `json:"id"`
	Numero         string    `json:"numero"`
	OrdineUscitaID int       `json:"ordine_uscita_id"`
	Data           time.Time `json:"data"`
	Mittente       string    `json:"mittente"`
	Destinatario   string    `json:"destinatario"`
	PathPDF        string    `json:"path_pdf"`
}

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

type ResoInputDTO struct {
	OrdineUscitaID int    `json:"ordine_uscita_id"`
	ProdottoID     int    `json:"prodotto_id"`
	Quantita       int    `json:"quantita"`
	Motivo         string `json:"motivo"`
	CreatedByID    int    `json:"created_by_id"`
}
