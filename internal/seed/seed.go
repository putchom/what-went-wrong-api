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

		if err := SeedUserPlans(tx); err != nil {
			return fmt.Errorf("failed to seed user plans: %w", err)
		}

		if err := SeedGoals(tx); err != nil {
			return fmt.Errorf("failed to seed goals: %w", err)
		}

		if err := SeedExcuseTemplates(tx); err != nil {
			return fmt.Errorf("failed to seed excuse templates: %w", err)
		}

		if err := SeedExcuseEntries(tx); err != nil {
			return fmt.Errorf("failed to seed excuse entries: %w", err)
		}

		log.Println("Database seeding completed successfully")
		return nil
	})
}
