package mcp

import (
	"fmt"
	"smart-task-planner/internal/modules/plan/repository"
	"smart-task-planner/internal/modules/plan/models"
)

func RunTool(tool string, params map[string]interface{}, repo *repository.PlanRepository) (interface{}, error) {
	switch tool {
	case "create_task_plan":
		return CreateTaskPlan(params, repo)
	case "get_goal_data":
		userID, ok := params["user_id"].(string)
		if !ok || userID == "" {
			return nil, fmt.Errorf("user_id parameter required")
		}
		return GetGoalData(userID, repo)
	case "refine_task":
		task, ok := params["task"].(models.Task)
		if !ok {
			return nil, fmt.Errorf("task parameter required")
		}
		return RefineTask(task)
	case "update_task_status":
		taskID, ok1 := params["task_id"].(string)
		status, ok2 := params["status"].(string)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("task_id and status required")
		}
		return UpdateTaskStatus(taskID, status, repo)
	case "get_task_details":
		taskID, ok := params["task_id"].(string)
		if !ok {
			return nil, fmt.Errorf("task_id required")
		}
		return GetTaskDetails(taskID, repo)

	case "interpret_user_message":
		return interpret_user_message(params)
	case "reschedule_plan":
		return reschedule_plan(params, repo)
	case "analyze_risks":
		return analyze_risks(params, repo)
	case "get_user_progress":
		return get_user_progress(params, repo)
	case "provide_feedback":
		return provide_feedback(params)
	case "generate_alternative_plans":
		return generate_alternative_plans(params)
	case "handle_general_query":
    	return handle_general_query(params, repo)

	default:
		return nil, fmt.Errorf("unknown MCP tool: %s", tool)
	}
}


