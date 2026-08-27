// Package catalogo raccoglie i test d'unità della sotto-pagina Catalogo.
package catalogo

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"

	tea "charm.land/bubbletea/v2"
)

// TestNew_StatoIniziale verifica che il catalogo parta in caricamento e senza
// alcun filtro categoria applicato.
func TestNew_StatoIniziale(t *testing.T) {
	m := New(nil)

	if !m.loading {
		t.Error("atteso loading=true all'avvio")
	}
	if m.filtroCategoria != nessunFiltro {
		t.Errorf("atteso nessun filtro (%d) all'avvio, ottenuto %d", nessunFiltro, m.filtroCategoria)
	}
}

// TestUpdate_CatalogoLoadedMsg verifica che il caricamento popoli prodotti e categorie.
func TestUpdate_CatalogoLoadedMsg(t *testing.T) {
	m := caricaCatalogoDiProva(t)

	if m.loading {
		t.Error("atteso loading=false dopo il caricamento")
	}
	if len(m.prodottiFiltrati()) != 2 {
		t.Errorf("senza filtro attesi tutti i 2 prodotti, ottenuti %d", len(m.prodottiFiltrati()))
	}
}

// TestUpdate_ErroreMsg verifica che un errore async venga registrato correttamente.
func TestUpdate_ErroreMsg(t *testing.T) {
	m := New(nil)
	res, _ := m.Update(ErroreMsg{Err: errors.New("catalogo non disponibile")})
	m = res.(Model)

	if m.err == nil || m.loading {
		t.Error("atteso err impostato e loading=false")
	}
	if v := m.View(); !strings.Contains(v.Content, "catalogo non disponibile") {
		t.Errorf("attesa la view di errore contenente il messaggio, ottenuta %q", v.Content)
	}
}

// TestUpdate_TabCambiaFiltro verifica che "tab" scorra ciclicamente tra le
// categorie disponibili e torni infine a "tutte le categorie".
func TestUpdate_TabCambiaFiltro(t *testing.T) {
	m := caricaCatalogoDiProva(t)

	res, _ := m.Update(tea.KeyPressMsg{Text: "tab"})
	m = res.(Model)
	if m.filtroCategoria != 0 {
		t.Errorf("atteso filtro sulla prima categoria (indice 0), ottenuto %d", m.filtroCategoria)
	}
	if prodotti := m.prodottiFiltrati(); len(prodotti) != 1 || prodotti[0].Codice != "P1" {
		t.Errorf("filtro 'Ferramenta' deve restituire solo P1, ottenuti %+v", prodotti)
	}

	res, _ = m.Update(tea.KeyPressMsg{Text: "tab"})
	m = res.(Model)
	if m.filtroCategoria != 1 {
		t.Errorf("atteso filtro sulla seconda categoria (indice 1), ottenuto %d", m.filtroCategoria)
	}

	res, _ = m.Update(tea.KeyPressMsg{Text: "tab"})
	m = res.(Model)
	if m.filtroCategoria != nessunFiltro {
		t.Errorf("atteso ritorno a 'tutte le categorie' (%d) dopo l'ultimo tab, ottenuto %d", nessunFiltro, m.filtroCategoria)
	}
}

// TestUpdate_TabSenzaCategorie verifica che "tab" su un catalogo privo di
// categorie non porti il filtro fuori dai limiti della slice.
func TestUpdate_TabSenzaCategorie(t *testing.T) {
	m := New(nil)
	res, _ := m.Update(CatalogoLoadedMsg{
		Prodotti:  []api.Prodotto{{ID: 1, Codice: "P1", Nome: "Vite M6"}},
		Categorie: nil,
	})
	m = res.(Model)

	res, _ = m.Update(tea.KeyPressMsg{Text: "tab"})
	m = res.(Model)

	if m.filtroCategoria != nessunFiltro {
		t.Errorf("senza categorie il filtro deve restare a %d, ottenuto %d", nessunFiltro, m.filtroCategoria)
	}
	if len(m.prodottiFiltrati()) != 1 {
		t.Errorf("atteso l'intero catalogo senza categorie, ottenuti %d prodotti", len(m.prodottiFiltrati()))
	}
}

// TestUpdate_FiltroNonValidoDopoRicaricamento verifica che un filtro rimasto
// puntato oltre la fine delle nuove categorie non faccia panic e venga trattato
// come "tutte le categorie".
func TestUpdate_FiltroNonValidoDopoRicaricamento(t *testing.T) {
	m := caricaCatalogoDiProva(t)

	// Filtro sulla seconda categoria ("Utensili", indice 1).
	for i := 0; i < 2; i++ {
		res, _ := m.Update(tea.KeyPressMsg{Text: "tab"})
		m = res.(Model)
	}
	if m.filtroCategoria != 1 {
		t.Fatalf("precondizione fallita: filtro atteso a 1, ottenuto %d", m.filtroCategoria)
	}

	// Il nuovo catalogo ha una sola categoria: l'indice 1 non esiste più.
	res, _ := m.Update(CatalogoLoadedMsg{
		Prodotti:  []api.Prodotto{{ID: 1, Codice: "P1", Nome: "Vite M6", Categoria: "Ferramenta"}},
		Categorie: []string{"Ferramenta"},
	})
	m = res.(Model)

	if len(m.prodottiFiltrati()) != 1 {
		t.Errorf("atteso l'intero catalogo con un filtro non più valido, ottenuti %d prodotti", len(m.prodottiFiltrati()))
	}
	if v := m.View(); !strings.Contains(v.Content, "Tutte le categorie") {
		t.Errorf("atteso il ritorno all'etichetta 'Tutte le categorie', ottenuta %q", v.Content)
	}
}

// TestUpdate_CursoreResettatoAlCambioFiltro verifica che cambiare filtro
// riporti il cursore a 0, per evitare indici fuori range sulla lista filtrata.
func TestUpdate_CursoreResettatoAlCambioFiltro(t *testing.T) {
	m := caricaCatalogoDiProva(t)

	res, _ := m.Update(tea.KeyPressMsg{Text: "down"})
	m = res.(Model)
	if m.cursore != 1 {
		t.Fatalf("precondizione fallita: cursore atteso a 1, ottenuto %d", m.cursore)
	}

	res, _ = m.Update(tea.KeyPressMsg{Text: "tab"})
	m = res.(Model)
	if m.cursore != 0 {
		t.Errorf("atteso reset del cursore a 0 dopo il cambio filtro, ottenuto %d", m.cursore)
	}
}

// TestUpdate_CursoreNonEsceDaiLimiti verifica che il cursore non vada mai sotto
// 0 né oltre l'ultimo prodotto della lista attualmente visibile.
func TestUpdate_CursoreNonEsceDaiLimiti(t *testing.T) {
	m := caricaCatalogoDiProva(t)

	res, _ := m.Update(tea.KeyPressMsg{Text: "up"})
	m = res.(Model)
	if m.cursore != 0 {
		t.Errorf("il cursore non deve andare sotto 0, ottenuto %d", m.cursore)
	}

	for i := 0; i < 5; i++ {
		res, _ = m.Update(tea.KeyPressMsg{Text: "down"})
		m = res.(Model)
	}
	if m.cursore != 1 {
		t.Errorf("con 2 prodotti il cursore deve fermarsi a 1, ottenuto %d", m.cursore)
	}
}

// TestView_FiltroApplicato verifica che la view mostri l'etichetta della
// categoria selezionata e nasconda i prodotti delle altre categorie.
func TestView_FiltroApplicato(t *testing.T) {
	m := caricaCatalogoDiProva(t)

	res, _ := m.Update(tea.KeyPressMsg{Text: "tab"})
	m = res.(Model)

	v := m.View()
	if !strings.Contains(v.Content, "Ferramenta") {
		t.Errorf("attesa l'etichetta della categoria filtrata nella view, ottenuta %q", v.Content)
	}
	if strings.Contains(v.Content, "Cacciavite") {
		t.Errorf("il prodotto di un'altra categoria non doveva comparire nella view, ottenuta %q", v.Content)
	}
}

// TestView_MostraPrezzoEStato verifica che la lista riporti il prezzo di listino
// e lo stato del prodotto, e non solo codice e nome.
func TestView_MostraPrezzoEStato(t *testing.T) {
	m := New(nil)
	res, _ := m.Update(CatalogoLoadedMsg{
		Prodotti: []api.Prodotto{
			{ID: 1, Codice: "P1", Nome: "Vite M6", Categoria: "Ferramenta", Prezzo: 0.10, Attivo: true},
			{ID: 2, Codice: "P2", Nome: "Cacciavite", Categoria: "Utensili", Prezzo: 4.50, Attivo: false},
		},
		Categorie: []string{"Ferramenta", "Utensili"},
	})
	m = res.(Model)

	v := m.View()
	if !strings.Contains(v.Content, "€ 0.10") {
		t.Errorf("atteso il prezzo di listino nella view, ottenuta %q", v.Content)
	}
	if !strings.Contains(v.Content, "non attivo") {
		t.Errorf("atteso lo stato del prodotto disattivato nella view, ottenuta %q", v.Content)
	}
}

// TestLoadCatalogoCmd_Successo esegue il Cmd di Init() contro un backend di
// prova e verifica che prodotti e categorie arrivino uniti nel messaggio,
// con la categoria risolta dal proprio identificativo.
func TestLoadCatalogoCmd_Successo(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/api/v1/products":
			json.NewEncoder(w).Encode([]map[string]any{{
				"id":             7,
				"codice":         "P1",
				"nome":           "Vite M6",
				"categoria_id":   1,
				"prezzo_vendita": 0.10,
			}})
		case "/api/v1/categories":
			json.NewEncoder(w).Encode([]map[string]any{{"id": 1, "nome": "Ferramenta"}})
		default:
			t.Errorf("endpoint non previsto dal test: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	m := New(api.New(srv.URL))

	msg := m.Init()()
	loaded, ok := msg.(CatalogoLoadedMsg)
	if !ok {
		t.Fatalf("atteso CatalogoLoadedMsg, ottenuto %T (%v)", msg, msg)
	}
	if len(loaded.Prodotti) != 1 || len(loaded.Categorie) != 1 {
		t.Fatalf("attesi 1 prodotto e 1 categoria, ottenuti %d e %d", len(loaded.Prodotti), len(loaded.Categorie))
	}
	if loaded.Prodotti[0].Categoria != "Ferramenta" {
		t.Errorf("attesa la categoria risolta dal categoria_id, ottenuta %q", loaded.Prodotti[0].Categoria)
	}
	if loaded.Prodotti[0].Prezzo != 0.10 {
		t.Errorf("atteso il prezzo di vendita come prezzo di listino, ottenuto %v", loaded.Prodotti[0].Prezzo)
	}
}

// TestLoadCatalogoCmd_ErroreServer verifica che un errore HTTP diventi un
// ErroreMsg e chiuda lo stato di caricamento.
func TestLoadCatalogoCmd_ErroreServer(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	m := New(api.New(srv.URL))

	msg := m.Init()()
	if _, ok := msg.(ErroreMsg); !ok {
		t.Fatalf("atteso ErroreMsg su risposta 500, ottenuto %T", msg)
	}

	res, _ := m.Update(msg)
	m = res.(Model)
	if m.loading || m.err == nil {
		t.Error("atteso loading=false ed errore memorizzato dopo il fallimento del caricamento")
	}
}

// caricaCatalogoDiProva restituisce un Model già popolato con due prodotti in
// due categorie diverse, per i test sul filtro e sulla navigazione.
func caricaCatalogoDiProva(t *testing.T) Model {
	t.Helper()
	m := New(nil)
	res, _ := m.Update(CatalogoLoadedMsg{
		Prodotti: []api.Prodotto{
			{ID: 1, Codice: "P1", Nome: "Vite M6", Categoria: "Ferramenta", Prezzo: 0.10, Attivo: true},
			{ID: 2, Codice: "P2", Nome: "Cacciavite", Categoria: "Utensili", Prezzo: 4.50, Attivo: true},
		},
		Categorie: []string{"Ferramenta", "Utensili"},
	})
	return res.(Model)
}
