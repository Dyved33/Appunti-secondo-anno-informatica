package DTO

type CategoryInputDTO struct {
	Nome        string `json:"nome"`
	Descrizione string `json:"descrizione"`
}

type CategoryOutputDTO struct {
	ID          int    `json:"id"`
	Nome        string `json:"nome"`
	Descrizione string `json:"descrizione"`
}