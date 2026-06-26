package user

import (
	"context"
	"gialosai/backend/domain"
	"golang.org/x/crypto/bcrypt"
)

// la struct Service serve a richiamare le funzioni del Repository e a gestire la logica di business relativa agli utenti
type Service struct {
	Repo *Repository
	Tx   *domain.TransactionHelper
}

// NewService crea una nuova istanza del service per gli utenti.
func NewService(repo *Repository, tx *domain.TransactionHelper) *Service {
	return &Service{Repo: repo, Tx: tx}
}

// Login verifica le credenziali di un utente e gestisce il meccanismo di blocco dopo 3 tentativi falliti. Se il login ha successo, ripristina i tentativi falliti, se fallisce ripetutamente, blocca l'account.
func (s *Service) Login(ctx context.Context, username, password string) (*domain.Utente, error) {
	var user *domain.Utente
	var loginErr error

	// Esegue tutte le operazioni di login entro una transazione per garantire la consistenza
	err := s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		// Recupera l'utente dal database
		u, err := s.Repo.GetByUsername(txCtx, username)
		if err != nil {
			return err
		}

		// Controlla se l'account è già bloccato o se ha superato il limite di tentativi
		if u.Bloccato || u.TentativiFalliti >= domain.MaxTentativiLogin {
			u.Bloccato = true
			_ = s.Repo.Update(txCtx, u)
			return domain.ErrAccountBloccato
		}

		// Verifica la password usando bcrypt
		err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
		if err != nil {
			// Password errata: incrementa il contatore dei tentativi falliti
			u.TentativiFalliti++
			if u.TentativiFalliti >= domain.MaxTentativiLogin {
				u.Bloccato = true
				_ = s.Repo.Update(txCtx, u)
				loginErr = domain.ErrAccountBloccato
				return nil
			}
			// Persiste il nuovo contatore sul database
			_ = s.Repo.Update(txCtx, u)
			loginErr = domain.ErrCredenzaliErrate
			return nil
		}

		// Login riuscito: resetta i tentativi falliti e sblocca l'account se era bloccato
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

// GetSecurityQuestion recupera la domanda di sicurezza associata a un utente.
func (s *Service) GetSecurityQuestion(ctx context.Context, username string) (string, error) {
	u, err := s.Repo.GetByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	return u.DomandaSicurezza, nil
}

// VerifySecurityAnswerAndResetPassword verifica la risposta alla domanda di sicurezza e, se corretta, resetta la password.
func (s *Service) VerifySecurityAnswerAndResetPassword(ctx context.Context, username, answer, newPassword string) error {
	// Esegue l'intera operazione in una transazione per coerenza
	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {

		u, err := s.Repo.GetByUsername(txCtx, username)
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
		return s.Repo.Update(txCtx, u)
	})
}
