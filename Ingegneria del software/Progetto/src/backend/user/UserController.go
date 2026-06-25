package user

import (
	"context"
)

// Trasla le richieste in ingresso (DTO) alle chiamate al Service e formatta le risposte.
type Controller struct {
	Service *Service
}

// NewController crea una nuova istanza del controller per gli utenti.
func NewController(svc *Service) *Controller {
	return &Controller{Service: svc}
}

// Serve per l'autenticazione di un utente, accetta una richiesta di login, chiama il Service e restituisce la risposta formattata.
func (c *Controller) HandleLogin(ctx context.Context, req LoginRequestDTO) LoginResponseDTO {
	u, err := c.Service.Login(ctx, req.Username, req.Password)
	if err != nil {
		return LoginResponseDTO{
			Success:      false,
			ErrorMessage: err.Error(),
		}
	}

	return LoginResponseDTO{
		Success:  true,
		Username: u.Username,
		Ruolo:    string(u.Ruolo),
	}
}

// HandleGetSecurityQuestion recupera la domanda di sicurezza per un utente.
func (c *Controller) HandleGetSecurityQuestion(ctx context.Context, req RecoveryRequestDTO) (RecoveryResponseDTO, error) {
	q, err := c.Service.GetSecurityQuestion(ctx, req.Username)
	if err != nil {
		return RecoveryResponseDTO{}, err
	}
	return RecoveryResponseDTO{DomandaSicurezza: q}, nil
}

// Chiama il Service che verifica la risposta e, se corretta, resetta la password.
func (c *Controller) HandleResetPassword(ctx context.Context, req ResetPasswordRequestDTO) error {
	return c.Service.VerifySecurityAnswerAndResetPassword(ctx, req.Username, req.Risposta, req.NewPassword)
}
