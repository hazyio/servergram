package db

import (
	"log/slog"
	"time"

	"gorm.io/gorm"
)

type CronJobStatus string

const (
	CronJobStatusInactive CronJobStatus = "inactive"
	CronJobStatusActive   CronJobStatus = "active"
)

type CronJob struct {
	ID              uint64 `gorm:"primaryKey"`
	Name            string `gorm:"size:255;not null"`
	Hash            string        `gorm:"uniqueIndex:idx_server_hash"`
	Managed         bool          // This might be updated on demand
	Status          CronJobStatus `gorm:"type:varchar(20)"`
	RunLogKeepCount int32
	LastRunAt       time.Time
	NextRunAt       time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Runs            []CronRun `gorm:"foreignKey:CronJobID;constraint:OnDelete:CASCADE"`
}

func FindCronJobByID(db *gorm.DB, id uint64) (*CronJob, error) {
	var job CronJob
	if err := db.First(&job, id).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

func FindCronJobByHash(db *gorm.DB, hash string) (*CronJob, error) {
	slog.Debug("finding cron job by hash", "hash", hash)

	var job CronJob
	if err := db.Where("hash = ?", hash).First(&job).Error; err != nil {
		return nil, err
	}
	return &job, nil
}
