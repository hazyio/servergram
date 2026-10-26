package db

import "time"

type Sessions struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;index:idx_sessions_user_created,priority:1"`
	Token     string    `gorm:"size:64;uniqueIndex;not null"`
	Device    string    `gorm:"size:255;not null"`
	Os        string    `gorm:"size:255;not null"`
	IpAddress string    `gorm:"size:45;not null"`
	CreatedAt time.Time `gorm:"index:idx_sessions_user_created,priority:2"`
	UpdatedAt time.Time // set automatically
}
