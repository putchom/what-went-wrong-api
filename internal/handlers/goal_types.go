package handlers

import "time"

type GoalResponse struct {
	ID                  string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Title               string    `json:"title" example:"本を10ページ読む"`
	NotificationTime    *string   `json:"notificationTime,omitempty" example:"20:00"`
	NotificationEnabled bool      `json:"notificationEnabled" example:"true"`
	Order               int       `json:"order" example:"1"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
}

type GetGoalsResponse struct {
	Goals []GoalResponse `json:"goals"`
}

type CreateGoalRequest struct {
	Title               string  `json:"title" binding:"required,max=200" example:"本を10ページ読む"`
	NotificationTime    *string `json:"notificationTime" example:"20:00"`
	NotificationEnabled bool    `json:"notificationEnabled" example:"true"`
}

type CreateGoalResponse struct {
	Goal GoalResponse `json:"goal"`
}

type UpdateGoalRequest struct {
	Title               *string `json:"title" binding:"omitempty,max=200" example:"本を20ページ読む"`
	NotificationTime    *string `json:"notificationTime" example:"21:00"`
	NotificationEnabled *bool   `json:"notificationEnabled" example:"false"`
}

type GoalLimitReachedResponse struct {
	Error string `json:"error" example:"プランの目標作成数上限に達しました"`
}

type GoalNotFoundErrorResponse struct {
	Error string `json:"error" example:"目標が見つかりません"`
}

type GoalFetchErrorResponse struct {
	Error string `json:"error" example:"目標の取得に失敗しました"`
}

type GoalCreateErrorResponse struct {
	Error string `json:"error" example:"目標の作成に失敗しました"`
}

type GoalUpdateErrorResponse struct {
	Error string `json:"error" example:"目標の更新に失敗しました"`
}

type GoalDeleteErrorResponse struct {
	Error string `json:"error" example:"目標の削除に失敗しました"`
}

type GoalValidationErrorResponse struct {
	Error string `json:"error" example:"入力内容が正しくありません"`
}
