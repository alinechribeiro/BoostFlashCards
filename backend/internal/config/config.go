package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort   string
	DBUser       string
	DBPassword   string
	DBHost       string
	DBPort       string
	DBName       string
	OpenAIAPIKey string
	OpenAIModel  string
	// Auth
	JWTSecret        string
	AuthCookieName   string
	ServerURL        string // e.g. http://localhost:8080 for OAuth redirects
	GoogleClientID   string
	GoogleSecret     string
	FacebookClientID string
	FacebookSecret   string
	LinkedInClientID string
	LinkedInSecret   string
	InstagramClientID string
	InstagramSecret   string
	FrontendURL       string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	c := &Config{
		ServerPort:        getEnv("SERVER_PORT", "8080"),
		DBUser:            getEnv("DB_USER", "boostflash"),
		DBPassword:        getEnv("DB_PASSWORD", "boostflash"),
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "3306"),
		DBName:            getEnv("DB_NAME", "boostflashcards"),
		OpenAIAPIKey:      getEnv("OPENAI_API_KEY", ""),
		OpenAIModel:       getEnv("OPENAI_MODEL", "gpt-4o-mini"),
		JWTSecret:         getEnv("JWT_SECRET", ""),
		AuthCookieName:   getEnv("AUTH_COOKIE_NAME", "session"),
		ServerURL:         getEnv("SERVER_URL", "http://localhost:8080"),
		GoogleClientID:    getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleSecret:      getEnv("GOOGLE_CLIENT_SECRET", ""),
		FacebookClientID:  getEnv("FACEBOOK_CLIENT_ID", ""),
		FacebookSecret:    getEnv("FACEBOOK_CLIENT_SECRET", ""),
		LinkedInClientID:  getEnv("LINKEDIN_CLIENT_ID", ""),
		LinkedInSecret:    getEnv("LINKEDIN_CLIENT_SECRET", ""),
		InstagramClientID: getEnv("INSTAGRAM_CLIENT_ID", ""),
		InstagramSecret:   getEnv("INSTAGRAM_CLIENT_SECRET", ""),
		FrontendURL:       getEnv("FRONTEND_URL", "http://localhost:3000"),
	}
	return c, nil
}

func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
