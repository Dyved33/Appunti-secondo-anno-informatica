package style

import (
	"testing"
)

// TestPaginaCalcola verifica la funzione PaginaCalcola per il calcolo dell'offset
// in tabelle paginate, coprendo tutti i casi limite (totale zero, scorrimento verso
// il basso oltre la pagina visibile, scorrimento verso l'alto e permanenza nella finestra).
func TestPaginaCalcola(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		cursore     int
		offset      int
		totale      int
		expectedOff int
	}{
		{
			name:        "Totale zero restituisce offset zero",
			cursore:     5,
			offset:      10,
			totale:      0,
			expectedOff: 0,
		},
		{
			name:        "Cursore dentro la pagina mantiene l'offset invariato",
			cursore:     10,
			offset:      0,
			totale:      100,
			expectedOff: 0,
		},
		{
			name:        "Cursore sull'ultimo elemento della prima pagina esatta",
			cursore:     RighePerPagina - 1, // 19
			offset:      0,
			totale:      100,
			expectedOff: 0,
		},
		{
			name:        "Cursore oltre la pagina visibile verso il basso (avanzamento offset)",
			cursore:     RighePerPagina, // 20
			offset:      0,
			totale:      100,
			expectedOff: 1, // 20 - 20 + 1 = 1
		},
		{
			name:        "Cursore molto avanti verso il basso",
			cursore:     45,
			offset:      10,
			totale:      100,
			expectedOff: 26, // 45 - 20 + 1 = 26
		},
		{
			name:        "Cursore prima dell'offset corrente (scorrimento verso l'alto)",
			cursore:     5,
			offset:      10,
			totale:      100,
			expectedOff: 5,
		},
		{
			name:        "Cursore a zero con offset maggiore di zero",
			cursore:     0,
			offset:      15,
			totale:      100,
			expectedOff: 0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := PaginaCalcola(tt.cursore, tt.offset, tt.totale)
			if got != tt.expectedOff {
				t.Errorf("PaginaCalcola(%d, %d, %d) = %d, atteso %d",
					tt.cursore, tt.offset, tt.totale, got, tt.expectedOff)
			}
		})
	}
}

// TestNewWindowSize verifica il costruttore NewWindowSize e le dimensioni iniziali.
func TestNewWindowSize(t *testing.T) {
	t.Parallel()

	ws := NewWindowSize(80, 24)
	if ReadWidth(ws) != 80 {
		t.Errorf("ReadWidth atteso 80, ottenuto %d", ReadWidth(ws))
	}
	if ReadHeight(ws) != 24 {
		t.Errorf("ReadHeight atteso 24, ottenuto %d", ReadHeight(ws))
	}
}

// TestWindowSize_ChangeSize verifica il metodo ChangeSize per la mutazione delle dimensioni.
func TestWindowSize_ChangeSize(t *testing.T) {
	t.Parallel()

	ws := NewWindowSize(80, 24)
	ws.ChangeSize(120, 60)

	if got := ReadWidth(ws); got != 120 {
		t.Errorf("ReadWidth dopo ChangeSize atteso 120, ottenuto %d", got)
	}
	if got := ReadHeight(ws); got != 60 {
		t.Errorf("ReadHeight dopo ChangeSize atteso 60, ottenuto %d", got)
	}
}

// TestReadWidth_ReadHeight verifica le funzioni getter ReadWidth e ReadHeight.
func TestReadWidth_ReadHeight(t *testing.T) {
	t.Parallel()

	ws := WindowSize{width: 150, height: 75}
	if got := ReadWidth(ws); got != 150 {
		t.Errorf("ReadWidth atteso 150, ottenuto %d", got)
	}
	if got := ReadHeight(ws); got != 75 {
		t.Errorf("ReadHeight atteso 75, ottenuto %d", got)
	}
}

// TestSpacing_SetSpacing verifica l'impostazione dei quattro margini dello Spacing.
func TestSpacing_SetSpacing(t *testing.T) {
	t.Parallel()

	var s Spacing
	s.SetSpacing(10, 20, 30, 40)

	if s.top != 10 || s.bottom != 20 || s.left != 30 || s.right != 40 {
		t.Errorf("SetSpacing fallito: top=%d, bottom=%d, left=%d, right=%d", s.top, s.bottom, s.left, s.right)
	}
}

// TestGetInnerSpace verifica il calcolo dello spazio interno della finestra sottraendo
// margini, bordi e padding su entrambi gli assi orizzontale e verticale.
func TestGetInnerSpace(t *testing.T) {
	t.Parallel()

	ws := NewWindowSize(200, 100)
	ws.Margin.SetSpacing(5, 5, 10, 10)   // orizzontale: 20, verticale: 10
	ws.Border.SetSpacing(1, 1, 2, 2)     // orizzontale: 4, verticale: 2
	ws.Padding.SetSpacing(3, 3, 4, 4)    // orizzontale: 8, verticale: 6

	innerW, innerH := GetInnerSpace(ws)
	expectedW := 200 - 20 - 4 - 8 // 168
	expectedH := 100 - 10 - 2 - 6 // 82

	if innerW != expectedW {
		t.Errorf("GetInnerSpace larghezza attesa %d, ottenuta %d", expectedW, innerW)
	}
	if innerH != expectedH {
		t.Errorf("GetInnerSpace altezza attesa %d, ottenuta %d", expectedH, innerH)
	}
}

// TestStylesDefinitions verifica che tutti gli stili Lipgloss definiti nel package siano validi.
func TestStylesDefinitions(t *testing.T) {
	t.Parallel()

	stili := []struct {
		nome string
		out  string
	}{
		{"Tab", Tab.Render("Test")},
		{"ActiveTab", ActiveTab.Render("Test")},
		{"TabGap", TabGap.Render("Test")},
		{"FinestraStile", FinestraStile.Render("Test")},
		{"FinestraInternaStile", FinestraInternaStile.Render("Test")},
		{"HeaderStile", HeaderStile.Render("Test")},
		{"FreeHeaderStile", FreeHeaderStile.Render("Test")},
		{"FooterStile", FooterStile.Render("Test")},
		{"TitleStile", TitleStile.Render("Test")},
		{"SubTitleStile", SubTitleStile.Render("Test")},
		{"ErrorStile", ErrorStile.Render("Test")},
		{"SuccessStile", SuccessStile.Render("Test")},
		{"InfoStile", InfoStile.Render("Test")},
		{"FocusStile", FocusStile.Render("Test")},
		{"BoxStile", BoxStile.Render("Test")},
		{"LoginBoxStile", LoginBoxStile.Render("Test")},
	}

	for _, s := range stili {
		if s.out == "" {
			t.Errorf("stile %s ha prodotto una stringa vuota", s.nome)
		}
	}

	if len(WindowBorderColors) != 3 {
		t.Errorf("attesi 3 colori in WindowBorderColors, trovati %d", len(WindowBorderColors))
	}
}
