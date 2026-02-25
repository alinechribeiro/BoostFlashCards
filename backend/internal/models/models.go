package models

import "time"

type Subject struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	LastReviewedAt *time.Time `json:"last_reviewed_at,omitempty"`
}

type Topic struct {
	ID        int64     `json:"id"`
	SubjectID int64     `json:"subject_id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}

type FlashcardStatus string

const (
	FlashcardStatusNotYet    FlashcardStatus = "not_yet"
	FlashcardStatusConfident FlashcardStatus = "confident"
)

type Flashcard struct {
	ID             int64           `json:"id"`
	TopicID        int64           `json:"topic_id"`
	Front          string          `json:"front"`
	Back           string          `json:"back"`
	Status         FlashcardStatus `json:"status"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	LastReviewedAt *time.Time      `json:"last_reviewed_at,omitempty"`
	NextDueAt      *time.Time      `json:"next_due_at,omitempty"`
}

type CreateFlashcardRequest struct {
	TopicID int64  `json:"topic_id"`
	Front   string `json:"front"`
	Back    string `json:"back"`
}

type UpdateFlashcardRequest struct {
	Front *string `json:"front,omitempty"`
	Back  *string `json:"back,omitempty"`
}

type AnswerAttempt struct {
	ID             int64     `json:"id"`
	SubjectID      int64     `json:"subject_id"`
	TopicID        *int64    `json:"topic_id,omitempty"`
	Question       string    `json:"question"`
	StudentAnswer  string    `json:"student_answer"`
	PredictedScore float64   `json:"predicted_score"`
	MaxScore       float64   `json:"max_score"`
	PredictedGrade string    `json:"predicted_grade"`
	Feedback       string    `json:"feedback"`
	CreatedAt      time.Time `json:"created_at"`
}

type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // never expose
	Name         string    `json:"name"`
	Role         string    `json:"role"` // student, tutor, admin
	AvatarURL    string    `json:"avatar_url"`
	CreatedAt    time.Time `json:"created_at"`
}

type OAuthIdentity struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"user_id"`
	Provider       string    `json:"provider"`
	ProviderUserID string    `json:"provider_user_id"`
	CreatedAt      time.Time `json:"created_at"`
}

type TutorProfile struct {
	UserID                 int64     `json:"user_id"`
	Bio                    string    `json:"bio"`
	Headline               string    `json:"headline"`
	HourlyRateCents        int       `json:"hourly_rate_cents"`
	StripeConnectAccountID string    `json:"stripe_connect_account_id"`
	IsListed               bool      `json:"is_listed"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

type TutorSubject struct {
	ID          int64     `json:"id"`
	TutorUserID int64     `json:"tutor_user_id"`
	SubjectName string    `json:"subject_name"`
	Level       string    `json:"level"`
	CreatedAt   time.Time `json:"created_at"`
}

type EssayBundleStatus string

const (
	EssayBundleStatusPendingPayment EssayBundleStatus = "pending_payment"
	EssayBundleStatusActive         EssayBundleStatus = "active"
	EssayBundleStatusExhausted      EssayBundleStatus = "exhausted"
	EssayBundleStatusExpired        EssayBundleStatus = "expired"
)

type EssayBundle struct {
	ID                    int64             `json:"id"`
	StudentID             int64             `json:"student_id"`
	TutorID               *int64            `json:"tutor_id,omitempty"`
	TotalEssays           int               `json:"total_essays"`
	UsedEssays            int               `json:"used_essays"`
	PriceCents            int               `json:"price_cents"`
	Status                EssayBundleStatus `json:"status"`
	StripePaymentIntentID string            `json:"stripe_payment_intent_id"`
	CreatedAt             time.Time         `json:"created_at"`
	ExpiresAt             *time.Time        `json:"expires_at,omitempty"`
}

type EssayRequestStatus string

const (
	EssayStatusPendingPayment   EssayRequestStatus = "pending_payment"
	EssayStatusAwaitingTutor    EssayRequestStatus = "awaiting_tutor"
	EssayStatusInReview         EssayRequestStatus = "in_review"
	EssayStatusSubmitted        EssayRequestStatus = "submitted"
	EssayStatusViewedByStudent  EssayRequestStatus = "viewed_by_student"
	EssayStatusCompleted        EssayRequestStatus = "completed"
)

type EssayTier string

const (
	EssayTierQuick    EssayTier = "1_quick"
	EssayTierStandard EssayTier = "2_standard"
	EssayTierPremium  EssayTier = "3_premium"
)

type EssayRequest struct {
	ID                    int64             `json:"id"`
	StudentID             int64             `json:"student_id"`
	TutorID               int64             `json:"tutor_id"`
	Tier                  EssayTier         `json:"tier"`
	BundleID              *int64            `json:"bundle_id,omitempty"`
	Status                EssayRequestStatus `json:"status"`
	Subject               string            `json:"subject"`
	QuestionPrompt        string            `json:"question_prompt"`
	StudentAnswer         string            `json:"student_answer"`
	AnswerFileURL         string            `json:"answer_file_url"`
	StripePaymentIntentID string            `json:"stripe_payment_intent_id"`
	CreatedAt             time.Time         `json:"created_at"`
	UpdatedAt             time.Time         `json:"updated_at"`
}

type EssayReview struct {
	ID                 int64      `json:"id"`
	EssayRequestID     int64      `json:"essay_request_id"`
	Grade              string     `json:"grade"`
	QuickComments      string     `json:"quick_comments"`
	MarkSchemeRef      string     `json:"mark_scheme_ref"`
	Strengths          string     `json:"strengths"`
	Improvements       string     `json:"improvements"`
	ImprovedParagraph  string     `json:"improved_paragraph"`
	AudioVideoURL      string     `json:"audio_video_url"`
	ImprovementPlanURL string     `json:"improvement_plan_url"`
	SubmittedAt        *time.Time `json:"submitted_at,omitempty"`
	ViewedAt           *time.Time `json:"viewed_at,omitempty"`
}


