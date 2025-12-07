package models

import (
	"time"

	"github.com/google/uuid"
)

type ExcuseEntry struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID     string    `gorm:"size:255;not null;index;uniqueIndex:idx_user_goal_date"`
	GoalID     uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_user_goal_date"`
	Date       string    `gorm:"type:date;not null;uniqueIndex:idx_user_goal_date"` // YYYY-MM-DD
	ExcuseText string    `gorm:"type:text;not null"`
	TemplateID *string   `gorm:"size:255"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
