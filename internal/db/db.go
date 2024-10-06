package db

import (
	"github.com/glebarez/sqlite"
	"github.com/nodeboxhq/nodebox-dashboard/internal"
	"github.com/nodeboxhq/nodebox-dashboard/internal/db/models"
	"github.com/nodeboxhq/nodebox-dashboard/internal/logger"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func SetupDatabase(cfg *internal.NodeboxConfig) *gorm.DB {
	ormConfig := &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	}

	db, err := gorm.Open(sqlite.Open(cfg.DataPath+"/nodebox.db"), ormConfig)

	if err != nil {
		logger.L.Fatal().Msgf("Error connecting to database: %v", err)
	}

	db.Exec("PRAGMA foreign_keys = ON")
	db.Exec("PRAGMA journal_mode=WAL")
	db.Exec("PRAGMA synchronous=NORMAL")
	db.Exec("PRAGMA auto_vacuum=FULL")

	err = db.AutoMigrate(
		&models.Stats{},
		&models.Host{},
	)

	if err != nil {
		logger.L.Fatal().Msgf("Error migrating database: %v", err)
	}

	return db
}
