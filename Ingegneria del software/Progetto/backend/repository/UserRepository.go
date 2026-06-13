package Repository

import (
	"context"
	"database/sql"
	"errors"
	"gialosai/backend/domain"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.Utente, error) {
	exec := GetExecutor(ctx, r.DB)
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
	u.Ruolo = domain.Ruolo(ruoloStr)
	return &u, nil
}

func (r *UserRepository) Update(ctx context.Context, u *domain.Utente) error {
	exec := GetExecutor(ctx, r.DB)
	query := `UPDATE Utenti SET password_hash = ?, tentativi_falliti = ?, bloccato = ? WHERE id = ?`
	_, err := exec.ExecContext(ctx, query, u.PasswordHash, u.TentativiFalliti, u.Bloccato, u.ID)
	return err
}

func (r *UserRepository) Create(ctx context.Context, u *domain.Utente) error {
	exec := GetExecutor(ctx, r.DB)
	query := `INSERT INTO Utenti (username, password_hash, ruolo, domanda_sicurezza, risposta_hash, tentativi_falliti, bloccato, created_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	res, err := exec.ExecContext(ctx, query, u.Username, u.PasswordHash, string(u.Ruolo), u.DomandaSicurezza, u.RispostaHash, u.TentativiFalliti, u.Bloccato, u.CreatedAt)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		u.ID = int(id)
	}
	return nil
}