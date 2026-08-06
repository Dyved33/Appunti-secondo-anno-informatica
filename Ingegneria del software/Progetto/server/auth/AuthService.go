package auth

import (
	"context"
	"time"

	"github.com/tommasogarofani/ingsw_25-26_group-11/server/domain"
)

type userService interface {
	UserService_Login(ctx context.Context, username, password string) (*domain.Utente, error)
}

type sessionRepository interface {
	SessionRepository_Create(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) error
}

type AuthService struct {
	UserService userService
	SessionRepo sessionRepository
}

func NewAuthService(us userService, sr sessionRepository) *AuthService {
	return &AuthService{
		UserService: us,
		SessionRepo: sr,
	}
}

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