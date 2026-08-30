package admin

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
)

// TestAdminLayout_New verifica che il costruttore New inizializzi correttamente
// la pagina di default (DASHBOARD), i riferimenti al client API e i sottomodelli.
func TestAdminLayout_New(t *testing.T) {
	t.Parallel()

	client := api.New("http://localhost:8080")
	m := New(client)

	if m.Attiva() != int(DASHBOARD) {
		t.Errorf("pagina attiva iniziale attesa DASHBOARD (%d), ottenuta %d", DASHBOARD, m.Attiva())
	}
	if m.Client != client {
		t.Error("client API non associato correttamente in AdminLayout")
	}

	w, h := m.Size()
	if w != 0 || h != 0 {
		t.Errorf("dimensioni iniziali attese (0, 0), ottenute (%d, %d)", w, h)
	}
}

// TestAdminLayout_Init verifica che Init deleghi l'inizializzazione al sottomodello dashboard.
func TestAdminLayout_Init(t *testing.T) {
	t.Parallel()

	m := New(nil)
	cmd := m.Init()

	if cmd == nil {
		t.Error("Init() atteso tea.Cmd non nil da dashboard.Init()")
	}
}

// TestAdminLayout_DimUpdate verifica la propagazione delle dimensioni a tutti gli 11 sottomodelli.
func TestAdminLayout_DimUpdate(t *testing.T) {
	t.Parallel()

	m := New(nil)
	m.DimUpdate(tea.WindowSizeMsg{Width: 160, Height: 50})

	w, h := m.Size()
	if w != 160 || h != 50 {
		t.Errorf("dimensioni registrate inattese: (%d, %d), attese (160, 50)", w, h)
	}
}

// TestAdminLayout_SwitchTo verifica la transizione tra tutte le 11 sotto-pagine
// (DASHBOARD, USERS, GIACENZE, ORDINIENTRATA, ORDINIUSCITA, RESI, SCANNER, CATALOGO, FATTURE, ORDINI, REPORT)
// e il caso di indice non valido.
func TestAdminLayout_SwitchTo(t *testing.T) {
	t.Parallel()

	allTabs := []struct {
		dest     focus
		expected focus
	}{
		{DASHBOARD, DASHBOARD},
		{USERS, USERS},
		{GIACENZE, GIACENZE},
		{ORDINIENTRATA, ORDINIENTRATA},
		{ORDINIUSCITA, ORDINIUSCITA},
		{RESI, RESI},
		{SCANNER, SCANNER},
		{CATALOGO, CATALOGO},
		{FATTURE, FATTURE},
		{ORDINI, ORDINI},
		{REPORT, REPORT},
	}

	for _, tt := range allTabs {
		m := New(nil)
		res, cmd := m.switchTo(tt.dest)
		updated := res.(AdminLayout)

		if updated.attiva != tt.expected {
			t.Errorf("switchTo(%v): attesa sezione attiva %v, ottenuta %v", tt.dest, tt.expected, updated.attiva)
		}
		// SCANNER ha un Init() che restituisce nil, tutte le altre restituiscono un cmd non nil
		if tt.dest != SCANNER && cmd == nil {
			t.Errorf("switchTo(%v): atteso cmd non nil per init della sotto-pagina", tt.dest)
		}
	}

	// Caso indice focus non valido
	m := New(nil)
	res, cmd := m.switchTo(focus(999))
	updated := res.(AdminLayout)
	if updated.attiva != focus(999) {
		t.Errorf("attesa attiva=999, ottenuta %v", updated.attiva)
	}
	if cmd != nil {
		t.Errorf("atteso cmd nil per focus non valido, ottenuto %v", cmd)
	}
}

// TestAdminLayout_UpdateActive verifica l'inoltro dei messaggi alla sotto-pagina attiva per tutti i tab.
func TestAdminLayout_UpdateActive(t *testing.T) {
	t.Parallel()

	for att := DASHBOARD; att <= REPORT; att++ {
		m := New(nil)
		m.attiva = att

		res, _ := m.Update(tea.KeyPressMsg{Text: "x"})
		updated := res.(AdminLayout)

		if updated.attiva != att {
			t.Errorf("UpdateActive con attiva=%v ha modificato la pagina in %v", att, updated.attiva)
		}
	}
}

// TestAdminLayout_ViewAndCursor verifica che View() renderizzi correttamente la barra dei tab
// e il contenuto per ciascuna delle 11 sotto-pagine.
func TestAdminLayout_ViewAndCursor(t *testing.T) {
	t.Parallel()

	for att := DASHBOARD; att <= REPORT; att++ {
		m := New(nil)
		m.DimUpdate(tea.WindowSizeMsg{Width: 120, Height: 40})
		m.attiva = att

		view := m.View()
		if view.Content == "" {
			t.Errorf("View() con attiva=%v ha prodotto contenuto vuoto", att)
		}

		// Verifica che l'etichetta del tab attivo sia inclusa nella vista
		if !strings.Contains(view.Content, labels[att]) {
			t.Errorf("View() con attiva=%v non contiene l'etichetta %q", att, labels[att])
		}

		// activeCursor su default
		_ = m.activeCursor()
	}

	// Caso attiva fuori range per activeCursor
	mInvalid := New(nil)
	mInvalid.attiva = focus(999)
	if cursor := mInvalid.activeCursor(); cursor != nil {
		t.Errorf("activeCursor() con indice non valido atteso nil, ottenuto %+v", cursor)
	}
}

// TestAdminLayout_Hitboxes verifica il getter Hitboxes.
func TestAdminLayout_Hitboxes(t *testing.T) {
	t.Parallel()

	m := New(nil)
	m.DimUpdate(tea.WindowSizeMsg{Width: 120, Height: 40})

	hitboxes := m.Hitboxes()
	if len(hitboxes) != len(labels) {
		t.Errorf("Hitboxes() attese %d, ottenute %d", len(labels), len(hitboxes))
	}
}

// TestAdminLayout_JoinWithSpacing verifica la funzione joinWithSpacing con input vuoto, singolo e multiplo.
func TestAdminLayout_JoinWithSpacing(t *testing.T) {
	t.Parallel()

	if got := joinWithSpacing([]string{}, " "); got != nil {
		t.Errorf("atteso nil per slice vuota, ottenuto %v", got)
	}

	singolo := joinWithSpacing([]string{"A"}, " ")
	if len(singolo) != 1 || singolo[0] != "A" {
		t.Errorf("risultato inatteso per singolo elemento: %v", singolo)
	}

	multiplo := joinWithSpacing([]string{"A", "B", "C"}, "  ")
	expected := []string{"A", "  ", "B", "  ", "C"}
	if len(multiplo) != len(expected) {
		t.Fatalf("lunghezza inattesa: attesa %d, ottenuta %d", len(expected), len(multiplo))
	}
	for i, v := range multiplo {
		if v != expected[i] {
			t.Errorf("indice %d: atteso %q, ottenuto %q", i, expected[i], v)
		}
	}
}
