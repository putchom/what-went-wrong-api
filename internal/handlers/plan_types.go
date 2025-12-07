package handlers

import "what-went-wrong-api/internal/services"

type GetMePlanResponse struct {
	Plan         string                `json:"plan" example:"premium"`
	Entitlements services.Entitlements `json:"entitlements"` // Entitlements struct might need examples in its own definition if not here
}

type PostMePlanRequest struct {
	Plan string `json:"plan" binding:"required" example:"premium"`
}

type PostMePlanResponse struct {
	Plan         string                `json:"plan" example:"premium"`
	Entitlements services.Entitlements `json:"entitlements"`
}

type PlanInternalErrorResponse struct {
	Error string `json:"error" example:"Failed to get plan"`
}

type PlanValidationErrorResponse struct {
	Error string `json:"error" example:"Key: 'PostMePlanRequest.Plan' Error:Field validation for 'Plan' failed on the 'required' tag"`
}
