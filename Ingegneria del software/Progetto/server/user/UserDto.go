package user

// LoginRequestDTO rappresenta la richiesta di login da parte di un utente.
// Contiene le credenziali base: username e password.
type LoginRequestDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponseDTO rappresenta la risposta dopo un tentativo di login.
// Se Success è true, contiene il nome utente e il ruolo; altrimenti ErrorMessage spiega il motivo del fallimento.
type LoginResponseDTO struct {
	Success      bool   `json:"success"`
	Username     string `json:"username"`
	Ruolo        string `json:"ruolo"`
	ErrorMessage string `json:"error_message,omitempty"`
}
