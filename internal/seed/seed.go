package seed

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

func Run(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := SeedUsers(tx); err != nil {
			return fmt.Errorf("failed to seed users: %w", err)
		}

		log.Println("Database seeding completed successfully")
		return nil
	})
}
