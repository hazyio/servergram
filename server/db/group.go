package db

import "time"

type Group struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:100;not null;UniqueIndex"`
	Servers   []Server `gorm:"foreignKey:GroupID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
