package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/boostflashcards/backend/internal/ai"
	"github.com/boostflashcards/backend/internal/models"
	"github.com/boostflashcards/backend/internal/repository"
	"github.com/gorilla/mux"
)

type Handlers struct {
	Repo *repository.Repository
	AI   *ai.Client
	Auth *AuthConfig
}

func New(repo *repository.Repository, aiClient *ai.Client, auth *AuthConfig) *Handlers {
	return &Handlers{Repo: repo, AI: aiClient, Auth: auth}
}

func (h *Handlers) ListSubjects(w http.ResponseWriter, r *http.Request) {
	subjects, err := h.Repo.ListSubjects(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, subjects)
}

func (h *Handlers) GetSubject(w http.ResponseWriter, r *http.Request) {
	id, err := idFromVars(r, "id")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	subject, err := h.Repo.GetSubjectByID(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if subject == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "subject not found"})
		return
	}
	writeJSON(w, http.StatusOK, subject)
}

func (h *Handlers) ListTopics(w http.ResponseWriter, r *http.Request) {
	subjectID, err := idFromVars(r, "subjectId")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid subject id"})
		return
	}
	topics, err := h.Repo.ListTopicsBySubjectID(r.Context(), subjectID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, topics)
}

func (h *Handlers) GetTopic(w http.ResponseWriter, r *http.Request) {
	id, err := idFromVars(r, "id")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	topic, err := h.Repo.GetTopicByID(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if topic == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "topic not found"})
		return
	}
	writeJSON(w, http.StatusOK, topic)
}

func (h *Handlers) ListFlashcards(w http.ResponseWriter, r *http.Request) {
	topicID, err := idFromVars(r, "topicId")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid topic id"})
		return
	}
	// Touch subject last_reviewed_at whenever a topic's flashcards are accessed.
	_, _ = h.Repo.DB.ExecContext(r.Context(), `
		UPDATE subjects
		   SET last_reviewed_at = NOW()
		 WHERE id = (SELECT subject_id FROM topics WHERE id = ?)`,
		topicID)
	cards, err := h.Repo.ListFlashcardsByTopicID(r.Context(), topicID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, cards)
}

func (h *Handlers) GetFlashcard(w http.ResponseWriter, r *http.Request) {
	id, err := idFromVars(r, "id")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	card, err := h.Repo.GetFlashcardByID(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if card == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "flashcard not found"})
		return
	}
	writeJSON(w, http.StatusOK, card)
}

func (h *Handlers) CreateFlashcard(w http.ResponseWriter, r *http.Request) {
	var req models.CreateFlashcardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.Front == "" || req.Back == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "front and back are required"})
		return
	}
	card, err := h.Repo.CreateFlashcard(r.Context(), req)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, card)
}

func (h *Handlers) UpdateFlashcard(w http.ResponseWriter, r *http.Request) {
	id, err := idFromVars(r, "id")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	var req models.UpdateFlashcardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	card, err := h.Repo.UpdateFlashcard(r.Context(), id, req.Front, req.Back)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if card == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "flashcard not found"})
		return
	}
	writeJSON(w, http.StatusOK, card)
}

func (h *Handlers) DeleteFlashcard(w http.ResponseWriter, r *http.Request) {
	id, err := idFromVars(r, "id")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	ok, err := h.Repo.DeleteFlashcard(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "flashcard not found"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type updateFlashcardStatusRequest struct {
	Status string `json:"status"`
}

// UpdateFlashcardStatus allows the student to mark a card as "not_yet" or "confident",
// updating spaced‑repetition scheduling metadata.
func (h *Handlers) UpdateFlashcardStatus(w http.ResponseWriter, r *http.Request) {
	id, err := idFromVars(r, "id")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	var body updateFlashcardStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	var status models.FlashcardStatus
	switch body.Status {
	case string(models.FlashcardStatusConfident):
		status = models.FlashcardStatusConfident
	default:
		status = models.FlashcardStatusNotYet
	}
	card, err := h.Repo.UpdateFlashcardStatus(r.Context(), id, status)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if card == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "flashcard not found"})
		return
	}
	writeJSON(w, http.StatusOK, card)
}

type generateAIFlashcardsRequest struct {
	NumCards int    `json:"num_cards"`
	Topic    string `json:"topic_name"`
}

type aiGenerationResponse struct {
	SubjectID         int64   `json:"subject_id"`
	TopicIDs          []int64 `json:"topic_ids"`
	FlashcardsCreated int     `json:"flashcards_created"`
	Message           string  `json:"message"`
}

type aiCreateSubjectRequest struct {
	Prompt string `json:"prompt"`
}

type aiCreateSubjectResponse struct {
	SubjectID         int64   `json:"subject_id"`
	TopicIDs          []int64 `json:"topic_ids"`
	FlashcardsCreated int     `json:"flashcards_created"`
	Message           string  `json:"message"`
}

type practiceQuestionResponse struct {
	SubjectID int64  `json:"subject_id"`
	Question  string `json:"question"`
}

type practiceGradeRequest struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
	TopicID  *int64 `json:"topic_id,omitempty"`
}

type practiceGradeResponse struct {
	AttemptID       int64   `json:"attempt_id"`
	SubjectID       int64   `json:"subject_id"`
	TopicID         *int64  `json:"topic_id,omitempty"`
	Question        string  `json:"question"`
	StudentAnswer   string  `json:"student_answer"`
	Score           float64 `json:"score"`
	MaxScore        float64 `json:"max_score"`
	Grade           string  `json:"grade"`
	GradeBand       string  `json:"grade_band,omitempty"`
	Feedback        string  `json:"feedback"`
	Strengths       string  `json:"strengths,omitempty"`
	Improvements    string  `json:"improvements,omitempty"`
	ScorePercentage float64 `json:"score_percentage"`
}

type subjectProgressAttempt struct {
	ID              int64   `json:"id"`
	CreatedAt       string  `json:"created_at"`
	ScorePercentage float64 `json:"score_percentage"`
	Grade           string  `json:"grade"`
}

type subjectProgressResponse struct {
	SubjectID          int64                    `json:"subject_id"`
	SubjectName        string                   `json:"subject_name"`
	Attempts           []subjectProgressAttempt `json:"attempts"`
	LatestGrade        string                   `json:"latest_grade"`
	AverageScore       float64                  `json:"average_score"`
	AttemptsCount      int                      `json:"attempts_count"`
	EncouragementBlurb string                   `json:"encouragement_blurb"`
}

type extractInsightsRequest struct {
	Text        string `json:"text"`
	MaxInsights int    `json:"max_insights"`
}

type extractInsightsResponse struct {
	Insights []ai.Insight `json:"insights"`
}

type insightFlashcardInput struct {
	Text  string `json:"text"`
	Style string `json:"style"` // "qa" or "true_false"
}

type generateFromInsightsRequest struct {
	Insights []insightFlashcardInput `json:"insights"`
}

type generateFromInsightsResponse struct {
	Created int `json:"created"`
}

func (h *Handlers) GenerateAIFlashcardsForSubject(w http.ResponseWriter, r *http.Request) {
	if h.AI == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "AI integration is not configured"})
		return
	}

	subjectID, err := idFromVars(r, "subjectId")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid subject id"})
		return
	}

	var body generateAIFlashcardsRequest
	_ = json.NewDecoder(r.Body).Decode(&body)
	if body.NumCards <= 0 {
		body.NumCards = 10
	}

	ctx := r.Context()
	subject, err := h.Repo.GetSubjectByID(ctx, subjectID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if subject == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "subject not found"})
		return
	}

	aiResp, err := h.AI.GenerateFlashcardsForSubject(ctx, subject.Name, body.NumCards)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if len(aiResp.Topics) == 0 {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "AI did not return any topics"})
		return
	}

	var topicIDs []int64
	flashcardsCreated := 0

	for _, t := range aiResp.Topics {
		topicName := strings.TrimSpace(t.Name)
		if topicName == "" {
			if body.Topic != "" {
				topicName = body.Topic
			} else {
				topicName = "AI Practice"
			}
		}
		if body.Topic != "" {
			topicName = body.Topic
		}

		topic, err := h.Repo.CreateTopic(ctx, subject.ID, topicName, slugify(topicName))
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		topicIDs = append(topicIDs, topic.ID)

		for _, c := range t.Flashcards {
			front := strings.TrimSpace(c.Front)
			back := strings.TrimSpace(c.Back)
			if front == "" || back == "" {
				continue
			}
			_, err := h.Repo.CreateFlashcard(ctx, models.CreateFlashcardRequest{
				TopicID: topic.ID,
				Front:   front,
				Back:    back,
			})
			if err != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}
			flashcardsCreated++
		}
	}

	msg := "AI generated flashcards"
	if flashcardsCreated > 0 {
		msg = "AI generated flashcards and saved them to new topics"
	}

	writeJSON(w, http.StatusCreated, aiGenerationResponse{
		SubjectID:         subject.ID,
		TopicIDs:          topicIDs,
		FlashcardsCreated: flashcardsCreated,
		Message:           msg,
	})
}

func (h *Handlers) CreateSubjectWithAI(w http.ResponseWriter, r *http.Request) {
	if h.AI == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "AI integration is not configured"})
		return
	}

	var body aiCreateSubjectRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if strings.TrimSpace(body.Prompt) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "prompt is required"})
		return
	}

	ctx := r.Context()
	aiResp, err := h.AI.GenerateSubjectFromDescription(ctx, body.Prompt)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	subjectName := strings.TrimSpace(aiResp.Subject.Name)
	if subjectName == "" {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "AI did not return a subject name"})
		return
	}

	// Reuse existing subject with the same name if it already exists.
	subject, err := h.Repo.GetSubjectByName(ctx, subjectName)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if subject == nil {
		subject, err = h.Repo.CreateSubject(ctx, subjectName, slugify(subjectName))
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	}

	var topicIDs []int64
	flashcardsCreated := 0

	for _, t := range aiResp.Subject.Topics {
		topicName := strings.TrimSpace(t.Name)
		if topicName == "" {
			continue
		}
		topic, err := h.Repo.CreateTopic(ctx, subject.ID, topicName, slugify(topicName))
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		topicIDs = append(topicIDs, topic.ID)

		for _, c := range t.Flashcards {
			front := strings.TrimSpace(c.Front)
			back := strings.TrimSpace(c.Back)
			if front == "" || back == "" {
				continue
			}
			_, err := h.Repo.CreateFlashcard(ctx, models.CreateFlashcardRequest{
				TopicID: topic.ID,
				Front:   front,
				Back:    back,
			})
			if err != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}
			flashcardsCreated++
		}
	}

	message := "AI created a new subject"
	if flashcardsCreated > 0 {
		message = "AI created a new subject, topics, and flashcards"
	}

	writeJSON(w, http.StatusCreated, aiCreateSubjectResponse{
		SubjectID:         subject.ID,
		TopicIDs:          topicIDs,
		FlashcardsCreated: flashcardsCreated,
		Message:           message,
	})
}

// Practice endpoints: generate question and grade answer with predicted GCSE grade.

func (h *Handlers) GeneratePracticeQuestion(w http.ResponseWriter, r *http.Request) {
	if h.AI == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "AI integration is not configured"})
		return
	}

	subjectID, err := idFromVars(r, "subjectId")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid subject id"})
		return
	}

	ctx := r.Context()
	subject, err := h.Repo.GetSubjectByID(ctx, subjectID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if subject == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "subject not found"})
		return
	}

	q, err := h.AI.GeneratePracticeQuestionForSubject(ctx, subject.Name)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, practiceQuestionResponse{
		SubjectID: subject.ID,
		Question:  q,
	})
}

func (h *Handlers) GradePracticeAnswer(w http.ResponseWriter, r *http.Request) {
	if h.AI == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "AI integration is not configured"})
		return
	}

	subjectID, err := idFromVars(r, "subjectId")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid subject id"})
		return
	}

	var body practiceGradeRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if strings.TrimSpace(body.Question) == "" || strings.TrimSpace(body.Answer) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "question and answer are required"})
		return
	}

	ctx := r.Context()
	subject, err := h.Repo.GetSubjectByID(ctx, subjectID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if subject == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "subject not found"})
		return
	}

	gr, err := h.AI.GradeAnswerForSubject(ctx, subject.Name, body.Question, body.Answer)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	attempt, err := h.Repo.CreateAnswerAttempt(ctx, models.AnswerAttempt{
		SubjectID:      subject.ID,
		TopicID:        body.TopicID,
		Question:       body.Question,
		StudentAnswer:  body.Answer,
		PredictedScore: gr.Score,
		MaxScore:       gr.MaxScore,
		PredictedGrade: strings.TrimSpace(gr.Grade),
		Feedback:       gr.Feedback,
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	scorePct := 0.0
	if gr.MaxScore > 0 {
		scorePct = (gr.Score / gr.MaxScore) * 100
	}

	writeJSON(w, http.StatusCreated, practiceGradeResponse{
		AttemptID:       attempt.ID,
		SubjectID:       attempt.SubjectID,
		TopicID:         attempt.TopicID,
		Question:        attempt.Question,
		StudentAnswer:   attempt.StudentAnswer,
		Score:           gr.Score,
		MaxScore:        gr.MaxScore,
		Grade:           strings.TrimSpace(gr.Grade),
		GradeBand:       gr.GradeBand,
		Feedback:        gr.Feedback,
		Strengths:       gr.StrengthsText(),
		Improvements:    gr.ImprovementsText(),
		ScorePercentage: scorePct,
	})
}

// GetSubjectProgress aggregates past AI-marked attempts into a simple
// progress view for graphs and predicted grade per subject.
func (h *Handlers) GetSubjectProgress(w http.ResponseWriter, r *http.Request) {
	subjectID, err := idFromVars(r, "subjectId")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid subject id"})
		return
	}

	ctx := r.Context()
	subject, err := h.Repo.GetSubjectByID(ctx, subjectID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if subject == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "subject not found"})
		return
	}

	attempts, err := h.Repo.ListAnswerAttemptsBySubject(ctx, subjectID, 100)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if len(attempts) == 0 {
		writeJSON(w, http.StatusOK, subjectProgressResponse{
			SubjectID:          subject.ID,
			SubjectName:        subject.Name,
			Attempts:           []subjectProgressAttempt{},
			LatestGrade:        "",
			AverageScore:       0,
			AttemptsCount:      0,
			EncouragementBlurb: "Answer your first AI-marked question to start building your GCSE prediction.",
		})
		return
	}

	// attempts are newest-first; reverse for chronological order.
	for i, j := 0, len(attempts)-1; i < j; i, j = i+1, j-1 {
		attempts[i], attempts[j] = attempts[j], attempts[i]
	}

	var totalPct float64
	outAttempts := make([]subjectProgressAttempt, 0, len(attempts))
	for _, a := range attempts {
		pct := 0.0
		if a.MaxScore > 0 {
			pct = (a.PredictedScore / a.MaxScore) * 100
		}
		totalPct += pct
		outAttempts = append(outAttempts, subjectProgressAttempt{
			ID:              a.ID,
			CreatedAt:       a.CreatedAt.Format(time.RFC3339),
			ScorePercentage: pct,
			Grade:           a.PredictedGrade,
		})
	}

	avg := totalPct / float64(len(attempts))
	latest := attempts[len(attempts)-1]

	var blurb string
	if avg >= 80 {
		blurb = "You are tracking at a high GCSE grade here. Keep practising to lock this in."
	} else if avg >= 60 {
		blurb = "You are building a solid foundation. A bit more focused practice can push you into top grades."
	} else if avg >= 40 {
		blurb = "You’re on the journey – focus on the feedback and you’ll see this graph climb."
	} else {
		blurb = "Every attempt nudges your grade up. Keep going – progress comes from consistent practice."
	}

	writeJSON(w, http.StatusOK, subjectProgressResponse{
		SubjectID:          subject.ID,
		SubjectName:        subject.Name,
		Attempts:           outAttempts,
		LatestGrade:        latest.PredictedGrade,
		AverageScore:       avg,
		AttemptsCount:      len(attempts),
		EncouragementBlurb: blurb,
	})
}

// ExtractInsightsForTopic lets the user paste text and get back AI-extracted insights
// for a given topic/subject, without saving any flashcards yet.
func (h *Handlers) ExtractInsightsForTopic(w http.ResponseWriter, r *http.Request) {
	if h.AI == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "AI integration is not configured"})
		return
	}

	topicID, err := idFromVars(r, "topicId")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid topic id"})
		return
	}

	var body extractInsightsRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if strings.TrimSpace(body.Text) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "text is required"})
		return
	}

	ctx := r.Context()
	topic, err := h.Repo.GetTopicByID(ctx, topicID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if topic == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "topic not found"})
		return
	}

	subject, err := h.Repo.GetSubjectByID(ctx, topic.SubjectID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if subject == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "subject not found"})
		return
	}

	insights, err := h.AI.ExtractInsightsFromText(ctx, subject.Name, body.Text, body.MaxInsights)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, extractInsightsResponse{
		Insights: insights,
	})
}

// CreateFlashcardsFromInsights turns user-approved insights into saved flashcards
// for a given topic.
func (h *Handlers) CreateFlashcardsFromInsights(w http.ResponseWriter, r *http.Request) {
	if h.AI == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "AI integration is not configured"})
		return
	}

	topicID, err := idFromVars(r, "topicId")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid topic id"})
		return
	}

	var body generateFromInsightsRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if len(body.Insights) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "no insights provided"})
		return
	}

	ctx := r.Context()
	topic, err := h.Repo.GetTopicByID(ctx, topicID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if topic == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "topic not found"})
		return
	}

	subject, err := h.Repo.GetSubjectByID(ctx, topic.SubjectID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if subject == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "subject not found"})
		return
	}

	specs := make([]ai.InsightFlashcardSpec, 0, len(body.Insights))
	for _, in := range body.Insights {
		text := strings.TrimSpace(in.Text)
		if text == "" {
			continue
		}
		style := strings.ToLower(strings.TrimSpace(in.Style))
		if style != "qa" && style != "true_false" {
			style = "qa"
		}
		specs = append(specs, ai.InsightFlashcardSpec{
			Text:  text,
			Style: style,
		})
	}
	if len(specs) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "no valid insights provided"})
		return
	}

	generated, err := h.AI.GenerateFlashcardsFromInsights(ctx, subject.Name, specs)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	created := 0
	for _, fc := range generated {
		front := strings.TrimSpace(fc.Front)
		back := strings.TrimSpace(fc.Back)
		if front == "" || back == "" {
			continue
		}
		_, err := h.Repo.CreateFlashcard(ctx, models.CreateFlashcardRequest{
			TopicID: topic.ID,
			Front:   front,
			Back:    back,
		})
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		created++
	}

	writeJSON(w, http.StatusCreated, generateFromInsightsResponse{
		Created: created,
	})
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			b.WriteRune(r)
		}
	}
	out := b.String()
	if out == "" {
		return "subject"
	}
	return out
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func idFromVars(r *http.Request, key string) (int64, error) {
	vars := mux.Vars(r)
	s := vars[key]
	if s == "" {
		return 0, strconv.ErrSyntax
	}
	return strconv.ParseInt(s, 10, 64)
}
