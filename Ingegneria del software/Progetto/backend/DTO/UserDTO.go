package DTO

type LoginRequestDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponseDTO struct {
	Success      bool   `json:"success"`
	Username     string `json:"username"`
	Ruolo        string `json:"ruolo"`
	ErrorMessage string `json:"error_message,omitempty"`
}

type RecoveryRequestDTO struct {
	Username string `json:"username"`
}

type RecoveryResponseDTO struct {
	DomandaSicurezza string `json:"domanda_sicurezza"`
}

type ResetPasswordRequestDTO struct {
	Username    string `json:"username"`
	Risposta    string `json:"risposta"`
	NewPassword string `json:"new_password"`
}