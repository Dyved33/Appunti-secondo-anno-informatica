package auth

type LoginRequestDTO struct {// Cosa viene inviato dal client per effettuare il login
    Username string `json:"username"`
    Password string `json:"password"`
}

type LoginResponseDTO struct {// Cosa viene restituito dal server dopo il login
    Token string `json:"token"`
}