package auth

import (
	"context"
	"database/sql"
	"time"

	"github.com/tommasogarofani/ingsw_25-26_group-11/server/domain"
)

// SessionRepository gestisce la persistenza delle sessioni di login nella
// tabella Sessioni del database.
type SessionRepository struct {
	DB *sql.DB
}

// NewSessionRepository crea una nuova istanza del repository delle sessioni.
func NewSessionRepository(database *sql.DB) *SessionRepository {
	return &SessionRepository{DB: database}
}

// SessionRepository_Create inserisce una nuova sessione associata all'utente,
// memorizzando l'hash del token e la relativa data di scadenza.
func (r *SessionRepository) SessionRepository_Create(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) error {
	exec := domain.GetExecutor(ctx, r.DB)
	query := `INSERT INTO Sessioni (user_id, token_hash, expires_at) VALUES (?, ?, ?)`
	_, err := exec.ExecContext(ctx, query, userID, tokenHash, expiresAt)
	return err
}

// SessionRepository_IsValid verifica che esista una sessione con l'hash del
// token fornito e che non sia ancora scaduta. Restituisce true se la sessione
// è valida, insieme a un eventuale errore del database.
func (r *SessionRepository) SessionRepository_IsValid(ctx context.Context, tokenHash string) (bool, error) {
	exec := domain.GetExecutor(ctx, r.DB)
	query := `SELECT COUNT(*) FROM Sessioni WHERE token_hash = ? AND expires_at > NOW()`
	var count int
	err := exec.QueryRowContext(ctx, query, tokenHash).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// SessionRepository_DeleteByTokenHash elimina la sessione associata all'hash
// del token fornito; viene usata ad esempio in caso di logout.
func (r *SessionRepository) SessionRepository_DeleteByTokenHash(ctx context.Context, tokenHash string) error {
	exec := domain.GetExecutor(ctx, r.DB)
	query := `DELETE FROM Sessioni WHERE token_hash = ?`
	_, err := exec.ExecContext(ctx, query, tokenHash)
	return err
}
