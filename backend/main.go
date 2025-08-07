package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/Rastaiha/rasta-1404-contest/api/handlers"
	"github.com/Rastaiha/rasta-1404-contest/internal/repository"
	"github.com/Rastaiha/rasta-1404-contest/internal/service"
)

func main() {
	// Initialize repository (currently JSON file-based with embedded files)
	territoryRepo := repository.NewJSONTerritoryRepository()

	// Initialize service
	territoryService := service.NewTerritoryService(territoryRepo)

	// Initialize handlers
	territoryHandler := handlers.NewTerritoryHandler(territoryService)

	// Setup router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(corsMiddleware)

	// Routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/territories/{territoryID}", territoryHandler.GetTerritory)
	})

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

// Simple CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")

		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
	})
}
