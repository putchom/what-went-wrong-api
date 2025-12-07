package handlers

import "time"

type ExcuseTemplateResponse struct {
	ID         string    `json:"id"`
	PackID     string    `json:"packId"`
	ExcuseText string    `json:"excuseText"`
	Tags       []string  `json:"tags"`
	IsPremium  bool      `json:"isPremium"`
	CreatedAt  time.Time `json:"createdAt"`
}

type GetExcuseTemplatesResponse struct {
	Templates []ExcuseTemplateResponse `json:"templates"`
}
