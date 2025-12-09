package handlers

import (
	"errors"
	"net/http"
	"what-went-wrong-api/internal/models"
	"what-went-wrong-api/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ExcuseTemplateHandler struct {
	db *gorm.DB
}

func NewExcuseTemplateHandler(db *gorm.DB) *ExcuseTemplateHandler {
	return &ExcuseTemplateHandler{db: db}
}

// GetTemplates godoc
// @Summary List excuse templates
// @Description Get excuse templates. Can filter by pack_id. Premium users can access all. Free users restricted from premium packs.
// @Tags excuse-templates
// @Accept json
// @Produce json
// @Param pack_id query string false "Pack ID to filter"
// @Success 200 {object} GetExcuseTemplatesResponse
// @Security BearerAuth
// @Router /excuse-templates [get]
func (h *ExcuseTemplateHandler) GetExcuseTemplates(c *gin.Context) {
	packID := c.Query("pack_id")

	// Entitlement check
	entitlementsInterface, exists := c.Get("entitlements")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "プラン情報の取得に失敗しました"})
		return
	}
	entitlements := entitlementsInterface.(services.Entitlements)

	query := h.db.Model(&models.ExcuseTemplate{})

	if packID != "" {
		// If requesting a premium pack, check entitlement
		// For now, let's assume packs starting with "premium-" are premium, or we check specific IDs.
		query = query.Where("pack_id = ?", packID)
	}

	// Filter premium if not entitled
	// Actually spec says "free users restricted from premium packs".
	// Implementation: list all but maybe show isPremium? Or hide premium templates?
	// Let's assume hiding premium templates for free users if that's the requirement, or just returning all so they can see what they are missing.
	// Re-reading docstring: "Free users restricted from premium packs".
	if !entitlements.CanUsePremiumTemplates {
		query = query.Where("is_premium = ?", false)
	}

	var templates []models.ExcuseTemplate
	if err := query.Find(&templates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "テンプレートの取得に失敗しました"})
		return
	}

	res := GetExcuseTemplatesResponse{Templates: make([]ExcuseTemplateResponse, len(templates))}
	for i, t := range templates {
		res.Templates[i] = ExcuseTemplateResponse{
			ID:         t.ID,
			PackID:     t.PackID,
			ExcuseText: t.Text,
			Tags:       t.Tags,
			IsPremium:  t.IsPremium,
			CreatedAt:  t.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, res)
}

// GetExcuseTemplate godoc
// @Summary Get template details
// @Tags excuse-templates
// @Accept json
// @Produce json
// @Param id path string true "Template ID"
// @Success 200 {object} ExcuseTemplateResponse
// @Failure 401 {object} TemplateUnauthorizedResponse
// @Failure 404 {object} TemplateNotFoundResponse
// @Failure 500 {object} TemplateInternalErrorResponse
// @Security BearerAuth
// @Router /excuse-templates/{id} [get]
func (h *ExcuseTemplateHandler) GetExcuseTemplate(c *gin.Context) {
	id := c.Param("id")

	entitlementsInterface, exists := c.Get("entitlements")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "プラン情報の取得に失敗しました"})
		return
	}
	entitlements := entitlementsInterface.(services.Entitlements)

	var t models.ExcuseTemplate
	if err := h.db.First(&t, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "テンプレートが見つかりません"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "テンプレートの取得に失敗しました"})
		return
	}

	if t.IsPremium && !entitlements.CanUsePremiumTemplates {
		c.JSON(http.StatusForbidden, gin.H{"error": "プレミアムテンプレートを利用するにはプレミアムプランが必要です"})
		return
	}

	res := ExcuseTemplateResponse{
		ID:         t.ID,
		PackID:     t.PackID,
		ExcuseText: t.Text,
		Tags:       t.Tags,
		IsPremium:  t.IsPremium,
		CreatedAt:  t.CreatedAt,
	}
	c.JSON(http.StatusOK, res)
}
