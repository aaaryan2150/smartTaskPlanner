package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"smart-task-planner/config"           // ← Check this
	"smart-task-planner/internal/database" // ← Check this
	authRoutes "smart-task-planner/internal/modules/auth/routes"
	"smart-task-planner/internal/modules/auth/service"
	
)

func main() {
	config.Load()

	service.InitGoogleOAuth()


	if err := database.Connect(); err != nil {
		log.Fatal("❌ Database connection failed:", err)
	}
	defer database.Disconnect()

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":   "healthy",
			"message":  "Server is running",
			"database": "connected",
		})
	})

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Smart Task Planner API",
			"version": "1.0.0",
		})
	})

	authRoutes.RegisterAuthRoutes(router)
	authRoutes.RegisterOAuthRoutes(router)

	log.Println("🚀 Starting server...")
	log.Printf("🌐 Server running on http://localhost:%s", config.AppConfig.Port)
	log.Printf("📝 Health check: http://localhost:%s/health", config.AppConfig.Port)

	go func() {
		if err := router.Run(":" + config.AppConfig.Port); err != nil {
			log.Fatal("❌ Server failed:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down...")
	log.Println("✅ Server stopped")
}