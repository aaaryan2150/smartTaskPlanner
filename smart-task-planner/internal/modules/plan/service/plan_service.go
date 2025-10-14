package service

import (
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

func (s *PlanService) GenerateTasksFromGoal(goal string) []models.Task {
	plan, err := mcp.RunTool("create_task_plan", map[string]interface{}{"goal": goal})
	if err != nil {
		// fallback mock tasks
		return []models.Task{
			{Title: "Research " + goal, Description: "Understand " + goal, Status: "Pending"},
			{Title: "Plan milestones", Description: "Break goal into achievable parts", Status: "Pending"},
			{Title: "Execute tasks", Description: "Start working on subtasks", Status: "Pending"},
		}
	}
	return plan.Tasks
}

// CreatePlan links plan to a user
func (s *PlanService) CreatePlan(userID string, goal string) (*models.Plan, error) {
	tasks := s.GenerateTasksFromGoal(goal)

	plan := &models.Plan{
		UserID: userID,
		Goal:   goal,
		Tasks:  tasks,
	}

	err := s.Repo.Create(plan)
	return plan, err
}

// GetAllPlans fetches plans for a user
func (s *PlanService) GetAllPlans(userID string) ([]models.Plan, error) {
	return s.Repo.GetAllByUser(userID)
}
