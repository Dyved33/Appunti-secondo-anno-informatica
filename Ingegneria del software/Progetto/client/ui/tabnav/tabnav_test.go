package tabnav

import (
	"testing"
)

// TestRect_Contains verifica che il metodo Contains restituisca true per punti interni o sui bordi
// del rettangolo e false per punti esterni.
func TestRect_Contains(t *testing.T) {
	t.Parallel()

	r := Rect{X0: 10, Y0: 5, X1: 20, Y1: 15}

	tests := []struct {
		name string
		x, y int
		want bool
	}{
		{"Punto interno", 15, 10, true},
		{"Bordo superiore-sinistro", 10, 5, true},
		{"Bordo inferiore-destro", 20, 15, true},
		{"Bordo superiore", 15, 5, true},
		{"Bordo inferiore", 15, 15, true},
		{"Bordo sinistro", 10, 10, true},
		{"Bordo destro", 20, 10, true},
		{"Fuori a sinistra", 9, 10, false},
		{"Fuori a destra", 21, 10, false},
		{"Fuori in alto", 15, 4, false},
		{"Fuori in basso", 15, 16, false},
		{"Fuori diagonale", 5, 0, false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := r.Contains(tt.x, tt.y)
			if got != tt.want {
				t.Errorf("Rect.Contains(%d, %d) = %v, atteso %v", tt.x, tt.y, got, tt.want)
			}
		})
	}
}

// TestHitAt verifica la funzione HitAt per il rilevamento dell'indice del rettangolo cliccato.
func TestHitAt(t *testing.T) {
	t.Parallel()

	rects := []Rect{
		{X0: 0, Y0: 0, X1: 10, Y1: 5},
		{X0: 15, Y0: 0, X1: 25, Y1: 5},
		{X0: 30, Y0: 0, X1: 40, Y1: 5},
	}

	// Click sul primo rect
	if idx := HitAt(rects, 5, 2); idx != 0 {
		t.Errorf("HitAt(5, 2) atteso 0, ottenuto %d", idx)
	}

	// Click sul secondo rect
	if idx := HitAt(rects, 20, 3); idx != 1 {
		t.Errorf("HitAt(20, 3) atteso 1, ottenuto %d", idx)
	}

	// Click sul terzo rect
	if idx := HitAt(rects, 35, 1); idx != 2 {
		t.Errorf("HitAt(35, 1) atteso 2, ottenuto %d", idx)
	}

	// Click nello spazio vuoto tra rect 0 e 1
	if idx := HitAt(rects, 12, 2); idx != -1 {
		t.Errorf("HitAt(12, 2) atteso -1, ottenuto %d", idx)
	}

	// Click con slice vuota
	if idx := HitAt(nil, 5, 2); idx != -1 {
		t.Errorf("HitAt con slice nil atteso -1, ottenuto %d", idx)
	}
}

// TestPlacePage verifica che PlacePage componga il layout centrato con header, logo e contenuto.
func TestPlacePage(t *testing.T) {
	t.Parallel()

	res := PlacePage(100, 30, "Corpo di test")
	if res == "" {
		t.Fatal("PlacePage ha prodotto una stringa vuota")
	}

	// Con dimensioni nulle o negative viene comunque generata la composizione
	resZero := PlacePage(0, 0, "Corpo zero")
	if resZero == "" {
		t.Fatal("PlacePage(0, 0) ha prodotto una stringa vuota")
	}
}

// TestRects verifica il calcolo delle hitbox cliccabili per vari scenari e dimensioni.
func TestRects(t *testing.T) {
	t.Parallel()

	labels := []string{"Tab1", "Tab2", "Tab3"}
	page := "Tab1  Tab2  Tab3\n\nContenuto pagina"

	t.Run("Dimensioni standard", func(t *testing.T) {
		rects := Rects(120, 40, labels, 0, page)
		if len(rects) != len(labels) {
			t.Fatalf("attesi %d rects, ottenuti %d", len(labels), len(rects))
		}
		for i, r := range rects {
			if r.X1 < r.X0 || r.Y1 < r.Y0 {
				t.Errorf("rect %d con coordinate non valide: %+v", i, r)
			}
		}
	})

	t.Run("Dimensioni <= 0 usano fallback 80x24", func(t *testing.T) {
		rects := Rects(0, 0, labels, 0, page)
		if len(rects) != len(labels) {
			t.Fatalf("attesi %d rects con dimensioni fallback, ottenuti %d", len(labels), len(rects))
		}
	})

	t.Run("Labels vuote restituisce nil", func(t *testing.T) {
		rects := Rects(100, 30, nil, 0, page)
		if rects != nil {
			t.Errorf("atteso nil per labels vuote, ottenuto %+v", rects)
		}
	})

	t.Run("ActiveIdx fuori range viene normalizzato a 0", func(t *testing.T) {
		rects := Rects(120, 40, labels, -5, page)
		if len(rects) != len(labels) {
			t.Fatalf("attesi %d rects con activeIdx negativo, ottenuti %d", len(labels), len(rects))
		}
		rectsHigh := Rects(120, 40, labels, 999, page)
		if len(rectsHigh) != len(labels) {
			t.Fatalf("attesi %d rects con activeIdx elevato, ottenuti %d", len(labels), len(rectsHigh))
		}
	})

	t.Run("Etichetta non presente nel rendering restituisce nil", func(t *testing.T) {
		missingLabels := []string{"NON_ESISTE_1", "NON_ESISTE_2"}
		rects := Rects(100, 30, missingLabels, 0, page)
		if rects != nil {
			t.Errorf("atteso nil per etichette non trovate, ottenuto %+v", rects)
		}
	})
}

// TestFindLabels verifica la ricerca delle etichette sulle righe del rendering.
func TestFindLabels(t *testing.T) {
	t.Parallel()

	lines := []string{
		"  [ Tab1 ]    Tab2    Tab3  ",
		"  ────────    ────    ────  ",
	}
	labels := []string{"Tab1", "Tab2", "Tab3"}

	spans, ok := findLabels(lines, labels)
	if !ok {
		t.Fatal("findLabels doveva trovare tutte le etichette")
	}
	if len(spans) != 3 {
		t.Fatalf("attesi 3 span, ottenuti %d", len(spans))
	}

	// Caso etichetta mancante
	_, okMissing := findLabels(lines, []string{"Tab1", "Inesistente"})
	if okMissing {
		t.Error("findLabels non doveva trovare l'etichetta inesistente")
	}

	// Caso righe insufficienti (< 2 righe)
	_, okFewLines := findLabels([]string{"Solo una riga"}, labels)
	if okFewLines {
		t.Error("findLabels doveva fallire con meno di 2 righe")
	}
}

// TestMax verifica la funzione di supporto max.
func TestMax(t *testing.T) {
	t.Parallel()

	if got := max(10, 20); got != 20 {
		t.Errorf("max(10, 20) = %d, atteso 20", got)
	}
	if got := max(30, 5); got != 30 {
		t.Errorf("max(30, 5) = %d, atteso 30", got)
	}
	if got := max(7, 7); got != 7 {
		t.Errorf("max(7, 7) = %d, atteso 7", got)
	}
}
