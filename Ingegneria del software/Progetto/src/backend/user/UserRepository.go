package user

import (
	"context"
	"database/sql"
	"errors"
	"github.com/tommasogarofani/ingsw_25-26_group-11/src/backend/domain"
)

// questa struct serve per accedere al database e gestire le operazioni CRUD sugli utenti.
type UserRepository struct {
	DB *sql.DB
}

// NewUserRepository crea una nuova istanza del repository per gli utenti.
func NewUserRepository(database *sql.DB) *UserRepository {
	return &UserRepository{DB: database}
}

// GetByUsername recupera un utente dal database usando il suo username come chiave. Restituisce l'oggetto Utente se trovato, altrimenti un errore di credenziali errate.
// è l'equivalente di Read in un CRUD convenzionale
func (r *UserRepository) UserRepository_GetByUsername(ctx context.Context, username string) (*domain.Utente, error) {
	exec := domain.GetExecutor(ctx, r.DB) // serve per ottenere il corretto executor (DB o transazione) dal contesto
	query := `SELECT id, username, password_hash, ruolo, domanda_sicurezza, risposta_hash, tentativi_falliti, bloccato, created_at 
	          FROM Utenti WHERE username = ?`

	var u domain.Utente // crea un'istanza vuota di Utente per popolare i dati dal database
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

// Update modifica i dati di un utente esistente nel database
func (r *UserRepository) UserRepository_Update(ctx context.Context, u *domain.Utente) error {
	exec := domain.GetExecutor(ctx, r.DB)
	query := `UPDATE Utenti SET password_hash = ?, tentativi_falliti = ?, bloccato = ? WHERE id = ?`
	_, err := exec.ExecContext(ctx, query, u.PasswordHash, u.TentativiFalliti, u.Bloccato, u.ID)
	return err
}

// Create inserisce un nuovo utente nel database (assegna automaticamente l'ID generato dal database all'oggetto Utente passato)
func (r *UserRepository) UserRepository_Create(ctx context.Context, u *domain.Utente) error {
	exec := domain.GetExecutor(ctx, r.DB)
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

func(r *UserRepository) UserRepository_Delete(ctx context.Context, u *domain.Utente) error{
	exec := domain.GetExecutor(ctx, r.DB)
	query := `DELETE FROM Utenti WHERE id = ?`
	_, err := exec.ExecContext(ctx, query, u.ID)
	return err
}