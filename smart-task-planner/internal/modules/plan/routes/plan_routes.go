package routes

import (
	"smart-task-planner/internal/modules/plan/handlers"
	"smart-task-planner/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterPlanRoutes(router *gin.Engine, handler *handlers.PlanHandler) {
	api := router.Group("/api/plan")
	api.Use(middleware.JWTAuth())
	{
		api.POST("/generate", handler.CreatePlan)
		api.GET("/", handler.GetPlans)
	}
}

