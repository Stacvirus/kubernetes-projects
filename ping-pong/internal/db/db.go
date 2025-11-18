package db

import (
	"context"
	"database/sql"
	"ping-pong/internal/config"
	"time"

	_ "github.com/lib/pq"
)

func New(config *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", config.Addr)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetMaxOpenConns(config.MaxOpenConns)

	duration, err := time.ParseDuration(config.MaxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
