package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"post_service/internal/entity"
	"post_service/internal/usecase"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

// MockPostRepository implements PostRepository interface for testing
type MockPostRepository struct{}

func (m *MockPostRepository) GetPosts() ([]entity.Post, error) {
	return []entity.Post{}, nil
}

func (m *MockPostRepository) CreatePost(post entity.Post) (entity.Post, error) {
	return post, nil
}

func (m *MockPostRepository) UpdatePost(post entity.Post) (entity.Post, error) {
	return post, nil
}

func (m *MockPostRepository) DeletePost(post entity.Post) (entity.Post, error) {
	return post, nil
}

func (m *MockPostRepository) GetPost() (entity.Post, error) {
	return entity.Post{
		ID:           1,
		AuthorID:     123,
		Title:        "Test Post",
		Content:      "Test Content",
		CreationTime: time.Now(),
	}, nil
}

func TestHealthCheck(t *testing.T) {
	// Create a new handler with mock repository
	useCase := &usecase.PostUsecase{
		PostRepository: &MockPostRepository{},
	}
	handler := NewHandler(useCase)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	handler.HealthCheck(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			w.Code, http.StatusOK)
	}

	expected := "OK"
	if w.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			w.Body.String(), expected)
	}
}

func TestGetPosts(t *testing.T) {
	// Create a new handler with mock repository
	useCase := &usecase.PostUsecase{
		PostRepository: &MockPostRepository{},
	}
	handler := NewHandler(useCase)

	req := httptest.NewRequest(http.MethodGet, "/posts", nil)
	w := httptest.NewRecorder()

	handler.GetPosts(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			w.Code, http.StatusOK)
	}
}

func TestGetPost(t *testing.T) {
	// Create a new handler with mock repository
	useCase := &usecase.PostUsecase{
		PostRepository: &MockPostRepository{},
	}
	handler := NewHandler(useCase)

	tests := []struct {
		name           string
		postID         string
		expectedStatus int
		expectedError  bool
	}{
		{
			name:           "valid post id",
			postID:         "1",
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "invalid post id",
			postID:         "invalid",
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/posts/"+tt.postID, nil)
			req = mux.SetURLVars(req, map[string]string{
				"id": tt.postID,
			})
			w := httptest.NewRecorder()

			handler.GetPost(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					w.Code, tt.expectedStatus)
			}

			if tt.expectedError {
				// For error responses, we expect a plain text error message
				errorMsg := w.Body.String()
				if errorMsg == "" {
					t.Error("expected error message in response")
				}
			} else {
				var post entity.Post
				err := json.NewDecoder(w.Body).Decode(&post)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if post.ID != 1 {
					t.Errorf("expected post ID 1, got %v", post.ID)
				}
			}
		})
	}
}

func TestCreatePost(t *testing.T) {
	// Create a new handler with mock repository
	useCase := &usecase.PostUsecase{
		PostRepository: &MockPostRepository{},
	}
	handler := NewHandler(useCase)

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "valid post data",
			payload: map[string]interface{}{
				"title":    "Test Post",
				"content":  "Test Content",
				"authorId": float64(123),
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name:           "empty payload",
			payload:        map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "missing required fields",
			payload: map[string]interface{}{
				"title": "Test Post",
				// missing content and authorId
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.CreatePost(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					w.Code, tt.expectedStatus)
			}

			if tt.expectedError {
				var response map[string]string
				err := json.NewDecoder(w.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if response["error"] == "" {
					t.Error("expected error message in response")
				}
			} else {
				var post entity.Post
				err := json.NewDecoder(w.Body).Decode(&post)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if post.Title != tt.payload["title"] {
					t.Errorf("expected title %v, got %v", tt.payload["title"], post.Title)
				}
				if post.Content != tt.payload["content"] {
					t.Errorf("expected content %v, got %v", tt.payload["content"], post.Content)
				}
				if post.AuthorID != int32(tt.payload["authorId"].(float64)) {
					t.Errorf("expected authorId %v, got %v", tt.payload["authorId"], post.AuthorID)
				}
				if post.CreationTime.IsZero() {
					t.Error("expected non-zero creation time")
				}
			}
		})
	}
}

func TestUpdatePost(t *testing.T) {
	// Create a new handler with mock repository
	useCase := &usecase.PostUsecase{
		PostRepository: &MockPostRepository{},
	}
	handler := NewHandler(useCase)

	tests := []struct {
		name           string
		postID         string
		payload        map[string]interface{}
		expectedStatus int
		expectedError  bool
	}{
		{
			name:   "valid update data",
			postID: "1",
			payload: map[string]interface{}{
				"title":   "Updated Post",
				"content": "Updated Content",
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "invalid post id",
			postID:         "invalid",
			payload:        map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPut, "/posts/"+tt.postID, bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			req = mux.SetURLVars(req, map[string]string{
				"id": tt.postID,
			})
			w := httptest.NewRecorder()

			handler.UpdatePost(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					w.Code, tt.expectedStatus)
			}

			if !tt.expectedError {
				var response map[string]interface{}
				err := json.NewDecoder(w.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if response["id"] != float64(1) && tt.postID == "1" {
					t.Errorf("expected id 1 in response, got %v", response["id"])
				}
			} else {
				var response map[string]string
				err := json.NewDecoder(w.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if response["error"] != "Invalid post ID" {
					t.Errorf("expected error message 'Invalid post ID', got %v", response["error"])
				}
			}
		})
	}
}

func TestDeletePost(t *testing.T) {
	// Create a new handler with mock repository
	useCase := &usecase.PostUsecase{
		PostRepository: &MockPostRepository{},
	}
	handler := NewHandler(useCase)

	tests := []struct {
		name           string
		postID         string
		expectedStatus int
		expectedError  bool
	}{
		{
			name:           "valid post id",
			postID:         "1",
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "invalid post id",
			postID:         "invalid",
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/posts/"+tt.postID, nil)
			req = mux.SetURLVars(req, map[string]string{
				"id": tt.postID,
			})
			w := httptest.NewRecorder()

			handler.DeletePost(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					w.Code, tt.expectedStatus)
			}

			if !tt.expectedError {
				var response map[string]interface{}
				err := json.NewDecoder(w.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if response["id"] != float64(1) && tt.postID == "1" {
					t.Errorf("expected id 1 in response, got %v", response["id"])
				}
			} else {
				var response map[string]string
				err := json.NewDecoder(w.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if response["error"] != "Invalid post ID" {
					t.Errorf("expected error message 'Invalid post ID', got %v", response["error"])
				}
			}
		})
	}
}
