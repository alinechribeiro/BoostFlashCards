package repository

import (
	"context"

	"github.com/boostflashcards/backend/internal/models"
)

type TutorWithProfile struct {
	User     models.User
	Profile  models.TutorProfile
	Subjects []models.TutorSubject
}

func (r *Repository) GetTutorWithProfile(ctx context.Context, userID int64) (*TutorWithProfile, error) {
	user, err := r.GetUserByID(ctx, userID)
	if err != nil || user == nil {
		return nil, err
	}

	var p models.TutorProfile
	err = r.DB.QueryRowContext(ctx,
		`SELECT user_id, bio, headline, hourly_rate_cents, stripe_connect_account_id, is_listed, created_at, updated_at
		 FROM tutor_profiles WHERE user_id = ?`,
		userID).
		Scan(&p.UserID, &p.Bio, &p.Headline, &p.HourlyRateCents, &p.StripeConnectAccountID, &p.IsListed, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}

	rows, err := r.DB.QueryContext(ctx,
		`SELECT id, tutor_user_id, subject_name, level, created_at
		 FROM tutor_subjects WHERE tutor_user_id = ? ORDER BY subject_name`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subjects []models.TutorSubject
	for rows.Next() {
		var s models.TutorSubject
		if err := rows.Scan(&s.ID, &s.TutorUserID, &s.SubjectName, &s.Level, &s.CreatedAt); err != nil {
			return nil, err
		}
		subjects = append(subjects, s)
	}

	return &TutorWithProfile{
		User:     *user,
		Profile:  p,
		Subjects: subjects,
	}, rows.Err()
}

func (r *Repository) ListPublicTutors(ctx context.Context) ([]TutorWithProfile, error) {
	rows, err := r.DB.QueryContext(ctx,
		`SELECT u.id, u.email, u.name, u.role, u.avatar_url, u.created_at,
		        p.bio, p.headline, p.hourly_rate_cents, p.stripe_connect_account_id, p.is_listed, p.created_at, p.updated_at
		 FROM users u
		 JOIN tutor_profiles p ON p.user_id = u.id
		 WHERE p.is_listed = 1
		 ORDER BY p.headline, u.name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tutors []TutorWithProfile
	for rows.Next() {
		var tw TutorWithProfile
		if err := rows.Scan(
			&tw.User.ID, &tw.User.Email, &tw.User.Name, &tw.User.Role, &tw.User.AvatarURL, &tw.User.CreatedAt,
			&tw.Profile.Bio, &tw.Profile.Headline, &tw.Profile.HourlyRateCents, &tw.Profile.StripeConnectAccountID, &tw.Profile.IsListed,
			&tw.Profile.CreatedAt, &tw.Profile.UpdatedAt,
		); err != nil {
			return nil, err
		}
		tw.Profile.UserID = tw.User.ID

		// Fetch subjects for each tutor
		subRows, err := r.DB.QueryContext(ctx,
			`SELECT id, tutor_user_id, subject_name, level, created_at
			 FROM tutor_subjects WHERE tutor_user_id = ? ORDER BY subject_name`, tw.User.ID)
		if err != nil {
			return nil, err
		}
		var subjects []models.TutorSubject
		for subRows.Next() {
			var s models.TutorSubject
			if err := subRows.Scan(&s.ID, &s.TutorUserID, &s.SubjectName, &s.Level, &s.CreatedAt); err != nil {
				subRows.Close()
				return nil, err
			}
			subjects = append(subjects, s)
		}
		subRows.Close()
		tw.Subjects = subjects

		tutors = append(tutors, tw)
	}
	return tutors, rows.Err()
}

func (r *Repository) UpdateTutorStripeAccountID(ctx context.Context, userID int64, accountID string) error {
	_, err := r.DB.ExecContext(ctx, `
		UPDATE tutor_profiles
		   SET stripe_connect_account_id = ?
		 WHERE user_id = ?
	`, accountID, userID)
	return err
}


