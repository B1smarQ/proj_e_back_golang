package http

import (
	"encoding/json"
	"net/http"

	"auth/internal/domain"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	userUseCase domain.UserUseCase
}

func NewUserHandler(router *mux.Router, useCase domain.UserUseCase) {
	handler := &UserHandler{
		userUseCase: useCase,
	}

	router.HandleFunc("/auth/register", handler.Register).Methods("POST")
	router.HandleFunc("/auth/login", handler.Login).Methods("POST")
	router.HandleFunc("/auth/validate", handler.ValidateToken).Methods("GET")
}

type registerRequest struct {
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

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := h.userUseCase.Register(req.Email, req.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		Status: "success",
		Data:   user,
	})
}

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
