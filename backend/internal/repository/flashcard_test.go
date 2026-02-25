package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/boostflashcards/backend/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateFlashcard(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	defer db.Close()
	repo := &Repository{DB: db}

	mock.ExpectExec("INSERT INTO flashcards").
		WithArgs(1, "What is 2+2?", "4").
		WillReturnResult(sqlmock.NewResult(10, 1))

	mock.ExpectQuery("SELECT .+ FROM flashcards WHERE id = ?").
		WithArgs(10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "topic_id", "front", "back", "created_at", "updated_at"}).
			AddRow(10, 1, "What is 2+2?", "4", nil, nil))

	ctx := context.Background()
	card, err := repo.CreateFlashcard(ctx, models.CreateFlashcardRequest{
		TopicID: 1,
		Front:   "What is 2+2?",
		Back:    "4",
	})
	if err != nil {
		t.Fatalf("CreateFlashcard: %v", err)
	}
	if card == nil || card.ID != 10 || card.Front != "What is 2+2?" || card.Back != "4" {
		t.Errorf("unexpected card: %+v", card)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations: %v", err)
	}
}

func TestGetFlashcard_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	defer db.Close()
	repo := &Repository{DB: db}

	mock.ExpectQuery("SELECT .+ FROM flashcards WHERE id = ?").
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	ctx := context.Background()
	card, err := repo.GetFlashcardByID(ctx, 999)
	if err != nil {
		t.Fatalf("GetFlashcardByID: %v", err)
	}
	if card != nil {
		t.Errorf("expected nil card, got %+v", card)
	}
}
