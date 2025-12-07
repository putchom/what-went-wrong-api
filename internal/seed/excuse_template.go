package seed

import (
	"what-went-wrong-api/internal/models"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

func SeedExcuseTemplates(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.ExcuseTemplate{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	templates := []models.ExcuseTemplate{
		{
			ID:        "gravity-strong",
			Text:      "Today, gravity was unusually strong.",
			PackID:    "core",
			IsActive:  true,
			IsPremium: false,
			Tags:      pq.StringArray{"physics", "funny"},
		},
		{
			ID:        "cat-monitor",
			Text:      "My cat fell asleep on my monitor.",
			PackID:    "core",
			IsActive:  true,
			IsPremium: false,
			Tags:      pq.StringArray{"cat", "cute"},
		},
		{
			ID:        "coffee-spill",
			Text:      "I spilled coffee on my keyboard.",
			PackID:    "core",
			IsActive:  true,
			IsPremium: false,
			Tags:      pq.StringArray{"accident", "coffee"},
		},
		{
			ID:        "aliens",
			Text:      "Aliens abducted my motivation.",
			PackID:    "surreal",
			IsActive:  true,
			IsPremium: true,
			Tags:      pq.StringArray{"scifi", "alien"},
		},
	}

	return db.Create(&templates).Error
}
