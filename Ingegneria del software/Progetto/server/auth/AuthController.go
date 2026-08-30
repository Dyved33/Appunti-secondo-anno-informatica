package auth

import "context"

// authService definisce il sottoinsieme di metodi di AuthService necessario
// al controller.
type authService interface {
	AuthService_Login(ctx context.Context, username, password string) (string, string, error)
}

// AuthController riceve le richieste di login, le inoltra al servizio di
// autenticazione e formatta la risposta da restituire al client.
type AuthController struct {
	Service authService
}

// NewAuthController crea una nuova istanza del controller di autenticazione.
func NewAuthController(svc authService) *AuthController {
	return &AuthController{Service: svc}
}

// AuthController_Login gestisce una richiesta di login, delegando al servizio
// la verifica delle credenziali, e restituisce un DTO contenente token e ruolo
// oppure un errore.
func (c *AuthController) AuthController_Login(ctx context.Context, req LoginRequestDTO) (LoginResponseDTO, error) {
	token, role, err := c.Service.AuthService_Login(ctx, req.Username, req.Password)
	if err != nil {
		return LoginResponseDTO{}, err
	}
	return LoginResponseDTO{Token: token, Role: role}, nil
}
