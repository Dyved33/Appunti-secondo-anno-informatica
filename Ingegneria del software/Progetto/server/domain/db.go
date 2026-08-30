package domain

import (
	"context"
	"database/sql"
)

// per prima cosa ho bisogno di memorizzare una transazione nel context (inserisce il puntatore alla transazione dentro il context)
type contextKey string

const txKey contextKey = "tx"

// QueryExecutor definisce i metodi necessari per eseguire query sia sulla
// connessione al database sia su una transazione, così da poter usare lo
// stesso codice in entrambe le situazioni.
type QueryExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// GetExecutor restituisce la transazione presente nel context, se ne esiste
// una; altrimenti restituisce la connessione al database di default. Consente
// di eseguire le query nel contesto di esecuzione corretto.
func GetExecutor(ctx context.Context, defaultDB *sql.DB) QueryExecutor {
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		return tx
	}
	return defaultDB
}

// TransactionHelper incapsula una connessione al database e fornisce un
// mezzo per eseguire funzioni all'interno di una transazione.
type TransactionHelper struct {
	DB *sql.DB
}

// NewTransactionHelper crea un nuovo TransactionHelper che utilizza la
// connessione al database fornita.
func NewTransactionHelper(db *sql.DB) *TransactionHelper {
	return &TransactionHelper{DB: db}
}

// WithinTransaction esegue la funzione fn all'interno di una transazione.
// Se nel context è già presente una transazione, la funzione viene eseguita
// con quella esistente senza crearne una nuova; altrimenti viene avviata una
// nuova transazione, eseguita fn nel relativo contesto e, in caso di errore,
// viene eseguito un ROLLBACK; in caso di successo la transazione viene
// salvata con COMMIT.
func (th *TransactionHelper) WithinTransaction(ctx context.Context, fn func(context.Context) error) error {
	if _, ok := ctx.Value(txKey).(*sql.Tx); ok { // se già esiste una transazione nel contesto
		return fn(ctx) // non ne creo una nuova, ma eseguo la funzione con il contesto esistente
	}

	tx, err := th.DB.BeginTx(ctx, nil) // creo una nuova transazione
	if err != nil {
		return err
	}

	txCtx := context.WithValue(ctx, txKey, tx) // creo un nuovo contesto con la transazione appena creata, così da poterla usare all'interno della funzione passata come parametro
	err = fn(txCtx)
	if err != nil {
		_ = tx.Rollback() // annullo modifiche
		return err
	}

	return tx.Commit() // modifiche permanenti
}
