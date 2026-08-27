package domain

import (
	"errors"
	"fmt"
)

// ─────────────────────────────────────────────
// ERRORI SENTINELLA
// Usati con errors.Is() per identificare il tipo di errore
// senza dipendere dal testo del messaggio.
// ─────────────────────────────────────────────

var (
	// ErrGiacenzaInsufficient è sollevato quando la quantità richiesta
	// per un ordine di uscita supera la giacenza disponibile (TC-01, RF-07).
	// Il service layer esegue ROLLBACK prima di restituire questo errore.
	ErrGiacenzaInsufficient = errors.New("giacenza insufficiente")

	// ErrProdottoNonTrovato è restituito quando il barcode scansionato
	// non corrisponde ad alcun prodotto in anagrafica (RF-15, criterio 3).
	ErrProdottoNonTrovato = errors.New("prodotto non trovato")

	// ErrCategoriaNonTrovata è restituito quando una categoria richiesta
	// non esiste più nel database.
	ErrCategoriaNonTrovata = errors.New("categoria non trovata")

	// ErrAccountBloccato è restituito quando l'utente ha superato
	// il numero massimo di tentativi di login falliti (RF-02).
	ErrAccountBloccato = errors.New("account bloccato")

	// ErrCredenzaliErrate è restituito quando username o password
	// non corrispondono a nessun utente attivo (RF-01).
	ErrCredenzaliErrate = errors.New("credenziali errate")

	// ErrPermessoNegato è restituito quando un utente tenta di eseguire
	// un'operazione non consentita dal suo ruolo (RF-01 RBAC).
	ErrPermessoNegato = errors.New("permesso negato: ruolo insufficiente")

	// ErrOrdineGiaEvaso è restituito quando si tenta di modificare
	// un ordine già nello stato "evaso" (RF-09).
	ErrOrdineGiaEvaso = errors.New("ordine già evaso: operazione non consentita")

	// ErrDDTGiaGenerato è restituito se si tenta di generare un secondo
	// DDT per un ordine di uscita già evaso (RF-08).
	ErrDDTGiaGenerato = errors.New("DDT già generato per questo ordine")

	// ErrQuantitaNonValida è restituito se la quantità fornita è <= 0.
	// Usato come guardia prima di qualsiasi operazione su giacenza.
	ErrQuantitaNonValida = errors.New("la quantità deve essere maggiore di zero")

	// ErrBarcodeNonValido è restituito se il barcode letto non rispetta
	// il formato EAN-13 (RF-15, criterio 1).
	ErrBarcodeNonValido = errors.New("barcode non valido: formato EAN-13 atteso")
)

// ─────────────────────────────────────────────
// ERRORI CONTESTUALI
// Wrappano gli errori sentinella aggiungendo dati specifici.
// Usati con errors.As() per estrarre i dettagli.
// ─────────────────────────────────────────────

// GiacenzaInsufficienteError contiene i dettagli del fallimento TC-01.
// Permette alla TUI di mostrare un messaggio preciso all'operatore.
type GiacenzaInsufficienteError struct {
	ProdottoID   int
	NomeProdotto string
	Disponibile  int
	Richiesta    int
}

func (e *GiacenzaInsufficienteError) Error() string {
	return fmt.Sprintf(
		"giacenza insufficiente per '%s' (id=%d): disponibili %d, richiesti %d",
		e.NomeProdotto, e.ProdottoID, e.Disponibile, e.Richiesta,
	)
}

// Unwrap consente a errors.Is(err, ErrGiacenzaInsufficient) di funzionare
// anche quando l'errore è wrapped in GiacenzaInsufficienteError.
func (e *GiacenzaInsufficienteError) Unwrap() error {
	return ErrGiacenzaInsufficient
}

// NewGiacenzaInsufficienteError costruisce l'errore contestuale per TC-01.
func NewGiacenzaInsufficienteError(id int, nome string, disponibile, richiesta int) error {
	return &GiacenzaInsufficienteError{
		ProdottoID:   id,
		NomeProdotto: nome,
		Disponibile:  disponibile,
		Richiesta:    richiesta,
	}
}

// ─────────────────────────────────────────────
// HELPER
// ─────────────────────────────────────────────

// IsGiacenzaInsufficient è un helper leggibile per il controllo nel service layer.
//
//	if domain.IsGiacenzaInsufficient(err) { /* esegui ROLLBACK */ }
func IsGiacenzaInsufficient(err error) bool {
	return errors.Is(err, ErrGiacenzaInsufficient)
}

// IsNotFound è un helper per il caso barcode non trovato (RF-15, criterio 3).
func IsNotFound(err error) bool {
	return errors.Is(err, ErrProdottoNonTrovato)
}

// IsCategoryNotFound è un helper per il caso categoria assente.
func IsCategoryNotFound(err error) bool {
	return errors.Is(err, ErrCategoriaNonTrovata)
}

// IsPermissionDenied è un helper per i controlli RBAC nei controller (RF-01).
func IsPermissionDenied(err error) bool {
	return errors.Is(err, ErrPermessoNegato)
}
