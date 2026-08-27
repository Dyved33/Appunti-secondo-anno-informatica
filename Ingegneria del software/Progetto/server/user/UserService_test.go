package user

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/tommasogarofani/ingsw_25-26_group-11/server/domain"
)

// mockUserRepository implementa l'interfaccia userRepository (vedi 00_PATCH_SUGGERITA.md)
// usando testify/mock, così possiamo controllare in ogni test cosa restituisce
// senza toccare un vero database.
type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) UserRepository_GetByUsername(ctx context.Context, username string) (*domain.Utente, error) {
	args := m.Called(ctx, username)
	u, _ := args.Get(0).(*domain.Utente)
	return u, args.Error(1)
}

func (m *mockUserRepository) UserRepository_Update(ctx context.Context, u *domain.Utente) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *mockUserRepository) UserRepository_Create(ctx context.Context, u *domain.Utente) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *mockUserRepository) UserRepository_Delete(ctx context.Context, u *domain.Utente) error {
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

func TestUserService_Login(t *testing.T) {
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
				repo.On("UserRepository_GetByUsername", mock.Anything, "mario").Return(u, nil)
				repo.On("UserRepository_Update", mock.Anything, mock.MatchedBy(func(u *domain.Utente) bool {
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
				repo.On("UserRepository_GetByUsername", mock.Anything, "mario").Return(u, nil)
				repo.On("UserRepository_Update", mock.Anything, mock.MatchedBy(func(u *domain.Utente) bool {
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
				repo.On("UserRepository_GetByUsername", mock.Anything, "mario").Return(u, nil)
				repo.On("UserRepository_Update", mock.Anything, mock.MatchedBy(func(u *domain.Utente) bool {
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
				repo.On("UserRepository_GetByUsername", mock.Anything, "mario").Return(u, nil)
				repo.On("UserRepository_Update", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: domain.ErrAccountBloccato,
			wantNil: true,
		},
		{
			name:     "utente inesistente propaga l'errore del repository",
			username: "fantasma",
			password: "qualsiasi",
			mockSetup: func(repo *mockUserRepository) {
				repo.On("UserRepository_GetByUsername", mock.Anything, "fantasma").Return(nil, domain.ErrCredenzaliErrate)
			},
			wantErr: domain.ErrCredenzaliErrate,
			wantNil: true,
		},
		{
			name:     "admin con password errata oltre la soglia non viene bloccato",
			username: "admin",
			password: "passwordSbagliata",
			mockSetup: func(repo *mockUserRepository) {
				u := &domain.Utente{ID: 1, Username: "admin", PasswordHash: correctHash, Ruolo: domain.RuoloAdmin, TentativiFalliti: domain.MaxTentativiLogin - 1, Bloccato: false}
				repo.On("UserRepository_GetByUsername", mock.Anything, "admin").Return(u, nil)
				repo.On("UserRepository_Update", mock.Anything, mock.MatchedBy(func(u *domain.Utente) bool {
					return !u.Bloccato
				})).Return(nil)
			},
			wantErr: domain.ErrCredenzaliErrate,
			wantNil: true,
		},
		{
			name:     "admin già bloccato può comunque effettuare il login",
			username: "admin",
			password: "passwordCorretta",
			mockSetup: func(repo *mockUserRepository) {
				u := &domain.Utente{ID: 1, Username: "admin", PasswordHash: correctHash, Ruolo: domain.RuoloAdmin, Bloccato: true, TentativiFalliti: 5}
				repo.On("UserRepository_GetByUsername", mock.Anything, "admin").Return(u, nil)
				repo.On("UserRepository_Update", mock.Anything, mock.MatchedBy(func(u *domain.Utente) bool {
					return u.TentativiFalliti == 0 && !u.Bloccato
				})).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mockUserRepository)
			tt.mockSetup(repo)

			svc := NewUserService(repo, fakeTxRunner{})
			got, err := svc.UserService_Login(context.Background(), tt.username, tt.password)

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
