package handlers

import (
	"errors"
	"net/http"
	"time"
	"what-went-wrong-api/internal/models"
	"what-went-wrong-api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExcuseHandler struct {
	db *gorm.DB
}

func NewExcuseHandler(db *gorm.DB) *ExcuseHandler {
	return &ExcuseHandler{db: db}
}

// GetExcuses godoc
// @Summary List excuses for a goal
// @Description List excuses, filters by retention days if strictly limited by entitlement.
// @Tags excuses
// @Accept json
// @Produce json
// @Param goal_id path string true "Goal ID"
// @Param from query string false "From Date (YYYY-MM-DD)"
// @Param to query string false "To Date (YYYY-MM-DD)"
// @Success 200 {object} GetExcusesResponse
// @Failure 400 {object} ExcuseValidationErrorResponse "Invalid Goal ID"
// @Failure 401 {object} ExcuseUnauthorizedResponse
// @Failure 500 {object} ExcuseFetchErrorResponse
// @Security BearerAuth
// @Router /goals/{goal_id}/excuses [get]
func (h *ExcuseHandler) GetExcuses(c *gin.Context) {
	userIdStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証されていません"})
		return
	}
	userID := userIdStr.(string)

	goalIDStr := c.Param("id")
	goalID, err := uuid.Parse(goalIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "入力内容が正しくありません"})
		return
	}

	entitlementsInterface, exists := c.Get("entitlements")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "プラン情報の取得に失敗しました"})
		return
	}
	entitlements := entitlementsInterface.(services.Entitlements)

	query := h.db.Model(&models.ExcuseEntry{}).Where("user_id = ? AND goal_id = ?", userID, goalID)

	// Entitlement: logRetentionDays
	if entitlements.LogRetentionDays != nil {
		retentionDate := time.Now().AddDate(0, 0, -*entitlements.LogRetentionDays).Format("2006-01-02")
		// Force filter: date >= retentionDate
		query = query.Where("date >= ?", retentionDate)
	}

	// Manual Filters (if they don't violate retention)
	from := c.Query("from")
	if from != "" {
		query = query.Where("date >= ?", from)
	}
	to := c.Query("to")
	if to != "" {
		query = query.Where("date <= ?", to)
	}

	var excuses []models.ExcuseEntry
	if err := query.Order("date desc").Find(&excuses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "言い訳の取得に失敗しました"})
		return
	}

	res := GetExcusesResponse{Excuses: make([]ExcuseResponse, len(excuses))}
	for i, e := range excuses {
		res.Excuses[i] = ExcuseResponse{
			ID:         e.ID,
			GoalID:     e.GoalID,
			Date:       e.Date,
			ExcuseText: e.ExcuseText,
			TemplateID: e.TemplateID,
			CreatedAt:  e.CreatedAt,
			UpdatedAt:  e.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, res)
}

// GetExcuseToday godoc
// @Summary Get today's excuse for a goal
// @Tags excuses
// @Accept json
// @Produce json
// @Param goal_id path string true "Goal ID"
// @Success 200 {object} ExcuseResponse
// @Failure 400 {object} ExcuseValidationErrorResponse "Invalid Goal ID"
// @Failure 401 {object} ExcuseUnauthorizedResponse
// @Failure 404 {object} ExcuseNotFoundResponse
// @Failure 500 {object} ExcuseFetchErrorResponse
// @Security BearerAuth
// @Router /goals/{goal_id}/excuses/today [get]
func (h *ExcuseHandler) GetExcuseToday(c *gin.Context) {
	userIdStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証されていません"})
		return
	}
	userID := userIdStr.(string)

	goalIDStr := c.Param("id")
	goalID, err := uuid.Parse(goalIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "入力内容が正しくありません"})
		return
	}

	today := time.Now().Format("2006-01-02")
	var excuse models.ExcuseEntry
	if err := h.db.Where("user_id = ? AND goal_id = ? AND date = ?", userID, goalID, today).First(&excuse).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "本日の言い訳が見つかりません"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "言い訳の取得に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, ExcuseResponse{
		ID:         excuse.ID,
		GoalID:     excuse.GoalID,
		Date:       excuse.Date,
		ExcuseText: excuse.ExcuseText,
		TemplateID: excuse.TemplateID,
		CreatedAt:  excuse.CreatedAt,
		UpdatedAt:  excuse.UpdatedAt,
	})
}

// PostExcuse godoc
// @Summary Create or update an excuse
// @Description Upsert excuse for a date. Checks entitlement if using premium template.
// @Tags excuses
// @Accept json
// @Produce json
// @Param goal_id path string true "Goal ID"
// @Param request body CreateExcuseRequest true "Excuse Data"
// @Success 201 {object} ExcuseResponse
// @Failure 400 {object} ExcuseValidationErrorResponse
// @Failure 401 {object} ExcuseUnauthorizedResponse
// @Failure 403 {object} ExcuseForbiddenResponse
// @Failure 500 {object} ExcuseCreateErrorResponse
// @Security BearerAuth
// @Router /goals/{goal_id}/excuses [post]
func (h *ExcuseHandler) PostExcuse(c *gin.Context) {
	userIdStr, _ := c.Get("userID")
	userID := userIdStr.(string)
	goalIDStr := c.Param("id")
	goalID, err := uuid.Parse(goalIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "入力内容が正しくありません"})
		return
	}

	var req CreateExcuseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "入力内容が正しくありません"})
		return
	}

	entitlementsInterface, _ := c.Get("entitlements")
	entitlements := entitlementsInterface.(services.Entitlements)

	// Verify template if provided
	if req.TemplateID != "" {
		var tmpl models.ExcuseTemplate
		if err := h.db.First(&tmpl, "id = ?", req.TemplateID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "入力内容が正しくありません"})
			return
		}
		if tmpl.IsPremium && !entitlements.CanUsePremiumTemplates {
			c.JSON(http.StatusForbidden, gin.H{"error": "プレミアムテンプレートを利用するにはプレミアムプランが必要です"})
			return
		}
	}

	// Upsert Logic
	var excuse models.ExcuseEntry
	err = h.db.Where("user_id = ? AND goal_id = ? AND date = ?", userID, goalID, req.Date).First(&excuse).Error
	if err == nil {
		// Update
		excuse.ExcuseText = req.ExcuseText
		if req.TemplateID != "" {
			excuse.TemplateID = &req.TemplateID
		} else {
			excuse.TemplateID = nil
		}
		h.db.Save(&excuse)
		c.JSON(http.StatusOK, mapToResponse(excuse))
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create
		excuse = models.ExcuseEntry{
			UserID:     userID,
			GoalID:     goalID,
			Date:       req.Date,
			ExcuseText: req.ExcuseText,
		}
		if req.TemplateID != "" {
			excuse.TemplateID = &req.TemplateID
		}
		if err := h.db.Create(&excuse).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "言い訳の作成に失敗しました"})
			return
		}
		c.JSON(http.StatusCreated, mapToResponse(excuse))
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "言い訳の作成に失敗しました"})
	}
}

// PatchExcuse godoc
// @Summary Update an excuse
// @Tags excuses
// @Accept json
// @Produce json
// @Param id path string true "Excuse ID"
// @Param request body UpdateExcuseRequest true "Update Data"
// @Success 200 {object} ExcuseResponse
// @Failure 400 {object} ExcuseValidationErrorResponse
// @Failure 401 {object} ExcuseUnauthorizedResponse
// @Failure 403 {object} ExcuseForbiddenResponse
// @Failure 404 {object} ExcuseNotFoundResponse
// @Failure 500 {object} ExcuseUpdateErrorResponse
// @Security BearerAuth
// @Router /excuses/{id} [patch]
func (h *ExcuseHandler) PatchExcuse(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "入力内容が正しくありません"})
		return
	}
	userIdStr, _ := c.Get("userID")
	userID := userIdStr.(string)

	var req UpdateExcuseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "入力内容が正しくありません"})
		return
	}

	var excuse models.ExcuseEntry
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&excuse).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "言い訳が見つかりません"})
		return
	}

	entitlementsInterface, _ := c.Get("entitlements")
	entitlements := entitlementsInterface.(services.Entitlements)

	if req.ExcuseText != "" {
		excuse.ExcuseText = req.ExcuseText
	}

	// If TemplateID is updated (checked if present in request via pointer usually, but here string empty assumes no change or unset?
	// Spec says "Update specific excuse... check template entitlement".
	// For simplicity, if ID is provided, check it.
	if req.TemplateID != "" {
		var tmpl models.ExcuseTemplate
		if err := h.db.First(&tmpl, "id = ?", req.TemplateID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "入力内容が正しくありません"})
			return
		}
		if tmpl.IsPremium && !entitlements.CanUsePremiumTemplates {
			c.JSON(http.StatusForbidden, gin.H{"error": "プレミアムテンプレートを利用するにはプレミアムプランが必要です"})
			return
		}
		excuse.TemplateID = &req.TemplateID
	}

	if err := h.db.Save(&excuse).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "言い訳の更新に失敗しました"})
		return
	}
	c.JSON(http.StatusOK, mapToResponse(excuse))
}

// DeleteExcuse godoc
// @Summary Delete an excuse
// @Tags excuses
// @Param id path string true "Excuse ID"
// @Success 204 "No Content"
// @Failure 400 {object} ExcuseValidationErrorResponse "Invalid Excuse ID"
// @Failure 401 {object} ExcuseUnauthorizedResponse
// @Failure 404 {object} ExcuseNotFoundResponse
// @Failure 500 {object} ExcuseDeleteErrorResponse
// @Security BearerAuth
// @Router /excuses/{id} [delete]
func (h *ExcuseHandler) DeleteExcuse(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "入力内容が正しくありません"})
		return
	}
	userIdStr, _ := c.Get("userID")
	userID := userIdStr.(string)

	result := h.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.ExcuseEntry{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "言い訳の削除に失敗しました"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "言い訳が見つかりません"})
		return
	}

	c.Status(http.StatusNoContent)
}

func mapToResponse(e models.ExcuseEntry) ExcuseResponse {
	return ExcuseResponse{
		ID:         e.ID,
		GoalID:     e.GoalID,
		Date:       e.Date,
		ExcuseText: e.ExcuseText,
		TemplateID: e.TemplateID,
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
	}
}
