package main

import (
	"github.com/Rastaiha/bermudia/api/handler"
	"github.com/Rastaiha/bermudia/internal/config"
	"github.com/Rastaiha/bermudia/internal/mock"
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
	territoryRepo, err := repository.NewSqlTerritoryRepository(db)
	if err != nil {
		log.Fatal(err)
	}
	islandRepo, err := repository.NewSqlIslandRepository(db)
	if err != nil {
		log.Fatal(err)
	}
	userRepo, err := repository.NewSqlUser(db)
	if err != nil {
		log.Fatal(err)
	}
	playerRepo, err := repository.NewSqlPlayerRepository(db)
	if err != nil {
		log.Fatal(err)
	}

	authService := service.NewAuth(cfg, userRepo)
	territoryService := service.NewTerritory(territoryRepo)
	islandService := service.NewIsland(islandRepo)
	playerService := service.NewPlayer(playerRepo, territoryRepo)
	adminService := service.NewAdmin(territoryRepo, islandRepo, userRepo, playerRepo)

	err = mock.CreateMockData(adminService, cfg.MockUsersPassword)
	if err != nil {
		log.Fatal("failed to create mock data: ", err)
	}

	h := handler.New(authService, territoryService, islandService, playerService)

	h.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	<-c
	slog.Info("Got signal, shutting down...")
	h.Stop()
}
