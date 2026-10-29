package db

import (
	"time"

	"gorm.io/gorm"
)

type Group struct {
	ID        uint     `gorm:"primaryKey"`
	Name      string   `gorm:"size:100;not null;UniqueIndex"`
	Servers   []Server `gorm:"foreignKey:GroupID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func FindGroupByID(id uint, db *gorm.DB) (*Group, error) {
	var group Group
	if err := db.First(&group, id).Error; err != nil {
		return nil, err
	}
	return &group, nil
}
