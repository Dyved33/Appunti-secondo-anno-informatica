package user

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"gialosai/backend/domain"
)

// mockUserRepository implementa l'interfaccia userRepository (vedi 00_PATCH_SUGGERITA.md)
// usando testify/mock, così possiamo controllare in ogni test cosa restituisce
// senza toccare un vero database.
type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) GetByUsername(ctx context.Context, username string) (*domain.Utente, error) {
	args := m.Called(ctx, username)
	u, _ := args.Get(0).(*domain.Utente)
	return u, args.Error(1)
}

func (m *mockUserRepository) Update(ctx context.Context, u *domain.Utente) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *mockUserRepository) Create(ctx context.Context, u *domain.Utente) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *mockUserRepository) Delete(ctx context.Context, u *domain.Utente) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

// fakeTxRunner esegue semplicemente la funzione passata, senza apertura di una
// vera transazione: nei test di Service non ci interessa testare la gestione
// delle transazioni (quella appartiene ai test di db.TransactionHelper), solo
// la logica di business.
type fakeTxRunner struct{}

func (fakeTxRunner) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func hashPassword(t *testing.T, pw string) string {
	t.Helper()
	h, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	require.NoError(t, err)
	return string(h)
}

func TestService_Login(t *testing.T) {
	correctHash := hashPassword(t, "passwordCorretta")

	tests := []struct {
		name      string
		username  string
		password  string
		mockSetup func(repo *mockUserRepository)
		wantErr   error
		wantNil   bool
	}{
		{
			name:     "login riuscito resetta i tentativi falliti",
			username: "mario",
			password: "passwordCorretta",
			mockSetup: func(repo *mockUserRepository) {
				u := &domain.Utente{ID: 1, Username: "mario", PasswordHash: correctHash, TentativiFalliti: 2, Bloccato: false}
				repo.On("GetByUsername", mock.Anything, "mario").Return(u, nil)
				repo.On("Update", mock.Anything, mock.MatchedBy(func(u *domain.Utente) bool {
					return u.TentativiFalliti == 0 && !u.Bloccato
				})).Return(nil)
			},
		},
		{
			name:     "password errata sotto la soglia incrementa il contatore",
			username: "mario",
			password: "passwordSbagliata",
			mockSetup: func(repo *mockUserRepository) {
				u := &domain.Utente{ID: 1, Username: "mario", PasswordHash: correctHash, TentativiFalliti: 0, Bloccato: false}
				repo.On("GetByUsername", mock.Anything, "mario").Return(u, nil)
				repo.On("Update", mock.Anything, mock.MatchedBy(func(u *domain.Utente) bool {
					return u.TentativiFalliti == 1 && !u.Bloccato
				})).Return(nil)
			},
			wantErr: domain.ErrCredenzaliErrate,
			wantNil: true,
		},
		{
			name:     "password errata al limite blocca l'account",
			username: "mario",
			password: "passwordSbagliata",
			mockSetup: func(repo *mockUserRepository) {
				u := &domain.Utente{ID: 1, Username: "mario", PasswordHash: correctHash, TentativiFalliti: domain.MaxTentativiLogin - 1, Bloccato: false}
				repo.On("GetByUsername", mock.Anything, "mario").Return(u, nil)
				repo.On("Update", mock.Anything, mock.MatchedBy(func(u *domain.Utente) bool {
					return u.Bloccato
				})).Return(nil)
			},
			wantErr: domain.ErrAccountBloccato,
			wantNil: true,
		},
		{
			name:     "account già bloccato viene rifiutato senza controllare la password",
			username: "mario",
			password: "qualsiasi",
			mockSetup: func(repo *mockUserRepository) {
				u := &domain.Utente{ID: 1, Username: "mario", PasswordHash: correctHash, Bloccato: true}
				repo.On("GetByUsername", mock.Anything, "mario").Return(u, nil)
				repo.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: domain.ErrAccountBloccato,
			wantNil: true,
		},
		{
			name:     "utente inesistente propaga l'errore del repository",
			username: "fantasma",
			password: "qualsiasi",
			mockSetup: func(repo *mockUserRepository) {
				repo.On("GetByUsername", mock.Anything, "fantasma").Return(nil, domain.ErrCredenzaliErrate)
			},
			wantErr: domain.ErrCredenzaliErrate,
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mockUserRepository)
			tt.mockSetup(repo)

			svc := NewService(repo, fakeTxRunner{})
			got, err := svc.Login(context.Background(), tt.username, tt.password)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
			if tt.wantNil {
				assert.Nil(t, got)
			} else {
				require.NotNil(t, got)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestService_GetSecurityQuestion(t *testing.T) {
	t.Run("restituisce la domanda dell'utente", func(t *testing.T) {
		repo := new(mockUserRepository)
		repo.On("GetByUsername", mock.Anything, "mario").
			Return(&domain.Utente{Username: "mario", DomandaSicurezza: "Il nome del tuo primo animale?"}, nil)

		svc := NewService(repo, fakeTxRunner{})
		q, err := svc.GetSecurityQuestion(context.Background(), "mario")

		require.NoError(t, err)
		assert.Equal(t, "Il nome del tuo primo animale?", q)
		repo.AssertExpectations(t)
	})

	t.Run("propaga l'errore se l'utente non esiste", func(t *testing.T) {
		repo := new(mockUserRepository)
		repo.On("GetByUsername", mock.Anything, "fantasma").
			Return(nil, domain.ErrCredenzaliErrate)

		svc := NewService(repo, fakeTxRunner{})
		_, err := svc.GetSecurityQuestion(context.Background(), "fantasma")

		require.ErrorIs(t, err, domain.ErrCredenzaliErrate)
		repo.AssertExpectations(t)
	})
}

func TestService_VerifySecurityAnswerAndResetPassword(t *testing.T) {
	answerHash := hashPassword(t, "rispostaCorretta")

	t.Run("risposta corretta resetta la password e sblocca l'account", func(t *testing.T) {
		repo := new(mockUserRepository)
		u := &domain.Utente{ID: 1, Username: "mario", RispostaHash: answerHash, TentativiFalliti: 2, Bloccato: true}
		repo.On("GetByUsername", mock.Anything, "mario").Return(u, nil)
		repo.On("Update", mock.Anything, mock.MatchedBy(func(u *domain.Utente) bool {
			return u.TentativiFalliti == 0 && !u.Bloccato && u.PasswordHash != ""
		})).Return(nil)

		svc := NewService(repo, fakeTxRunner{})
		err := svc.VerifySecurityAnswerAndResetPassword(context.Background(), "mario", "rispostaCorretta", "nuovaPassword")

		require.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("risposta errata restituisce ErrRispostaSecretErrata e non aggiorna nulla", func(t *testing.T) {
		repo := new(mockUserRepository)
		u := &domain.Utente{ID: 1, Username: "mario", RispostaHash: answerHash}
		repo.On("GetByUsername", mock.Anything, "mario").Return(u, nil)

		svc := NewService(repo, fakeTxRunner{})
		err := svc.VerifySecurityAnswerAndResetPassword(context.Background(), "mario", "rispostaSbagliata", "nuovaPassword")

		require.ErrorIs(t, err, domain.ErrRispostaSecretErrata)
		repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("utente inesistente propaga l'errore del repository", func(t *testing.T) {
		repo := new(mockUserRepository)
		repo.On("GetByUsername", mock.Anything, "fantasma").Return(nil, errors.New("not found"))

		svc := NewService(repo, fakeTxRunner{})
		err := svc.VerifySecurityAnswerAndResetPassword(context.Background(), "fantasma", "risposta", "nuovaPassword")

		require.Error(t, err)
		repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})
}
