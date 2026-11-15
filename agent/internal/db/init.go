package db

import (
	"log"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func InitDatabase() (*gorm.DB, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(homeDir, ".servergram")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	url := filepath.Join(dir, "servergram.db")
	database, err := gorm.Open(sqlite.Open(url), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, err
	}
	sqlDB, err := database.DB()
	if err != nil {
		return nil, err
	}
	// SQLite serializes writes: one connection avoids "database is locked"
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetConnMaxLifetime(0) // never expire; an expired in-memory DB would vanish
	// WAL mode and a busy timeout so concurrent readers and a writer don't collide
	database.Exec("PRAGMA journal_mode = WAL")
	database.Exec("PRAGMA busy_timeout = 5000")
	database.Exec("PRAGMA foreign_keys = ON")
	if err := database.AutoMigrate(&CronJob{}); err != nil {
		log.Fatal(err)
	}
	return database, nil
}
