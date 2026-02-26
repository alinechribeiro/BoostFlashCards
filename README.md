# BoostFlashCards – GCSE UK study app

BoostFlashCards project is aiming to help students revise for the **GCSE UK** exams using **flashcards** and AI‑assisted content.

- Backend: **Go (Golang)** + **MySQL**
- Frontend: **React + TypeScript (Vite)**

---

## Project structure

- `backend/` – Go REST API, MySQL
- `frontend/` – React SPA (Vite + TypeScript)
- `backend/migrations/` – MySQL schema and seed data (GCSE subjects/topics)

---

## Prerequisites

- Go **1.21+**
- Node.js **18+**
- MySQL **8+**

---

## Setting up MySQL

You can either create the database directly in MySQL, or run it in Docker.

### Option A – Local MySQL

Create the database and user (or adjust to your own credentials):

```bash
mysql -u root -e "
  CREATE DATABASE IF NOT EXISTS boostflashcards;
  CREATE USER IF NOT EXISTS 'boostflash'@'%' IDENTIFIED BY 'boostflash';
  GRANT ALL ON boostflashcards.* TO 'boostflash'@'%';
  FLUSH PRIVILEGES;
"
```

Apply the initial migration:

```bash
mysql -u boostflash -pboostflash boostflashcards < backend/migrations/001_schema.sql
```

You can then apply the rest of the migrations in order if needed.

### Option B – MySQL in Docker

```bash
docker run -d --name boostflash-mysql \
  -e MYSQL_ROOT_PASSWORD=root \
  -e MYSQL_DATABASE=boostflashcards \
  -e MYSQL_USER=boostflash \
  -e MYSQL_PASSWORD=boostflash \
  -p 3306:3306 \
  mysql:8
```

After the container is running, execute the migrations:

```bash
mysql -h 127.0.0.1 -P 3306 -u boostflash -pboostflash boostflashcards < backend/migrations/001_schema.sql
```

The script creates the tables and seeds GCSE subjects/topics.

---

## Backend (Go API)

From the project root:

```bash
cd backend
go mod download
```

### Environment variables

These are the main variables the backend understands (with defaults):

| Variable        | Default          |
|----------------|------------------|
| `SERVER_PORT`  | `8080`           |
| `DB_USER`      | `boostflash`     |
| `DB_PASSWORD`  | `boostflash`     |
| `DB_HOST`      | `localhost`      |
| `DB_PORT`      | `3306`           |
| `DB_NAME`      | `boostflashcards`|
| `OPENAI_API_KEY` | _(empty)_      |
| `OPENAI_MODEL` | `gpt-4o-mini`    |

Typical `.env` in `backend/`:

```env
SERVER_PORT=8080
DB_USER=boostflash
DB_PASSWORD=boostflash
DB_HOST=localhost
DB_PORT=3306
DB_NAME=boostflashcards
OPENAI_API_KEY=sk-...
OPENAI_MODEL=gpt-4o-mini
```

### Running the API

From `backend/`:

```bash
go run ./cmd/server
```

The API will be available at `http://localhost:8080`.

Some of the main endpoints:

- `GET /api/subjects` – list GCSE subjects
- `GET /api/subjects/:subjectId/topics` – topics for a subject
- `GET /api/topics/:topicId/flashcards` – flashcards for a topic
- `POST /api/flashcards` – create a flashcard (`topic_id`, `front`, `back`)
- `GET/PUT/DELETE /api/flashcards/:id`
- `GET /health` – simple health check
- `POST /api/subjects/{subjectId}/ai/flashcards` – generate topics + flashcards for a subject using OpenAI
- `POST /api/ai/subjects` – AI assistant to create a brand‑new subject with topics/flashcards

### Backend tests

Unit tests (no MySQL required):

```bash
cd backend
go test ./...
```

Integration tests (hit a real MySQL instance; optional):

```bash
cd backend
RUN_INTEGRATION=1 go test -tags=integration ./testutil/...
```

---

## Frontend (React + Vite)

From the project root:

```bash
cd frontend
npm install
npm run dev
```

The app runs at `http://localhost:3000`.  
Vite is configured to proxy `/api` and `/health` to `http://localhost:8080`, so you can run frontend and backend together without dealing with CORS in local development.

### Frontend tests

```bash
cd frontend
npm run test
```

---

## Test command summary

| Area                      | Command                                                    |
|---------------------------|------------------------------------------------------------|
| Backend (unit)            | `cd backend && go test ./...`                             |
| Backend (integration)     | `cd backend && RUN_INTEGRATION=1 go test -tags=integration ./testutil/...` |
| Frontend                  | `cd frontend && npm run test`                             |

---

## License / usage

This is a learning / demo project for educational purposes.  
Feel free to explore, tweak, and adapt it to your own GCSE or flashcard‑style experiments.
