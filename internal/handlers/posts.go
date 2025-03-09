package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"olbcloud.com/webapi/internal/models"
	"olbcloud.com/webapi/internal/services"
)

var validate = validator.New()

type Handlers struct {
	PostService services.PostService
}

func NewHandlers(ps services.PostService) Handlers {
	return Handlers{PostService: ps}
}

func (h *Handlers) GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := h.PostService.GetPosts()
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to get posts"})
		return
	}

	writeResponse(w, http.StatusOK, map[string]interface{}{"posts": posts})
}

func (h *Handlers) GetPostByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	post, err := h.PostService.GetPostByID(id)
	if err != nil {
		if errors.Is(err, services.ErrPostNotFound) {
			writeResponse(w, http.StatusNotFound, map[string]string{"error": "post not found"})
			return
		}
		writeResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeResponse(w, http.StatusOK, map[string]interface{}{"post": post})
}

func (h *Handlers) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid request payload"})
		return
	}

	if err := validate.Struct(post); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]string{"error": "missing required fields"})
		return
	}

	post, err := h.PostService.CreatePost(post)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to create post"})
		return
	}
	writeResponse(w, http.StatusCreated, map[string]interface{}{"message": "post created", "status": "success", "post": post})
}

func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}
