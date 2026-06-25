package db

import (
	"context"
	"database/sql"
)

// per prima cosa ho bisogno di memorizzare una transazione nel context (inserisce il puntatore alla transazione dentro il context)
type contextKey string
const txKey contextKey = "tx"

// questa interfaccia mi serve per definire i metodi necessari per eseguire query sia per la connessione al database che  per la transazione, così da poter usare lo stesso codice per entrambe le situazioni.
type QueryExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// controlla se nel context è presente una transazione (se c'è, la restituisce; altrimenti restituisce la connessione al database). Mi serve per eseguire le query nel contesto corretto
func GetExecutor(ctx context.Context, defaultDB *sql.DB) QueryExecutor {
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		return tx
	}
	return defaultDB
}

// si tratta di una struttura che contiene il database
type TransactionHelper struct {
	DB *sql.DB
}

// crea un nuovo TransactionHelper che utilizza la connessione al database fornita
func NewTransactionHelper(db *sql.DB) *TransactionHelper {
	return &TransactionHelper{DB: db}
}

// esegue la funzione passata come parametro all'interno del contesto 
func (th *TransactionHelper) WithinTransaction(ctx context.Context, fn func(context.Context) error) error {
	if _, ok := ctx.Value(txKey).(*sql.Tx); ok {// se già esiste una transazione nel contesto
		return fn(ctx)// non ne creo una nuova, ma eseguo la funzione con il contesto esistente
	}

	tx, err := th.DB.BeginTx(ctx, nil)// creo una nuova transazione
	if err != nil {
		return err
	}

	txCtx := context.WithValue(ctx, txKey, tx)// creo un nuovo contesto con la transazione appena creata, così da poterla usare all'interno della funzione passata come parametro
	err = fn(txCtx)
	if err != nil {
		_ = tx.Rollback()// annullo modifiche
		return err
	}

	return tx.Commit()// modifiche permanenti
}
