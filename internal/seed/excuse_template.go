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
			Text:      "今日は重力が強かった。",
			PackID:    "core",
			IsActive:  true,
			IsPremium: false,
			Tags:      pq.StringArray{"物理", "面白い"},
		},
		{
			ID:        "cat-monitor",
			Text:      "猫がモニターの上で寝てしまった。",
			PackID:    "core",
			IsActive:  true,
			IsPremium: false,
			Tags:      pq.StringArray{"猫", "かわいい"},
		},
		{
			ID:        "coffee-spill",
			Text:      "キーボードにコーヒーをこぼした。",
			PackID:    "core",
			IsActive:  true,
			IsPremium: false,
			Tags:      pq.StringArray{"事故", "コーヒー"},
		},
		{
			ID:        "aliens",
			Text:      "エイリアンにやる気を奪われた。",
			PackID:    "surreal",
			IsActive:  true,
			IsPremium: true,
			Tags:      pq.StringArray{"SF", "エイリアン"},
		},
	}

	return db.Create(&templates).Error
}
