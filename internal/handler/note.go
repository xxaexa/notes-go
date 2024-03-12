package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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




func GetNote(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) != 3 {
			http.Error(w, "Invalid URL format", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(parts[2])
		if err != nil {
			http.Error(w, "Invalid note ID", http.StatusBadRequest)
			return
		}

		var note models.Note
		err = db.QueryRow("SELECT id, title, content FROM notes WHERE id = $1", id).Scan(&note.ID, &note.Title, &note.Content)
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		} else if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching note: %v", err), http.StatusInternalServerError)
			return
		}

		SendJSONResponse(w, http.StatusOK, "success", "Notes fetched successfully", note)
	}
}

func CreateNote(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newNote models.Note

		if err := json.NewDecoder(r.Body).Decode(&newNote); err != nil {
			SendJSONResponse(w, http.StatusBadRequest, "error", "Invalid request payload", nil)
			return
		}

		var noteID int
		err := db.QueryRow("INSERT INTO notes (title, content) VALUES ($1, $2) RETURNING id", newNote.Title, newNote.Content).Scan(&noteID)
		if err != nil {
			SendJSONResponse(w, http.StatusInternalServerError, "error", "Failed to insert new note", nil)
			return
		}

		
		SendJSONResponse(w, http.StatusCreated, "success", "Note created successfully", map[string]int{"id": noteID})
	}
}

func DeleteNote(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 3 {
			http.Error(w, "Invalid URL or Note ID", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(parts[2])
		if err != nil {
			http.Error(w, "Invalid note ID", http.StatusBadRequest)
			return
		}

		_, err = db.Exec("DELETE FROM notes WHERE id = $1", id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error deleting note: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		
	}
}

func UpdateNote(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 3 {
			http.Error(w, "Invalid URL or Note ID", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(parts[2])
		if err != nil {
			http.Error(w, "Invalid note ID", http.StatusBadRequest)
			return
		}

		var updatedNote models.Note
		if err := json.NewDecoder(r.Body).Decode(&updatedNote); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		_, err = db.Exec("UPDATE notes SET title = $1, content = $2 WHERE id = $3", updatedNote.Title, updatedNote.Content, id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error updating note: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedNote)
	}
}