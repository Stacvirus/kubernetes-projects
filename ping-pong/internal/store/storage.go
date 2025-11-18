package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Lines interface {
		Create(context.Context, *Line) error
		ReadLatest(context.Context) (*Line, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Lines: &LineStore{db},
	}
}
