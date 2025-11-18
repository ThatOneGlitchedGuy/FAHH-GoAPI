package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"golang-app/internal/config"
	"golang-app/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db   *gorm.DB
	once sync.Once
)

func InitDB() (*gorm.DB, error) {
	var err error
	once.Do(func() {
		cfg, cfgErr := config.GetSettings()
		if cfgErr != nil {
			err = fmt.Errorf("failed to load configuration: %w", cfgErr)
			return
		}

		gormConfig := &gorm.Config{}
		if cfg.Debug {
			gormConfig.Logger = logger.Default.LogMode(logger.Info)
		} else {
			gormConfig.Logger = logger.Default.LogMode(logger.Silent)
		}

		dbPath := strings.TrimPrefix(cfg.DatabaseURL, "sqlite://")
		dbPath = strings.TrimPrefix(dbPath, "sqlite+aiosqlite://")
		dbPath = strings.TrimPrefix(dbPath, "./")

		dbDir := filepath.Dir(dbPath)
		if dbDir != "" && dbDir != "." {
			if err = os.MkdirAll(dbDir, os.ModePerm); err != nil {
				err = fmt.Errorf("failed to create database directory %s: %w", dbDir, err)
				return
			}
		}

		db, err = gorm.Open(sqlite.Open(dbPath), gormConfig)
		if err != nil {
			err = fmt.Errorf("failed to connect to database: %w", err)
			return
		}

		err = db.AutoMigrate(
			&models.User{},
			&models.Product{},
			&models.Order{},
			&models.OrderItem{},
			&models.Message{},
			&models.Review{},
		)
		if err != nil {
			err = fmt.Errorf("failed to auto migrate database: %w", err)
			return
		}

		log.Println("Database initialized and migrated successfully.")
	})
	return db, err
}

func GetDB() *gorm.DB {
	if db == nil {
		log.Fatal("Database not initialized. Call InitDB() first.")
	}
	return db
}