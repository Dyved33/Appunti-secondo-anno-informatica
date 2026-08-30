package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tommasogarofani/ingsw_25-26_group-11/server/domain"
)

type mockUserService struct {
	mock.Mock
}

func (m *mockUserService) UserService_Login(ctx context.Context, username, password string) (*domain.Utente, error) {
	args := m.Called(ctx, username, password)
	u, _ := args.Get(0).(*domain.Utente)
	return u, args.Error(1)
}

type mockSessionRepository struct {
	mock.Mock
}

func (m *mockSessionRepository) SessionRepository_Create(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) error {
	args := m.Called(ctx, userID, tokenHash, expiresAt)
	return args.Error(0)
}

type mockAuthService struct {
	mock.Mock
}

func (m *mockAuthService) AuthService_Login(ctx context.Context, username, password string) (string, string, error) {
	args := m.Called(ctx, username, password)
	return args.String(0), args.String(1), args.Error(2)
}

func TestPasswordHashing(t *testing.T) {
	password := "segretissimo"
	hash, err := HashPassword(password)
	require.NoError(t, err)
	assert.NotEmpty(t, hash)

	err = CheckPassword(hash, password)
	assert.NoError(t, err)

	err = CheckPassword(hash, "sbagliata")
	assert.Error(t, err)
}

func TestJWTTokens(t *testing.T) { // testa sia la generazione che la validazione dei token JWT
	userID := 42
	username := "mario"
	role := "Magazziniere"

	token, err := GenerateToken(userID, username, role)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := ValidateToken(token)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, role, claims.Role)

	_, err = ValidateToken("invalid-token")
	assert.Error(t, err)
}

func TestHashToken(t *testing.T) {
	h1 := HashToken("my-token")
	h2 := HashToken("my-token")
	h3 := HashToken("different-token")

	assert.Equal(t, h1, h2, "same input should produce same hash")
	assert.NotEqual(t, h1, h3, "different input should produce different hash")
	assert.Len(t, h1, 64, "SHA-256 hex digest should be 64 characters")
}

func TestAuthService_Login(t *testing.T) {
	t.Run("login con successo restituisce token e ruolo", func(t *testing.T) {
		us := new(mockUserService)
		us.On("UserService_Login", mock.Anything, "mario", "password").Return(&domain.Utente{ID: 42, Username: "mario", Ruolo: domain.RuoloMagazziniere}, nil)

		svc := NewAuthService(us, nil)
		token, role, err := svc.AuthService_Login(context.Background(), "mario", "password")

		require.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.Equal(t, string(domain.RuoloMagazziniere), role)
		us.AssertExpectations(t)
	})

	t.Run("login fallito propaga errore", func(t *testing.T) {
		us := new(mockUserService)
		us.On("UserService_Login", mock.Anything, "mario", "sbagliata").Return(nil, errors.New("credenziali errate"))

		svc := NewAuthService(us, nil)
		_, _, err := svc.AuthService_Login(context.Background(), "mario", "sbagliata")

		assert.Error(t, err)
		us.AssertExpectations(t)
	})

	t.Run("login con sessione: crea la sessione nel DB", func(t *testing.T) {
		us := new(mockUserService)
		us.On("UserService_Login", mock.Anything, "mario", "password").
			Return(&domain.Utente{ID: 42, Username: "mario", Ruolo: domain.RuoloMagazziniere}, nil)

		sr := new(mockSessionRepository)
		sr.On("SessionRepository_Create", mock.Anything, 42, mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).
			Return(nil)

		svc := NewAuthService(us, sr)
		token, role, err := svc.AuthService_Login(context.Background(), "mario", "password")

		require.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.Equal(t, string(domain.RuoloMagazziniere), role)
		sr.AssertExpectations(t)
	})

	t.Run("login con sessione: errore nella creazione sessione", func(t *testing.T) {
		us := new(mockUserService)
		us.On("UserService_Login", mock.Anything, "mario", "password").
			Return(&domain.Utente{ID: 42, Username: "mario", Ruolo: domain.RuoloContabile}, nil)

		sr := new(mockSessionRepository)
		sr.On("SessionRepository_Create", mock.Anything, 42, mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).
			Return(errors.New("db error"))

		svc := NewAuthService(us, sr)
		_, _, err := svc.AuthService_Login(context.Background(), "mario", "password")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})
}

func TestAuthController_Login(t *testing.T) {
	t.Run("login con successo restituisce DTO con token e ruolo", func(t *testing.T) {
		as := new(mockAuthService)
		as.On("AuthService_Login", mock.Anything, "mario", "password").Return("mock-token", "Magazziniere", nil)

		c := NewAuthController(as)
		res, err := c.AuthController_Login(context.Background(), LoginRequestDTO{Username: "mario", Password: "password"})

		require.NoError(t, err)
		assert.Equal(t, "mock-token", res.Token)
		assert.Equal(t, "Magazziniere", res.Role)
		as.AssertExpectations(t)
	})

	t.Run("login fallito restituisce errore", func(t *testing.T) {
		as := new(mockAuthService)
		as.On("AuthService_Login", mock.Anything, "mario", "sbagliata").Return("", "", errors.New("auth error"))

		c := NewAuthController(as)
		_, err := c.AuthController_Login(context.Background(), LoginRequestDTO{Username: "mario", Password: "sbagliata"})

		assert.Error(t, err)
		as.AssertExpectations(t)
	})
}
