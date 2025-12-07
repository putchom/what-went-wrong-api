package models

import (
	"time"

	"github.com/google/uuid"
)

type ExcuseEntry struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;index"`
	GoalID     uuid.UUID `gorm:"type:uuid;not null;index"`
	Date       string    `gorm:"type:date;not null;index:idx_user_goal_date,unique"` // YYYY-MM-DD
	ExcuseText string    `gorm:"type:text;not null"`
	TemplateID *string   `gorm:"size:255"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
