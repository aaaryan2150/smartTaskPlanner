package handlers

import (
	"net/http"
	"smart-task-planner/internal/modules/plan/dto"
	"smart-task-planner/internal/modules/plan/models"
	"smart-task-planner/internal/modules/plan/service"

	"github.com/gin-gonic/gin"
)

type PlanHandler struct {
	service *service.PlanService
}

func NewPlanHandler(svc *service.PlanService) *PlanHandler {
	return &PlanHandler{service: svc}
}

// GenerateDraftPlan generates an AI-based draft plan without saving to DB
func (h *PlanHandler) GenerateDraftPlan(c *gin.Context) {
	var req dto.CreatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tasks, err := h.service.GenerateDraftPlan(req.Goal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate draft plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"goal":  req.Goal,
		"tasks": tasks,
	})
}

// ConfirmPlanRequest is used when user confirms the draft plan
type ConfirmPlanRequest struct {
	Goal  string       `json:"goal" binding:"required"`
	Tasks []models.Task `json:"tasks" binding:"required"`
}

// ConfirmPlan saves the confirmed plan to the database
func (h *PlanHandler) ConfirmPlan(c *gin.Context) {
	userID := c.GetString("user_id") // from JWT middleware

	var req ConfirmPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plan, err := h.service.ConfirmPlan(userID, req.Goal, req.Tasks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm plan"})
		return
	}

	c.JSON(http.StatusCreated, plan)
}

// GetPlans fetches all plans for the current user
func (h *PlanHandler) GetPlans(c *gin.Context) {
	userID := c.GetString("user_id")

	plans, err := h.service.GetAllPlans(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch plans"})
		return
	}

	c.JSON(http.StatusOK, plans)
}

func (h *PlanHandler) RefineTask(c *gin.Context) {
	var req struct {
		TaskID string `json:"task_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.service.GetTaskDetails(req.TaskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	subtasks, err := h.service.RefineTask(*task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"subtasks": subtasks})
}

func (h *PlanHandler) UpdateTaskStatus(c *gin.Context) {
	var req struct {
		TaskID string `json:"task_id" binding:"required"`
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.service.UpdateTaskStatus(req.TaskID, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *PlanHandler) GetTaskDetails(c *gin.Context) {
	taskID := c.Query("task_id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task_id query param required"})
		return
	}

	task, err := h.service.GetTaskDetails(taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *PlanHandler) AddSubTasks(c *gin.Context) {
	var req struct {
		PlanID   string         `json:"plan_id" binding:"required"`
		TaskID   string         `json:"task_id" binding:"required"`
		SubTasks []models.Task  `json:"sub_tasks" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedPlan, err := h.service.AddSubTasks(req.PlanID, req.TaskID, req.SubTasks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedPlan)
}

