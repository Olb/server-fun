package handlers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"olbcloud.com/webapi/internal/apiserver"
	"olbcloud.com/webapi/internal/database"
	"olbcloud.com/webapi/internal/database/mocks"
	"olbcloud.com/webapi/internal/handlers"
	"olbcloud.com/webapi/internal/models"
	"olbcloud.com/webapi/internal/services"
)

var mockPosts = []models.Post{
	{ID: 1, Title: "Post 1", Body: "Content 1", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 2, Title: "Post 2", Body: "Content 2", CreatedAt: time.Now(), UpdatedAt: time.Now()},
}

func setupTest(mockDB database.DB) (*httptest.ResponseRecorder, http.Handler) {
	postService := services.NewPostService(mockDB)
	h := handlers.NewHandlers(postService)
	mux := apiserver.NewServer(h)
	return httptest.NewRecorder(), mux
}

func validateResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int, expectedBody string) {
	t.Helper()
	assert.Equal(t, expectedStatus, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var actual map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &actual)
	assert.NoError(t, err)

	var expected map[string]interface{}
	err = json.Unmarshal([]byte(expectedBody), &expected)
	assert.NoError(t, err)

	// get rid of those annoying timestamps
	removeTimestamps(actual)
	removeTimestamps(expected)

	assert.Equal(t, expected, actual)
}

func removeTimestamps(data map[string]interface{}) {
	if post, ok := data["post"].(map[string]interface{}); ok {
		delete(post, "created_at")
		delete(post, "updated_at")
	}
	if posts, ok := data["posts"].([]interface{}); ok {
		for _, p := range posts {
			if postMap, ok := p.(map[string]interface{}); ok {
				delete(postMap, "created_at")
				delete(postMap, "updated_at")
			}
		}
	}
}

func TestHandlers(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		url            string
		body           string
		mockSetup      func(mockDB *mocks.DB)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Get all posts successfully",
			method: http.MethodGet,
			url:    "/posts",
			mockSetup: func(m *mocks.DB) {
				m.On("GetPosts").Return(mockPosts, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"posts":[{"id":1,"title":"Post 1","body":"Content 1"},{"id":2,"title":"Post 2","body":"Content 2"}]}`,
		},
		{
			name:   "Get all posts - DB failure",
			method: http.MethodGet,
			url:    "/posts",
			mockSetup: func(m *mocks.DB) {
				m.On("GetPosts").Return(nil, database.ErrFailedConnection)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"failed to get posts"}`,
		},
		{
			name:   "Get post by ID successfully",
			method: http.MethodGet,
			url:    "/posts/1",
			mockSetup: func(m *mocks.DB) {
				m.On("GetPostByID", "1").Return(mockPosts[0], nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"post":{"id":1,"title":"Post 1","body":"Content 1"}}`,
		},
		{
			name:   "Get post by ID - Not Found",
			method: http.MethodGet,
			url:    "/posts/9",
			mockSetup: func(m *mocks.DB) {
				m.On("GetPostByID", "9").Return(models.Post{}, database.ErrNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"post not found"}`,
		},
		{
			name:   "Get post by ID - DB failure",
			method: http.MethodGet,
			url:    "/posts/9",
			mockSetup: func(m *mocks.DB) {
				m.On("GetPostByID", "9").Return(models.Post{}, database.ErrFailedConnection)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"failed to connect to the database"}`,
		},
		{
			name:   "Creates a post returns invalid payload",
			method: http.MethodPost,
			url:    "/posts",
			body:   `{"some":1}`,
			mockSetup: func(m *mocks.DB) {
				m.AssertNotCalled(t, "CreatePost")
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"missing required fields"}`,
		},
		{
			name:   "Creates a post returns invalid payload with missing title",
			method: http.MethodPost,
			url:    "/posts",
			body:   `{"body":"a post"}`,
			mockSetup: func(m *mocks.DB) {
				m.AssertNotCalled(t, "CreatePost")
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"missing required fields"}`,
		},
		{
			name:   "Creates a post returns invalid payload with missing body",
			method: http.MethodPost,
			url:    "/posts",
			body:   `{"title":"a title"}`,
			mockSetup: func(m *mocks.DB) {
				m.AssertNotCalled(t, "CreatePost")
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"missing required fields"}`,
		},
		{
			name:   "Creates a post returns empty body",
			method: http.MethodPost,
			url:    "/posts",
			body:   ``,
			mockSetup: func(m *mocks.DB) {
				m.AssertNotCalled(t, "CreatePost")
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid request payload"}`,
		},
		{
			name:   "Creates a post successfully",
			method: http.MethodPost,
			url:    "/posts",
			body:   `{"title": "New Post", "body": "New Content"}`,
			mockSetup: func(m *mocks.DB) {
				m.On("CreatePost", mock.AnythingOfType("models.Post")).Return(models.Post{
					ID:        9,
					Title:     "New Post",
					Body:      "New Content",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"message":"post created","status":"success","post":{"id":9,"title":"New Post","body":"New Content"}}`,
		},
		{
			name:   "Fails to create post due to DB error",
			method: http.MethodPost,
			url:    "/posts",
			body:   `{"title": "New Post", "body": "New Content"}`,
			mockSetup: func(m *mocks.DB) {
				m.On("CreatePost", mock.AnythingOfType("models.Post")).Return(models.Post{}, errors.New("failed to create post"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"failed to create post"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(mocks.DB)
			tt.mockSetup(mockDB)
			w, mux := setupTest(mockDB)

			req := httptest.NewRequest(tt.method, tt.url, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			mux.ServeHTTP(w, req)

			validateResponse(t, w, tt.expectedStatus, tt.expectedBody)
			mockDB.AssertExpectations(t)
		})
	}
}
