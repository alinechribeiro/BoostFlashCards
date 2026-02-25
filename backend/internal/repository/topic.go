package repository

import (
	"context"
	"database/sql"
	"strings"

	"github.com/boostflashcards/backend/internal/models"
)

func (r *Repository) ListTopicsBySubjectID(ctx context.Context, subjectID int64) ([]models.Topic, error) {
	rows, err := r.DB.QueryContext(ctx,
		`SELECT id, subject_id, name, slug, created_at FROM topics WHERE subject_id = ? ORDER BY name`,
		subjectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Topic
	for rows.Next() {
		var t models.Topic
		if err := rows.Scan(&t.ID, &t.SubjectID, &t.Name, &t.Slug, &t.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, t)
	}
	return list, rows.Err()
}

func (r *Repository) GetTopicByID(ctx context.Context, id int64) (*models.Topic, error) {
	var t models.Topic
	err := r.DB.QueryRowContext(ctx,
		`SELECT id, subject_id, name, slug, created_at FROM topics WHERE id = ?`, id).
		Scan(&t.ID, &t.SubjectID, &t.Name, &t.Slug, &t.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *Repository) CreateTopic(ctx context.Context, subjectID int64, name, slug string) (*models.Topic, error) {
	baseSlug := slug
	trySlug := slug
	for i := 0; i < 5; i++ {
		res, err := r.DB.ExecContext(ctx,
			`INSERT INTO topics (subject_id, name, slug, created_at) VALUES (?, ?, ?, NOW())`,
			subjectID, name, trySlug)
		if err != nil {
			// If duplicate key, tweak slug and retry a few times.
			if strings.Contains(err.Error(), "Duplicate entry") {
				trySlug = baseSlug + "-" + string('a'+i)
				continue
			}
			return nil, err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}
		return r.GetTopicByID(ctx, id)
	}
	return nil, sql.ErrNoRows
}

