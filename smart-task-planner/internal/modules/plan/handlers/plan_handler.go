package handlers

import (
	"net/http"
	"smart-task-planner/internal/modules/plan/dto"
	"smart-task-planner/internal/modules/plan/service"

	"github.com/gin-gonic/gin"
)

type PlanHandler struct {
	service *service.PlanService
}

func NewPlanHandler(svc *service.PlanService) *PlanHandler {
	return &PlanHandler{service: svc}
}

func (h *PlanHandler) CreatePlan(c *gin.Context) {
	userID := c.GetString("user_id") // from JWT middleware

	var req dto.CreatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plan, err := h.service.CreatePlan(userID, req.Goal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create plan"})
		return
	}

	c.JSON(http.StatusCreated, plan)
}

func (h *PlanHandler) GetPlans(c *gin.Context) {
	userID := c.GetString("user_id")

	plans, err := h.service.GetAllPlans(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch plans"})
		return
	}

	c.JSON(http.StatusOK, plans)
}
