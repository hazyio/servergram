package db

import (
	"fmt"
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
	if err := database.AutoMigrate(&CronJob{}, &CronRun{}); err != nil {
		return nil, err
	}
	if err := migrateCronRuns(database); err != nil {
		return nil, fmt.Errorf("migrate cron runs: %w", err)
	}
	return database, nil
}

// migrateCronRuns upgrades databases created before the CronRun.JobID field was
// renamed to CronRun.CronJobID. Older schemas stored the reference in job_id
// without ON DELETE CASCADE, which blocks deleting a CronJob that still has runs.
func migrateCronRuns(database *gorm.DB) error {
	var columns []struct {
		Name string
	}
	if err := database.Raw("PRAGMA table_info(cron_runs)").Scan(&columns).Error; err != nil {
		return err
	}

	var hasJobID, hasCronJobID bool
	for _, column := range columns {
		switch column.Name {
		case "job_id":
			hasJobID = true
		case "cron_job_id":
			hasCronJobID = true
		}
	}
	if !hasJobID {
		return nil
	}

	jobIDExpr := "job_id"
	if hasCronJobID {
		jobIDExpr = "COALESCE(cron_job_id, job_id)"
	}

	return database.Transaction(func(tx *gorm.DB) error {
		stmts := []string{
			`DROP TABLE IF EXISTS cron_runs_migrated`,
			`CREATE TABLE cron_runs_migrated (
				id integer PRIMARY KEY AUTOINCREMENT,
				cron_job_id integer,
				status text,
				trigger text,
				started_at datetime,
				finished_at datetime,
				exit_code integer,
				CONSTRAINT fk_cron_jobs_runs FOREIGN KEY (cron_job_id) REFERENCES cron_jobs(id) ON DELETE CASCADE
			)`,
			`INSERT INTO cron_runs_migrated (id, cron_job_id, status, trigger, started_at, finished_at, exit_code)
				SELECT id, ` + jobIDExpr + `, status, trigger, started_at, finished_at, exit_code FROM cron_runs`,
			`DROP TABLE cron_runs`,
			`ALTER TABLE cron_runs_migrated RENAME TO cron_runs`,
		}
		for _, stmt := range stmts {
			if err := tx.Exec(stmt).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
