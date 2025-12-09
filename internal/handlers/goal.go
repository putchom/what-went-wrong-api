package handlers

import (
	"errors"
	"net/http"
	"time"
	"what-went-wrong-api/internal/models"
	"what-went-wrong-api/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GoalHandler struct {
	db *gorm.DB
}

func NewGoalHandler(db *gorm.DB) *GoalHandler {
	return &GoalHandler{db: db}
}

// GetGoals godoc
// @Summary List goals
// @Description Get all goals for the current user
// @Tags goals
// @Accept json
// @Produce json
// @Success 200 {object} GetGoalsResponse
// @Failure 401 {object} GoalUnauthorizedResponse
// @Failure 500 {object} GoalFetchErrorResponse
// @Security BearerAuth
// @Router /goals [get]
func (h *GoalHandler) GetGoals(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証されていません"})
		return
	}
	userID := userIDStr.(string)

	var goals []models.Goal
	if err := h.db.Where("user_id = ?", userID).Order("\"order\" asc, created_at desc").Find(&goals).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "目標の取得に失敗しました"})
		return
	}

	res := GetGoalsResponse{Goals: make([]GoalResponse, len(goals))}
	for i, g := range goals {
		res.Goals[i] = GoalResponse{
			ID:                  g.ID.String(),
			Title:               g.Title,
			NotificationTime:    g.NotificationTime,
			NotificationEnabled: g.NotificationEnabled,
			Order:               g.Order,
			CreatedAt:           g.CreatedAt,
			UpdatedAt:           g.UpdatedAt,
		}
	}
	c.JSON(http.StatusOK, res)
}

// PostGoals godoc
// @Summary Create goal
// @Description Create a new goal. Checks for plan limits.
// @Tags goals
// @Accept json
// @Produce json
// @Param request body CreateGoalRequest true "Request body"
// @Success 201 {object} CreateGoalResponse
// @Failure 400 {object} GoalValidationErrorResponse
// @Failure 401 {object} GoalUnauthorizedResponse
// @Failure 403 {object} GoalLimitReachedResponse "Forbidden if max goals reached"
// @Failure 500 {object} GoalCreateErrorResponse
// @Security BearerAuth
// @Router /goals [post]
func (h *GoalHandler) PostGoals(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証されていません"})
		return
	}
	userID := userIDStr.(string)

	// Entitlement check
	entitlementsInterface, exists := c.Get("entitlements")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "プラン情報の取得に失敗しました"})
		return
	}
	entitlements := entitlementsInterface.(services.Entitlements)

	var currentCount int64
	if err := h.db.Model(&models.Goal{}).Where("user_id = ?", userID).Count(&currentCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "目標数の取得に失敗しました"})
		return
	}

	if int(currentCount) >= entitlements.MaxGoals {
		c.JSON(http.StatusForbidden, gin.H{"error": "プランの目標作成数上限に達しました"})
		return
	}

	var req CreateGoalRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "入力内容が正しくありません"})
		return
	}

	// Calculate order (simple: put at end)
	// Or just default to 0

	newGoal := models.Goal{
		UserID:              userID,
		Title:               req.Title,
		NotificationTime:    req.NotificationTime,
		NotificationEnabled: req.NotificationEnabled,
		Order:               int(currentCount) + 1, // Simple default order
	}

	if err := h.db.Create(&newGoal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "目標の作成に失敗しました"})
		return
	}

	c.JSON(http.StatusCreated, CreateGoalResponse{
		Goal: GoalResponse{
			ID:                  newGoal.ID.String(),
			Title:               newGoal.Title,
			NotificationTime:    newGoal.NotificationTime,
			NotificationEnabled: newGoal.NotificationEnabled,
			Order:               newGoal.Order,
			CreatedAt:           newGoal.CreatedAt,
			UpdatedAt:           newGoal.UpdatedAt,
		},
	})
}

// PatchGoal godoc
// @Summary Update goal
// @Tags goals
// @Accept json
// @Produce json
// @Param id path string true "Goal ID"
// @Param request body UpdateGoalRequest true "Request body"
// @Success 200 {object} CreateGoalResponse
// @Failure 400 {object} GoalValidationErrorResponse
// @Failure 401 {object} GoalUnauthorizedResponse
// @Failure 404 {object} GoalNotFoundErrorResponse
// @Failure 500 {object} GoalUpdateErrorResponse
// @Security BearerAuth
// @Router /goals/{id} [patch]
func (h *GoalHandler) PatchGoal(c *gin.Context) {
	goalID := c.Param("id")
	userIDStr, _ := c.Get("userID")
	userID := userIDStr.(string)

	var goal models.Goal
	if err := h.db.First(&goal, "id = ? AND user_id = ?", goalID, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "目標が見つかりません"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "目標の取得に失敗しました"})
		return
	}

	var req UpdateGoalRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "入力内容が正しくありません"})
		return
	}

	if req.Title != nil {
		goal.Title = *req.Title
	}
	if req.NotificationTime != nil {
		goal.NotificationTime = req.NotificationTime
	}
	if req.NotificationEnabled != nil {
		goal.NotificationEnabled = *req.NotificationEnabled
	}
	goal.UpdatedAt = time.Now()

	if err := h.db.Save(&goal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "目標の更新に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, CreateGoalResponse{
		Goal: GoalResponse{
			ID:                  goal.ID.String(),
			Title:               goal.Title,
			NotificationTime:    goal.NotificationTime,
			NotificationEnabled: goal.NotificationEnabled,
			Order:               goal.Order,
			CreatedAt:           goal.CreatedAt,
			UpdatedAt:           goal.UpdatedAt,
		},
	})
}

// DeleteGoal godoc
// @Summary Delete goal
// @Tags goals
// @Accept json
// @Produce json
// @Param id path string true "Goal ID"
// @Success 204 "No Content"
// @Failure 401 {object} GoalUnauthorizedResponse
// @Failure 404 {object} GoalNotFoundErrorResponse
// @Failure 500 {object} GoalDeleteErrorResponse
// @Security BearerAuth
// @Router /goals/{id} [delete]
func (h *GoalHandler) DeleteGoal(c *gin.Context) {
	goalID := c.Param("id")
	userIDStr, _ := c.Get("userID")
	userID := userIDStr.(string)

	// Transaction to delete associated excuses if necessary
	// Assuming cascade delete is configured in DB or we handle it here
	// Spec says: Goal削除＋紐づくExcuseEntry削除

	err := h.db.Transaction(func(tx *gorm.DB) error {
		// Verify ownership
		var goal models.Goal
		if err := tx.First(&goal, "id = ? AND user_id = ?", goalID, userID).Error; err != nil {
			return err
		}

		// Delete excuses
		if err := tx.Where("goal_id = ?", goalID).Delete(&models.ExcuseEntry{}).Error; err != nil {
			return err
		}

		// Delete goal
		return tx.Delete(&goal).Error
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "目標が見つかりません"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "目標の削除に失敗しました"})
		return
	}

	c.Status(http.StatusNoContent)
}
