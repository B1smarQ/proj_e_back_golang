package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"post_service/config"
	"post_service/internal/repo"
	"post_service/internal/usecase"
	"post_service/pkg/httpserver/handlers"
	"post_service/pkg/logger"
	"post_service/pkg/mysql"
	"time"

	"post_service/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func init() {
	cfg := config.NewConfig()
	l = *logger.New(cfg.App.LogLevel)
}

var server mux.Router
var l logger.Logger

type Router struct {
	router mux.Router
}

func (r *Router) Run() {
	l.Info("Running server")
	http.Handle("/", &r.router)

	// Start the server and listen for connections
	l.Info("Server listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		l.Fatal(err)
	}
}

func Create() Router {
	l.Info("Creating server")
	server = mux.Router{}

	// Initialize repositories
	postRepo := repo.NewPostRepository(&mysql.MySql)
	threadRepo := repo.NewThreadRepository(&mysql.MySql)

	// Initialize usecases
	postUseCase := usecase.NewPostUsecase(postRepo)
	threadUseCase := usecase.NewThreadUsecase(threadRepo)

	// Initialize handlers
	h := handlers.NewHandler(postUseCase, threadUseCase)

	// Register routes
	server.HandleFunc("/health", h.HealthCheck).Methods(http.MethodGet)

	// Post routes
	server.HandleFunc("/posts", h.GetPosts).Methods(http.MethodGet)
	server.HandleFunc("/posts/{id}", h.GetPost).Methods(http.MethodGet)
	server.HandleFunc("/posts", h.CreatePost).Methods(http.MethodPost)
	server.HandleFunc("/posts/{id}", h.UpdatePost).Methods(http.MethodPut)
	server.HandleFunc("/posts/{id}", h.DeletePost).Methods(http.MethodDelete)

	// Thread routes
	server.HandleFunc("/threads", h.GetThreads).Methods(http.MethodGet)
	server.HandleFunc("/threads/{id}", h.GetThread).Methods(http.MethodGet)
	server.HandleFunc("/threads", h.CreateThread).Methods(http.MethodPost)
	server.HandleFunc("/threads/{id}", h.UpdateThread).Methods(http.MethodPut)
	server.HandleFunc("/threads/{id}", h.DeleteThread).Methods(http.MethodDelete)
	server.HandleFunc("/threads/{id}/posts", h.GetThreadPosts).Methods(http.MethodGet)

	// Swagger documentation
	server.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))

	l.Info("Server created")
	return Router{server}
}

func Run() {
	l.Info("Running server")
	http.Handle("/", &server)
	l.Info("Server running")
}

type Server struct {
	httpServer *http.Server
	router     *mux.Router
	handler    *handlers.Handler
}

func NewServer(port int, postUseCase *usecase.PostUsecase, threadUseCase *usecase.ThreadUsecase) *Server {
	router := mux.NewRouter()
	handler := handlers.NewHandler(postUseCase, threadUseCase)

	// Initialize Swagger
	docs.SwaggerInfo.Title = "Post Service API"
	docs.SwaggerInfo.Description = "A service for managing blog posts and threads"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%d", port)
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}

	// Register routes
	router.HandleFunc("/health", handler.HealthCheck).Methods(http.MethodGet)

	// Post routes
	router.HandleFunc("/posts", handler.GetPosts).Methods(http.MethodGet)
	router.HandleFunc("/posts", handler.CreatePost).Methods(http.MethodPost)
	router.HandleFunc("/posts/{id}", handler.GetPost).Methods(http.MethodGet)
	router.HandleFunc("/posts/{id}", handler.UpdatePost).Methods(http.MethodPut)
	router.HandleFunc("/posts/{id}", handler.DeletePost).Methods(http.MethodDelete)

	// Thread routes
	router.HandleFunc("/threads", handler.GetThreads).Methods(http.MethodGet)
	router.HandleFunc("/threads/{id}", handler.GetThread).Methods(http.MethodGet)
	router.HandleFunc("/threads", handler.CreateThread).Methods(http.MethodPost)
	router.HandleFunc("/threads/{id}", handler.UpdateThread).Methods(http.MethodPut)
	router.HandleFunc("/threads/{id}", handler.DeleteThread).Methods(http.MethodDelete)
	router.HandleFunc("/threads/{id}/posts", handler.GetThreadPosts).Methods(http.MethodGet)

	// Swagger UI
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", port)),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
	))

	return &Server{
		httpServer: &http.Server{
			Addr:         fmt.Sprintf(":%d", port),
			Handler:      router,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		router:  router,
		handler: handler,
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
