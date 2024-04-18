package main

import (
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/config"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/initialization"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/server"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/router"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
)

func main() {
	logger := logger.Logger()
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load config", "error", err)
		return
	}
	logger.Info("Start config", "config", cfg)

	pdb, rdb, err := initialization.DataBaseIniInitialization()
	if err != nil {
		logger.Error("DataBase initialization error", "error", err)
	}

	storages := initialization.StorageInit(pdb, rdb)
	usecases := initialization.UseCaseInit(storages)
	handlers := initialization.HandlerInit(usecases)

	router := router.SetupRouter(cfg, handlers)

	if err := server.StartServer(router, cfg); err != nil {
		logger.Error("Failed to start server", "error", err)
		return
	}
}
