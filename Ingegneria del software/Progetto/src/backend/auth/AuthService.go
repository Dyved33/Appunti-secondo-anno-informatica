package auth

import (
	"context"

	"github.com/tommasogarofani/ingsw_25-26_group-11/src/backend/domain"
)

type userService interface {
	UserService_Login(ctx context.Context, username, password string) (*domain.Utente, error)
}

type AuthService struct {
	UserService userService
}

func NewAuthService(us userService) *AuthService {
	return &AuthService{UserService: us}
}

func (s *AuthService) AuthService_Login(ctx context.Context, username, password string) (string, error) {
	u, err := s.UserService.UserService_Login(ctx, username, password)
	if err != nil {
		return "", err
	}
	return GenerateToken(u.ID, u.Username) // verifica che domain.Utente abbia un campo ID
}