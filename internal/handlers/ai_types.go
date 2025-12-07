package handlers

type CreateAiExcuseRequest struct {
	GoalID  string `json:"goalId" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Date    string `json:"date" binding:"required" example:"2023-10-27"`
	Tone    string `json:"tone" example:"serious"`
	Context string `json:"context" example:"I had a lot of meetings"`
}

type CreateAiExcuseResponse struct {
	Candidates []string `json:"candidates" example:"I had unexpected urgent work.,I wasn't feeling well."`
}

type ValidationErrorResponse struct {
	Error string `json:"error" example:"Key: 'CreateAiExcuseRequest.GoalID' Error:Field validation for 'GoalID' failed on the 'required' tag"`
}

type PremiumRequiredResponse struct {
	Error string `json:"error" example:"This feature requires a premium plan"`
}

type InternalErrorResponse struct {
	Error string `json:"error" example:"Failed to generate excuses"`
}
