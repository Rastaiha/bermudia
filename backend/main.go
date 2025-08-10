package main

import (
	"github.com/Rastaiha/bermudia/api/handler"
	"github.com/Rastaiha/bermudia/internal/repository"
	"github.com/Rastaiha/bermudia/internal/service"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	territoryRepo := repository.NewJSONTerritoryRepository()
	islandRepo := repository.NewJSONIslandRepository()

	territoryService := service.NewTerritory(territoryRepo)
	islandService := service.NewIsland(islandRepo)

	h := handler.New(territoryService, islandService)

	h.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	h.Stop()
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
