package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VitaliySynytskyi/pollpulse/pkg/common/config"
	"github.com/VitaliySynytskyi/pollpulse/pkg/common/database"
	"github.com/VitaliySynytskyi/pollpulse/pkg/common/logging"
	"github.com/VitaliySynytskyi/pollpulse/services/user-service/handler"
	"github.com/VitaliySynytskyi/pollpulse/services/user-service/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	// Initialize logger
	logger := logging.NewLogger(&logging.Config{
		Level:       config.GetEnv("LOG_LEVEL", "info"),
		ServiceName: "user-service",
		Environment: config.GetEnv("ENV", "development"),
	})
	logger.Info("Starting user service")

	// Database configuration
	dbConfig := &database.Config{
		Host:     config.GetEnv("DB_HOST", "localhost"),
		Port:     config.GetEnvInt("DB_PORT", 5432),
		User:     config.GetEnv("DB_USER", "postgres"),
		Password: config.GetEnv("DB_PASSWORD", "postgres"),
		DBName:   config.GetEnv("DB_NAME", "pollpulse_users"),
		SSLMode:  config.GetEnv("DB_SSLMODE", "disable"),
	}

	// Connect to the database
	db, err := database.Connect(dbConfig)
	if err != nil {
		logger.Fatal("Failed to connect to database", "error", err)
	}
	defer database.Close(db)

	// Create repository
	userRepo := repository.NewUserRepository(db)

	// Initialize router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Create handler
	jwtSecret := config.GetEnv("JWT_SECRET", "dev_secret_key")
	userHandler := handler.NewUserHandler(userRepo, logger, jwtSecret)

	// Register routes
	r.Route("/api/v1", func(r chi.Router) {
		userHandler.RegisterRoutes(r)
	})

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start server
	port := config.GetEnvInt("PORT", 8081)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	// Create a channel to listen for errors from the server
	serverErrors := make(chan error, 1)

	// Start the server in a goroutine
	go func() {
		logger.Info("Starting server", "port", port)
		serverErrors <- server.ListenAndServe()
	}()

	// Create a channel to listen for an interrupt or terminate signal
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Block until an error or shutdown signal is received
	select {
	case err := <-serverErrors:
		logger.Error("Server error", "error", err)

	case <-shutdown:
		logger.Info("Shutting down server")

		// Create a context with a timeout to give outstanding requests a chance to complete
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Ask the server to shutdown gracefully
		if err := server.Shutdown(ctx); err != nil {
			// If graceful shutdown fails, forcefully close
			logger.Error("Server forced to shutdown", "error", err)
			if err := server.Close(); err != nil {
				logger.Error("Server close error", "error", err)
			}
		}
	}

	logger.Info("Server stopped")
}
