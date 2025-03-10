package postgresql

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/lib/pq"
	"olbcloud.com/webapi/internal/database"
	"olbcloud.com/webapi/internal/models"
)

// PostgreSQL struct
type PostgreSQL struct {
	conn *sql.DB
}

func NewPostgreSQL(dsn string) (database.DB, error) {
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Println("PostgreSQL connection failed:", err)
		return nil, database.ErrFailedConnection
	}

	if err := conn.Ping(); err != nil {
		log.Println("PostgreSQL ping failed:", err)
		return nil, database.ErrFailedConnection
	}

	log.Println("Connected to PostgreSQL")
	return &PostgreSQL{conn: conn}, nil
}

func (p *PostgreSQL) GetPosts() ([]models.Post, error) {
	rows, err := p.conn.Query("SELECT id, title, body, created_at, updated_at FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Body, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	if len(posts) == 0 {
		return nil, database.ErrNotFound
	}

	return posts, nil
}

func (p *PostgreSQL) GetPostByID(id string) (models.Post, error) {
	var post models.Post
	err := p.conn.QueryRow("SELECT id, title, body, created_at, updated_at FROM posts WHERE id = $1", id).
		Scan(&post.ID, &post.Title, &post.Body, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Post{}, database.ErrNotFound
		}
		return models.Post{}, err
	}
	return post, nil
}

func (p *PostgreSQL) CreatePost(post models.Post) (models.Post, error) {
	err := p.conn.QueryRow(
		`INSERT INTO posts (title, body) 
		 VALUES ($1, $2) 
		 RETURNING id, title, body, created_at, updated_at`,
		post.Title, post.Body,
	).Scan(&post.ID, &post.Title, &post.Body, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return models.Post{}, err
	}
	return post, nil
}

func (p *PostgreSQL) UpdatePost(post models.Post) (models.Post, error) {
	err := p.conn.QueryRow(
		`UPDATE posts SET title = $1, body = $2, updated_at = NOW() 
		 WHERE id = $3 
		 RETURNING id, title, body, created_at, updated_at`,
		post.Title, post.Body, post.ID,
	).Scan(&post.ID, &post.Title, &post.Body, &post.CreatedAt, &post.UpdatedAt)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Post{}, database.ErrNotFound
		}
		return models.Post{}, err
	}
	return post, nil
}

func (p *PostgreSQL) Close() error {
	return p.conn.Close()
}
