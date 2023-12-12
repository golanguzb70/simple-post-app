package main

import (
	"log"

	"github.com/golanguzb70/simple-post-app/api"
	"github.com/golanguzb70/simple-post-app/config"
	"github.com/golanguzb70/simple-post-app/pkg/db"
	"github.com/golanguzb70/simple-post-app/pkg/logger"
	"github.com/golanguzb70/simple-post-app/storage"
)

func main() {
	cfg := config.Load()
	logger := logger.New(cfg.LogLevel)

	db, err := db.New(cfg)
	if err != nil {
		logger.Error("Error while connecting to database", err)
	} else {
		logger.Info("Successfully connected to database")
	}

	router := api.New(logger, cfg, storage.New(db, logger, cfg))

	if err := router.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatal("failed to run http server", err)
	}
}
