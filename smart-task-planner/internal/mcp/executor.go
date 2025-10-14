package mcp

import (
	"fmt"
	"smart-task-planner/internal/modules/plan/repository"
)

// RunTool executes an MCP tool by name
func RunTool(tool string, params map[string]interface{}, repo *repository.PlanRepository) (interface{}, error) {
	switch tool {
	case "create_task_plan":
		return CreateTaskPlan(params)
	case "get_goal_data":
		userID, ok := params["user_id"].(string)
		if !ok || userID == "" {
			return nil, fmt.Errorf("user_id parameter required")
		}
		return GetGoalData(userID, repo)
	default:
		return nil, fmt.Errorf("unknown MCP tool: %s", tool)
	}
}
