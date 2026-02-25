package handlers

import (
	"net/http"

	"github.com/boostflashcards/backend/internal/repository"
)

type TutorSummary struct {
	ID           int64    `json:"id"`
	Name         string   `json:"name"`
	AvatarURL    string   `json:"avatar_url"`
	Headline     string   `json:"headline"`
	Bio          string   `json:"bio"`
	HourlyRate   int      `json:"hourly_rate_cents"`
	Subjects     []string `json:"subjects"`
}

type TutorDetail struct {
	ID           int64             `json:"id"`
	Name         string            `json:"name"`
	AvatarURL    string            `json:"avatar_url"`
	Headline     string            `json:"headline"`
	Bio          string            `json:"bio"`
	HourlyRate   int               `json:"hourly_rate_cents"`
	Subjects     []TutorDetailSubject `json:"subjects"`
}

type TutorDetailSubject struct {
	Name  string `json:"name"`
	Level string `json:"level"`
}

func (h *Handlers) ListTutors(w http.ResponseWriter, r *http.Request) {
	tutors, err := h.Repo.ListPublicTutors(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	out := make([]TutorSummary, 0, len(tutors))
	for _, t := range tutors {
		var subjects []string
		for _, s := range t.Subjects {
			if s.Level != "" {
				subjects = append(subjects, s.SubjectName+" ("+s.Level+")")
			} else {
				subjects = append(subjects, s.SubjectName)
			}
		}
		out = append(out, TutorSummary{
			ID:         t.User.ID,
			Name:       t.User.Name,
			AvatarURL:  t.User.AvatarURL,
			Headline:   t.Profile.Headline,
			Bio:        t.Profile.Bio,
			HourlyRate: t.Profile.HourlyRateCents,
			Subjects:   subjects,
		})
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *Handlers) GetTutor(w http.ResponseWriter, r *http.Request) {
	id, err := idFromVars(r, "id")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	t, err := h.Repo.GetTutorWithProfile(r.Context(), id)
	if err != nil {
		if err == repository.ErrNotFound {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "tutor not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	var subjects []TutorDetailSubject
	for _, s := range t.Subjects {
		subjects = append(subjects, TutorDetailSubject{
			Name:  s.SubjectName,
			Level: s.Level,
		})
	}

	out := TutorDetail{
		ID:         t.User.ID,
		Name:       t.User.Name,
		AvatarURL:  t.User.AvatarURL,
		Headline:   t.Profile.Headline,
		Bio:        t.Profile.Bio,
		HourlyRate: t.Profile.HourlyRateCents,
		Subjects:   subjects,
	}
	writeJSON(w, http.StatusOK, out)
}

