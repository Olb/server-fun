package models

import "time"

// Post represents a blog post.
type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title" validate:"required"`
	Body      string    `json:"body" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
