package main

import (
	"github.com/Rastaiha/bermudia/api/handler"
	"github.com/Rastaiha/bermudia/internal/repository"
	"github.com/Rastaiha/bermudia/internal/service"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	territoryRepo := repository.NewJSONTerritoryRepository()
	islandRepo := repository.NewJSONIslandRepository()

	territoryService := service.NewTerritory(territoryRepo)
	islandService := service.NewIsland(islandRepo)

	h := handler.New(territoryService, islandService)

	h.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	<-c
	slog.Info("Got signal, shutting down...")
	h.Stop()
}
