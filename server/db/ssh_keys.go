package db

import (
	"time"

	"gorm.io/gorm"
)

type SSHKey struct {
	ID        uint     `gorm:"primaryKey"`
	Name      string   `gorm:"size:100;not null;UniqueIndex"`
	Key       string   `gorm:"size:16384;not null"`
	Servers   []Server `gorm:"foreignKey:SSHKeyID"`
	Encrypted bool     `gorm:"not null"`
	CreatedAt time.Time
}

func FindSSHKeyByID(id uint, db *gorm.DB) (*SSHKey, error) {
	var key SSHKey
	if err := db.First(&key, id).Error; err != nil {
		return nil, err
	}
	return &key, nil
}
