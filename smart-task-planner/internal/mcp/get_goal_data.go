package mcp

import (
	"fmt"
	"smart-task-planner/internal/modules/plan/models"
	"smart-task-planner/internal/modules/plan/repository"
)

// GetGoalData fetches all plans and tasks for a given user
func GetGoalData(userID string, repo *repository.PlanRepository) (UserGoals, error) {
	if userID == "" {
		return UserGoals{}, fmt.Errorf("user_id is required")
	}

	plans, err := repo.GetAllByUser(userID)
	if err != nil {
		return UserGoals{}, fmt.Errorf("failed to fetch plans: %w", err)
	}

	return UserGoals{
		UserID: userID,
		Plans:  plans,
	}, nil
}

// UserGoals represents the user + their plans/tasks
type UserGoals struct {
	UserID string        `json:"user_id"`
	Plans  []models.Plan `json:"plans"`
}
