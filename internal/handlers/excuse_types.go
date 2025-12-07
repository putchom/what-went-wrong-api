package handlers

import (
	"time"

	"github.com/google/uuid"
)

type CreateExcuseRequest struct {
	Date       string `json:"date" binding:"required" example:"2023-10-27"` // YYYY-MM-DD
	ExcuseText string `json:"excuseText" binding:"required,max=500" example:"寝坊しました。"`
	TemplateID string `json:"templateId" example:"template_123"`
}

type UpdateExcuseRequest struct {
	ExcuseText string `json:"excuseText" binding:"max=500" example:"盛大に寝坊しました。"`
	TemplateID string `json:"templateId" example:"template_123"`
}

type ExcuseResponse struct {
	ID         uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	GoalID     uuid.UUID `json:"goalId" example:"550e8400-e29b-41d4-a716-446655440001"`
	Date       string    `json:"date" example:"2023-10-27"`
	ExcuseText string    `json:"excuseText" example:"寝坊しました。"`
	TemplateID *string   `json:"templateId,omitempty" example:"template_123"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type GetExcusesResponse struct {
	Excuses []ExcuseResponse `json:"excuses"`
}

type ExcuseValidationErrorResponse struct {
	Error string `json:"error" example:"入力内容が正しくありません"`
}

type ExcuseFetchErrorResponse struct {
	Error string `json:"error" example:"言い訳の取得に失敗しました"`
}

type ExcuseCreateErrorResponse struct {
	Error string `json:"error" example:"言い訳の作成に失敗しました"`
}

type ExcuseUpdateErrorResponse struct {
	Error string `json:"error" example:"言い訳の更新に失敗しました"`
}

type ExcuseDeleteErrorResponse struct {
	Error string `json:"error" example:"言い訳の削除に失敗しました"`
}

type ExcuseUnauthorizedResponse struct {
	Error string `json:"error" example:"認証されていません"`
}

type ExcuseForbiddenResponse struct {
	Error string `json:"error" example:"プレミアムテンプレートを利用するにはプレミアムプランが必要です"`
}

type ExcuseNotFoundResponse struct {
	Error string `json:"error" example:"言い訳が見つかりません"`
}
