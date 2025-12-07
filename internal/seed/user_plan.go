package seed

import (
	"what-went-wrong-api/internal/models"

	"gorm.io/gorm"
)

func SeedUserPlans(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.UserPlan{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		return err
	}

	var plans []models.UserPlan
	for _, user := range users {
		planType := "free"
		if user.Email == "jiro@example.com" {
			planType = "premium"
		}

		plans = append(plans, models.UserPlan{
			UserID: user.ID,
			Plan:   planType,
		})
	}

	if len(plans) > 0 {
		return db.Create(&plans).Error
	}

	return nil
}
