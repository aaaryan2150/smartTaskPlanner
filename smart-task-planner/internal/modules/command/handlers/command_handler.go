package handlers

import (
	"net/http"
	"smart-task-planner/internal/modules/command/service"

	"github.com/gin-gonic/gin"
)

type CommandHandler struct {
	Service *services.CommandService
}

func NewCommandHandler(service *services.CommandService) *CommandHandler {
	return &CommandHandler{Service: service}
}

// POST /api/command
func (h *CommandHandler) ProcessCommand(c *gin.Context) {
	// Extract userID from JWT (set by middleware)
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// Bind message from request body
	var req struct {
		Message string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the command service (now with automatic chaining support)
	response, err := h.Service.HandleCommand(userID, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}