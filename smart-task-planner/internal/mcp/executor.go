package mcp

import (
	"errors"
	"smart-task-planner/internal/modules/plan/models"
)

func RunTool(toolName string, params map[string]interface{}) (models.Plan, error) {
	switch toolName {
	case "create_task_plan":
		return CreateTaskPlan(params)
	default:
		return models.Plan{}, errors.New("unknown tool")
	}
}
