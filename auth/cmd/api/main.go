package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	handler "auth/internal/delivery/http/handler"
	userrepo "auth/internal/repository/mongo"
	"auth/internal/usecase"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Auth Service API
// @version 1.0
// @description This is an authentication microservice using JWT and MongoDB.
// @host localhost:8080
// @BasePath /
func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Get environment variables
	mongoURI := getEnv("MONGO_URI", "mongodb://localhost:27017")
	jwtKey := getEnv("JWT_KEY", "your-secret-key")
	port := getEnv("PORT", "8080")

	// Connect to MongoDB
	client, err := mongoDriver.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	// Ping the database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	// Initialize dependencies
	db := client.Database("auth_service")
	userRepo := userrepo.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepo, jwtKey)

	// Initialize router
	router := mux.NewRouter()
	handler.NewUserHandler(router, userUseCase)

	// Swagger docs route
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Add middleware
	router.Use(loggingMiddleware)

	// Start server
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
	})
}
