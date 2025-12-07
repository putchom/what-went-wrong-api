package handlers

import "what-went-wrong-api/internal/services"

type GetMePlanResponse struct {
	Plan         string                `json:"plan"`
	Entitlements services.Entitlements `json:"entitlements"`
}

type PostMePlanRequest struct {
	Plan string `json:"plan" binding:"required"`
}

type PostMePlanResponse struct {
	Plan         string                `json:"plan"`
	Entitlements services.Entitlements `json:"entitlements"`
}
