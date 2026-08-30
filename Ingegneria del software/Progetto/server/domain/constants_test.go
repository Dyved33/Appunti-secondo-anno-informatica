package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRuolo_IsValid(t *testing.T) {
	tests := []struct {
		ruolo   Ruolo
		isValid bool
	}{
		{RuoloMagazziniere, true},
		{RuoloContabile, true},
		{RuoloAdmin, true},
		{"non_esistente", false},
		{"", false},
	}
	for _, tt := range tests {
		t.Run(string(tt.ruolo), func(t *testing.T) {
			assert.Equal(t, tt.isValid, tt.ruolo.IsValid())
		})
	}
}

func TestRuolo_IsAdmin(t *testing.T) {
	assert.True(t, RuoloAdmin.IsAdmin())
	assert.False(t, RuoloContabile.IsAdmin())
	assert.False(t, RuoloMagazziniere.IsAdmin())
}

func TestRuolo_Label(t *testing.T) {
	tests := []struct {
		ruolo Ruolo
		label string
	}{
		{RuoloMagazziniere, "Magazziniere"},
		{RuoloContabile, "Contabile"},
		{RuoloAdmin, "Admin"},
		{"unknown", "Sconosciuto"},
	}
	for _, tt := range tests {
		t.Run(string(tt.ruolo), func(t *testing.T) {
			assert.Equal(t, tt.label, tt.ruolo.Label())
		})
	}
}

func TestBuildMySQLDSN(t *testing.T) {
	dsn := BuildMySQLDSN("root", "secret")
	assert.Contains(t, dsn, "root:secret")
	assert.Contains(t, dsn, "127.0.0.1:3306")
	assert.Contains(t, dsn, "gialosai")
}

func TestGiacenza_SottoSoglia(t *testing.T) {
	t.Run("con prodotto sotto soglia", func(t *testing.T) {
		g := &Giacenza{Quantita: 5, Prodotto: &Prodotto{SogliaMin: 10}}
		assert.True(t, g.SottoSoglia())
	})

	t.Run("con prodotto sopra soglia", func(t *testing.T) {
		g := &Giacenza{Quantita: 20, Prodotto: &Prodotto{SogliaMin: 10}}
		assert.False(t, g.SottoSoglia())
	})

	t.Run("con prodotto uguale a soglia", func(t *testing.T) {
		g := &Giacenza{Quantita: 10, Prodotto: &Prodotto{SogliaMin: 10}}
		assert.False(t, g.SottoSoglia())
	})

	t.Run("senza prodotto associato", func(t *testing.T) {
		g := &Giacenza{Quantita: 0, Prodotto: nil}
		assert.False(t, g.SottoSoglia())
	})
}
