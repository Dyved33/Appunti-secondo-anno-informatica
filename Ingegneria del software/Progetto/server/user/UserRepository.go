package user

import (
	"context"
	"database/sql"
	"errors"
	"github.com/tommasogarofani/ingsw_25-26_group-11/server/domain"
)

// UserRepository accede al database e gestisce le operazioni CRUD sugli utenti.
type UserRepository struct {
	DB *sql.DB
}

// NewUserRepository crea una nuova istanza del repository per gli utenti.
func NewUserRepository(database *sql.DB) *UserRepository {
	return &UserRepository{DB: database}
}

// UserRepository_GetByUsername recupera un utente dal database usando il suo
// username come chiave (equivalente del Read in un CRUD convenzionale).
// Restituisce l'oggetto Utente se trovato, altrimenti un errore.
func (r *UserRepository) UserRepository_GetByUsername(ctx context.Context, username string) (*domain.Utente, error) {
	exec := domain.GetExecutor(ctx, r.DB) // serve per ottenere il corretto executor (DB o transazione) dal contesto
	query := `SELECT id, username, password_hash, ruolo, tentativi_falliti, bloccato, created_at 
	          FROM Utenti WHERE username = ?`

	var u domain.Utente // crea un'istanza vuota di Utente per popolare i dati dal database
	var ruoloStr string
	err := exec.QueryRowContext(ctx, query, username).Scan(
		&u.ID, &u.Username, &u.PasswordHash, &ruoloStr,
		&u.TentativiFalliti, &u.Bloccato, &u.CreatedAt,
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

// UserRepository_Update modifica i dati di un utente esistente nel database.
func (r *UserRepository) UserRepository_Update(ctx context.Context, u *domain.Utente) error {
	exec := domain.GetExecutor(ctx, r.DB)
	query := `UPDATE Utenti SET username = ?, password_hash = ?, ruolo = ?, tentativi_falliti = ?, bloccato = ? WHERE id = ?`
	_, err := exec.ExecContext(ctx, query, u.Username, u.PasswordHash, string(u.Ruolo), u.TentativiFalliti, u.Bloccato, u.ID)
	return err
}

// UserRepository_Create inserisce un nuovo utente nel database e assegna
// all'oggetto Utente l'ID auto-generato dal database.
func (r *UserRepository) UserRepository_Create(ctx context.Context, u *domain.Utente) error {
	exec := domain.GetExecutor(ctx, r.DB)
	query := `INSERT INTO Utenti (username, password_hash, ruolo, tentativi_falliti, bloccato, created_at)
	          VALUES (?, ?, ?, ?, ?, ?)`
	res, err := exec.ExecContext(ctx, query, u.Username, u.PasswordHash, string(u.Ruolo), u.TentativiFalliti, u.Bloccato, u.CreatedAt)
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

// UserRepository_Delete elimina dal database l'utente fornito, identificandolo
// tramite il suo ID.
func (r *UserRepository) UserRepository_Delete(ctx context.Context, u *domain.Utente) error {
	exec := domain.GetExecutor(ctx, r.DB)
	query := `DELETE FROM Utenti WHERE id = ?`
	_, err := exec.ExecContext(ctx, query, u.ID)
	return err
}
