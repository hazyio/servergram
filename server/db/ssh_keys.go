package db

import "time"

type SSHKey struct {
	ID      uint     `gorm:"primaryKey"`
	Name    string   `gorm:"size:100;not null;UniqueIndex"`
	Key     string   `gorm:"size:16384;not null"`
	Servers []Server `gorm:"foreignKey:SSHKeyID"`
	Encrypted bool `gorm:"not null"`
	CreatedAt time.Time
}
