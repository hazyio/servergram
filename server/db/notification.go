package db

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	ID        int    `gorm:"primaryKey"`
	Title     string `gorm:"not null"`
	Message   string `gorm:"not null"`
	Read      bool   `gorm:"not null"`
	Url       *string
	CreatedAt time.Time
}

func NewNotificationAvailable(database *gorm.DB) (bool, error) {
	var db Notification
	err := database.Model(&Notification{}).Where("read = ?", false).First(&db).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return false, err
	}
	return db.ID != 0, nil
}

func MarkNotificationsAsRead(database *gorm.DB, ids []int) error {
	if len(ids) == 0 {
		return nil
	}
	err := database.Model(&Notification{}).Where("id IN ?", ids).Update("read", true).Error
	if err != nil {
		return err
	}
	return nil
}
