package user

import (
	"context"
	"github.com/tommasogarofani/ingsw_25-26_group-11/src/backend/domain"
)

// userService è il sottoinsieme di UserService di cui Controller ha bisogno.
// *UserService lo implementa già implicitamente.
type userService interface {
	UserRepository_Login(ctx context.Context, username, password string) (*domain.Utente, error)
	UserRepository_GetSecurityQuestion(ctx context.Context, username string) (string, error)
	UserRepository_VerifySecurityAnswerAndResetPassword(ctx context.Context, username, answer, newPassword string) error
}

// Trasla le richieste in ingresso (DTO) alle chiamate al Service e formatta le risposte.
type UserController struct {
	Service userService
}

// NewUserController crea una nuova istanza del controller per gli utenti.
func NewUserController(svc userService) *UserController {
	return &UserController{Service: svc}
}

// Serve per l'autenticazione di un utente, accetta una richiesta di login, chiama il Service e restituisce la risposta formattata.
func (c *UserController) UserController_Login(ctx context.Context, req LoginRequestDTO) LoginResponseDTO {
	u, err := c.Service.UserRepository_Login(ctx, req.Username, req.Password)
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
func (c *UserController) UserController_GetSecurityQuestion(ctx context.Context, req RecoveryRequestDTO) (RecoveryResponseDTO, error) {
	q, err := c.Service.UserRepository_GetSecurityQuestion(ctx, req.Username)
	if err != nil {
		return RecoveryResponseDTO{}, err
	}
	return RecoveryResponseDTO{DomandaSicurezza: q}, nil
}

// Chiama il Service che verifica la risposta e, se corretta, resetta la password.
func (c *UserController) UserController_ResetPassword(ctx context.Context, req ResetPasswordRequestDTO) error {
	return c.Service.UserRepository_VerifySecurityAnswerAndResetPassword(ctx, req.Username, req.Risposta, req.NewPassword)
}
