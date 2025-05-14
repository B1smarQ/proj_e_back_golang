package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreate(t *testing.T) {
	router := Create()

	// Test cases for each endpoint
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "health check",
			method:         http.MethodGet,
			path:           "/health",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "get posts",
			method:         http.MethodGet,
			path:           "/posts",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "get single post",
			method:         http.MethodGet,
			path:           "/posts/1",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "create post",
			method:         http.MethodPost,
			path:           "/posts",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "update post",
			method:         http.MethodPut,
			path:           "/posts/1",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "delete post",
			method:         http.MethodDelete,
			path:           "/posts/1",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			router.router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					w.Code, tt.expectedStatus)
			}
		})
	}
}

func TestInvalidMethods(t *testing.T) {
	router := Create()

	// Test cases for invalid methods
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "post to health endpoint",
			method:         http.MethodPost,
			path:           "/health",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "get to posts with invalid method",
			method:         http.MethodPatch,
			path:           "/posts",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "invalid method for single post",
			method:         http.MethodPatch,
			path:           "/posts/1",
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			router.router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					w.Code, tt.expectedStatus)
			}
		})
	}
}

func TestInvalidPaths(t *testing.T) {
	router := Create()

	// Test cases for invalid paths
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "non-existent path",
			method:         http.MethodGet,
			path:           "/nonexistent",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "invalid post id format",
			method:         http.MethodGet,
			path:           "/posts/invalid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			router.router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					w.Code, tt.expectedStatus)
			}
		})
	}
}
