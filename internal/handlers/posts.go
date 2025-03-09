package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"olbcloud.com/webapi/internal/services"
)

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

	writeResponse(w, http.StatusOK, posts)
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

	writeResponse(w, http.StatusOK, post)
}

func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}
