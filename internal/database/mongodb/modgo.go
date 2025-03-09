package mongodb

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"olbcloud.com/webapi/internal/database"
	"olbcloud.com/webapi/internal/models"
)

// MongoDB struct
type MongoDB struct {
	client *mongo.Client
	posts  *mongo.Collection
}

// NewMongoDB initializes the connection
func NewMongoDB(uri string) (database.DB, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Println("MongoDB connection failed:", err)
		return nil, database.ErrFailedConnection
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Println("MongoDB connect failed:", err)
		return nil, database.ErrFailedConnection
	}

	log.Println("Connected to MongoDB")
	db := client.Database("blog")
	return &MongoDB{client: client, posts: db.Collection("posts")}, nil
}

// GetPosts retrieves all posts from MongoDB
func (m *MongoDB) GetPosts() ([]models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := m.posts.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []models.Post
	for cursor.Next(ctx) {
		var post models.Post
		if err := cursor.Decode(&post); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if len(posts) == 0 {
		return nil, database.ErrNotFound
	}

	return posts, nil
}

// GetPostByID retrieves a single post by ID
func (m *MongoDB) GetPostByID(id string) (models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var post models.Post
	err := m.posts.FindOne(ctx, bson.M{"id": id}).Decode(&post)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Post{}, database.ErrNotFound
		}
		return models.Post{}, err
	}
	return post, nil
}

func (m *MongoDB) CreatePost(post models.Post) (models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.posts.InsertOne(ctx, post)
	if err != nil {
		return models.Post{}, err
	}
	return post, nil
}

// Close closes the MongoDB connection
func (m *MongoDB) Close() error {
	return m.client.Disconnect(context.Background())
}
