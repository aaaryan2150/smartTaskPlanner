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
		// Generate AI draft plan (not saved to DB)
		api.POST("/draft", handler.GenerateDraftPlan)

		// Confirm & save a plan to DB
		api.POST("/confirm", handler.ConfirmPlan)

		// Get all plans for current user
		api.GET("/", handler.GetPlans)

		api.POST("/refine-task", handler.RefineTask)
		api.POST("/update-task-status", handler.UpdateTaskStatus)
		api.GET("/task-details", handler.GetTaskDetails)

		api.POST("/add-subtasks", handler.AddSubTasks)

	}
}
