package routes

import (
	"github.com/gin-gonic/gin"
	"smart-task-planner/internal/modules/command/handlers"
	"smart-task-planner/internal/middleware"
)


func RegisterCommandRoutes(router *gin.Engine, handler *handlers.CommandHandler) {
	api := router.Group("/api/command")
	api.Use(middleware.JWTAuth())
	{
		api.POST("/", handler.ProcessCommand)
	}
}
