package handlers

import (
	"net/http"
	"what-went-wrong-api/internal/services"

	"github.com/gin-gonic/gin"
)

type AIHandler struct {
	aiService services.AIService
}

func NewAIHandler(aiService services.AIService) *AIHandler {
	return &AIHandler{
		aiService: aiService,
	}
}

type CreateAiExcuseRequest struct {
	GoalID  string `json:"goalId" binding:"required"`
	Date    string `json:"date" binding:"required"`
	Tone    string `json:"tone"`
	Context string `json:"context"`
}

type CreateAiExcuseResponse struct {
	Candidates []string `json:"candidates"`
}

// PostAiExcuse godoc
// @Summary Generate AI excuses
// @Description Generate excuse candidates using AI. Requires premium plan.
// @Tags ai
// @Accept json
// @Produce json
// @Param request body CreateAiExcuseRequest true "Request body"
// @Success 200 {object} CreateAiExcuseResponse
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string "Forbidden if not premium"
// @Failure 500 {object} map[string]string
// @Router /ai-excuse [post]
func (h *AIHandler) PostAiExcuse(c *gin.Context) {
	// UserID extraction kept if needed for future logic (e.g. logging), otherwise remove or underscore
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	entitlementsInterface, exists := c.Get("entitlements")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Entitlements not found in context"})
		return
	}
	entitlements := entitlementsInterface.(services.Entitlements)

	if !entitlements.CanUseAiExcuse {
		c.JSON(http.StatusForbidden, gin.H{"error": "This feature requires a premium plan"})
		return
	}

	var req CreateAiExcuseRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	candidates, err := h.aiService.GenerateExcuse(req.Tone, req.Context)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate excuses"})
		return
	}

	c.JSON(http.StatusOK, CreateAiExcuseResponse{Candidates: candidates})
}
