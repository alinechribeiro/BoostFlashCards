package repository

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Repository struct {
	DB *sql.DB
}

var ErrNotFound = sql.ErrNoRows

func New(dsn string) (*Repository, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Repository{DB: db}, nil
}

func (r *Repository) Close() error {
	return r.DB.Close()
}
