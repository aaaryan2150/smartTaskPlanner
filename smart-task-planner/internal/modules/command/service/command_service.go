package services

import (
	"fmt"
	"smart-task-planner/internal/mcp"
	"smart-task-planner/internal/modules/plan/repository"
)

type CommandService struct {
	Repo *repository.PlanRepository
}

func NewCommandService(repo *repository.PlanRepository) *CommandService {
	return &CommandService{Repo: repo}
}

// HandleCommand interprets natural language and triggers MCP tools with smart chaining
func (s *CommandService) HandleCommand(userID, message string) (map[string]interface{}, error) {
	// Step 1: Interpret the user message
	intent, err := mcp.RunTool("interpret_user_message", map[string]interface{}{
		"user_id": userID,
		"message": message,
	}, s.Repo)
	if err != nil {
		return nil, fmt.Errorf("interpretation failed: %v", err)
	}

	// Step 1b: Type assert
	intentMap, ok := intent.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid intent format")
	}

	toolName, ok := intentMap["tool"].(string)
	if !ok {
		return nil, fmt.Errorf("tool name not found")
	}

	params, ok := intentMap["params"].(map[string]interface{})
	if !ok {
		params = make(map[string]interface{})
	}

	// Check if chaining is needed
	needsChaining, _ := intentMap["needs_chaining"].(bool)

	// Step 2: Execute the primary tool
	result, err := mcp.RunTool(toolName, params, s.Repo)
	if err != nil {
		return nil, fmt.Errorf("tool execution failed: %v", err)
	}

	// Step 3: Handle smart chaining for progress → feedback
	if needsChaining && toolName == "get_user_progress" {
		progressData, ok := result.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid progress data format")
		}

		// Chain with provide_feedback
		feedback, err := mcp.RunTool("provide_feedback", map[string]interface{}{
			"progress_data": progressData,
		}, s.Repo)

		if err != nil {
			return nil, fmt.Errorf("feedback generation failed: %v", err)
		}

		feedbackData, ok := feedback.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid feedback format")
		}

		// Return clean feedback response
		return feedbackData, nil
	}

	// Step 4: Return result for all other tools (including handle_general_query)
	return map[string]interface{}{
		"interpreted_action": toolName,
		"result":             result,
	}, nil
}