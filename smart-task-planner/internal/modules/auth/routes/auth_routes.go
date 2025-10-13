package routes

import (
	"github.com/gin-gonic/gin"
	"smart-task-planner/internal/modules/auth/handlers"
)

// RegisterAuthRoutes registers auth endpoints
func RegisterAuthRoutes(router *gin.Engine) {
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}
}
