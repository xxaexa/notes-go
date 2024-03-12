package main

import (
	"log"
	"net/http"

	"github.com/xxaexa/be-notes/internal/config"
	"github.com/xxaexa/be-notes/internal/routes"
)


func main() {
	dbConn, err := config.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	//routes
	routes.NoteRoutes(dbConn)
	
	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}