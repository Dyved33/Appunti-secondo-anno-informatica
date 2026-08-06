package user

import (
	"context"
	"github.com/tommasogarofani/ingsw_25-26_group-11/server/domain"
	"golang.org/x/crypto/bcrypt"
)

// userRepository è il sottoinsieme di Repository di cui Service ha bisogno.
// *Repository lo implementa già implicitamente (stessi metodi, stesse firme).
type userRepository interface {
	UserRepository_GetByUsername(ctx context.Context, username string) (*domain.Utente, error)
	UserRepository_Update(ctx context.Context, u *domain.Utente) error
	UserRepository_Create(ctx context.Context, u *domain.Utente) error
	UserRepository_Delete(ctx context.Context, u *domain.Utente) error
}

// transactionRunner è il sottoinsieme di domain.TransactionHelper di cui Service ha bisogno.
type transactionRunner interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// la struct UserService serve a richiamare le funzioni del UserRepository e a gestire la logica di business relativa agli utenti
type UserService struct {
	Repo userRepository
	Tx   transactionRunner
}

// UserRepository_GetSecurityQuestion implements [userService].
func (s *UserService) UserRepository_GetSecurityQuestion(ctx context.Context, username string) (string, error) {
	return s.UserService_GetSecurityQuestion(ctx, username)
}

// UserRepository_Login implements [userService].
func (s *UserService) UserRepository_Login(ctx context.Context, username string, password string) (*domain.Utente, error) {
	return s.UserService_Login(ctx, username, password)
}

// UserRepository_VerifySecurityAnswerAndResetPassword implements [userService].
func (s *UserService) UserRepository_VerifySecurityAnswerAndResetPassword(ctx context.Context, username string, answer string, newPassword string) error {
	return s.UserService_VerifySecurityAnswerAndResetPassword(ctx, username, answer, newPassword)
}

// NewUserService crea una nuova istanza del service per gli utenti.
func NewUserService(repo userRepository, tx transactionRunner) *UserService {
	return &UserService{Repo: repo, Tx: tx}
}

// UserService_Login verifica le credenziali di un utente e gestisce il meccanismo di blocco dopo 3 tentativi falliti. Se il login ha successo, ripristina i tentativi falliti, se fallisce ripetutamente, blocca l'account.
func (s *UserService) UserService_Login(ctx context.Context, username, password string) (*domain.Utente, error) {
	var user *domain.Utente
	var loginErr error

	// Esegue tutte le operazioni di login entro una transazione per garantire la consistenza
	err := s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		// Recupera l'utente dal database
		u, err := s.Repo.UserRepository_GetByUsername(txCtx, username)
		if err != nil {
			return err
		}

		// Controlla se l'account è già bloccato o se ha superato il limite di tentativi
		if u.Bloccato || u.TentativiFalliti >= domain.MaxTentativiLogin {
			u.Bloccato = true
			if err := s.Repo.UserRepository_Update(txCtx, u); err != nil {
 				return err
 			}
			return domain.ErrAccountBloccato
		}

		// Verifica la password usando bcrypt
		err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
		if err != nil {
			// Password errata: incrementa il contatore dei tentativi falliti
			u.TentativiFalliti++
			if u.TentativiFalliti >= domain.MaxTentativiLogin {
				u.Bloccato = true
				if err := s.Repo.UserRepository_Update(txCtx, u); err != nil {
 					return err
 				}
				loginErr = domain.ErrAccountBloccato
				return nil
			}
			// Persiste il nuovo contatore sul database
			if err := s.Repo.UserRepository_Update(txCtx, u); err != nil {
				return err
			}
			loginErr = domain.ErrCredenzaliErrate
			return nil
		}

		// Login riuscito: resetta i tentativi falliti e sblocca l'account se era bloccato
		u.TentativiFalliti = 0
		u.Bloccato = false
		err = s.Repo.UserRepository_Update(txCtx, u)
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

// GetSecurityQuestion recupera la domanda di sicurezza associata a un utente.
func (s *UserService) UserService_GetSecurityQuestion(ctx context.Context, username string) (string, error) {
	u, err := s.Repo.UserRepository_GetByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	return u.DomandaSicurezza, nil
}

// VerifySecurityAnswerAndResetPassword verifica la risposta alla domanda di sicurezza e, se corretta, resetta la password.
func (s *UserService) UserService_VerifySecurityAnswerAndResetPassword(ctx context.Context, username, answer, newPassword string) error {
	// Esegue l'intera operazione in una transazione per coerenza
	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {

		u, err := s.Repo.UserRepository_GetByUsername(txCtx, username)
		if err != nil {
			return err
		}

		// Verifica la risposta alla domanda di sicurezza usando bcrypt
		err = bcrypt.CompareHashAndPassword([]byte(u.RispostaHash), []byte(answer))
		if err != nil {
			return domain.ErrRispostaSecretErrata
		}

		// Genera l'hash della nuova password
		newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		// Aggiorna l'utente: nuova password, nessun tentativo fallito, account sbloccato
		u.PasswordHash = string(newHash)
		u.TentativiFalliti = 0
		u.Bloccato = false
		return s.Repo.UserRepository_Update(txCtx, u)
	})
}
