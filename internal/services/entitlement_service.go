package services

import (
	"errors"
	"what-went-wrong-api/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EntitlementService struct {
	db *gorm.DB
}

func NewEntitlementService(db *gorm.DB) *EntitlementService {
	return &EntitlementService{db: db}
}

func (s *EntitlementService) CanUseAiExcuse(userID uuid.UUID) (bool, error) {
	var plan models.UserPlan
	if err := s.db.First(&plan, "user_id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Default to free plan if not found
			return false, nil
		}
		return false, err
	}

	return plan.Plan == "premium", nil
}
