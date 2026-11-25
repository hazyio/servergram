package db

import (
	"log/slog"
	"time"

	"gorm.io/gorm"
)

type LogAction string

const (
	LogActionDelete LogAction = "delete"
	LogActionCreate LogAction = "create"
	LogActionAdd    LogAction = "add"
	LogActionUpdate LogAction = "update"
	LogActionRead   LogAction = "read"
)

type Log struct {
	ID        uint      `gorm:"primaryKey"`
	Action    LogAction `gorm:"not null;size:20;check:chk_log_action,action IN ('delete','create','add','update')"`
	Message   string    `gorm:"not null;size:255"`
	CreatedAt time.Time
}

// Add log asynchronously to the database
func AddLog(action LogAction, message string, database *gorm.DB) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("panic", "r", r)
		}
	}()
	go func() {
		log := &Log{
			Action:  action,
			Message: message,
		}
		err := database.Create(log).Error
		if err != nil {
			// handle error if needed
			slog.Error("failed to add log", "error", err)
		}
	}()
}
