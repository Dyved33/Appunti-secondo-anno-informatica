package style

// STRUCT PER I CALCOLI DI SPAZIATURA
type Spacing struct {
	top    int
	bottom int
	left   int
	right  int
}

// STRUCT PER LA DIMENSIONE FINESTRA
type WindowSize struct {
	width   int
	height  int
	Margin  Spacing
	Border  Spacing
	Padding Spacing
}

// COSTRUTTORE DIMENSIONI FINESTRA
func NewWindowSize(width, height int) WindowSize {
	return WindowSize{
		width:  width,
		height: height,
	}
}

// CAMBIO DIMENSIONI
func (w *WindowSize) ChangeSize(width, height int) {
	w.height = height
	w.width = width
}

// LEGGI LA LARGHEZZA
func ReadWidth(w WindowSize) int {
	return w.width
}

// LEGGI L'ALTEZZA
func ReadHeight(w WindowSize) int {
	return w.height
}

// RICEVE LE DIMENSIONI DI TUTTA LA FINESTRA
func GetInnerSpace(w WindowSize) (width, height int) {
	width = w.width - w.Margin.left - w.Margin.right - w.Border.right - w.Border.left - w.Padding.left - w.Padding.right
	height = w.height - w.Margin.top - w.Margin.bottom - w.Border.top - w.Border.bottom - w.Padding.top - w.Padding.bottom

	return width, height
}

// MODIFICO LO SPACING DEGLI OGGETTI INTERNI
func (s *Spacing) SetSpacing(t, b, l, r int) {
	s.top = t
	s.bottom = b
	s.left = l
	s.right = r
}
