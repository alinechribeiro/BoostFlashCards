package repository

import (
	"context"
	"database/sql"

	"github.com/boostflashcards/backend/internal/models"
)

func (r *Repository) ListSubjects(ctx context.Context) ([]models.Subject, error) {
	rows, err := r.DB.QueryContext(ctx,
		`SELECT id, name, slug, created_at, last_reviewed_at FROM subjects ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Subject
	for rows.Next() {
		var s models.Subject
		var lastReviewed sql.NullTime
		if err := rows.Scan(&s.ID, &s.Name, &s.Slug, &s.CreatedAt, &lastReviewed); err != nil {
			return nil, err
		}
		if lastReviewed.Valid {
			t := lastReviewed.Time
			s.LastReviewedAt = &t
		}
		list = append(list, s)
	}
	return list, rows.Err()
}

func (r *Repository) GetSubjectByID(ctx context.Context, id int64) (*models.Subject, error) {
	var s models.Subject
	var lastReviewed sql.NullTime
	err := r.DB.QueryRowContext(ctx,
		`SELECT id, name, slug, created_at, last_reviewed_at FROM subjects WHERE id = ?`, id).
		Scan(&s.ID, &s.Name, &s.Slug, &s.CreatedAt, &lastReviewed)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if lastReviewed.Valid {
		t := lastReviewed.Time
		s.LastReviewedAt = &t
	}
	return &s, nil
}

func (r *Repository) GetSubjectByName(ctx context.Context, name string) (*models.Subject, error) {
	var s models.Subject
	var lastReviewed sql.NullTime
	err := r.DB.QueryRowContext(ctx,
		`SELECT id, name, slug, created_at, last_reviewed_at FROM subjects WHERE name = ?`, name).
		Scan(&s.ID, &s.Name, &s.Slug, &s.CreatedAt, &lastReviewed)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if lastReviewed.Valid {
		t := lastReviewed.Time
		s.LastReviewedAt = &t
	}
	return &s, nil
}

func (r *Repository) CreateSubject(ctx context.Context, name, slug string) (*models.Subject, error) {
	res, err := r.DB.ExecContext(ctx,
		`INSERT INTO subjects (name, slug, created_at) VALUES (?, ?, NOW())`,
		name, slug)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return r.GetSubjectByID(ctx, id)
}

