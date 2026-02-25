package repository

import (
	"context"
	"database/sql"

	"github.com/boostflashcards/backend/internal/models"
)

func (r *Repository) CreateAnswerAttempt(ctx context.Context, a models.AnswerAttempt) (*models.AnswerAttempt, error) {
	res, err := r.DB.ExecContext(ctx,
		`INSERT INTO answer_attempts (subject_id, topic_id, question, student_answer, predicted_score, max_score, predicted_grade, feedback)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		a.SubjectID, a.TopicID, a.Question, a.StudentAnswer, a.PredictedScore, a.MaxScore, a.PredictedGrade, a.Feedback)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return r.GetAnswerAttemptByID(ctx, id)
}

func (r *Repository) GetAnswerAttemptByID(ctx context.Context, id int64) (*models.AnswerAttempt, error) {
	row := r.DB.QueryRowContext(ctx,
		`SELECT id, subject_id, topic_id, question, student_answer, predicted_score, max_score, predicted_grade, feedback, created_at
		 FROM answer_attempts WHERE id = ?`, id)
	var a models.AnswerAttempt
	var topicID sql.NullInt64
	if err := row.Scan(&a.ID, &a.SubjectID, &topicID, &a.Question, &a.StudentAnswer, &a.PredictedScore, &a.MaxScore, &a.PredictedGrade, &a.Feedback, &a.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if topicID.Valid {
		a.TopicID = &topicID.Int64
	}
	return &a, nil
}

func (r *Repository) ListAnswerAttemptsBySubject(ctx context.Context, subjectID int64, limit int) ([]models.AnswerAttempt, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := r.DB.QueryContext(ctx,
		`SELECT id, subject_id, topic_id, question, student_answer, predicted_score, max_score, predicted_grade, feedback, created_at
		 FROM answer_attempts
		 WHERE subject_id = ?
		 ORDER BY created_at DESC
		 LIMIT ?`, subjectID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.AnswerAttempt
	for rows.Next() {
		var a models.AnswerAttempt
		var topicID sql.NullInt64
		if err := rows.Scan(&a.ID, &a.SubjectID, &topicID, &a.Question, &a.StudentAnswer, &a.PredictedScore, &a.MaxScore, &a.PredictedGrade, &a.Feedback, &a.CreatedAt); err != nil {
			return nil, err
		}
		if topicID.Valid {
			a.TopicID = &topicID.Int64
		}
		list = append(list, a)
	}
	return list, rows.Err()
}

