package handlers

import "time"

type GoalResponse struct {
	ID                  string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Title               string    `json:"title" example:"Read 10 pages"`
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
	Title               string  `json:"title" binding:"required,max=200" example:"Read 10 pages"`
	NotificationTime    *string `json:"notificationTime" example:"20:00"`
	NotificationEnabled bool    `json:"notificationEnabled" example:"true"`
}

type CreateGoalResponse struct {
	Goal GoalResponse `json:"goal"`
}

type UpdateGoalRequest struct {
	Title               *string `json:"title" binding:"omitempty,max=200" example:"Read 20 pages"`
	NotificationTime    *string `json:"notificationTime" example:"21:00"`
	NotificationEnabled *bool   `json:"notificationEnabled" example:"false"`
}

type GoalLimitReachedResponse struct {
	Error string `json:"error" example:"Maximum goal limit reached for your plan"`
}

type GoalNotFoundErrorResponse struct {
	Error string `json:"error" example:"Goal not found"`
}

type GoalInternalErrorResponse struct {
	Error string `json:"error" example:"Failed to fetch goals"`
}

type GoalValidationErrorResponse struct {
	Error string `json:"error" example:"Key: 'CreateGoalRequest.Title' Error:Field validation for 'Title' failed on the 'required' tag"`
}
