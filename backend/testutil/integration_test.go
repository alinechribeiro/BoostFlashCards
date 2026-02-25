//go:build integration

package testutil

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/boostflashcards/backend/internal/config"
	"github.com/boostflashcards/backend/internal/handlers"
	"github.com/boostflashcards/backend/internal/repository"
	"github.com/boostflashcards/backend/internal/server"
)

func TestIntegration_SubjectsAndFlashcards(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION") == "" {
		t.Skip("set RUN_INTEGRATION=1 to run")
	}
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config: %v", err)
	}
	repo, err := repository.New(cfg.DSN())
	if err != nil {
		t.Skipf("database not available: %v", err)
	}
	defer repo.Close()

	h := handlers.New(repo)
	srv := server.New(h)

	req := httptest.NewRequest("GET", "/api/subjects", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("GET /api/subjects: %d", w.Code)
	}

	req = httptest.NewRequest("GET", "/api/subjects/1/topics", nil)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("GET /api/subjects/1/topics: %d", w.Code)
	}

	body := map[string]interface{}{
		"topic_id": 1,
		"front":    "What is the quadratic formula?",
		"back":     "x = (-b ± √(b²-4ac)) / 2a",
	}
	b, _ := json.Marshal(body)
	req = httptest.NewRequest("POST", "/api/flashcards", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("POST /api/flashcards: %d body %s", w.Code, w.Body.String())
	}
}
