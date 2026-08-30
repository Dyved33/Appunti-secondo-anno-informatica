package auth

// LoginRequestDTO rappresenta la richiesta di login inviata dal client:
// contiene username e password.
type LoginRequestDTO struct { // Cosa viene inviato dal client per effettuare il login
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponseDTO rappresenta la risposta restituita dal server dopo un login
// riuscito: contiene il token JWT e il ruolo dell'utente.
type LoginResponseDTO struct { // Cosa viene restituito dal server dopo il login
	Token string `json:"token"`
	Role  string `json:"role"`
}
