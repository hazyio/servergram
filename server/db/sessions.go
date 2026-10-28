package db

import "time"

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
