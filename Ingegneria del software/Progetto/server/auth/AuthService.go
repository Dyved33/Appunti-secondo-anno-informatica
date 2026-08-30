package auth

import (
	"context"
	"time"

	"github.com/tommasogarofani/ingsw_25-26_group-11/server/domain"
)

// userService definisce il sottoinsieme di metodi del servizio utenti
// necessario all'autenticazione.
type userService interface {
	UserService_Login(ctx context.Context, username, password string) (*domain.Utente, error)
}

// sessionRepository definisce il sottoinsieme di metodi del repository delle
// sessioni necessario per persistere le sessioni di login.
type sessionRepository interface {
	SessionRepository_Create(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) error
}

// AuthService coordina il flusso di login: valida le credenziali tramite il
// servizio utenti, genera il token JWT e, se configurato, registra la sessione
// nel repository.
type AuthService struct {
	UserService userService
	SessionRepo sessionRepository
}

// NewAuthService crea una nuova istanza di AuthService.
func NewAuthService(us userService, sr sessionRepository) *AuthService {
	return &AuthService{
		UserService: us,
		SessionRepo: sr,
	}
}

// AuthService_Login autentica un utente tramite username e password, genera un
// token JWT valido 24 ore e ne registra la sessione nel repository.
// Restituisce il token e il ruolo dell'utente, oppure un errore se le
// credenziali non sono valide o la registrazione della sessione fallisce.
func (s *AuthService) AuthService_Login(ctx context.Context, username, password string) (string, string, error) {
	u, err := s.UserService.UserService_Login(ctx, username, password)
	if err != nil {
		return "", "", err
	}

	roleStr := string(u.Ruolo)
	token, err := GenerateToken(u.ID, u.Username, roleStr)
	if err != nil {
		return "", "", err
	}

	if s.SessionRepo != nil {
		tokenHash := HashToken(token)
		expiresAt := time.Now().Add(24 * time.Hour)
		if err := s.SessionRepo.SessionRepository_Create(ctx, u.ID, tokenHash, expiresAt); err != nil {
			return "", "", err
		}
	}

	return token, roleStr, nil
}
