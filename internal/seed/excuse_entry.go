package seed

import (
	"time"
	"what-went-wrong-api/internal/models"

	"gorm.io/gorm"
)

func SeedExcuseEntries(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.ExcuseEntry{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	var goals []models.Goal
	if err := db.Preload("User").Find(&goals).Error; err != nil {
		return err
	}

	// Assuming templates are already seeded
	templateID := "gravity-strong"

	var entries []models.ExcuseEntry

	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	for i, goal := range goals {
		// Create entry for today for some goals
		if i%2 == 0 {
			entries = append(entries, models.ExcuseEntry{
				UserID:     goal.UserID,
				GoalID:     goal.ID,
				Date:       today,
				ExcuseText: "It was too hard.",
				TemplateID: nil,
			})
		} else {
			// Create entry for yesterday
			entries = append(entries, models.ExcuseEntry{
				UserID:     goal.UserID,
				GoalID:     goal.ID,
				Date:       yesterday,
				ExcuseText: "Gravity was unusually strong.",
				TemplateID: &templateID,
			})
		}
	}

	if len(entries) > 0 {
		return db.Create(&entries).Error
	}

	return nil
}
