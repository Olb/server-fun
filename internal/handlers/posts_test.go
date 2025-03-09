package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"olbcloud.com/webapi/internal/database"
	"olbcloud.com/webapi/internal/database/mocks"
	"olbcloud.com/webapi/internal/handlers"
	"olbcloud.com/webapi/internal/models"
	"olbcloud.com/webapi/internal/server"
	"olbcloud.com/webapi/internal/services"
)

var mockPosts = []models.Post{
	{ID: 1, Title: "Post 1", Body: "Content 1"},
	{ID: 2, Title: "Post 2", Body: "Content 2"},
}

func setupTest(mockDB database.DB) (*httptest.ResponseRecorder, http.Handler) {
	postService := services.NewPostService(mockDB)
	h := handlers.NewHandlers(postService)
	mux := server.NewServer(h)
	return httptest.NewRecorder(), mux
}

func validateResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int, expectedBody string) {
	t.Helper()
	assert.Equal(t, expectedStatus, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Body.String(), expectedBody)
}

func TestHandlers(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		url            string
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
			expectedBody:   "Post 1",
		},
		{
			name:   "Get all posts - DB failure",
			method: http.MethodGet,
			url:    "/posts",
			mockSetup: func(m *mocks.DB) {
				m.On("GetPosts").Return(mockPosts, database.ErrFailedConnection)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "failed to get posts",
		},
		{
			name:   "Get post by ID successfully",
			method: http.MethodGet,
			url:    "/posts/1",
			mockSetup: func(m *mocks.DB) {
				m.On("GetPostByID", "1").Return(mockPosts[0], nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "Post 1",
		},
		{
			name:   "Get post by ID - Not Found",
			method: http.MethodGet,
			url:    "/posts/9",
			mockSetup: func(m *mocks.DB) {
				m.On("GetPostByID", "9").Return(models.Post{}, database.ErrNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "post not found",
		},
		{
			name:   "Get post by ID - DB failure",
			method: http.MethodGet,
			url:    "/posts/9",
			mockSetup: func(m *mocks.DB) {
				m.On("GetPostByID", "9").Return(models.Post{}, database.ErrFailedConnection)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "failed to connect",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new mock for each subtest.
			mockDB := new(mocks.DB)
			tt.mockSetup(mockDB)
			w, mux := setupTest(mockDB)

			req := httptest.NewRequest(tt.method, tt.url, nil)
			mux.ServeHTTP(w, req)

			validateResponse(t, w, tt.expectedStatus, tt.expectedBody)
			mockDB.AssertExpectations(t)
		})
	}
}
