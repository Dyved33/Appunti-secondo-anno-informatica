package main

import (
    "log"
    "net/http"

    "database/sql"

    "github.com/tommasogarofani/ingsw_25-26_group-11/src/backend/router"
)

func main() {
    db, err := sql.Open("sqlite", "gialosai.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
    
    r := router.New(db)// Crea una nuova istanza del router (con tutte le rotte e i middleware già configurati)

    log.Println("Server in ascolto su :8080")// Avvia un server TCP in ascolto sulla porta 8080 e utilizza il router come handler per le richieste HTTP
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatal(err)
    }
}