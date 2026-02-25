module github.com/boostflashcards/backend

go 1.24.0

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/go-sql-driver/mysql v1.7.1
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/gorilla/mux v1.8.1
	github.com/joho/godotenv v1.5.1
	github.com/sashabaranov/go-openai v1.27.0
	github.com/stripe/stripe-go/v82 v82.0.0
	golang.org/x/crypto v0.48.0
	golang.org/x/oauth2 v0.35.0
)

require cloud.google.com/go/compute/metadata v0.3.0 // indirect
