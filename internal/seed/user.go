package seed

import (
	"what-went-wrong-api/internal/models"

	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.User{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	users := []models.User{
		{Name: "Taro Yamada", Email: "taro@example.com", Auth0ID: "auth0|dummy1"},
		{Name: "Hanako Suzuki", Email: "hanako@example.com", Auth0ID: "auth0|dummy2"},
		{Name: "Jiro Sato", Email: "jiro@example.com", Auth0ID: "auth0|dummy3"},
	}

	return db.Create(&users).Error
}
