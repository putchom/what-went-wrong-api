package models

import (
	"time"
)

type UserPlan struct {
	UserID    string `gorm:"size:255;primaryKey"`
	Plan      string `gorm:"size:50;not null;default:'free'"` // "free", "premium"
	ExpiresAt *time.Time
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
