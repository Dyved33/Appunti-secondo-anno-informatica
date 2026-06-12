package DTO

type ProductInputDTO struct {
	Codice      string  `json:"codice"`
	Barcode     string  `json:"barcode"`
	Nome        string  `json:"nome"`
	Descrizione string  `json:"descrizione"`
	CategoriaID int     `json:"categoria_id"`
	PrezzoAcq   float64 `json:"prezzo_acquisto"`
	PrezzoVend  float64 `json:"prezzo_vendita"`
	SogliaMin   int     `json:"soglia_minima"`
}

type ProductOutputDTO struct {
	ID          int     `json:"id"`
	Codice      string  `json:"codice"`
	Barcode     string  `json:"barcode"`
	Nome        string  `json:"nome"`
	Descrizione string  `json:"descrizione"`
	CategoriaID int     `json:"categoria_id"`
	PrezzoAcq   float64 `json:"prezzo_acquisto"`
	PrezzoVend  float64 `json:"prezzo_vendita"`
	SogliaMin   int     `json:"soglia_minima"`
}

type PriceUpdateInputDTO struct {
	PrezzoAcq  float64 `json:"prezzo_acquisto"`
	PrezzoVend float64 `json:"prezzo_vendita"`
}
