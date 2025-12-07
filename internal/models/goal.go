package models

import (
	"time"

	"github.com/google/uuid"
)

type Goal struct {
	ID                  uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID              string    `gorm:"size:255;not null;index"`
	Title               string    `gorm:"size:255;not null"`
	NotificationTime    *string   `gorm:"size:5"` // "HH:MM" format
	NotificationEnabled bool      `gorm:"default:false"`
	Order               int       `gorm:"default:0"`
	CreatedAt           time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt           time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
