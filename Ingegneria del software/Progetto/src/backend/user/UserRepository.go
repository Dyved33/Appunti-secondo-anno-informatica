package user

import (
	"context"
	"database/sql"
	"errors"
	"gialosai/backend/db"
	"gialosai/backend/domain"
)

// Repository incapsula tutte le operazioni di persistenza relative agli utenti.
// Comunica direttamente con il database MySQL attraverso il QueryExecutor fornito dal layer db.
type Repository struct {
	DB *sql.DB
}

// NewRepository crea una nuova istanza del repository per gli utenti.
func NewRepository(database *sql.DB) *Repository {
	return &Repository{DB: database}
}

// GetByUsername recupera un utente dal database usando il suo username come chiave.
// Restituisce l'oggetto Utente se trovato, altrimenti un errore di credenziali errate.
func (r *Repository) GetByUsername(ctx context.Context, username string) (*domain.Utente, error) {
	// Usa GetExecutor per ottenere il corretto executor (DB o transazione) dal contesto
	exec := db.GetExecutor(ctx, r.DB)
	query := `SELECT id, username, password_hash, ruolo, domanda_sicurezza, risposta_hash, tentativi_falliti, bloccato, created_at 
	          FROM Utenti WHERE username = ?`

	var u domain.Utente
	var ruoloStr string
	err := exec.QueryRowContext(ctx, query, username).Scan(
		&u.ID, &u.Username, &u.PasswordHash, &ruoloStr,
		&u.DomandaSicurezza, &u.RispostaHash, &u.TentativiFalliti,
		&u.Bloccato, &u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrCredenzaliErrate
		}
		return nil, err
	}
	// Converte la stringa del ruolo nel tipo enum Ruolo
	u.Ruolo = domain.Ruolo(ruoloStr)
	return &u, nil
}

// Update modifica i dati di un utente esistente nel database.
// Aggiorna la password (hash), il contatore dei tentativi falliti e lo stato di blocco dell'account.
func (r *Repository) Update(ctx context.Context, u *domain.Utente) error {
	exec := db.GetExecutor(ctx, r.DB)
	query := `UPDATE Utenti SET password_hash = ?, tentativi_falliti = ?, bloccato = ? WHERE id = ?`
	_, err := exec.ExecContext(ctx, query, u.PasswordHash, u.TentativiFalliti, u.Bloccato, u.ID)
	return err
}

// Create inserisce un nuovo utente nel database.
// Assegna automaticamente l'ID generato dal database all'oggetto Utente passato.
func (r *Repository) Create(ctx context.Context, u *domain.Utente) error {
	exec := db.GetExecutor(ctx, r.DB)
	query := `INSERT INTO Utenti (username, password_hash, ruolo, domanda_sicurezza, risposta_hash, tentativi_falliti, bloccato, created_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	res, err := exec.ExecContext(ctx, query, u.Username, u.PasswordHash, string(u.Ruolo), u.DomandaSicurezza, u.RispostaHash, u.TentativiFalliti, u.Bloccato, u.CreatedAt)
	if err != nil {
		return err
	}
	// Recupera l'ID auto-generato dal database e lo assegna all'utente
	id, err := res.LastInsertId()
	if err == nil {
		u.ID = int(id)
	}
	return nil
}
