package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"smart-task-planner/internal/modules/plan/models"
	"smart-task-planner/internal/modules/plan/repository"
)

// Temp struct for parsing AI JSON
type AITask struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	DeadlineStr string `json:"deadline"`
}

type TaskPlan struct {
	Tasks []models.Task
}

func CreateTaskPlan(params map[string]interface{}, repo *repository.PlanRepository) (TaskPlan, error) {
	// 1️⃣ Validate input
	userID, ok := params["user_id"].(string)
	if !ok || userID == "" {
		return TaskPlan{}, fmt.Errorf("user_id required")
	}

	goal, ok := params["goal"].(string)
	if !ok || goal == "" {
		return TaskPlan{}, fmt.Errorf("goal required")
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return TaskPlan{}, fmt.Errorf("OPENAI_API_KEY not set")
	}

	now := time.Now()

	// 2️⃣ Fetch existing risky dates for the user
	riskData, _ := analyze_risks(map[string]interface{}{
		"user_id":        userID,
		"threshold_days": 3,
	}, repo)

	riskyDates := make(map[string]bool)
	if risks, ok := riskData["risks"].([]RiskTask); ok {
		for _, r := range risks {
			riskyDates[r.Deadline] = true
		}
	}

	// 3️⃣ Build AI prompt
	prompt := fmt.Sprintf(`You are an expert AI task planner.
Generate at least 10 actionable, detailed tasks for this goal:
"%s"
Each task must have:
- title
- description
- deadline (YYYY-MM-DD, evenly distributed across goal duration)
Avoid scheduling tasks on dates that are already risky for the user:
%v
Return ONLY valid JSON:
[
  {"title": "...", "description": "...", "deadline": "..."}
]`, goal, riskyDates)

	// 4️⃣ Call OpenAI
	aiResp, err := CallOpenAIAPI(prompt)
	if err != nil {
		return TaskPlan{}, err
	}

	// 5️⃣ Parse AI response
	raw := strings.TrimSpace(aiResp)
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimSuffix(raw, "```")
	raw = strings.TrimSpace(raw)

	var aiTasks []AITask
	if err := json.Unmarshal([]byte(raw), &aiTasks); err != nil {
		return TaskPlan{}, fmt.Errorf("failed to parse AI JSON: %v\nRaw: %s", err, raw)
	}

	// 6️⃣ Convert AI tasks → models.Task, adjust deadlines
	var tasks []models.Task
	for i, t := range aiTasks {
		deadline, err := time.Parse("2006-01-02", t.DeadlineStr)
		if err != nil {
			// Fallback if AI gave invalid date
			deadline = now.AddDate(0, 0, i+1)
		} else if deadline.Year() < now.Year() {
			// Fix year if AI gave old year
			deadline = time.Date(now.Year(), deadline.Month(), deadline.Day(), 0, 0, 0, 0, time.Local)
		}

		// Push task to next safe date if it conflicts with risky date
		for riskyDates[deadline.Format("2006-01-02")] {
			deadline = deadline.AddDate(0, 0, 1)
		}

		tasks = append(tasks, models.Task{
			Title:       t.Title,
			Description: t.Description,
			Status:      "Pending",
			Deadline:    deadline,
		})
	}

	// 7️⃣ Save plan to DB
	plan := &models.Plan{
		UserID:    userID,
		Goal:      goal,
		Tasks:     tasks,
		
	}

	if err := repo.Create(plan); err != nil {
		return TaskPlan{}, fmt.Errorf("failed to save plan: %v", err)
	}

	return TaskPlan{Tasks: tasks}, nil
}
