package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/pollpulse/pkg/common/config"
	"github.com/pollpulse/pkg/common/logging"
)

// ServiceConfig represents the configuration for a service
type ServiceConfig struct {
	Name            string
	URL             string
	AuthRequired    bool
	PathPrefix      string
	StripPathPrefix bool
}

func main() {
	// Initialize logger
	logger := logging.NewLogger(&logging.Config{
		Level:       config.GetEnv("LOG_LEVEL", "info"),
		ServiceName: "api-gateway",
		Environment: config.GetEnv("ENV", "development"),
	})
	logger.Info("Starting API gateway")

	// Service configurations
	services := []ServiceConfig{
		{
			Name:            "user-service",
			URL:             config.GetEnv("USER_SERVICE_URL", "http://localhost:8081"),
			AuthRequired:    false, // Some endpoints require auth, handled by the service
			PathPrefix:      "/api/v1",
			StripPathPrefix: false,
		},
		{
			Name:            "survey-service",
			URL:             config.GetEnv("SURVEY_SERVICE_URL", "http://localhost:8082"),
			AuthRequired:    true, // Most endpoints require auth
			PathPrefix:      "/api/v1/surveys",
			StripPathPrefix: false,
		},
		{
			Name:            "result-service",
			URL:             config.GetEnv("RESULT_SERVICE_URL", "http://localhost:8083"),
			AuthRequired:    true, // All endpoints require auth
			PathPrefix:      "/api/v1/results",
			StripPathPrefix: false,
		},
	}

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

	// Set up proxy routes for each service
	for _, service := range services {
		targetURL, err := url.Parse(service.URL)
		if err != nil {
			logger.Fatal("Invalid service URL", "service", service.Name, "url", service.URL, "error", err)
		}

		// Create a proxy for the service
		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		// Set up the proxy director
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			req.Header.Add("X-Gateway", "PollPulse-API-Gateway")
			req.URL.Host = targetURL.Host
			req.URL.Scheme = targetURL.Scheme
			req.Host = targetURL.Host

			// Strip the path prefix if needed
			if service.StripPathPrefix {
				req.URL.Path = strings.TrimPrefix(req.URL.Path, service.PathPrefix)
				if !strings.HasPrefix(req.URL.Path, "/") {
					req.URL.Path = "/" + req.URL.Path
				}
			}
		}

		// Set up error handler
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			logger.Error("Proxy error", "service", service.Name, "error", err, "path", r.URL.Path)
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(fmt.Sprintf("Service %s is unavailable", service.Name)))
		}

		// Define the handler
		handler := func(w http.ResponseWriter, r *http.Request) {
			// Log the request
			logger.Info("Proxying request",
				"service", service.Name,
				"method", r.Method,
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
			)
			proxy.ServeHTTP(w, r)
		}

		// Register the route
		r.Handle(service.PathPrefix+"/*", http.HandlerFunc(handler))
	}

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Root endpoint with API information
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"name": "PollPulse API Gateway", "version": "1.0.0", "status": "running"}`))
	})

	// Start server
	port := config.GetEnvInt("PORT", 8080)
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
