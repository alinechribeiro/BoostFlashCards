package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/boostflashcards/backend/internal/auth"
	"github.com/boostflashcards/backend/internal/models"
)

type createEssayRequestBody struct {
	TutorID        int64   `json:"tutor_id"`
	Tier           string  `json:"tier"`
	BundleID       *int64  `json:"bundle_id,omitempty"`
	Subject        string  `json:"subject"`
	QuestionPrompt string  `json:"question_prompt"`
	StudentAnswer  string  `json:"student_answer"`
	AnswerFileURL  string  `json:"answer_file_url"`
}

// CreateEssayRequest lets a student create an essay review request.
// Stripe / bundle payment can be handled separately; this starts as pending_payment
// or attaches to an already-paid bundle.
func (h *Handlers) CreateEssayRequest(w http.ResponseWriter, r *http.Request) {
	if h.Auth == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "auth not configured"})
		return
	}
	user := auth.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	if user.Role != "student" {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "only students can create essay requests"})
		return
	}

	var body createEssayRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	if body.TutorID == 0 || body.QuestionPrompt == "" || strings.TrimSpace(body.Subject) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "tutor_id, subject and question_prompt are required"})
		return
	}

	var tier models.EssayTier
	switch body.Tier {
	case string(models.EssayTierQuick), "":
		tier = models.EssayTierQuick
	case string(models.EssayTierStandard):
		tier = models.EssayTierStandard
	case string(models.EssayTierPremium):
		tier = models.EssayTierPremium
	default:
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid tier"})
		return
	}

	status := models.EssayStatusPendingPayment
	var bundleID *int64
	if body.BundleID != nil {
		// If attached to a bundle, ensure it's active and has remaining essays
		b, err := h.Repo.GetEssayBundleByID(r.Context(), *body.BundleID)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid bundle id"})
			return
		}
		if b.Status != models.EssayBundleStatusActive || b.UsedEssays >= b.TotalEssays {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bundle is not active or exhausted"})
			return
		}
		bundleID = body.BundleID
		status = models.EssayStatusAwaitingTutor
		if err := h.Repo.IncrementBundleUsage(r.Context(), *bundleID); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to consume bundle credit"})
			return
		}
	}

	req := &models.EssayRequest{
		StudentID:     user.ID,
		TutorID:       body.TutorID,
		Tier:          tier,
		BundleID:      bundleID,
		Status:        status,
		Subject:       strings.TrimSpace(body.Subject),
		QuestionPrompt: body.QuestionPrompt,
		StudentAnswer:  body.StudentAnswer,
		AnswerFileURL:  body.AnswerFileURL,
	}
	created, err := h.Repo.CreateEssayRequest(r.Context(), req)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

type submitEssayReviewBody struct {
	Grade             string `json:"grade"`
	QuickComments     string `json:"quick_comments"`
	MarkSchemeRef     string `json:"mark_scheme_ref"`
	Strengths         string `json:"strengths"`
	Improvements      string `json:"improvements"`
	ImprovedParagraph string `json:"improved_paragraph"`
	AudioVideoURL     string `json:"audio_video_url"`
	ImprovementPlanURL string `json:"improvement_plan_url"`
}

// SubmitEssayReview lets the assigned tutor submit feedback for an essay request.
func (h *Handlers) SubmitEssayReview(w http.ResponseWriter, r *http.Request) {
	if h.Auth == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "auth not configured"})
		return
	}
	user := auth.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	if user.Role != "tutor" {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "only tutors can submit reviews"})
		return
	}

	id, err := idFromVars(r, "id")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	req, err := h.Repo.GetEssayRequestByID(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "essay request not found"})
		return
	}
	if req.TutorID != user.ID {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "not your essay request"})
		return
	}

	var body submitEssayReviewBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}

	now := time.Now()
	review := &models.EssayReview{
		EssayRequestID:     req.ID,
		Grade:              body.Grade,
		QuickComments:      body.QuickComments,
		MarkSchemeRef:      body.MarkSchemeRef,
		Strengths:          body.Strengths,
		Improvements:       body.Improvements,
		ImprovedParagraph:  body.ImprovedParagraph,
		AudioVideoURL:      body.AudioVideoURL,
		ImprovementPlanURL: body.ImprovementPlanURL,
		SubmittedAt:        &now,
	}
	created, err := h.Repo.CreateEssayReview(r.Context(), review)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	_ = h.Repo.UpdateEssayRequestStatus(r.Context(), req.ID, models.EssayStatusSubmitted)
	writeJSON(w, http.StatusCreated, created)
}

// MarkEssayViewed is called when the student opens the review; it updates timestamps and status.
func (h *Handlers) MarkEssayViewed(w http.ResponseWriter, r *http.Request) {
	if h.Auth == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "auth not configured"})
		return
	}
	user := auth.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	id, err := idFromVars(r, "id")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	req, err := h.Repo.GetEssayRequestByID(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "essay request not found"})
		return
	}
	if req.StudentID != user.ID {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "not your essay"})
		return
	}

	if err := h.Repo.MarkEssayReviewViewed(r.Context(), req.ID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to mark viewed"})
		return
	}
	if err := h.Repo.UpdateEssayRequestStatus(r.Context(), req.ID, models.EssayStatusViewedByStudent); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to update status"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// StudentEssays lists the current student's essay requests with tutor name and bundle/review info.
func (h *Handlers) StudentEssays(w http.ResponseWriter, r *http.Request) {
	if h.Auth == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "auth not configured"})
		return
	}
	user := auth.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	res, err := h.Repo.ListStudentEssays(r.Context(), user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, res)
}

// TutorEssays lists the current tutor's essay requests, to feed the tutor dashboard.
func (h *Handlers) TutorEssays(w http.ResponseWriter, r *http.Request) {
	if h.Auth == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "auth not configured"})
		return
	}
	user := auth.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	if user.Role != "tutor" {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "only tutors can view this dashboard"})
		return
	}
	res, err := h.Repo.ListTutorEssays(r.Context(), user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, res)
}

// GetEssayDetail returns the request + optional review so both student and tutor can drill into it.
func (h *Handlers) GetEssayDetail(w http.ResponseWriter, r *http.Request) {
	if h.Auth == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "auth not configured"})
		return
	}
	user := auth.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	id, err := idFromVars(r, "id")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	req, err := h.Repo.GetEssayRequestByID(r.Context(), id)
	if err != nil || req == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "essay request not found"})
		return
	}
	if user.ID != req.StudentID && user.ID != req.TutorID {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "forbidden"})
		return
	}
	// Try to load review (if exists)
	// We don't have a helper to get by request id, so a small query here:
	type response struct {
		Request *models.EssayRequest `json:"request"`
		Review  *models.EssayReview  `json:"review,omitempty"`
	}
	var rev *models.EssayReview
	row := h.Repo.DB.QueryRowContext(r.Context(), `
		SELECT id FROM essay_reviews WHERE essay_request_id = ? LIMIT 1
	`, req.ID)
	var reviewID int64
	if err := row.Scan(&reviewID); err == nil {
		if r2, err2 := h.Repo.GetEssayReviewByID(r.Context(), reviewID); err2 == nil {
			rev = r2
		}
	}
	writeJSON(w, http.StatusOK, response{Request: req, Review: rev})
}

