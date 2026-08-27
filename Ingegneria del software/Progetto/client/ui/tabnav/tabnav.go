// Package tabnav gestisce la barra dei tab cliccabile condivisa dai layout
// Magazzino, Contabilità e Admin.
//
// Strategia: la composizione della pagina è definita UNA sola volta qui
// (PlacePage) ed è usata sia da client/ui/app (View) sia dal calcolo delle
// aree cliccabili (Rects). Le hitbox non sono calcolate con formule geometriche
// ma MISURATE sul rendering reale: si compone la pagina, si rimuovono i codici
// ANSI e si cercano le etichette dei tab sulla riga di navigazione. Qualsiasi
// modifica a stili, bordi o centratura resta quindi automaticamente coerente.
package tabnav

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
	components "github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/components"
	style "github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/styles"

	"charm.land/lipgloss/v2"
)

// Rect è un'area rettangolare cliccabile, estremi inclusi, in coordinate schermo.
type Rect struct {
	X0, Y0, X1, Y1 int
}

// Contains indica se il punto (x, y) cade dentro il rettangolo.
func (r Rect) Contains(x, y int) bool {
	return x >= r.X0 && x <= r.X1 && y >= r.Y0 && y <= r.Y1
}

// PlacePage compone la pagina finale esattamente come la mostra il programma:
// header (logo + intestazione), riga vuota, blocco contenuti incorniciato,
// tutto centrato nel terminale. È l'unica definizione del layout d'insieme:
// client/ui/app/app.go la usa per il rendering e Rects la usa per le hitbox.
func PlacePage(winW, winH int, contentBlock string) string {
	headerBlock := lipgloss.JoinHorizontal(lipgloss.Left, "", components.Logo, "")
	headerBlock = lipgloss.JoinVertical(lipgloss.Center, headerBlock, components.Intestazione)

	finalPage := lipgloss.JoinVertical(lipgloss.Center, headerBlock, "", contentBlock)

	targetW := max(lipgloss.Width(finalPage), winW)
	targetH := max(lipgloss.Height(finalPage), winH)

	if targetW > 0 && targetH > 0 {
		finalPage = lipgloss.Place(targetW, targetH, lipgloss.Center, lipgloss.Center, finalPage)
	}
	return finalPage
}

// Rects restituisce, per ciascun tab, l'area cliccabile in coordinate schermo.
//
// winW e winH sono le dimensioni correnti del terminale (WindowSizeMsg);
// labels contiene le etichette dei tab nell'ordine visuale; activeIdx è
// l'indice del tab attivo; page è la pagina NON incorniciata ("tabs\n\ncorpo")
// così com'è prodotta da rawPage() del layout chiamante.
//
// Se le dimensioni non sono ancora note (valori <= 0) vengono usati 80x24.
// Restituisce nil se una o più etichette non sono reperibili sul rendering:
// in tal caso i click vengono semplicemente ignorati, senza errori.
func Rects(winW, winH int, labels []string, activeIdx int, page string) []Rect {
	if winW <= 0 || winH <= 0 {
		winW, winH = 80, 24
	}
	if len(labels) == 0 {
		return nil
	}
	if activeIdx < 0 || activeIdx >= len(labels) {
		activeIdx = 0
	}

	// Rendering finale identico a quello mostrato all'utente, senza codici colore.
	composed := ansi.Strip(PlacePage(winW, winH, style.FinestraStile.Render(page)))
	lines := strings.Split(composed, "\n")

	// La barra dei tab è individuata cercando tutte le etichette in ordine
	// su due righe consecutive: JoinHorizontal(Top) allinea i blocchi in alto,
	// quindi il testo dei tab inattivi può stare sulla riga del bordo
	// superiore del tab attivo.
	spans, ok := findLabels(lines, labels)
	if !ok {
		return nil
	}

	rects := make([]Rect, len(labels))
	for i, label := range labels {
		w := max(1, ansi.StringWidth(label))
		x0, x1 := spans[i].x, spans[i].x+w-1
		y := spans[i].y
		if i == activeIdx {
			// Il tab attivo è racchiuso da un bordo completo: rendilo cliccabile.
			x0--
			x1++
			y--
		}
		// Dal bordo superiore (se presente) fino alla riga sottostante
		// (bordo inferiore / gap).
		rects[i] = Rect{X0: x0, X1: x1, Y0: y, Y1: spans[i].y + 1}
	}
	return rects
}

// span è la posizione (x, y) dell'inizio di un'etichetta nel rendering privo di ANSI.
type span struct {
	x, y int
}

// findLabels cerca tutte le etichette, nell'ordine, su due righe consecutive
// e restituisce la posizione iniziale di ciascuna. Su ciascuna riga le
// occorrenze devono procedere da sinistra a destra; l'ordine può however
// alternarsi tra le due righe (il testo dei tab inattivi e quello del tab
// attivo stanno tipicamente su righe diverse).
func findLabels(lines []string, labels []string) ([]span, bool) {
	for y := 0; y+1 < len(lines); y++ {
		cur, nxt := lines[y], lines[y+1]
		curFrom := 0 // prossima colonna libera sulla riga corrente
		nxtFrom := 0 // prossima colonna libera sulla riga successiva
		spans := make([]span, 0, len(labels))
		found := true

		for _, label := range labels {
			if idx := strings.Index(cur[curFrom:], label); idx >= 0 {
				absByte := curFrom + idx
				spans = append(spans, span{x: ansi.StringWidth(cur[:absByte]), y: y})
				curFrom = absByte + len(label)
				continue
			}
			if idx := strings.Index(nxt[nxtFrom:], label); idx >= 0 {
				absByte := nxtFrom + idx
				spans = append(spans, span{x: ansi.StringWidth(nxt[:absByte]), y: y + 1})
				nxtFrom = absByte + len(label)
				continue
			}
			found = false
			break
		}
		if found && len(spans) == len(labels) {
			return spans, true
		}
	}
	return nil, false
}

// HitAt restituisce l'indice del tab cliccato alle coordinate (x, y), oppure -1.
func HitAt(rects []Rect, x, y int) int {
	for i, r := range rects {
		if r.Contains(x, y) {
			return i
		}
	}
	return -1
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
