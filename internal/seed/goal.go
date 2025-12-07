package seed

import (
	"what-went-wrong-api/internal/models"

	"gorm.io/gorm"
)

func SeedGoals(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.Goal{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		return err
	}

	var goals []models.Goal

	for _, user := range users {
		notifTime := "09:00"
		goals = append(goals, models.Goal{
			UserID:              user.ID,
			Title:               "早起きする",
			NotificationTime:    &notifTime,
			NotificationEnabled: true,
			Order:               1,
		})

		goals = append(goals, models.Goal{
			UserID:              user.ID,
			Title:               "Go言語を勉強する",
			NotificationEnabled: false,
			Order:               2,
		})
	}

	if len(goals) > 0 {
		return db.Create(&goals).Error
	}

	return nil
}
