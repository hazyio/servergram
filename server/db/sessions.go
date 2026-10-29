package db

import (
	"errors"
	"log/slog"
	"time"

	"gorm.io/gorm"
)

type Sessions struct {
	ID        uint      `gorm:"primaryKey"`
	Token     string    `gorm:"size:64;uniqueIndex;not null"`
	Os        string    `gorm:"size:255;not null"`
	UserAgent string    `gorm:"size:255;not null"`
	Expires   time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"index:idx_sessions_user_created,priority:2"`
	UpdatedAt time.Time // set automatically
}

func (s *Sessions) IsExpired() bool {
	return time.Now().After(s.Expires)
}

func TokenIsValid(hashed_token string, database *gorm.DB) bool {
	find_session := Sessions{}
	if err := database.Where("token = ?", hashed_token).First(&find_session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Warn("Session not found")
			return false
		}
		slog.Error("An internal error occurred while trying to check the session", "error", err)
		return false
	}
	if find_session.IsExpired() {
		slog.Warn("Session expired")
		return false
	}
	return !find_session.IsExpired()
}
