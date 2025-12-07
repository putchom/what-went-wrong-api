package models

import (
	"time"

	"github.com/google/uuid"
)

type UserPlan struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	Plan      string    `gorm:"size:50;not null;default:'free'"` // "free", "premium"
	ExpiresAt *time.Time
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
