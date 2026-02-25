# BoostFlashCards – GCSE UK

Projeto de educação para geração e estudo de **flashcards** para o **GCSE UK** (General Certificate of Secondary Education).  
Backend em **Golang**, base de dados **MySQL**, frontend **React** (TypeScript).

## Estrutura

- **backend/** – API REST em Go (Gorilla Mux, MySQL)
- **frontend/** – SPA em React (Vite, TypeScript)
- **backend/migrations/** – Schema e seeds MySQL (subjects/topics GCSE)

## Pré-requisitos

- Go 1.21+
- Node.js 18+
- MySQL 8+ (ou Docker)

## Base de dados MySQL

1. Criar utilizador e base (ou usar root):

```bash
mysql -u root -e "
  CREATE DATABASE IF NOT EXISTS boostflashcards;
  CREATE USER IF NOT EXISTS 'boostflash'@'%' IDENTIFIED BY 'boostflash';
  GRANT ALL ON boostflashcards.* TO 'boostflash'@'%';
  FLUSH PRIVILEGES;
"
```

2. Aplicar migrations:

```bash
mysql -u boostflash -pboostflash boostflashcards < backend/migrations/001_schema.sql
```

### Docker (MySQL apenas)

```bash
docker run -d --name boostflash-mysql \
  -e MYSQL_ROOT_PASSWORD=root \
  -e MYSQL_DATABASE=boostflashcards \
  -e MYSQL_USER=boostflash \
  -e MYSQL_PASSWORD=boostflash \
  -p 3306:3306 \
  mysql:8
```

Depois executar o SQL em `backend/migrations/001_schema.sql` (a base e o user já existem; as tabelas e seeds são criados pelo script).

## Backend (Go)

```bash
cd backend
go mod download
```

Variáveis de ambiente (opcional; valores por defeito em baixo):

| Variável      | Default        |
|---------------|----------------|
| SERVER_PORT   | 8080           |
| DB_USER       | boostflash     |
| DB_PASSWORD   | boostflash     |
| DB_HOST       | localhost      |
| DB_PORT       | 3306           |
| DB_NAME       | boostflashcards|
| OPENAI_API_KEY| _(vazio)_      |
| OPENAI_MODEL  | gpt-4o-mini    |

Ficheiro `.env` na pasta `backend` (exemplo):

```
SERVER_PORT=8080
DB_USER=boostflash
DB_PASSWORD=boostflash
DB_HOST=localhost
DB_PORT=3306
DB_NAME=boostflashcards
OPENAI_API_KEY=sk-...
OPENAI_MODEL=gpt-4o-mini
```

Executar API:

```bash
go run ./cmd/server
```

A API fica em `http://localhost:8080`. Endpoints principais:

- `GET /api/subjects` – lista de disciplinas GCSE
- `GET /api/subjects/:subjectId/topics` – tópicos da disciplina
- `GET /api/topics/:topicId/flashcards` – flashcards do tópico
- `POST /api/flashcards` – criar flashcard (body: `topic_id`, `front`, `back`)
- `GET/PUT/DELETE /api/flashcards/:id`
- `GET /health` – health check
– `POST /api/subjects/{subjectId}/ai/flashcards` – gera tópicos e flashcards para uma disciplina usando OpenAI
– `POST /api/ai/subjects` – chatbot para criar uma nova disciplina (e tópicos/flashcards) com OpenAI

### Testes (backend)

Testes unitários (sem MySQL):

```bash
cd backend
go test ./...
```

Testes de integração (com MySQL; opcional):

```bash
RUN_INTEGRATION=1 go test -tags=integration ./testutil/...
```

## Frontend (React)

```bash
cd frontend
npm install
npm run dev
```

A app corre em `http://localhost:3000` e o Vite faz proxy de `/api` e `/health` para `http://localhost:8080`.

### Testes (frontend)

```bash
cd frontend
npm run test
```

## Resumo de testes

| Onde        | Comando |
|------------|---------|
| Backend    | `cd backend && go test ./...` |
| Backend (integração) | `RUN_INTEGRATION=1 go test -tags=integration ./testutil/...` |
| Frontend   | `cd frontend && npm run test` |

## Licença

Uso educacional / projeto de demonstração.
