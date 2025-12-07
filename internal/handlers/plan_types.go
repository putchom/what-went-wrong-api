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

type PlanFetchErrorResponse struct {
	Error string `json:"error" example:"プランの取得に失敗しました"`
}

type PlanUpdateErrorResponse struct {
	Error string `json:"error" example:"プランの更新に失敗しました"`
}

type PlanValidationErrorResponse struct {
	Error string `json:"error" example:"入力内容が正しくありません"`
}
