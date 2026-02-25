package main

import (
	"log"
	"net/http"

	"github.com/boostflashcards/backend/internal/auth"
	"github.com/boostflashcards/backend/internal/config"
	"github.com/boostflashcards/backend/internal/ai"
	"github.com/boostflashcards/backend/internal/handlers"
	"github.com/boostflashcards/backend/internal/repository"
	"github.com/boostflashcards/backend/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	repo, err := repository.New(cfg.DSN())
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer repo.Close()

	var aiClient *ai.Client
	if cfg.OpenAIAPIKey != "" {
		aiClient = ai.New(cfg.OpenAIAPIKey, cfg.OpenAIModel)
	} else {
		log.Println("warning: OPENAI_API_KEY is not set; AI features will be disabled")
	}

	var authConfig *handlers.AuthConfig
	if cfg.JWTSecret != "" {
		authConfig = &handlers.AuthConfig{
			JWTSecret:    cfg.JWTSecret,
			CookieName:   cfg.AuthCookieName,
			ServerURL:    cfg.ServerURL,
			FrontendURL:  cfg.FrontendURL,
			OAuth: auth.NewOAuthConfig(
				cfg.ServerURL,
				cfg.GoogleClientID, cfg.GoogleSecret,
				cfg.FacebookClientID, cfg.FacebookSecret,
				cfg.LinkedInClientID, cfg.LinkedInSecret,
				cfg.InstagramClientID, cfg.InstagramSecret,
			),
		}
	} else {
		log.Println("warning: JWT_SECRET is not set; auth (signup, login, OAuth) will be disabled")
	}

	h := handlers.New(repo, aiClient, authConfig)
	srv := server.New(h)

	addr := ":" + cfg.ServerPort
	log.Printf("BoostFlashCards API listening on %s", addr)
	if err := http.ListenAndServe(addr, srv); err != nil {
		log.Fatalf("server: %v", err)
	}
}
