package handlers

type CreateAiExcuseRequest struct {
	GoalID  string `json:"goalId" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Date    string `json:"date" binding:"required" example:"2023-10-27"`
	Tone    string `json:"tone" example:"真面目"`
	Context string `json:"context" example:"会議が多すぎました。"`
}

type CreateAiExcuseResponse struct {
	Candidates []string `json:"candidates" example:"急な緊急対応が入りました。,体調が優れませんでした。"`
}

type AiUnauthorizedResponse struct {
	Error string `json:"error" example:"認証されていません"`
}

type ValidationErrorResponse struct {
	Error string `json:"error" example:"入力内容が正しくありません"`
}

type PremiumRequiredResponse struct {
	Error string `json:"error" example:"この機能を利用するにはプレミアムプランが必要です"`
}

type InternalErrorResponse struct {
	Error string `json:"error" example:"AI言い訳の生成に失敗しました"`
}
