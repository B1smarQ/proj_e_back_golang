package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

type AuthService interface {
	ValidateToken(token string) (map[string]interface{}, error)
}

type authMiddleware struct {
	authService AuthService
}

func NewAuthMiddleware(authService AuthService) *authMiddleware {
	return &authMiddleware{
		authService: authService,
	}
}

func (m *authMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondWithError(w, http.StatusUnauthorized, "Missing authorization token")
			return
		}

		// Remove "Bearer " prefix if present
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			respondWithError(w, http.StatusUnauthorized, "Invalid token format")
			return
		}

		claims, err := m.authService.ValidateToken(token)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Add user info to request context
		ctx := r.Context()
		for key, value := range claims {
			ctx = context.WithValue(ctx, key, value)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	response := map[string]interface{}{
		"status":  "error",
		"message": message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}
