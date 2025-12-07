package handlers

import (
	"time"

	"github.com/google/uuid"
)

type CreateExcuseRequest struct {
	Date       string `json:"date" binding:"required"` // YYYY-MM-DD
	ExcuseText string `json:"excuseText" binding:"required,max=500"`
	TemplateID string `json:"templateId"`
}

type UpdateExcuseRequest struct {
	ExcuseText string `json:"excuseText" binding:"max=500"`
	TemplateID string `json:"templateId"`
}

type ExcuseResponse struct {
	ID         uuid.UUID `json:"id"`
	GoalID     uuid.UUID `json:"goalId"`
	Date       string    `json:"date"`
	ExcuseText string    `json:"excuseText"`
	TemplateID *string   `json:"templateId,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type GetExcusesResponse struct {
	Excuses []ExcuseResponse `json:"excuses"`
}
