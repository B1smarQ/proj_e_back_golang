package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	"post_service/config"
	"post_service/internal/entity"
	"post_service/internal/middleware"
	"post_service/internal/repository/mysql"
	"post_service/internal/service"
	"post_service/pkg/logger"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize logger
	logger := logger.New("info")

	// Connect to database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Fatal("Failed to connect to database", err)
	}
	defer db.Close()

	// Initialize repositories
	postRepo := mysql.NewPostRepository(db)

	// Initialize services
	authClient := service.NewAuthClient(cfg.AuthURL)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(authClient)

	// Initialize router
	router := mux.NewRouter()

	// Apply middleware
	router.Use(authMiddleware.Authenticate)

	// Initialize and register handlers
	_ = NewPostHandler(router, postRepo)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	logger.Info("Starting server on " + addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Fatal("Failed to start server", err)
	}
}

type PostHandler struct {
	router   *mux.Router
	postRepo entity.PostRepository
}

func NewPostHandler(router *mux.Router, postRepo entity.PostRepository) *PostHandler {
	handler := &PostHandler{
		router:   router,
		postRepo: postRepo,
	}

	router.HandleFunc("/posts", handler.CreatePost).Methods("POST")
	router.HandleFunc("/posts", handler.GetPosts).Methods("GET")
	router.HandleFunc("/posts/{id}", handler.GetPost).Methods("GET")
	router.HandleFunc("/posts/{id}", handler.UpdatePost).Methods("PUT")
	router.HandleFunc("/posts/{id}", handler.DeletePost).Methods("DELETE")

	return handler
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	// Get user info from context
	userID := r.Context().Value("user_id").(string)
	username := r.Context().Value("username").(string)

	var post entity.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Set author info
	post.AuthorID = userID
	post.AuthorName = username
	post.CreationTime = time.Now()

	if err := post.Create(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func (h *PostHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.postRepo.GetPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(posts)
}

func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	post, err := h.postRepo.GetPost(vars["id"])
	if err != nil {
		if err == entity.ErrPostNotFound {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(post)
}

func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	// Get user info from context
	userID := r.Context().Value("user_id").(string)

	vars := mux.Vars(r)
	post, err := h.postRepo.GetPost(vars["id"])
	if err != nil {
		if err == entity.ErrPostNotFound {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if user is the author
	if post.AuthorID != userID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var updatedPost entity.Post
	if err := json.NewDecoder(r.Body).Decode(&updatedPost); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	updatedPost.ID = post.ID
	updatedPost.AuthorID = post.AuthorID
	updatedPost.AuthorName = post.AuthorName
	updatedPost.CreationTime = post.CreationTime

	if err := updatedPost.Update(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedPost)
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	// Get user info from context
	userID := r.Context().Value("user_id").(string)

	vars := mux.Vars(r)
	post, err := h.postRepo.GetPost(vars["id"])
	if err != nil {
		if err == entity.ErrPostNotFound {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if user is the author
	if post.AuthorID != userID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	deletePost := entity.Post{ID: post.ID}
	if err := deletePost.Delete(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
