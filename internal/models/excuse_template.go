package models

import (
	"time"

	"github.com/lib/pq"
)

type ExcuseTemplate struct {
	ID        string         `gorm:"primaryKey;size:255"` // "gravity-strong" etc
	Text      string         `gorm:"type:text;not null"`
	PackID    string         `gorm:"size:255;default:'core'"` // "core", "pack.surreal", etc
	IsActive  bool           `gorm:"default:true"`
	IsPremium bool           `gorm:"default:false"`
	Tags      pq.StringArray `gorm:"type:text[]"`
	CreatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
}
