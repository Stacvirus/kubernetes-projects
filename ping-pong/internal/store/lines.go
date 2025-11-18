package store

import (
	"context"
	"database/sql"
)

type Line struct {
	ID        int64  `json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type LineStore struct {
	db *sql.DB
}

func (s *LineStore) Create(ctx context.Context, line *Line) error {
	query := `
		INSERT INTO lines (content)
		VALUES ($1) RETURNING id, created_at
	`
	err := s.db.QueryRowContext(ctx, query, line.Content).Scan(&line.ID, &line.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *LineStore) ReadLatest(ctx context.Context) (*Line, error) {
	query := `
		SELECT * FROM lines ORDER BY created_at DESC LIMIT 1
	`
	line := &Line{}
	err := s.db.QueryRowContext(ctx, query).Scan(&line.ID, &line.Content, &line.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil // no lines yet
	}
	if err != nil {
		return nil, err
	}

	return line, nil
}
