package handlers

import (
	"encoding/json"
	"net/http"
	"post_service/internal/entity"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// @Summary     Get all threads
// @Description Retrieve a list of all threads
// @Tags        threads
// @Produce     json
// @Success     200 {array}  entity.ReturnThread
// @Failure     500 {object} map[string]string
// @Router      /threads [get]
func (h *Handler) GetThreads(w http.ResponseWriter, r *http.Request) {
	threads, err := h.threadUseCase.GetThreadsUsecase()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(threads)
}

// @Summary     Create a new thread
// @Description Create a new discussion thread
// @Tags        threads
// @Accept      json
// @Produce     json
// @Param       thread body entity.Thread true "Thread object"
// @Success     201 {object} entity.ReturnThread
// @Failure     400 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /threads [post]
func (h *Handler) CreateThread(w http.ResponseWriter, r *http.Request) {
	var thread entity.Thread
	if err := json.NewDecoder(r.Body).Decode(&thread); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	thread.CreatedAt = time.Now()
	thread.UpdatedAt = time.Now()

	createdThread, err := h.threadUseCase.CreateThreadUsecase(thread)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdThread)
}

// @Summary     Update a thread
// @Description Update an existing discussion thread
// @Tags        threads
// @Accept      json
// @Produce     json
// @Param       id     path int true "Thread ID"
// @Param       thread body entity.Thread true "Thread object"
// @Success     200 {object} entity.ReturnThread
// @Failure     400 {object} map[string]string
// @Failure     404 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /threads/{id} [put]
func (h *Handler) UpdateThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid thread ID", http.StatusBadRequest)
		return
	}

	var thread entity.Thread
	if err := json.NewDecoder(r.Body).Decode(&thread); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	thread.ID = int32(id)
	thread.UpdatedAt = time.Now()

	updatedThread, err := h.threadUseCase.UpdateThreadUsecase(thread)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedThread)
}

// @Summary     Delete a thread
// @Description Delete an existing discussion thread
// @Tags        threads
// @Produce     json
// @Param       id path int true "Thread ID"
// @Success     200 {object} entity.ReturnThread
// @Failure     400 {object} map[string]string
// @Failure     404 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /threads/{id} [delete]
func (h *Handler) DeleteThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid thread ID", http.StatusBadRequest)
		return
	}

	thread := entity.Thread{ID: int32(id)}
	deletedThread, err := h.threadUseCase.DeleteThreadUsecase(thread)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deletedThread)
}

// @Summary     Get a thread
// @Description Get a specific thread by ID
// @Tags        threads
// @Produce     json
// @Param       id path int true "Thread ID"
// @Success     200 {object} entity.ReturnThread
// @Failure     400 {object} map[string]string
// @Failure     404 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /threads/{id} [get]
func (h *Handler) GetThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid thread ID", http.StatusBadRequest)
		return
	}

	retrievedThread, err := h.threadUseCase.GetThreadUsecase(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(retrievedThread)
}

// @Summary     Get thread posts
// @Description Get all posts in a specific thread
// @Tags        threads
// @Produce     json
// @Param       id path int true "Thread ID"
// @Success     200 {array}  entity.ReturnPost
// @Failure     400 {object} map[string]string
// @Failure     404 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /threads/{id}/posts [get]
func (h *Handler) GetThreadPosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid thread ID", http.StatusBadRequest)
		return
	}

	posts, err := h.threadUseCase.GetThreadPostsUsecase(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
