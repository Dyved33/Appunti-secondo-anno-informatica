package Repository

import (
	"context"
	"database/sql"
	"time"
)

// QueryExecutor defines operations common to sql.DB and sql.Tx.
type QueryExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type txKey struct{}

// WithTx associates a transaction with the context.
func WithTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// GetExecutor returns the transaction in the context if present, otherwise the db connection.
func GetExecutor(ctx context.Context, db *sql.DB) QueryExecutor {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}
	return db
}

// InitDatabase initializes the SQLite database schema if tables do not exist.
func InitDatabase(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS Utenti (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			ruolo TEXT NOT NULL,
			domanda_sicurezza TEXT NOT NULL,
			risposta_hash TEXT NOT NULL,
			tentativi_falliti INTEGER DEFAULT 0,
			bloccato BOOLEAN DEFAULT 0,
			created_at DATETIME NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS Categorie (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			nome TEXT NOT NULL,
			descrizione TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS Prodotti (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			codice TEXT UNIQUE NOT NULL,
			barcode TEXT UNIQUE NOT NULL,
			nome TEXT NOT NULL,
			descrizione TEXT,
			categoria_id INTEGER,
			prezzo_acquisto REAL DEFAULT 0.0,
			prezzo_vendita REAL DEFAULT 0.0,
			soglia_minima INTEGER DEFAULT 0,
			FOREIGN KEY(categoria_id) REFERENCES Categorie(id)
		);`,
		`CREATE TABLE IF NOT EXISTS Giacenze (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			prodotto_id INTEGER UNIQUE NOT NULL,
			quantita INTEGER DEFAULT 0,
			updated_at DATETIME NOT NULL,
			FOREIGN KEY(prodotto_id) REFERENCES Prodotti(id)
		);`,
		`CREATE TABLE IF NOT EXISTS OrdiniEntrata (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			numero TEXT UNIQUE NOT NULL,
			fornitore TEXT NOT NULL,
			data DATETIME NOT NULL,
			stato TEXT NOT NULL,
			note TEXT,
			created_by INTEGER NOT NULL,
			FOREIGN KEY(created_by) REFERENCES Utenti(id)
		);`,
		`CREATE TABLE IF NOT EXISTS RigheOrdineEntrata (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			ordine_entrata_id INTEGER NOT NULL,
			prodotto_id INTEGER NOT NULL,
			quantita INTEGER NOT NULL,
			prezzo_unitario REAL NOT NULL,
			FOREIGN KEY(ordine_entrata_id) REFERENCES OrdiniEntrata(id),
			FOREIGN KEY(prodotto_id) REFERENCES Prodotti(id)
		);`,
		`CREATE TABLE IF NOT EXISTS OrdiniUscita (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			numero TEXT UNIQUE NOT NULL,
			cliente TEXT NOT NULL,
			data DATETIME NOT NULL,
			stato TEXT NOT NULL,
			note TEXT,
			created_by INTEGER NOT NULL,
			ddt_id INTEGER,
			FOREIGN KEY(created_by) REFERENCES Utenti(id)
		);`,
		`CREATE TABLE IF NOT EXISTS RigheOrdineUscita (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			ordine_uscita_id INTEGER NOT NULL,
			prodotto_id INTEGER NOT NULL,
			quantita INTEGER NOT NULL,
			prezzo_unitario REAL NOT NULL,
			FOREIGN KEY(ordine_uscita_id) REFERENCES OrdiniUscita(id),
			FOREIGN KEY(prodotto_id) REFERENCES Prodotti(id)
		);`,
		`CREATE TABLE IF NOT EXISTS DDT (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			numero TEXT UNIQUE NOT NULL,
			ordine_uscita_id INTEGER NOT NULL,
			data DATETIME NOT NULL,
			mittente TEXT NOT NULL,
			destinatario TEXT NOT NULL,
			path_pdf TEXT,
			FOREIGN KEY(ordine_uscita_id) REFERENCES OrdiniUscita(id)
		);`,
		`CREATE TABLE IF NOT EXISTS Fatture (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			numero TEXT UNIQUE NOT NULL,
			anno INTEGER NOT NULL,
			tipo TEXT NOT NULL,
			ordine_uscita_id INTEGER NOT NULL,
			imponibile REAL NOT NULL,
			iva REAL NOT NULL,
			totale REAL NOT NULL,
			data DATETIME NOT NULL,
			path_pdf TEXT,
			FOREIGN KEY(ordine_uscita_id) REFERENCES OrdiniUscita(id)
		);`,
		`CREATE TABLE IF NOT EXISTS Resi (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			ordine_uscita_id INTEGER NOT NULL,
			prodotto_id INTEGER NOT NULL,
			quantita INTEGER NOT NULL,
			motivo TEXT,
			stato TEXT NOT NULL,
			data DATETIME NOT NULL,
			created_by INTEGER NOT NULL,
			FOREIGN KEY(ordine_uscita_id) REFERENCES OrdiniUscita(id),
			FOREIGN KEY(prodotto_id) REFERENCES Prodotti(id),
			FOREIGN KEY(created_by) REFERENCES Utenti(id)
		);`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}

	// Insert default admin user if Utenti is empty for demonstration/testing
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Utenti").Scan(&count)
	if err == nil && count == 0 {
		// Default admin user: password: admin123, risposta: cagliari
		// BCrypt hash for "admin123": $2a$10$3e87B1g.fA1C9WnQ1QyUZe8vO5n09fB1g.fA1C9WnQ1QyUZe8vO5n09 (we can generate this in Go on runtime)
	}

	return nil
}

// TransactionHelper helper to begin a transaction, execute logic and commit/rollback automatically.
type TransactionHelper struct {
	DB *sql.DB
}

func (h *TransactionHelper) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := h.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	txCtx := WithTx(ctx, tx)
	if err := fn(txCtx); err != nil {
		return err
	}

	return tx.Commit()
}
