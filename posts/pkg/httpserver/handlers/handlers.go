package handlers

import (
	"encoding/json"
	"net/http"
	"post_service/internal/entity"
	"post_service/internal/usecase"
	"time"

	"github.com/gorilla/mux"
)

type Handler struct {
	useCase       *usecase.PostUsecase
	threadUseCase *usecase.ThreadUsecase
}

func NewHandler(useCase *usecase.PostUsecase, threadUseCase *usecase.ThreadUsecase) *Handler {
	return &Handler{
		useCase:       useCase,
		threadUseCase: threadUseCase,
	}
}

// @Summary     Health check endpoint
// @Description Get the health status of the service
// @Tags        health
// @Produce     plain
// @Success     200 {string} string "OK"
// @Router      /health [get]
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// @Summary     Get all posts
// @Description Retrieve a list of all posts
// @Tags        posts
// @Produce     json
// @Success     200 {array}  entity.ReturnPost
// @Failure     500 {object} map[string]string
// @Router      /posts [get]
func (h *Handler) GetPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.useCase.GetPostsUsecase()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// @Summary     Create a new post
// @Description Create a new blog post
// @Tags        posts
// @Accept      json
// @Produce     json
// @Param       post body entity.Post true "Post object"
// @Success     201 {object} entity.Post
// @Failure     400 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /posts [post]
func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var post entity.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if post.CreationTime.IsZero() {
		post.CreationTime = time.Now()
	}

	createdPost, err := h.useCase.CreatePostUsecase(post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdPost)
}

// @Summary     Update a post
// @Description Update an existing blog post
// @Tags        posts
// @Accept      json
// @Produce     json
// @Param       id   path string true "Post ID"
// @Param       post body entity.Post true "Post object"
// @Success     200 {object} entity.Post
// @Failure     400 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /posts/{id} [put]
func (h *Handler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var post entity.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	post.ID = id
	updatedPost, err := h.useCase.UpdatePostUsecase(post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedPost)
}

// @Summary     Delete a post
// @Description Delete an existing blog post
// @Tags        posts
// @Produce     json
// @Param       id path string true "Post ID"
// @Success     200 {object} entity.Post
// @Failure     400 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /posts/{id} [delete]
func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	post := entity.Post{ID: id}
	deletedPost, err := h.useCase.DeletePostUsecase(post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deletedPost)
}

// @Summary     Get a post
// @Description Get a specific blog post by ID
// @Tags        posts
// @Produce     json
// @Param       id path string true "Post ID"
// @Success     200 {object} entity.ReturnPost
// @Failure     400 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /posts/{id} [get]
func (h *Handler) GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	retrievedPost, err := h.useCase.GetPostUsecase(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(retrievedPost)
}
