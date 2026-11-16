package db

import (
	"time"

	"gorm.io/gorm"
)

type CronRunStatus string

const (
	RunStatusFailed    CronRunStatus = "failed"
	RunStatusSuccess   CronRunStatus = "success"
	RunStatusRunning   CronRunStatus = "running"
	RunStatusCancelled CronRunStatus = "cancelled"
)

type CronRunTrigger string

const (
	CronRunTriggerManual    CronRunTrigger = "manual"
	CronRunTriggerScheduled CronRunTrigger = "scheduled"
)

type CronRun struct {
	ID         int64 `gorm:"primaryKey"`
	CronJobID      int64
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

func (c *CronRun) UpdateStatus(database *gorm.DB, status CronRunStatus) error {
	c.Status = status
	return database.Save(c).Error
}
