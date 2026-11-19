package db

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type CronRunStatus string

const (
	RunStatusFailed    CronRunStatus = "failed"
	RunStatusSuccess   CronRunStatus = "success"
	RunStatusRunning   CronRunStatus = "running"
	RunStatusCancelled CronRunStatus = "cancelled"
	RunStatusTimedOut  CronRunStatus = "timed_out"
)

type CronRunTrigger string

const (
	CronRunTriggerManual    CronRunTrigger = "manual"
	CronRunTriggerScheduled CronRunTrigger = "scheduled"
)

type CronRun struct {
	ID         int64 `gorm:"primaryKey"`
	CronJobID  int64
	Status     CronRunStatus
	Trigger    CronRunTrigger
	StartedAt  time.Time
	FinishedAt time.Time
	ExitCode   int
	CronJob    CronJob `gorm:"constraint:OnDelete:CASCADE;"`
}

func (c *CronRun) RetrieveStatus(database *gorm.DB) error {

	return database.First(c).Error
}

func PruneCronRuns(db *gorm.DB, keep int32, cronJobID uint64) error {
	var cutoff CronRun
	err := db.Select("id").Where("status != ?", RunStatusRunning).Where("cron_job_id = ?", cronJobID).Order("id DESC").Offset(int(keep)).Limit(1).Take(&cutoff).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil // fewer than keep+1 rows, nothing to delete
	}
	if err != nil {
		return err
	}
	return db.Where("id <= ?", cutoff.ID).Delete(&CronRun{}).Error
}
