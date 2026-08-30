package style

// Spacing rappresenta i margini di spaziatura di una finestra su ciascun lato.
type Spacing struct {
	top    int
	bottom int
	left   int
	right  int
}

// WindowSize rappresenta le dimensioni della finestra di terminale insieme ai
// margini, ai bordi e ai padding degli elementi interni.
type WindowSize struct {
	width   int
	height  int
	Margin  Spacing
	Border  Spacing
	Padding Spacing
}

// NewWindowSize crea e restituisce una WindowSize con la larghezza e l'altezza date.
func NewWindowSize(width, height int) WindowSize {
	return WindowSize{
		width:  width,
		height: height,
	}
}

// ChangeSize aggiorna larghezza e altezza della finestra ai valori specificati.
func (w *WindowSize) ChangeSize(width, height int) {
	w.height = height
	w.width = width
}

// ReadWidth restituisce la larghezza della finestra.
func ReadWidth(w WindowSize) int {
	return w.width
}

// ReadHeight restituisce l'altezza della finestra.
func ReadHeight(w WindowSize) int {
	return w.height
}

// GetInnerSpace restituisce le dimensioni utili (larghezza, altezza) della
// finestra dopo aver sottratto margini, bordi e padding.
func GetInnerSpace(w WindowSize) (width, height int) {
	width = w.width - w.Margin.left - w.Margin.right - w.Border.right - w.Border.left - w.Padding.left - w.Padding.right
	height = w.height - w.Margin.top - w.Margin.bottom - w.Border.top - w.Border.bottom - w.Padding.top - w.Padding.bottom

	return width, height
}

// SetSpacing imposta la spaziatura del lato superiore, inferiore, sinistro e destro.
func (s *Spacing) SetSpacing(t, b, l, r int) {
	s.top = t
	s.bottom = b
	s.left = l
	s.right = r
}
