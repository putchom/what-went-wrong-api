package handlers

import "time"

type ExcuseTemplateResponse struct {
	ID         string    `json:"id" example:"template_123"`
	PackID     string    `json:"packId" example:"pack_abc"`
	ExcuseText string    `json:"excuseText" example:"My dog ate my homework."`
	Tags       []string  `json:"tags" example:"funny,classic"`
	IsPremium  bool      `json:"isPremium" example:"false"`
	CreatedAt  time.Time `json:"createdAt"`
}

type GetExcuseTemplatesResponse struct {
	Templates []ExcuseTemplateResponse `json:"templates"`
}

type TemplateInternalErrorResponse struct {
	Error string `json:"error" example:"Failed to fetch templates"`
}

type TemplateNotFoundResponse struct {
	Error string `json:"error" example:"Template not found"`
}
