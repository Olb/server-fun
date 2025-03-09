package database

import (
	"errors"

	"olbcloud.com/webapi/internal/models"
)

var ErrFailedConnection = errors.New("failed to connect to the database")
var ErrNotFound = errors.New("entity not found")

// DB interface defines database operations
type DB interface {
	GetPosts() ([]models.Post, error)
	GetPostByID(id string) (models.Post, error)
	Close() error
}
