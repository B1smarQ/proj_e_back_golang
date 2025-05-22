// Package handler provides HTTP handlers for authentication endpoints.
//
// @title Auth Service API
// @version 1.0
// @description This is an authentication microservice using JWT and MongoDB.
// @host localhost:8080
// @BasePath /
package handler

import (
	"encoding/json"
	"net/http"

	"auth/internal/domain"

	"github.com/gorilla/mux"
)

// UserHandler handles authentication endpoints.
type UserHandler struct {
	userUseCase domain.UserUseCase
}

// NewUserHandler registers authentication routes.
func NewUserHandler(router *mux.Router, useCase domain.UserUseCase) {
	handler := &UserHandler{
		userUseCase: useCase,
	}

	router.HandleFunc("/auth/register", handler.Register).Methods("POST")
	router.HandleFunc("/auth/login", handler.Login).Methods("POST")
	router.HandleFunc("/auth/validate", handler.ValidateToken).Methods("GET")
}

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with username, email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param registerRequest body registerRequest true "Register request"
// @Success 201 {object} response
// @Failure 400 {object} response
// @Router /auth/register [post]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Username == "" {
		respondWithError(w, http.StatusBadRequest, "Username is required")
		return
	}

	user, err := h.userUseCase.Register(req.Username, req.Email, req.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		Status: "success",
		Data:   user,
	})
}

// Login godoc
// @Summary Login
// @Description Login with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param loginRequest body loginRequest true "Login request"
// @Success 200 {object} response
// @Failure 400 {object} response
// @Failure 401 {object} response
// @Router /auth/login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	token, err := h.userUseCase.Login(req.Email, req.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Status: "success",
		Data: map[string]string{
			"token": token,
		},
	})
}

// ValidateToken godoc
// @Summary Validate JWT token
// @Description Validate a JWT token and return user info
// @Tags auth
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} response
// @Failure 401 {object} response
// @Router /auth/validate [get]
func (h *UserHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		respondWithError(w, http.StatusUnauthorized, "Missing authorization token")
		return
	}

	// Remove "Bearer " prefix if present
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	user, err := h.userUseCase.ValidateToken(token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Status: "success",
		Data:   user,
	})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, response{
		Status:  "error",
		Message: message,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","message":"Internal server error"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
