package test

import (
	"testing"

	models "github.com/xxaexa/be-notes/internal/model"
)


func TestCreateNote(t *testing.T) {
	note := models.Note{Title: "Test Title", Content: "Test Content"}

	if note.Title != "Test Title" {
		t.Errorf("Expected 'Test Title', got %s", note.Title )
	}

	if note.Content != "Test Content" {
		t.Errorf("Expected 'Test Title', got %s", note.Content )
	}

}