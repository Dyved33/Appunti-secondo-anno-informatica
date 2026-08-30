package user

import (
	"context"
	"github.com/tommasogarofani/ingsw_25-26_group-11/server/domain"
)

// userService è il sottoinsieme di UserService di cui Controller ha bisogno.
// *UserService lo implementa già implicitamente.
type userService interface {
	UserService_Login(ctx context.Context, username, password string) (*domain.Utente, error)
}

// UserController trasla le richieste in ingresso (DTO) alle chiamate al
// Service e formatta le risposte.
type UserController struct {
	Service userService
}

// NewUserController crea una nuova istanza del controller per gli utenti.
func NewUserController(svc userService) *UserController {
	return &UserController{Service: svc}
}

// UserController_Login gestisce l'autenticazione di un utente: accetta una
// richiesta di login, chiama il Service e restituisce la risposta formattata,
// indicando successo/fallimento ed eventuale messaggio di errore.
func (c *UserController) UserController_Login(ctx context.Context, req LoginRequestDTO) LoginResponseDTO {
	u, err := c.Service.UserService_Login(ctx, req.Username, req.Password)
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
