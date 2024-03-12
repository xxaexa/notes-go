package routes

import (
	"database/sql"
	"net/http"

	"github.com/xxaexa/be-notes/internal/handler"
)


func NoteRoutes(db *sql.DB) {
	http.HandleFunc("/notes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetNotes(db)(w, r)
		case http.MethodPost:
			handler.CreateNote(db)(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/notes/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetNote(db)(w, r)
		case http.MethodDelete:
			handler.DeleteNote(db)(w, r)
		case http.MethodPut:
			handler.UpdateNote(db)(w, r)
		default:
			http.Error(w, "Unsupported HTTP method", http.StatusMethodNotAllowed)
		}
	})
	
}