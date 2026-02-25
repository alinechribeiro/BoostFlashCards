package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/boostflashcards/backend/internal/models"
	"github.com/boostflashcards/backend/internal/repository"
	"github.com/boostflashcards/backend/internal/server"
	"github.com/gorilla/mux"
)

func TestCreateFlashcard_InvalidBody(t *testing.T) {
	repo := &repository.Repository{DB: nil}
	h := New(repo, nil, nil)
	srv := server.New(h)

	req := httptest.NewRequest("POST", "/api/flashcards", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreateFlashcard_EmptyFront(t *testing.T) {
	repo := &repository.Repository{DB: nil}
	h := New(repo, nil, nil)
	srv := server.New(h)

	body := models.CreateFlashcardRequest{TopicID: 1, Front: "", Back: "answer"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/flashcards", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreateFlashcard_EmptyBack(t *testing.T) {
	repo := &repository.Repository{DB: nil}
	h := New(repo, nil, nil)
	srv := server.New(h)

	body := models.CreateFlashcardRequest{TopicID: 1, Front: "question", Back: ""}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/flashcards", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestGetFlashcard_InvalidID(t *testing.T) {
	repo := &repository.Repository{DB: nil}
	h := New(repo, nil, nil)
	r := mux.NewRouter()
	r.HandleFunc("/api/flashcards/{id}", h.GetFlashcard).Methods("GET")

	req := httptest.NewRequest("GET", "/api/flashcards/abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid id, got %d", w.Code)
	}
}

func TestHealth(t *testing.T) {
	repo := &repository.Repository{DB: nil}
	h := New(repo, nil, nil)
	srv := server.New(h)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if body := w.Body.String(); body != "OK" {
		t.Errorf("expected body OK, got %q", body)
	}
}
