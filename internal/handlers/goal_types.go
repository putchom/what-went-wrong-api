package handlers

import "time"

type GoalResponse struct {
	ID                  string    `json:"id"`
	Title               string    `json:"title"`
	NotificationTime    *string   `json:"notificationTime,omitempty"`
	NotificationEnabled bool      `json:"notificationEnabled"`
	Order               int       `json:"order"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
}

type GetGoalsResponse struct {
	Goals []GoalResponse `json:"goals"`
}

type CreateGoalRequest struct {
	Title               string  `json:"title" binding:"required,max=200"`
	NotificationTime    *string `json:"notificationTime"`
	NotificationEnabled bool    `json:"notificationEnabled"`
}

type CreateGoalResponse struct {
	Goal GoalResponse `json:"goal"`
}

type UpdateGoalRequest struct {
	Title               *string `json:"title" binding:"omitempty,max=200"`
	NotificationTime    *string `json:"notificationTime"`
	NotificationEnabled *bool   `json:"notificationEnabled"`
}
