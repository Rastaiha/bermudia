package main

import (
	"github.com/Rastaiha/bermudia/api/handler"
	"github.com/Rastaiha/bermudia/internal/config"
	"github.com/Rastaiha/bermudia/internal/domain"
	"github.com/Rastaiha/bermudia/internal/repository"
	"github.com/Rastaiha/bermudia/internal/service"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.Load()

	db, err := repository.ConnectToSqlite()
	if err != nil {
		log.Fatal("failed to connect to sqlite", err)
	}
	territoryRepo := repository.NewJSONTerritoryRepository()
	islandRepo := repository.NewJSONIslandRepository()
	userRepo, err := repository.NewSqlUser(db)
	if err != nil {
		log.Fatal(err)
	}
	err = domain.CreateMockData(userRepo, cfg.MockUsersPassword)
	if err != nil {
		log.Fatal("failed to create mock data", err)
	}

	authService := service.NewAuth(cfg, userRepo)
	territoryService := service.NewTerritory(territoryRepo)
	islandService := service.NewIsland(islandRepo)

	h := handler.New(authService, territoryService, islandService)

	h.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	<-c
	slog.Info("Got signal, shutting down...")
	h.Stop()
}
