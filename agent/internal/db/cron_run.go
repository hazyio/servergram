package db

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"hazyio.com/servergram/agent/internal/common"
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
	ID         uint64 `gorm:"primaryKey"`
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
func FindCronJobRunByID(database *gorm.DB, id uint64) (*CronRun, error) {
	run := &CronRun{}
	err := database.First(run, id).Error
	if err != nil {
		return nil, err
	}
	return run, nil
}
func PruneCronRuns(db *gorm.DB, keep int32, cronJobID uint64) error {
	var cutoff CronRun
	err := db.Select("id").
		Where("cron_job_id = ? AND status != ?", cronJobID, RunStatusRunning).
		Order("id DESC").
		Offset(int(keep)).
		Limit(1).
		Take(&cutoff).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	deleted := make([]struct {
		StartedAt time.Time
		ID        uint64
	}, 0)
	err = db.
		Model(&CronRun{}).
		Where("cron_job_id = ? AND status != ? AND id <= ?", cronJobID, RunStatusRunning, cutoff.ID).
		Delete(&CronRun{}).Take(&deleted).Error
	if err != nil {
		return err
	}
	common.RemoveCronRunLog(deleted)

	return nil
}

func GetCronJobRunStatus(database *gorm.DB, id uint64) (CronRunStatus, error) {
	var status CronRunStatus
	err := database.Model(&CronRun{}).Select("status").Where("id = ?", id).Scan(&status).Error
	if err != nil {
		return "", err
	}
	return status, nil
}
