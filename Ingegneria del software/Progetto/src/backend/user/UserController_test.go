package user

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"gialosai/backend/domain"
)

// mockUserServiceForController implementa l'interfaccia userService
// (vedi 00_PATCH_SUGGERITA.md) per testare il Controller senza un vero Service.
type mockUserServiceForController struct {
	mock.Mock
}

func (m *mockUserServiceForController) Login(ctx context.Context, username, password string) (*domain.Utente, error) {
	args := m.Called(ctx, username, password)
	u, _ := args.Get(0).(*domain.Utente)
	return u, args.Error(1)
}

func (m *mockUserServiceForController) GetSecurityQuestion(ctx context.Context, username string) (string, error) {
	args := m.Called(ctx, username)
	return args.String(0), args.Error(1)
}

func (m *mockUserServiceForController) VerifySecurityAnswerAndResetPassword(ctx context.Context, username, answer, newPassword string) error {
	args := m.Called(ctx, username, answer, newPassword)
	return args.Error(0)
}

func TestController_HandleLogin(t *testing.T) {
	tests := []struct {
		name      string
		req       LoginRequestDTO
		mockSetup func(svc *mockUserServiceForController)
		want      LoginResponseDTO
	}{
		{
			name: "login riuscito mappa l'utente nella risposta",
			req:  LoginRequestDTO{Username: "mario", Password: "passwordCorretta"},
			mockSetup: func(svc *mockUserServiceForController) {
				svc.On("Login", mock.Anything, "mario", "passwordCorretta").
					Return(&domain.Utente{Username: "mario", Ruolo: domain.Ruolo("magazziniere")}, nil)
			},
			want: LoginResponseDTO{Success: true, Username: "mario", Ruolo: "magazziniere"},
		},
		{
			name: "login fallito propaga il messaggio d'errore",
			req:  LoginRequestDTO{Username: "mario", Password: "sbagliata"},
			mockSetup: func(svc *mockUserServiceForController) {
				svc.On("Login", mock.Anything, "mario", "sbagliata").
					Return(nil, domain.ErrCredenzaliErrate)
			},
			want: LoginResponseDTO{Success: false, ErrorMessage: domain.ErrCredenzaliErrate.Error()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockUserServiceForController)
			tt.mockSetup(svc)

			c := NewController(svc)
			got := c.HandleLogin(context.Background(), tt.req)

			assert.Equal(t, tt.want, got)
			svc.AssertExpectations(t)
		})
	}
}

func TestController_HandleGetSecurityQuestion(t *testing.T) {
	t.Run("restituisce la domanda di sicurezza", func(t *testing.T) {
		svc := new(mockUserServiceForController)
		svc.On("GetSecurityQuestion", mock.Anything, "mario").
			Return("Il nome del tuo primo animale?", nil)

		c := NewController(svc)
		resp, err := c.HandleGetSecurityQuestion(context.Background(), RecoveryRequestDTO{Username: "mario"})

		require.NoError(t, err)
		assert.Equal(t, RecoveryResponseDTO{DomandaSicurezza: "Il nome del tuo primo animale?"}, resp)
		svc.AssertExpectations(t)
	})

	t.Run("propaga l'errore se l'utente non esiste", func(t *testing.T) {
		svc := new(mockUserServiceForController)
		svc.On("GetSecurityQuestion", mock.Anything, "fantasma").
			Return("", errors.New("utente non trovato"))

		c := NewController(svc)
		_, err := c.HandleGetSecurityQuestion(context.Background(), RecoveryRequestDTO{Username: "fantasma"})

		require.Error(t, err)
		svc.AssertExpectations(t)
	})
}

func TestController_HandleResetPassword(t *testing.T) {
	t.Run("delega correttamente al service", func(t *testing.T) {
		svc := new(mockUserServiceForController)
		svc.On("VerifySecurityAnswerAndResetPassword", mock.Anything, "mario", "rispostaCorretta", "nuovaPassword").
			Return(nil)

		c := NewController(svc)
		err := c.HandleResetPassword(context.Background(), ResetPasswordRequestDTO{
			Username:    "mario",
			Risposta:    "rispostaCorretta",
			NewPassword: "nuovaPassword",
		})

		require.NoError(t, err)
		svc.AssertExpectations(t)
	})

	t.Run("propaga l'errore in caso di risposta errata", func(t *testing.T) {
		svc := new(mockUserServiceForController)
		svc.On("VerifySecurityAnswerAndResetPassword", mock.Anything, "mario", "rispostaSbagliata", "nuovaPassword").
			Return(domain.ErrRispostaSecretErrata)

		c := NewController(svc)
		err := c.HandleResetPassword(context.Background(), ResetPasswordRequestDTO{
			Username:    "mario",
			Risposta:    "rispostaSbagliata",
			NewPassword: "nuovaPassword",
		})

		require.ErrorIs(t, err, domain.ErrRispostaSecretErrata)
		svc.AssertExpectations(t)
	})
}
