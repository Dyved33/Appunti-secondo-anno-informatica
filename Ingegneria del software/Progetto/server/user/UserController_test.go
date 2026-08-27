package user

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/tommasogarofani/ingsw_25-26_group-11/server/domain"
)

// mockUserServiceForController implementa l'interfaccia userService
type mockUserServiceForController struct {
	mock.Mock
}

func (m *mockUserServiceForController) UserService_Login(ctx context.Context, username, password string) (*domain.Utente, error) {
	args := m.Called(ctx, username, password)
	u, _ := args.Get(0).(*domain.Utente)
	return u, args.Error(1)
}

func TestUserController_Login(t *testing.T) {
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
				svc.On("UserService_Login", mock.Anything, "mario", "passwordCorretta").
					Return(&domain.Utente{Username: "mario", Ruolo: domain.Ruolo("magazziniere")}, nil)
			},
			want: LoginResponseDTO{Success: true, Username: "mario", Ruolo: "magazziniere"},
		},
		{
			name: "login fallito propaga il messaggio d'errore",
			req:  LoginRequestDTO{Username: "mario", Password: "sbagliata"},
			mockSetup: func(svc *mockUserServiceForController) {
				svc.On("UserService_Login", mock.Anything, "mario", "sbagliata").
					Return(nil, domain.ErrCredenzaliErrate)
			},
			want: LoginResponseDTO{Success: false, ErrorMessage: domain.ErrCredenzaliErrate.Error()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockUserServiceForController)
			tt.mockSetup(svc)

			c := NewUserController(svc)
			got := c.UserController_Login(context.Background(), tt.req)

			assert.Equal(t, tt.want, got)
			svc.AssertExpectations(t)
		})
	}
}
