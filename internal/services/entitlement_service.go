package services

import (
	"errors"
	"time"
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

func (s *EntitlementService) GetPlan(userID uuid.UUID) (*models.UserPlan, error) {
	var plan models.UserPlan
	if err := s.db.First(&plan, "user_id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create default free plan if not exists
			newPlan := models.UserPlan{
				UserID: userID,
				Plan:   "free",
			}
			if err := s.db.Create(&newPlan).Error; err != nil {
				return nil, err
			}
			return &newPlan, nil
		}
		return nil, err
	}
	return &plan, nil
}

func (s *EntitlementService) UpdatePlan(userID uuid.UUID, planName string) (*models.UserPlan, error) {
	if planName != "free" && planName != "premium" {
		return nil, errors.New("invalid plan name")
	}

	var plan models.UserPlan
	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&plan, "user_id = ?", userID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				plan = models.UserPlan{
					UserID:    userID,
					Plan:      planName,
					UpdatedAt: time.Now(),
				}
				return tx.Create(&plan).Error
			}
			return err
		}

		plan.Plan = planName
		plan.UpdatedAt = time.Now()
		return tx.Save(&plan).Error
	})

	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (s *EntitlementService) GetEntitlements(planName string) Entitlements {
	if planName == "premium" {
		return Entitlements{
			MaxGoals:               100, // Practically unlimited
			LogRetentionDays:       nil, // Unlimited
			CanUseAiExcuse:         true,
			CanUsePremiumTemplates: true,
		}
	}

	// Default to free
	retention := 30
	return Entitlements{
		MaxGoals:               3,
		LogRetentionDays:       &retention,
		CanUseAiExcuse:         false,
		CanUsePremiumTemplates: false,
	}
}

func (s *EntitlementService) CanUseAiExcuse(userID uuid.UUID) (bool, error) {
	plan, err := s.GetPlan(userID)
	if err != nil {
		return false, err
	}
	entitlements := s.GetEntitlements(plan.Plan)
	return entitlements.CanUseAiExcuse, nil
}
