package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "post_service/docs" // This is required for swagger
	"post_service/internal/repo"
	"post_service/internal/usecase"
	"post_service/pkg/httpserver"
	"post_service/pkg/mysql"
)

// @title        Post Service API
// @version      1.0
// @description  A service for managing blog posts and threads
// @termsOfService http://swagger.io/terms/
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host          localhost:8080
// @BasePath      /
// @securityDefinitions.basic  BasicAuth
// @tag.name threads
// @tag.description Operations about threads
// @tag.name posts
// @tag.description Operations about posts
func Run() {
	// Initialize actual post repository
	postRepo := repo.NewPostRepository(&mysql.MySql)
	threadRepo := repo.NewThreadRepository(&mysql.MySql)

	// Initialize use case
	useCase := usecase.NewPostUsecase(postRepo)
	threadUseCase := usecase.NewThreadUsecase(threadRepo)

	// Initialize server
	server := httpserver.NewServer(8080, useCase, threadUseCase)

	// Start server in a goroutine
	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Stop(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
