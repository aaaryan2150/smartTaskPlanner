// mcp/refine_task.go
package mcp

import (
	"fmt"
	"smart-task-planner/internal/modules/plan/models"
	"time"
	"strings"
	"encoding/json"
)

// RefineTask generates subtasks for a given task
func RefineTask(task models.Task) ([]models.Task, error) {
	if task.Title == "" {
		return nil, fmt.Errorf("task title required")
	}

	prompt := fmt.Sprintf(`You are an AI task assistant.
Break the following task into detailed actionable subtasks:
Task: "%s"
Return ONLY valid JSON:
[
  {"title": "...", "description": "...", "deadline": "YYYY-MM-DD"}
]`, task.Title)

	resp, err := CallOpenAIAPI(prompt)
	if err != nil {
		return nil, err
	}

	// Clean markdown
	resp = TrimCodeBlock(resp)

	// Parse JSON
	var aiTasks []AITask
	if err := json.Unmarshal([]byte(resp), &aiTasks); err != nil {
		return nil, fmt.Errorf("failed to parse AI JSON: %v\nRaw: %s", err, resp)
	}

	var tasks []models.Task
	now := time.Now()
	for i, t := range aiTasks {
		deadline, _ := time.Parse("2006-01-02", t.DeadlineStr)
		if deadline.IsZero() {
			deadline = now.Add(time.Duration(i+1) * 24 * time.Hour)
		}
		tasks = append(tasks, models.Task{
			Title:       t.Title,
			Description: t.Description,
			Status:      "Pending",
			Deadline:    deadline,
		})
	}

	return tasks, nil
}

// Helper to trim ```json blocks
func TrimCodeBlock(s string) string {
	s = strings.TrimPrefix(s, "```json")
	s = strings.TrimSuffix(s, "```")
	return strings.TrimSpace(s)
}
