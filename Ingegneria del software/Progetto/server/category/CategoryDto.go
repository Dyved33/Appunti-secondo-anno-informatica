package category

// CategoryInputDTO rappresenta i dati in ingresso per la creazione o modifica di una categoria.
// Contiene il nome e una descrizione opzionale della categoria di prodotti.
type CategoryInputDTO struct {
	ID          int    `json:"id"`
	Nome        string `json:"nome"`
	Descrizione string `json:"descrizione"`
}

// CategoryOutputDTO rappresenta i dati in uscita di una categoria.
// Include l'ID assegnato dal database, il nome e la descrizione.
type CategoryOutputDTO struct {
	ID          int    `json:"id"`
	Nome        string `json:"nome"`
	Descrizione string `json:"descrizione"`
}
