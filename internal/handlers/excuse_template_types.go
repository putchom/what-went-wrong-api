package handlers

import "time"

type ExcuseTemplateResponse struct {
	ID         string    `json:"id" example:"template_123"`
	PackID     string    `json:"packId" example:"pack_abc"`
	ExcuseText string    `json:"excuseText" example:"宿題を犬に食べられました。"`
	Tags       []string  `json:"tags" example:"面白い,定番"`
	IsPremium  bool      `json:"isPremium" example:"false"`
	CreatedAt  time.Time `json:"createdAt"`
}

type GetExcuseTemplatesResponse struct {
	Templates []ExcuseTemplateResponse `json:"templates"`
}

type TemplateUnauthorizedResponse struct {
	Error string `json:"error" example:"認証されていません"`
}

type TemplateInternalErrorResponse struct {
	Error string `json:"error" example:"テンプレートの取得に失敗しました"`
}

type TemplateNotFoundResponse struct {
	Error string `json:"error" example:"テンプレートが見つかりません"`
}
