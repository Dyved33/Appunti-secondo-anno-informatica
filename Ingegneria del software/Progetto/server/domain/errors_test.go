package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGiacenzaInsufficienteError_Error(t *testing.T) {
	err := NewGiacenzaInsufficienteError(5, "Vite M6", 3, 10)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Vite M6")
	assert.Contains(t, err.Error(), "id=5")
	assert.Contains(t, err.Error(), "disponibili 3")
	assert.Contains(t, err.Error(), "richiesti 10")
}

func TestGiacenzaInsufficienteError_Unwrap(t *testing.T) {
	err := NewGiacenzaInsufficienteError(5, "Vite M6", 3, 10)
	assert.ErrorIs(t, err, ErrGiacenzaInsufficient)
}

func TestIsGiacenzaInsufficient(t *testing.T) {
	assert.True(t, IsGiacenzaInsufficient(ErrGiacenzaInsufficient))
	assert.True(t, IsGiacenzaInsufficient(NewGiacenzaInsufficienteError(1, "A", 0, 1)))
	assert.False(t, IsGiacenzaInsufficient(ErrProdottoNonTrovato))
	assert.False(t, IsGiacenzaInsufficient(nil))
}

func TestIsNotFound(t *testing.T) {
	assert.True(t, IsNotFound(ErrProdottoNonTrovato))
	assert.False(t, IsNotFound(ErrGiacenzaInsufficient))
	assert.False(t, IsNotFound(nil))
}

func TestIsCategoryNotFound(t *testing.T) {
	assert.True(t, IsCategoryNotFound(ErrCategoriaNonTrovata))
	assert.False(t, IsCategoryNotFound(ErrProdottoNonTrovato))
	assert.False(t, IsCategoryNotFound(nil))
}

func TestIsPermissionDenied(t *testing.T) {
	assert.True(t, IsPermissionDenied(ErrPermessoNegato))
	assert.False(t, IsPermissionDenied(ErrGiacenzaInsufficient))
	assert.False(t, IsPermissionDenied(nil))
}
