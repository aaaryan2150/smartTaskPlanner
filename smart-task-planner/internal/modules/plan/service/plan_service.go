package service

import (
	"fmt"
	"smart-task-planner/internal/mcp"
	"smart-task-planner/internal/modules/plan/models"
	"smart-task-planner/internal/modules/plan/repository"
)

type PlanService struct {
	Repo *repository.PlanRepository
}

func NewPlanService(repo *repository.PlanRepository) *PlanService {
	return &PlanService{Repo: repo}
}

// GenerateDraftPlan calls MCP to create an AI-generated draft plan without saving to DB
func (s *PlanService) GenerateDraftPlan(goal string) ([]models.Task, error) {
	// Pass the repository if your MCP tool needs it (optional for create_task_plan)
	result, err := mcp.RunTool("create_task_plan", map[string]interface{}{"goal": goal}, s.Repo)
	if err != nil {
		// fallback tasks
		return []models.Task{
			{Title: "Research " + goal, Description: "Understand " + goal, Status: "Pending"},
			{Title: "Plan milestones", Description: "Break goal into achievable parts", Status: "Pending"},
			{Title: "Execute tasks", Description: "Start working on subtasks", Status: "Pending"},
		}, nil
	}

	plan, ok := result.(mcp.TaskPlan)
	if !ok {
		return nil, fmt.Errorf("invalid data returned from MCP")
	}

	return plan.Tasks, nil
}

// ConfirmPlan saves the confirmed plan to the DB
func (s *PlanService) ConfirmPlan(userID string, goal string, tasks []models.Task) (*models.Plan, error) {
	plan := &models.Plan{
		UserID: userID,
		Goal:   goal,
		Tasks:  tasks,
	}

	if err := s.Repo.Create(plan); err != nil {
		return nil, err
	}
	return plan, nil
}

// GetAllPlans fetches all plans for a user
func (s *PlanService) GetAllPlans(userID string) ([]models.Plan, error) {
	return s.Repo.GetAllByUser(userID)
}

// GetUserGoals fetches all plans using the new MCP tool get_goal_data
func (s *PlanService) GetUserGoals(userID string) (mcp.UserGoals, error) {
	result, err := mcp.RunTool("get_goal_data", map[string]interface{}{"user_id": userID}, s.Repo)
	if err != nil {
		return mcp.UserGoals{}, err
	}

	goals, ok := result.(mcp.UserGoals)
	if !ok {
		return mcp.UserGoals{}, fmt.Errorf("invalid data returned from MCP")
	}

	return goals, nil
}
