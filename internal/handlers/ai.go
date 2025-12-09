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

// PostAiExcuse godoc
// @Summary Generate AI excuses
// @Description Generate excuse candidates using AI. Requires premium plan.
// @Tags ai
// @Accept json
// @Produce json
// @Param request body CreateAiExcuseRequest true "Request body"
// @Success 200 {object} CreateAiExcuseResponse
// @Failure 400 {object} ValidationErrorResponse
// @Failure 401 {object} AiUnauthorizedResponse
// @Failure 403 {object} PremiumRequiredResponse "Forbidden if not premium"
// @Failure 500 {object} InternalErrorResponse
// @Security BearerAuth
// @Router /ai-excuse [post]
func (h *AIHandler) PostAiExcuse(c *gin.Context) {
	// UserID extraction kept if needed for future logic (e.g. logging), otherwise remove or underscore
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証されていません"})
		return
	}

	entitlementsInterface, exists := c.Get("entitlements")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "プラン情報の取得に失敗しました"})
		return
	}
	entitlements := entitlementsInterface.(services.Entitlements)

	if !entitlements.CanUseAiExcuse {
		c.JSON(http.StatusForbidden, gin.H{"error": "この機能を利用するにはプレミアムプランが必要です"})
		return
	}

	var req CreateAiExcuseRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "入力内容が正しくありません"})
		return
	}

	candidates, err := h.aiService.GenerateExcuse(req.Tone, req.Context)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI言い訳の生成に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, CreateAiExcuseResponse{Candidates: candidates})
}
