package server

import (
	"net/http"

	"github.com/boostflashcards/backend/internal/auth"
	"github.com/boostflashcards/backend/internal/handlers"
	"github.com/gorilla/mux"
)

func New(h *handlers.Handlers) http.Handler {
	r := mux.NewRouter()

	// Auth
	r.HandleFunc("/api/auth/signup", h.Signup).Methods("POST")
	r.HandleFunc("/api/auth/login", h.Login).Methods("POST")
	r.HandleFunc("/api/auth/logout", h.Logout).Methods("POST")
	r.HandleFunc("/api/auth/me", h.Me).Methods("GET")
	r.HandleFunc("/api/auth/complete-signup", h.CompleteSignup).Methods("POST")
	r.HandleFunc("/api/auth/{provider}/redirect", h.OAuthRedirect).Methods("GET")
	r.HandleFunc("/api/auth/callback", h.OAuthCallback).Methods("GET")

	r.HandleFunc("/api/subjects", h.ListSubjects).Methods("GET")
	r.HandleFunc("/api/subjects/{id}", h.GetSubject).Methods("GET")
	r.HandleFunc("/api/subjects/{subjectId}/topics", h.ListTopics).Methods("GET")
	r.HandleFunc("/api/subjects/{subjectId}/practice/question", h.GeneratePracticeQuestion).Methods("POST")
	r.HandleFunc("/api/subjects/{subjectId}/practice/answer", h.GradePracticeAnswer).Methods("POST")
	r.HandleFunc("/api/subjects/{subjectId}/progress", h.GetSubjectProgress).Methods("GET")
	r.HandleFunc("/api/subjects/{subjectId}/ai/flashcards", h.GenerateAIFlashcardsForSubject).Methods("POST")
	r.HandleFunc("/api/tutors", h.ListTutors).Methods("GET")
	r.HandleFunc("/api/tutors/{id}", h.GetTutor).Methods("GET")
	r.HandleFunc("/api/topics/{id}", h.GetTopic).Methods("GET")
	r.HandleFunc("/api/topics/{topicId}/flashcards", h.ListFlashcards).Methods("GET")
	r.HandleFunc("/api/topics/{topicId}/ai/insights", h.ExtractInsightsForTopic).Methods("POST")
	r.HandleFunc("/api/topics/{topicId}/ai/flashcards-from-insights", h.CreateFlashcardsFromInsights).Methods("POST")
	r.HandleFunc("/api/flashcards", h.CreateFlashcard).Methods("POST")
	r.HandleFunc("/api/flashcards/{id}", h.GetFlashcard).Methods("GET")
	r.HandleFunc("/api/flashcards/{id}", h.UpdateFlashcard).Methods("PUT")
	r.HandleFunc("/api/flashcards/{id}", h.DeleteFlashcard).Methods("DELETE")
	r.HandleFunc("/api/flashcards/{id}/status", h.UpdateFlashcardStatus).Methods("POST")
	r.HandleFunc("/api/ai/subjects", h.CreateSubjectWithAI).Methods("POST")

	// Billing / essays (require auth if configured)
	if h.Auth != nil {
		r.HandleFunc("/api/billing/connect",
			auth.RequireAuth(h.Auth.JWTSecret, h.Auth.CookieName, h.BillingConnect)).Methods("GET")
		r.HandleFunc("/api/billing/bundles",
			auth.RequireAuth(h.Auth.JWTSecret, h.Auth.CookieName, h.CreateBundle)).Methods("POST")

		r.HandleFunc("/api/essays/requests",
			auth.RequireAuth(h.Auth.JWTSecret, h.Auth.CookieName, h.CreateEssayRequest)).Methods("POST")
		r.HandleFunc("/api/essays/{id:[0-9]+}/review",
			auth.RequireAuth(h.Auth.JWTSecret, h.Auth.CookieName, h.SubmitEssayReview)).Methods("POST")
		r.HandleFunc("/api/essays/{id:[0-9]+}/mark_viewed",
			auth.RequireAuth(h.Auth.JWTSecret, h.Auth.CookieName, h.MarkEssayViewed)).Methods("POST")
		r.HandleFunc("/api/students/me/essays",
			auth.RequireAuth(h.Auth.JWTSecret, h.Auth.CookieName, h.StudentEssays)).Methods("GET")
		r.HandleFunc("/api/tutors/me/essays",
			auth.RequireAuth(h.Auth.JWTSecret, h.Auth.CookieName, h.TutorEssays)).Methods("GET")
		r.HandleFunc("/api/essays/{id:[0-9]+}",
			auth.RequireAuth(h.Auth.JWTSecret, h.Auth.CookieName, h.GetEssayDetail)).Methods("GET")
	}

	// Stripe webhook (no auth; secured by Stripe signature)
	r.HandleFunc("/api/billing/webhook", h.StripeWebhook).Methods("POST")

	r.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	return corsMiddleware(r)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
