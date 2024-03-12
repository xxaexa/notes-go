package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	models "github.com/xxaexa/be-notes/internal/model"
)


func GetNotes(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var notes []models.Note
		rows, err := db.Query("SELECT id, title, content FROM notes")
		if err != nil {
			SendJSONResponse(w, http.StatusInternalServerError, "error", "Failed to fetch notes", nil)
			return
		}
		defer rows.Close()	

		for rows.Next() {
			var n models.Note
			if err := rows.Scan(&n.ID, &n.Title, &n.Content); err != nil {
				SendJSONResponse(w, http.StatusInternalServerError, "error", "Failed to scan notes", nil)
				return
			}
			notes = append(notes, n)
		}

		SendJSONResponse(w, http.StatusOK, "success", "Notes fetched successfully", notes)
	}
}

func CreateNote(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newNote models.Note

		// Menguraikan JSON dari body request ke struct Note
		if err := json.NewDecoder(r.Body).Decode(&newNote); err != nil {
			SendJSONResponse(w, http.StatusBadRequest, "error", "Invalid request payload", nil)
			return
		}

		// Menjalankan query INSERT ke database untuk membuat note baru
		var noteID int
		err := db.QueryRow("INSERT INTO notes (title, content) VALUES ($1, $2) RETURNING id", newNote.Title, newNote.Content).Scan(&noteID)
		if err != nil {
			SendJSONResponse(w, http.StatusInternalServerError, "error", "Failed to insert new note", nil)
			return
		}

		// Mengembalikan response sukses dengan ID note yang baru dibuat
		SendJSONResponse(w, http.StatusCreated, "success", "Note created successfully", map[string]int{"id": noteID})
	}
}