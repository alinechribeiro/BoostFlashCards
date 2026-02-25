package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/boostflashcards/backend/internal/models"
)

func (r *Repository) ListFlashcardsByTopicID(ctx context.Context, topicID int64) ([]models.Flashcard, error) {
	rows, err := r.DB.QueryContext(ctx,
		`SELECT id, topic_id, front, back, status, created_at, updated_at, last_reviewed_at, next_due_at
		   FROM flashcards
		  WHERE topic_id = ?
		    AND (
				status = 'not_yet'
				OR next_due_at IS NULL
				OR next_due_at <= NOW()
			)
		  ORDER BY created_at`,
		topicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Flashcard
	for rows.Next() {
		var f models.Flashcard
		var lastReviewed, nextDue sql.NullTime
		if err := rows.Scan(&f.ID, &f.TopicID, &f.Front, &f.Back, &f.Status, &f.CreatedAt, &f.UpdatedAt, &lastReviewed, &nextDue); err != nil {
			return nil, err
		}
		if lastReviewed.Valid {
			t := lastReviewed.Time
			f.LastReviewedAt = &t
		}
		if nextDue.Valid {
			t := nextDue.Time
			f.NextDueAt = &t
		}
		list = append(list, f)
	}
	return list, rows.Err()
}

func (r *Repository) GetFlashcardByID(ctx context.Context, id int64) (*models.Flashcard, error) {
	var f models.Flashcard
	var lastReviewed, nextDue sql.NullTime
	err := r.DB.QueryRowContext(ctx,
		`SELECT id, topic_id, front, back, status, created_at, updated_at, last_reviewed_at, next_due_at
		   FROM flashcards WHERE id = ?`, id).
		Scan(&f.ID, &f.TopicID, &f.Front, &f.Back, &f.Status, &f.CreatedAt, &f.UpdatedAt, &lastReviewed, &nextDue)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if lastReviewed.Valid {
		t := lastReviewed.Time
		f.LastReviewedAt = &t
	}
	if nextDue.Valid {
		t := nextDue.Time
		f.NextDueAt = &t
	}
	return &f, nil
}

func (r *Repository) CreateFlashcard(ctx context.Context, req models.CreateFlashcardRequest) (*models.Flashcard, error) {
	res, err := r.DB.ExecContext(ctx,
		`INSERT INTO flashcards (topic_id, front, back, status, created_at, updated_at) VALUES (?, ?, ?, 'not_yet', NOW(), NOW())`,
		req.TopicID, req.Front, req.Back)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return r.GetFlashcardByID(ctx, id)
}

func (r *Repository) UpdateFlashcard(ctx context.Context, id int64, front, back *string) (*models.Flashcard, error) {
	row := r.DB.QueryRowContext(ctx, `SELECT front, back FROM flashcards WHERE id = ?`, id)
	var currentFront, currentBack string
	if err := row.Scan(&currentFront, &currentBack); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	f, b := currentFront, currentBack
	if front != nil {
		f = *front
	}
	if back != nil {
		b = *back
	}

	_, err := r.DB.ExecContext(ctx,
		`UPDATE flashcards SET front = ?, back = ?, updated_at = NOW() WHERE id = ?`, f, b, id)
	if err != nil {
		return nil, err
	}
	return r.GetFlashcardByID(ctx, id)
}

func (r *Repository) DeleteFlashcard(ctx context.Context, id int64) (bool, error) {
	res, err := r.DB.ExecContext(ctx, `DELETE FROM flashcards WHERE id = ?`, id)
	if err != nil {
		return false, err
	}
	affected, _ := res.RowsAffected()
	return affected > 0, nil
}

// UpdateFlashcardStatus updates confidence status and scheduling, and touches the subject's last_reviewed_at.
func (r *Repository) UpdateFlashcardStatus(ctx context.Context, id int64, status models.FlashcardStatus) (*models.Flashcard, error) {
	var topicID int64
	if err := r.DB.QueryRowContext(ctx, `SELECT topic_id FROM flashcards WHERE id = ?`, id).Scan(&topicID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	now := time.Now()
	var nextDue interface{}
	if status == models.FlashcardStatusConfident {
		nextDue = now.Add(14 * 24 * time.Hour)
	} else {
		nextDue = nil
	}

	_, err := r.DB.ExecContext(ctx, `
		UPDATE flashcards
		   SET status = ?, last_reviewed_at = ?, next_due_at = ?
		 WHERE id = ?`,
		status, now, nextDue, id)
	if err != nil {
		return nil, err
	}

	// Touch subject last_reviewed_at
	_, _ = r.DB.ExecContext(ctx, `
		UPDATE subjects
		   SET last_reviewed_at = ?
		 WHERE id = (SELECT subject_id FROM topics WHERE id = ?)`,
		now, topicID)

	return r.GetFlashcardByID(ctx, id)
}
