# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download
COPY . .
RUN go mod tidy && CGO_ENABLED=0 go build -o /server ./cmd/server

# Run stage
FROM alpine:3.19

RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /server .

EXPOSE 8080

ENV SERVER_PORT=8080
ENV DB_USER=boostflash
ENV DB_PASSWORD=boostflash
ENV DB_HOST=mysql
ENV DB_PORT=3306
ENV DB_NAME=boostflashcards

CMD ["./server"]
