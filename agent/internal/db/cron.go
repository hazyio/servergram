package db

import (
	"time"

	"gorm.io/gorm"
)

type CronJobStatus string

const (
	CronJobStatusFailed    CronJobStatus = "failed"
	CronJobStatusSuccess   CronJobStatus = "success"
	CronJobStatusPaused    CronJobStatus = "paused"
	CronJobStatusRunning   CronJobStatus = "running"
	CronJobStatusScheduled CronJobStatus = "scheduled"
)

type CronJob struct {
	ID               uint64 `gorm:"primaryKey"`
	Name             string
	Hash             string `gorm:"uniqueIndex:idx_server_hash"`
	Managed          bool
	Status           CronJobStatus `gorm:"type:varchar(20)"`
	MetaKeepCount    int32
	RunLogKeepCount  int32
	ErrorLogKeepDays int32
	LastRunAt        time.Time
	NextRunAt        time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Runs             []CronRun `gorm:"foreignKey:JobID"`
}

func FindCronJobByID(db *gorm.DB, id uint64) (*CronJob, error) {
	var job CronJob
	if err := db.First(&job, id).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

func FindCronJobByHash(db *gorm.DB, hash string) (*CronJob, error) {
	var job CronJob
	if err := db.Where("hash = ?", hash).First(&job).Error; err != nil {
		return nil, err
	}
	return &job, nil
}
