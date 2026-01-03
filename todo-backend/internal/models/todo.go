package models

import "time"

type Todo struct {
	ID        int64     `json:"id"`
	Task      string    `json:"task"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
