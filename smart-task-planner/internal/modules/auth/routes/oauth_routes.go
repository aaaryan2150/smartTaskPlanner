package routes

import (
	"github.com/gin-gonic/gin"
	"smart-task-planner/internal/modules/auth/handlers"
)

// RegisterOAuthRoutes registers Google OAuth routes
func RegisterOAuthRoutes(router *gin.Engine) {
	auth := router.Group("/api/auth")
	{
		auth.GET("/google/login", handlers.GoogleLogin)
		auth.GET("/google/callback", handlers.GoogleCallback)
	}
}
