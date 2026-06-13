package Service

import (
	"context"
	"errors"
	"gialosai/backend/domain"
	"gialosai/backend/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repo *Repository.UserRepository
	Tx   *Repository.TransactionHelper
}

func NewUserService(repo *Repository.UserRepository, tx *Repository.TransactionHelper) *UserService {
	return &UserService{Repo: repo, Tx: tx}
}

// Login verifies the user credentials. Blocks user after 3 failed attempts.
func (s *UserService) Login(ctx context.Context, username, password string) (*domain.Utente, error) {
	var user *domain.Utente
	var loginErr error

	err := s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		u, err := s.Repo.GetByUsername(txCtx, username)
		if err != nil {
			return err
		}

		if u.Bloccato || u.TentativiFalliti >= domain.MaxTentativiLogin {
			u.Bloccato = true
			_ = s.Repo.Update(txCtx, u)
			return domain.ErrAccountBloccato
		}

		err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
		if err != nil {
			u.TentativiFalliti++
			if u.TentativiFalliti >= domain.MaxTentativiLogin {
				u.Bloccato = true
				_ = s.Repo.Update(txCtx, u)
				loginErr = domain.ErrAccountBloccato
				return nil
			}
			_ = s.Repo.Update(txCtx, u)
			loginErr = domain.ErrCredenzaliErrate
			return nil
		}

		// Success: reset attempts
		u.TentativiFalliti = 0
		u.Bloccato = false
		err = s.Repo.Update(txCtx, u)
		if err != nil {
			return err
		}

		user = u
		return nil
	})

	if err != nil {
		return nil, err
	}
	if loginErr != nil {
		return nil, loginErr
	}
	return user, nil
}

// GetSecurityQuestion retrieves the security question for password recovery
func (s *UserService) GetSecurityQuestion(ctx context.Context, username string) (string, error) {
	u, err := s.Repo.GetByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	return u.DomandaSicurezza, nil
}

// VerifySecurityAnswerAndResetPassword validates the security answer and resets the password
func (s *UserService) VerifySecurityAnswerAndResetPassword(ctx context.Context, username, answer, newPassword string) error {
	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		u, err := s.Repo.GetByUsername(txCtx, username)
		if err != nil {
			return err
		}

		err = bcrypt.CompareHashAndPassword([]byte(u.RispostaHash), []byte(answer))
		if err != nil {
			return domain.ErrRispostaSecretErrata
		}

		newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		u.PasswordHash = string(newHash)
		u.TentativiFalliti = 0
		u.Bloccato = false
		return s.Repo.Update(txCtx, u)
	})
}