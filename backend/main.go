package main

import (
	"database/sql"
	"fmt"
	"github.com/Rastaiha/bermudia/adminbot"
	"github.com/Rastaiha/bermudia/api/handler"
	"github.com/Rastaiha/bermudia/internal/config"
	"github.com/Rastaiha/bermudia/internal/domain"
	"github.com/Rastaiha/bermudia/internal/mock"
	"github.com/Rastaiha/bermudia/internal/repository"
	"github.com/Rastaiha/bermudia/internal/service"
	"github.com/go-telegram/bot"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func main() {
	cfg := config.Load()

	domain.ApplyConfig(cfg)

	theBot, err := bot.New(cfg.BotToken, bot.WithServerURL("https://tapi.bale.ai"))
	if err != nil {
		log.Fatal("failed to connect to bot api: ", err)
	}

	var db *sql.DB
	if cfg.Postgres.Enable {
		db, err = repository.ConnectToPostgres(cfg.Postgres)
		if err != nil {
			log.Fatal("failed to connect to postgres: ", err)
		}
	} else {
		db, err = repository.ConnectToSqlite()
		if err != nil {
			log.Fatal("failed to connect to sqlite: ", err)
		}
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
	questionStore, err := repository.NewSqlQuestionRepository(db)
	if err != nil {
		log.Fatal(err)
	}
	treasureRepo, err := repository.NewSqlTreasureRepository(db)
	if err != nil {
		log.Fatal(err)
	}
	marketRepo, err := repository.NewSqlMarketRepository(db)
	if err != nil {
		log.Fatal(err)
	}
	inboxRepo, err := repository.NewSqlInboxRepository(db)
	if err != nil {
		log.Fatal(err)
	}

	authService := service.NewAuth(cfg, userRepo)
	territoryService := service.NewTerritory(territoryRepo)
	islandService := service.NewIsland(theBot, islandRepo, questionStore, playerRepo, treasureRepo)
	playerService := service.NewPlayer(cfg, db, userRepo, playerRepo, territoryRepo, questionStore, islandRepo, treasureRepo, marketRepo, inboxRepo)
	correctionService := service.NewCorrection(cfg, questionStore)
	adminService := service.NewAdmin(cfg, territoryRepo, islandRepo, userRepo, playerRepo, questionStore, treasureRepo)

	islandService.OnNewPortableIsland(playerService.HandleNewPortableIsland)

	if cfg.DevMode && cfg.CreateMock {
		writeBackDir := filepath.Join(os.TempDir(), fmt.Sprintf("data_%d", time.Now().Unix()-1758018000))
		err = mock.SetGameContent(adminService, mock.DataFiles, writeBackDir, cfg.MockUsersPassword)
		if err != nil {
			log.Fatal("failed to create mock data: ", err)
		}
	}

	adminBot := adminbot.NewBot(cfg, theBot, islandService, correctionService, adminService, userRepo)

	h := handler.New(authService, territoryService, islandService, playerService)

	playerService.Start()
	adminBot.Start()
	h.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	<-c
	slog.Info("Got signal, shutting down...")

	h.Stop()
	adminBot.Stop()
	playerService.Stop()
}
