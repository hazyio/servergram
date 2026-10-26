package db

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey"`
	FirstName    string    `gorm:"size:100;not null"`
	LastName     string    `gorm:"size:100;not null"`
	Email        string    `gorm:"size:255;uniqueIndex;not null"`
	PasswordHash string    `gorm:"size:255;not null"`
	CreatedAt    time.Time // set automatically
	UpdatedAt    time.Time // set automatically
}
