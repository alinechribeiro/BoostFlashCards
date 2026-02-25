package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/boostflashcards/backend/internal/models"
)

type EssayTierPricing struct {
	StudentPriceCents int
	TutorShareCents   int
	PlatformCutCents  int
}

var TierPricing = map[models.EssayTier]EssayTierPricing{
	models.EssayTierQuick: {
		StudentPriceCents: 100,  // £1.00
		TutorShareCents:   75,   // £0.75
		PlatformCutCents:  25,   // £0.25
	},
	models.EssayTierStandard: {
		StudentPriceCents: 350,  // £3.50
		TutorShareCents:   300,  // £3.00
		PlatformCutCents:  50,   // £0.50
	},
	models.EssayTierPremium: {
		StudentPriceCents: 1000, // £10.00
		TutorShareCents:   900,  // £9.00
		PlatformCutCents:  100,  // £1.00
	},
}

// MinimumBundleCents is the minimum Stripe charge for bundles (e.g. £20).
const MinimumBundleCents = 2000

func (r *Repository) CreateEssayBundle(ctx context.Context, b *models.EssayBundle) (*models.EssayBundle, error) {
	res, err := r.DB.ExecContext(ctx, `
		INSERT INTO essay_bundles (
			student_id, tutor_id, total_essays, used_essays, price_cents,
			status, stripe_payment_intent_id, created_at, expires_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), ?)
	`, b.StudentID, b.TutorID, b.TotalEssays, b.UsedEssays, b.PriceCents,
		b.Status, b.StripePaymentIntentID, b.ExpiresAt)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return r.GetEssayBundleByID(ctx, id)
}

func (r *Repository) GetEssayBundleByID(ctx context.Context, id int64) (*models.EssayBundle, error) {
	var b models.EssayBundle
	var tutorID sql.NullInt64
	var expiresAt sql.NullTime
	err := r.DB.QueryRowContext(ctx, `
		SELECT id, student_id, tutor_id, total_essays, used_essays, price_cents,
		       status, stripe_payment_intent_id, created_at, expires_at
		  FROM essay_bundles
		 WHERE id = ?
	`, id).Scan(
		&b.ID,
		&b.StudentID,
		&tutorID,
		&b.TotalEssays,
		&b.UsedEssays,
		&b.PriceCents,
		&b.Status,
		&b.StripePaymentIntentID,
		&b.CreatedAt,
		&expiresAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if tutorID.Valid {
		id := tutorID.Int64
		b.TutorID = &id
	}
	if expiresAt.Valid {
		t := expiresAt.Time
		b.ExpiresAt = &t
	}
	return &b, nil
}

func (r *Repository) IncrementBundleUsage(ctx context.Context, bundleID int64) error {
	_, err := r.DB.ExecContext(ctx, `
		UPDATE essay_bundles
		   SET used_essays = used_essays + 1,
		       status = CASE
		         WHEN used_essays + 1 >= total_essays THEN 'exhausted'
		         ELSE status
		       END
		 WHERE id = ?
	`, bundleID)
	return err
}

func (r *Repository) GetEssayBundleByPaymentIntent(ctx context.Context, paymentIntentID string) (*models.EssayBundle, error) {
	var b models.EssayBundle
	var tutorID sql.NullInt64
	var expiresAt sql.NullTime
	err := r.DB.QueryRowContext(ctx, `
		SELECT id, student_id, tutor_id, total_essays, used_essays, price_cents,
		       status, stripe_payment_intent_id, created_at, expires_at
		  FROM essay_bundles
		 WHERE stripe_payment_intent_id = ?
	`, paymentIntentID).Scan(
		&b.ID,
		&b.StudentID,
		&tutorID,
		&b.TotalEssays,
		&b.UsedEssays,
		&b.PriceCents,
		&b.Status,
		&b.StripePaymentIntentID,
		&b.CreatedAt,
		&expiresAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if tutorID.Valid {
		id := tutorID.Int64
		b.TutorID = &id
	}
	if expiresAt.Valid {
		t := expiresAt.Time
		b.ExpiresAt = &t
	}
	return &b, nil
}

func (r *Repository) UpdateEssayBundleStatus(ctx context.Context, id int64, status models.EssayBundleStatus) error {
	_, err := r.DB.ExecContext(ctx, `
		UPDATE essay_bundles
		   SET status = ?
		 WHERE id = ?
	`, status, id)
	return err
}

func (r *Repository) SetRequestsAwaitingTutorForBundle(ctx context.Context, bundleID int64) error {
	_, err := r.DB.ExecContext(ctx, `
		UPDATE essay_requests
		   SET status = 'awaiting_tutor', updated_at = NOW()
		 WHERE bundle_id = ? AND status = 'pending_payment'
	`, bundleID)
	return err
}

type StudentEssaySummary struct {
	ID          int64              `json:"id"`
	TutorName   string             `json:"tutor_name"`
	Tier        models.EssayTier   `json:"tier"`
	Status      models.EssayRequestStatus `json:"status"`
	Subject     string             `json:"subject"`
	CreatedAt   time.Time          `json:"created_at"`
	BundleID    *int64             `json:"bundle_id,omitempty"`
	HasReview   bool               `json:"has_review"`
	ReviewID    *int64             `json:"review_id,omitempty"`
}

func (r *Repository) ListStudentEssays(ctx context.Context, studentID int64) ([]StudentEssaySummary, error) {
	rows, err := r.DB.QueryContext(ctx, `
		SELECT er.id, u.name, er.tier, er.status, er.subject,
		       er.created_at, er.bundle_id,
		       rev.id
		  FROM essay_requests er
		  JOIN users u ON u.id = er.tutor_id
		  LEFT JOIN essay_reviews rev ON rev.essay_request_id = er.id
		 WHERE er.student_id = ?
		 ORDER BY er.created_at DESC
	`, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []StudentEssaySummary
	for rows.Next() {
		var s StudentEssaySummary
		var bundleID sql.NullInt64
		var reviewID sql.NullInt64
		if err := rows.Scan(
			&s.ID,
			&s.TutorName,
			&s.Tier,
			&s.Status,
			&s.Subject,
			&s.CreatedAt,
			&bundleID,
			&reviewID,
		); err != nil {
			return nil, err
		}
		if bundleID.Valid {
			id := bundleID.Int64
			s.BundleID = &id
		}
		if reviewID.Valid {
			id := reviewID.Int64
			s.ReviewID = &id
			s.HasReview = true
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

type TutorEssaySummary struct {
	ID          int64              `json:"id"`
	StudentName string             `json:"student_name"`
	Subject     string             `json:"subject"`
	Tier        models.EssayTier   `json:"tier"`
	Status      models.EssayRequestStatus `json:"status"`
	CreatedAt   time.Time          `json:"created_at"`
	ViewedAt    *time.Time         `json:"viewed_at,omitempty"`
}

func (r *Repository) ListTutorEssays(ctx context.Context, tutorID int64) ([]TutorEssaySummary, error) {
	rows, err := r.DB.QueryContext(ctx, `
		SELECT er.id, u.name, er.subject, er.tier, er.status, er.created_at, rev.viewed_at
		  FROM essay_requests er
		  JOIN users u ON u.id = er.student_id
		  LEFT JOIN essay_reviews rev ON rev.essay_request_id = er.id
		 WHERE er.tutor_id = ?
		 ORDER BY er.created_at DESC
	`, tutorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []TutorEssaySummary
	for rows.Next() {
		var s TutorEssaySummary
		var viewedAt sql.NullTime
		if err := rows.Scan(
			&s.ID,
			&s.StudentName,
			&s.Subject,
			&s.Tier,
			&s.Status,
			&s.CreatedAt,
			&viewedAt,
		); err != nil {
			return nil, err
		}
		if viewedAt.Valid {
			t := viewedAt.Time
			s.ViewedAt = &t
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func (r *Repository) CreateEssayRequest(ctx context.Context, req *models.EssayRequest) (*models.EssayRequest, error) {
	var bundleID interface{}
	if req.BundleID != nil {
		bundleID = *req.BundleID
	}
	res, err := r.DB.ExecContext(ctx, `
		INSERT INTO essay_requests (
			student_id, tutor_id, tier, bundle_id, status, subject,
			question_prompt, student_answer, answer_file_url,
			stripe_payment_intent_id, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`, req.StudentID, req.TutorID, req.Tier, bundleID, req.Status, req.Subject,
		req.QuestionPrompt, req.StudentAnswer, req.AnswerFileURL, req.StripePaymentIntentID)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return r.GetEssayRequestByID(ctx, id)
}

func (r *Repository) GetEssayRequestByID(ctx context.Context, id int64) (*models.EssayRequest, error) {
	var req models.EssayRequest
	var bundleID sql.NullInt64
	err := r.DB.QueryRowContext(ctx, `
		SELECT id, student_id, tutor_id, tier, bundle_id, status, subject,
		       question_prompt, student_answer, answer_file_url,
		       stripe_payment_intent_id, created_at, updated_at
		  FROM essay_requests
		 WHERE id = ?
	`, id).Scan(
		&req.ID,
		&req.StudentID,
		&req.TutorID,
		&req.Tier,
		&bundleID,
		&req.Status,
		&req.Subject,
		&req.QuestionPrompt,
		&req.StudentAnswer,
		&req.AnswerFileURL,
		&req.StripePaymentIntentID,
		&req.CreatedAt,
		&req.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if bundleID.Valid {
		id := bundleID.Int64
		req.BundleID = &id
	}
	return &req, nil
}

func (r *Repository) UpdateEssayRequestStatus(ctx context.Context, id int64, status models.EssayRequestStatus) error {
	_, err := r.DB.ExecContext(ctx, `
		UPDATE essay_requests
		   SET status = ?, updated_at = NOW()
		 WHERE id = ?
	`, status, id)
	return err
}

func (r *Repository) CreateEssayReview(ctx context.Context, rev *models.EssayReview) (*models.EssayReview, error) {
	var submittedAt, viewedAt interface{}
	if rev.SubmittedAt != nil {
		submittedAt = *rev.SubmittedAt
	}
	if rev.ViewedAt != nil {
		viewedAt = *rev.ViewedAt
	}
	res, err := r.DB.ExecContext(ctx, `
		INSERT INTO essay_reviews (
			essay_request_id, grade, quick_comments, mark_scheme_ref,
			strengths, improvements, improved_paragraph,
			audio_video_url, improvement_plan_url,
			submitted_at, viewed_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, rev.EssayRequestID, rev.Grade, rev.QuickComments, rev.MarkSchemeRef,
		rev.Strengths, rev.Improvements, rev.ImprovedParagraph,
		rev.AudioVideoURL, rev.ImprovementPlanURL,
		submittedAt, viewedAt,
	)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return r.GetEssayReviewByID(ctx, id)
}

func (r *Repository) GetEssayReviewByID(ctx context.Context, id int64) (*models.EssayReview, error) {
	var rev models.EssayReview
	var submittedAt, viewedAt sql.NullTime
	err := r.DB.QueryRowContext(ctx, `
		SELECT id, essay_request_id, grade, quick_comments, mark_scheme_ref,
		       strengths, improvements, improved_paragraph,
		       audio_video_url, improvement_plan_url,
		       submitted_at, viewed_at
		  FROM essay_reviews
		 WHERE id = ?
	`, id).Scan(
		&rev.ID,
		&rev.EssayRequestID,
		&rev.Grade,
		&rev.QuickComments,
		&rev.MarkSchemeRef,
		&rev.Strengths,
		&rev.Improvements,
		&rev.ImprovedParagraph,
		&rev.AudioVideoURL,
		&rev.ImprovementPlanURL,
		&submittedAt,
		&viewedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if submittedAt.Valid {
		t := submittedAt.Time
		rev.SubmittedAt = &t
	}
	if viewedAt.Valid {
		t := viewedAt.Time
		rev.ViewedAt = &t
	}
	return &rev, nil
}

func (r *Repository) MarkEssayReviewViewed(ctx context.Context, essayRequestID int64) error {
	now := time.Now()
	_, err := r.DB.ExecContext(ctx, `
		UPDATE essay_reviews
		   SET viewed_at = ?
		 WHERE essay_request_id = ?
	`, now, essayRequestID)
	return err
}

