package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"smart-task-planner/config"
	"smart-task-planner/internal/database"

	authRoutes "smart-task-planner/internal/modules/auth/routes"
	authService "smart-task-planner/internal/modules/auth/service"

	planHandlers "smart-task-planner/internal/modules/plan/handlers"
	planRoutes "smart-task-planner/internal/modules/plan/routes"
	planRepository "smart-task-planner/internal/modules/plan/repository"
	planService "smart-task-planner/internal/modules/plan/service"

	commandHandlers "smart-task-planner/internal/modules/command/handlers"
	commandRoutes "smart-task-planner/internal/modules/command/routes"
	commandService "smart-task-planner/internal/modules/command/service"
)

func main() {
	// -------------------------------
	// Load configuration
	// -------------------------------
	config.Load()

	// -------------------------------
	// Initialize Google OAuth
	// -------------------------------
	authService.InitGoogleOAuth()

	// -------------------------------
	// Connect to MongoDB
	// -------------------------------
	if err := database.Connect(); err != nil {
		log.Fatal("❌ Database connection failed:", err)
	}
	defer database.Disconnect()

	// -------------------------------
	// Create Gin router
	// -------------------------------
	router := gin.Default()

	// -------------------------------
	// Health & root endpoints
	// -------------------------------
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

	// -------------------------------
	// Auth routes
	// -------------------------------
	authRoutes.RegisterAuthRoutes(router)
	authRoutes.RegisterOAuthRoutes(router)

	// -------------------------------
	// Plan module (Phase 1 & 2)
	// -------------------------------
	db := database.DB                                  // *mongo.Database
	planRepo := planRepository.NewPlanRepository(db)  // repository
	planSvc := planService.NewPlanService(planRepo)   // service
	planHandler := planHandlers.NewPlanHandler(planSvc) // handler
	planRoutes.RegisterPlanRoutes(router, planHandler) // plan routes

	// -------------------------------
	// Command module (Phase 3 & 4)
	// -------------------------------
	cmdSvc := commandService.NewCommandService(planRepo) // pass PlanRepository
	cmdHandler := commandHandlers.NewCommandHandler(cmdSvc)   // handler
	commandRoutes.RegisterCommandRoutes(router, cmdHandler)  // register /api/command

	// -------------------------------
	// Start server
	// -------------------------------
	log.Println("🚀 Starting server...")
	log.Printf("🌐 Server running on http://localhost:%s", config.AppConfig.Port)
	log.Printf("📝 Health check: http://localhost:%s/health", config.AppConfig.Port)

	go func() {
		if err := router.Run(":" + config.AppConfig.Port); err != nil {
			log.Fatal("❌ Server failed:", err)
		}
	}()

	// -------------------------------
	// Graceful shutdown
	// -------------------------------
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down...")
	log.Println("✅ Server stopped")
}
