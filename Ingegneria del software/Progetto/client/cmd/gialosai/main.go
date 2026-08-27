package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/app"
)

func main() {
	serverURL := os.Getenv("GIALOSAI_SERVER")
	if serverURL == "" {
		serverURL = "http://localhost:8080"
	}

	clientAPI := api.New(serverURL)
	app := app.New(clientAPI)

	p := tea.NewProgram(app)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Errore durante l'esecuzione dell'applicazione: %v\n", err)
		os.Exit(1)
	}
}
