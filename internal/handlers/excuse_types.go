package handlers

import (
	"time"

	"github.com/google/uuid"
)

type CreateExcuseRequest struct {
	Date       string `json:"date" binding:"required" example:"2023-10-27"` // YYYY-MM-DD
	ExcuseText string `json:"excuseText" binding:"required,max=500" example:"I overslept."`
	TemplateID string `json:"templateId" example:"template_123"`
}

type UpdateExcuseRequest struct {
	ExcuseText string `json:"excuseText" binding:"max=500" example:"I overslept a lot."`
	TemplateID string `json:"templateId" example:"template_123"`
}

type ExcuseResponse struct {
	ID         uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	GoalID     uuid.UUID `json:"goalId" example:"550e8400-e29b-41d4-a716-446655440001"`
	Date       string    `json:"date" example:"2023-10-27"`
	ExcuseText string    `json:"excuseText" example:"I overslept."`
	TemplateID *string   `json:"templateId,omitempty" example:"template_123"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type GetExcusesResponse struct {
	Excuses []ExcuseResponse `json:"excuses"`
}

type ExcuseValidationErrorResponse struct {
	Error string `json:"error" example:"Key: 'CreateExcuseRequest.Date' Error:Field validation for 'Date' failed on the 'required' tag"`
}

type ExcuseInternalErrorResponse struct {
	Error string `json:"error" example:"Failed to create excuse"`
}

type ExcuseForbiddenResponse struct {
	Error string `json:"error" example:"Premium template entitlement required"`
}

type ExcuseNotFoundResponse struct {
	Error string `json:"error" example:"Excuse not found"`
}
