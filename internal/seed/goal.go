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

	users := []string{"auth0|dummy1", "auth0|dummy2", "auth0|dummy3"}

	var goals []models.Goal

	for _, userID := range users {
		notifTime := "09:00"
		goals = append(goals, models.Goal{
			UserID:              userID,
			Title:               "早起きする",
			NotificationTime:    &notifTime,
			NotificationEnabled: true,
			Order:               1,
		})

		goals = append(goals, models.Goal{
			UserID:              userID,
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
