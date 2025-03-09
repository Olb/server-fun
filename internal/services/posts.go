package services

import (
	"errors"

	"olbcloud.com/webapi/internal/database"
	"olbcloud.com/webapi/internal/models"
)

var ErrPostNotFound = errors.New("the requested post was not found")

type PostService interface {
	GetPosts() ([]models.Post, error)
	GetPostByID(id string) (models.Post, error)
	CreatePost(post models.Post) (models.Post, error)
}

type postService struct {
	db database.DB
}

func NewPostService(db database.DB) PostService {
	return &postService{db: db}
}

// GetPosts returns all posts
func (ps *postService) GetPosts() ([]models.Post, error) {
	posts, err := ps.db.GetPosts()
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (ps *postService) GetPostByID(id string) (models.Post, error) {
	post, err := ps.db.GetPostByID(id)
	if err != nil {
		if err == database.ErrNotFound {
			return models.Post{}, ErrPostNotFound
		}
		return models.Post{}, err
	}

	return post, nil
}

func (ps *postService) CreatePost(post models.Post) (models.Post, error) {
	return ps.db.CreatePost(post)
}
