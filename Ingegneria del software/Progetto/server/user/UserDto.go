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

// RecoveryRequestDTO rappresenta la richiesta per il recupero delle credenziali.
// L'utente fornisce il suo username per ricevere la domanda di sicurezza.
type RecoveryRequestDTO struct {
	Username string `json:"username"`
}

// RecoveryResponseDTO contiene la domanda di sicurezza associata all'utente.
// Viene restituita quando si richiede il recupero delle credenziali.
type RecoveryResponseDTO struct {
	DomandaSicurezza string `json:"domanda_sicurezza"`
}

// ResetPasswordRequestDTO rappresenta la richiesta di reset della password.
// L'utente fornisce il suo username, la risposta alla domanda di sicurezza e la nuova password.
type ResetPasswordRequestDTO struct {
	Username    string `json:"username"`
	Risposta    string `json:"risposta"`
	NewPassword string `json:"new_password"`
}
