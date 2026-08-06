package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/tommasogarofani/ingsw_25-26_group-11/server/router"
)

func main() {
	dsn := os.Getenv("DB_DSN") // Recupera la variabile d'ambiente DB_DSN per ottenere il Data Source Name (DSN) del database MySQL
	if dsn == "" {
		// DSN di default per connettersi al container MySQL definito nel docker-compose.yml
		dsn = "root:rootpassword@tcp(db:3306)/gialosai?parseTime=true"
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Errore durante l'apertura del database: %v", err)
	}
	defer db.Close()

	// Verifica la connessione al database
	if err := db.Ping(); err != nil {
		log.Fatalf("Impossibile connettersi al database MySQL: %v", err)
	}

	r := router.NewRouter(db) // Crea una nuova istanza del router (con tutte le rotte e i middleware già configurati)

	log.Println("Server in ascolto su :8080") // Avvia un server TCP in ascolto sulla porta 8080 e utilizza il router come handler per le richieste HTTP
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
