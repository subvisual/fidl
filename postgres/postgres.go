package postgres

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type Config struct {
	Dsn          string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

type DB struct {
	*sqlx.DB
}

func Connect(cfg Config) *DB {
	dbx, err := sqlx.Connect("postgres", cfg.Dsn)

	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	db := &DB{DB: dbx}

	duration, err := time.ParseDuration(cfg.MaxIdleTime)
	if err != nil {
		log.Fatalf("Failed to parse duration: %v", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxIdleTime(duration)

	return db
}
