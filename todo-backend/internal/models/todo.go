package models

import "time"

type Todo struct {
	ID        int64     `json:"id"`
	Task      string    `json:"task"`
	CreatedAt time.Time `json:"created_at"`
}
