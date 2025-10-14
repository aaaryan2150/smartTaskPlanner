package mcp

import (
	"encoding/json"
	"fmt"
	"time"
	"smart-task-planner/internal/modules/plan/models"
)

// CreateTaskPlan calls the MCP/Gemini API with a structured prompt and returns the Plan
func CreateTaskPlan(params map[string]interface{}) (models.Plan, error) {
	goal := params["goal"].(string)

	// Optional: user can pass a deadline, else default to 2 weeks from today
	var deadline string
	if d, ok := params["deadline"].(string); ok && d != "" {
		deadline = d
	} else {
		deadline = time.Now().AddDate(0, 0, 14).Format("2006-01-02") // 2 weeks from today
	}

	today := time.Now().Format("2006-01-02")

	// Construct the detailed task planning prompt
	prompt := fmt.Sprintf(`You are an expert task planner. Today's date is %s.

Break down the following goal into a detailed, actionable task plan:

GOAL: "%s"
FINAL DEADLINE: %s

REQUIREMENTS:
1. Create SPECIFIC, ACTIONABLE tasks (not vague descriptions)
2. Number of tasks should be proportional to the time available until the deadline
3. Each task must be completable within 1-3 days
4. Include concrete deliverables for each task
5. Assign realistic deadlines to each task, distributed between today and the final deadline
6. Tasks should follow a logical sequence
7. Break down complex activities into smaller, measurable steps

OUTPUT FORMAT (respond with ONLY valid JSON, no markdown or additional text):
[
  {
    "id": "task-1",
    "name": "Specific task description with clear deliverable",
    "deadline": "YYYY-MM-DD"
  }
]

EXAMPLE:
If goal is "Launch a blog" with 10 days deadline:
- Don't say: "Set up website" (too vague)
- Do say: "Register domain name and set up hosting on Vercel/Netlify" (specific)
- Don't say: "Create content" (too vague)
- Do say: "Write and edit 5 blog posts (500-800 words each) on chosen topics" (specific)

Generate the task breakdown now:`, today, goal, deadline)

	// Call the Gemini API
	aiResponse, err := CallGeminiAPI(prompt)
	if err != nil {
		return models.Plan{}, err
	}

	// Parse the JSON response into tasks
	var tasks []models.Task
	if err := json.Unmarshal([]byte(aiResponse), &tasks); err != nil {
		// Try to extract JSON if extra text is present
		start, end := 0, len(aiResponse)
		for i, c := range aiResponse {
			if c == '[' {
				start = i
				break
			}
		}
		for i := len(aiResponse) - 1; i >= 0; i-- {
			if aiResponse[i] == ']' {
				end = i + 1
				break
			}
		}
		cleanJSON := aiResponse[start:end]
		_ = json.Unmarshal([]byte(cleanJSON), &tasks)
	}

	return models.Plan{
		Goal:  goal,
		Tasks: tasks,
	}, nil
}
