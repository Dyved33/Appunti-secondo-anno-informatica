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

// UserService espone la logica di business relativa agli utenti e delega le
// operazioni di persistenza al repository.
type UserService struct {
	Repo userRepository
	Tx   transactionRunner
}

// NewUserService crea una nuova istanza del service per gli utenti.
func NewUserService(repo userRepository, tx transactionRunner) *UserService {
	return &UserService{Repo: repo, Tx: tx}
}

// isAdminAccount indica se l'account è quello dell'amministratore principale
// (username riservato) o ha ruolo amministrativo: questi account non vengono
// mai bloccati automaticamente.
func isAdminAccount(u *domain.Utente) bool {
	return u.Username == domain.NomeUtenteAdmin || u.Ruolo.IsAdmin()
}

// UserService_Login verifica le credenziali di un utente e gestisce il meccanismo di blocco dopo 3 tentativi falliti. Se il login ha successo, ripristina i tentativi falliti, se fallisce ripetutamente, blocca l'account. L'account admin non viene mai bloccato.
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

		admin := isAdminAccount(u)

		// Controlla se l'account è già bloccato o se ha superato il limite di tentativi.
		// Gli account amministratore non sono soggetti a blocco.
		if !admin && (u.Bloccato || u.TentativiFalliti >= domain.MaxTentativiLogin) {
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
			if u.TentativiFalliti >= domain.MaxTentativiLogin && !admin {
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
