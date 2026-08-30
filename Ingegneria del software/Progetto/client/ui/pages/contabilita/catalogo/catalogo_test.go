package catalogo

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
	style "github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/styles"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/tabnav"
)

// helperNewModel crea un Model in stato di base (non in caricamento) con dati di test.
func helperNewModel() Model {
	m := New(nil)
	m.loading = false
	m.prodotti = []api.Prodotto{
		{ID: 1, Codice: "P01", Nome: "Laptop", Categoria: "Elettronica", Prezzo: 1000.0, Attivo: true},
		{ID: 2, Codice: "P02", Nome: "Mouse", Categoria: "Elettronica", Prezzo: 25.0, Attivo: true},
		{ID: 3, Codice: "P03", Nome: "Sedia", Categoria: "Arredamento", Prezzo: 150.0, Attivo: false},
	}
	m.categorie = []string{"Elettronica", "Arredamento"}
	m.categorieID = map[string]int{"Elettronica": 1, "Arredamento": 2}
	m.winH = 40 // Altezza di default sufficiente ad evitare paginazione spinta
	return m
}

// ── Test Inizializzazione e Utility ───────────────────────────────────────

func TestNew(t *testing.T) {
	m := New(nil)

	if !m.loading {
		t.Errorf("expected loading to be true initially, got false")
	}
	if m.mode != modeView {
		t.Errorf("expected modeView, got %v", m.mode)
	}
	if m.filtroCategoria != nessunFiltro {
		t.Errorf("expected filtroCategoria to be %d, got %d", nessunFiltro, m.filtroCategoria)
	}
	if m.categorieID == nil {
		t.Errorf("expected categorieID map to be initialized")
	}
}

func TestTroncaStringa(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"TestoLungo", 5, "Te..."},
		{"TestoLungo", 10, "TestoLungo"},
		{"ABC", 2, "AB"},
		{"Ciao", 4, "Ciao"},
	}

	for _, tt := range tests {
		got := troncaStringa(tt.input, tt.maxLen)
		if got != tt.expected {
			t.Errorf("troncaStringa(%q, %d) = %q; want %q", tt.input, tt.maxLen, got, tt.expected)
		}
	}
}

// ── Test Calcolo Righe e Paginazione ──────────────────────────────────────

func TestRigheVisibili(t *testing.T) {
	t.Run("Fallback per altezza <= 0", func(t *testing.T) {
		m := Model{winH: 0}
		if got := m.righeVisibili(); got != righeFallback {
			t.Errorf("expected fallback %d, got %d", righeFallback, got)
		}
	})

	t.Run("Calcolo standard con altezza finestra", func(t *testing.T) {
		m := Model{winH: 40}
		// m.winH - overheadChrome (21) - 11 = 40 - 32 = 8
		expected := 8
		if got := m.righeVisibili(); got != expected {
			t.Errorf("expected %d righe visibili, got %d", expected, got)
		}
	})

	t.Run("Valore minimo clampato a 1", func(t *testing.T) {
		m := Model{winH: 25} // 25 - 32 = -7 -> deve diventare 1
		if got := m.righeVisibili(); got != 1 {
			t.Errorf("expected minimum 1, got %d", got)
		}
	})
}

func TestScorriEMostraCursore(t *testing.T) {
	m := helperNewModel()
	m.winH = 33 // righeVisibili() = 1

	// Scorrimento verso il basso
	m.scorri(1)
	if m.cursore != 1 || m.offset != 1 {
		t.Errorf("expected cursore=1, offset=1; got cursore=%d, offset=%d", m.cursore, m.offset)
	}

	// Scorrimento oltre il limite massimo
	m.scorri(10)
	maxIndex := len(m.prodotti) - 1
	if m.cursore != maxIndex {
		t.Errorf("expected cursore clamped to %d, got %d", maxIndex, m.cursore)
	}

	// Scorrimento all'inizio (negativo)
	m.scorri(-10)
	if m.cursore != 0 || m.offset != 0 {
		t.Errorf("expected cursore=0, offset=0; got cursore=%d, offset=%d", m.cursore, m.offset)
	}
}

// ── Test Filtri Categoria ──────────────────────────────────────────────────

func TestProdottiFiltratiEEtichetta(t *testing.T) {
	m := helperNewModel()

	// 1. Nessun filtro
	if label := m.etichettaFiltro(); label != "Tutte le categorie" {
		t.Errorf("expected 'Tutte le categorie', got %q", label)
	}
	if len(m.prodottiFiltrati()) != 3 {
		t.Errorf("expected 3 products, got %d", len(m.prodottiFiltrati()))
	}

	// 2. Filtro attivo: Elettronica (indice 0)
	m.filtroCategoria = 0
	if label := m.etichettaFiltro(); label != "Elettronica" {
		t.Errorf("expected 'Elettronica', got %q", label)
	}
	filtrati := m.prodottiFiltrati()
	if len(filtrati) != 2 {
		t.Errorf("expected 2 products for Elettronica, got %d", len(filtrati))
	}

	// 3. Filtro fuori indice
	m.filtroCategoria = 99
	if label := m.etichettaFiltro(); label != "Tutte le categorie" {
		t.Errorf("expected default label on out of bounds index, got %q", label)
	}
}

// ── Test Gestione Eventi Update (Messaggi Asincroni) ──────────────────────

func TestUpdateAsyncMessages(t *testing.T) {
	m := New(nil)

	t.Run("CatalogoLoadedMsg", func(t *testing.T) {
		msg := CatalogoLoadedMsg{
			Prodotti:     []api.Prodotto{{ID: 10, Nome: "Test"}},
			Categorie:    []string{"CatTest"},
			CategorieIDs: map[string]int{"CatTest": 1},
		}
		updatedModel, cmd := m.Update(msg)
		mod := updatedModel.(Model)

		if cmd != nil {
			t.Errorf("expected nil cmd")
		}
		if mod.loading {
			t.Errorf("expected loading=false")
		}
		if len(mod.prodotti) != 1 || mod.prodotti[0].Nome != "Test" {
			t.Errorf("prodotti not correctly set")
		}
	})

	t.Run("ErroreMsg", func(t *testing.T) {
		errTest := errors.New("errore di rete")
		updatedModel, _ := m.Update(ErroreMsg{Err: errTest})
		mod := updatedModel.(Model)

		if mod.loading {
			t.Errorf("expected loading=false")
		}
		if mod.err != errTest {
			t.Errorf("expected error %v, got %v", errTest, mod.err)
		}
	})

	t.Run("prodottoCreatoMsg", func(t *testing.T) {
		m.mode = modeCreateProduct
		updatedModel, cmd := m.Update(prodottoCreatoMsg{Nome: "Nuovo"})
		mod := updatedModel.(Model)

		if mod.mode != modeView {
			t.Errorf("expected mode to reset to modeView, got %v", mod.mode)
		}
		if !strings.Contains(mod.success, "Nuovo") {
			t.Errorf("expected success message containing product name")
		}
		if cmd == nil {
			t.Errorf("expected command to reload catalog")
		}
	})

	t.Run("tea.WindowSizeMsg", func(t *testing.T) {
		updatedModel, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
		mod := updatedModel.(Model)

		if mod.winW != 100 || mod.winH != 50 {
			t.Errorf("expected window dimensions 100x50, got %dx%d", mod.winW, mod.winH)
		}
	})
}

// ── Test Tastiera e Navigazione in Modalità View ──────────────────────────

func TestHandleViewInput(t *testing.T) {
	m := helperNewModel()

	// Navigation Down (j/down)
	updatedModel, _ := m.Update(tea.KeyPressMsg{Code: 'j'})
	mod := updatedModel.(Model)
	if mod.cursore != 1 {
		t.Errorf("expected cursore 1 after 'j', got %d", mod.cursore)
	}

	// Navigation Up (k/up)
	updatedModel, _ = mod.Update(tea.KeyPressMsg{Code: 'k'})
	mod = updatedModel.(Model)
	if mod.cursore != 0 {
		t.Errorf("expected cursore 0 after 'k', got %d", mod.cursore)
	}

	// Tab per ciclare tra le categorie
	updatedModel, _ = mod.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	mod = updatedModel.(Model)
	if mod.filtroCategoria != 0 { // prima categoria: Elettronica
		t.Errorf("expected filtroCategoria 0 after tab, got %d", mod.filtroCategoria)
	}

	// Entrare in Modalità Creazione Prodotto (ctrl+n)
	updatedModel, _ = mod.Update(tea.KeyPressMsg{Code: 'n', Mod: tea.ModCtrl})
	mod = updatedModel.(Model)
	if mod.mode != modeCreateProduct {
		t.Errorf("expected modeCreateProduct after ctrl+n, got %v", mod.mode)
	}

	// Annullare con Esc per tornare a View
	updatedModel, _ = mod.Update(tea.KeyPressMsg{Code: tea.KeyEsc})
	mod = updatedModel.(Model)
	if mod.mode != modeView {
		t.Errorf("expected modeView after esc, got %v", mod.mode)
	}
}

// ── Test Form Creazione Prodotto e Validazione ───────────────────────────

func TestHandleCreateProductInput(t *testing.T) {
	m := helperNewModel()
	m.mode = modeCreateProduct

	t.Run("Spostamento Focus con Tab", func(t *testing.T) {
		updatedModel, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
		mod := updatedModel.(Model)
		if mod.createFocus != 1 {
			t.Errorf("expected focus on field 1, got %d", mod.createFocus)
		}
	})

	t.Run("Errore Validazione Campi Vuoti", func(t *testing.T) {
		// Tenta invio senza compilare i campi obbligatori
		updatedModel, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
		mod := updatedModel.(Model)

		if cmd != nil {
			t.Errorf("expected no cmd on validation failure")
		}
		if mod.err == nil || !strings.Contains(mod.err.Error(), "compilare codice e nome") {
			t.Errorf("expected validation error 'compilare codice e nome', got %v", mod.err)
		}
	})

	t.Run("Cambio Categoria Selezionata tramite freccia Destra/Sinistra", func(t *testing.T) {
		m.createFocus = 7 // Focus sul selettore di categoria
		m.createCategoriaIdx = 0

		updatedModel, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyRight})
		mod := updatedModel.(Model)

		if mod.createCategoriaIdx != 1 {
			t.Errorf("expected categoriaIdx 1 after 'right', got %d", mod.createCategoriaIdx)
		}
	})
}

// ── Test Conferma Eliminazione ────────────────────────────────────────────

func TestConfirmDelete(t *testing.T) {
	m := helperNewModel()

	// Apertura conferma eliminazione prodotto con ctrl+d
	updatedModel, _ := m.Update(tea.KeyPressMsg{Code: 'd', Mod: tea.ModCtrl})
	mod := updatedModel.(Model)

	if mod.mode != modeConfirmDeleteProduct {
		t.Fatalf("expected modeConfirmDeleteProduct, got %v", mod.mode)
	}
	if mod.confirmDeleteID != 1 || mod.confirmDeleteNome != "Laptop" {
		t.Errorf("expected delete target ID=1, Nome='Laptop'; got ID=%d, Nome=%s", mod.confirmDeleteID, mod.confirmDeleteNome)
	}

	// Annullamento con ESC
	updatedModel, _ = mod.Update(tea.KeyPressMsg{Code: tea.KeyEsc})
	mod = updatedModel.(Model)

	if mod.mode != modeView {
		t.Errorf("expected reset to modeView after esc, got %v", mod.mode)
	}
	if mod.confirmDeleteID != 0 {
		t.Errorf("expected reset of confirmDeleteID")
	}
}

// ── Test Rendering View ───────────────────────────────────────────────────

func TestViewRendering(t *testing.T) {
	t.Run("Stato di caricamento", func(t *testing.T) {
		m := New(nil)
		m.loading = true
		v := m.View()

		if !strings.Contains(v.Content, "Caricamento catalogo in corso...") {
			t.Errorf("expected loading message in View")
		}
	})

	t.Run("Stato di Errore", func(t *testing.T) {
		m := New(nil)
		m.loading = false
		m.err = errors.New("Errore Connessione")
		v := m.View()

		if !strings.Contains(v.Content, "Errore Connessione") {
			t.Errorf("expected error message in View output")
		}
	})

	t.Run("Tabella Prodotti Standard", func(t *testing.T) {
		m := helperNewModel()
		v := m.View()

		if !strings.Contains(v.Content, "Laptop") || !strings.Contains(v.Content, "CODICE") {
			t.Errorf("expected table headers and product names in View")
		}
	})

	t.Run("Form Creazione Categoria", func(t *testing.T) {
		m := helperNewModel()
		m.mode = modeCreateCategory
		v := m.View()

		if !strings.Contains(v.Content, "CREA NUOVA CATEGORIA") {
			t.Errorf("expected create category form title in View")
		}
	})
}

// TestCatalogo_FooterSempreVisibile verifica la regressione del "footer invisibile":
// il numero di righe della tabella si adegua all'altezza del terminale e la pagina
// composta non deve superare lo schermo.
func TestCatalogo_FooterSempreVisibile(t *testing.T) {
	m := New(nil)
	var prodotti []api.Prodotto
	for i := 1; i <= 26; i++ {
		prodotti = append(prodotti, api.Prodotto{
			ID: i, Codice: fmt.Sprintf("PRD-%03d", i), Nome: fmt.Sprintf("Prodotto %d", i),
			Categoria: "Ferramenta", Prezzo: 1.5, Attivo: true,
		})
	}
	res, _ := m.Update(CatalogoLoadedMsg{Prodotti: prodotti, Categorie: []string{"Ferramenta"}})
	m = res.(Model)

	for _, winH := range []int{40, 50} {
		upd, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: winH})
		m = upd.(Model)

		view := m.View().Content
		if !strings.Contains(view, "[pgup/pgdown o rotellina] Scorri") {
			t.Errorf("winH=%d: le scorciatoie del footer non sono presenti nella vista", winH)
		}

		tabs := lipgloss.JoinHorizontal(lipgloss.Top,
			style.ActiveTab.Render("ORDINI"), "  ", style.TabGap.Render("CATALOGO"))
		composed := tabnav.PlacePage(120, winH,
			style.FinestraStile.Render(tabs+"\n\n"+view))
		if strings.Count(composed, "\n")+1 > winH {
			t.Errorf("winH=%d: la pagina composta occupa %d righe (footer fuori schermo)",
				winH, strings.Count(composed, "\n")+1)
		}
	}

	// Lo scroll con la rotellina deve poter raggiungere l'ultimo prodotto.
	upd, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = upd.(Model)
	for i := 0; i < 40; i++ {
		upd, _ = m.Update(tea.MouseWheelMsg{Button: tea.MouseWheelDown})
		m = upd.(Model)
	}
	if m.cursore != len(prodotti)-1 {
		t.Errorf("Dopo lo scroll il cursore dovrebbe essere sull'ultimo prodotto (%d), trovato %d",
			len(prodotti)-1, m.cursore)
	}
}
