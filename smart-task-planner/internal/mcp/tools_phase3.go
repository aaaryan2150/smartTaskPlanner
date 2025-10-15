package mcp

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sort"
	"time"

	"smart-task-planner/internal/modules/plan/repository"
	"smart-task-planner/internal/modules/plan/models"
)

type RiskTask struct {
	Goal     string `json:"goal"`
	TaskName string `json:"task_name"`
	Deadline string `json:"deadline"`
	DaysLeft int    `json:"days_left"`
}

func interpret_user_message(params map[string]interface{}) (map[string]interface{}, error) {
	message := params["message"].(string)
	userID := params["user_id"].(string)

	switch {
	case contains(message, "behind"):
		return map[string]interface{}{
			"tool": "reschedule_plan",
			"params": map[string]interface{}{
				"user_id": userID,
				"message": message,
			},
		}, nil
	case contains(message, "risk"):
		return map[string]interface{}{
			"tool": "analyze_risks",
			"params": map[string]interface{}{
				"user_id": userID,
			},
		}, nil
	case contains(message, "faster"):
		return map[string]interface{}{
			"tool": "generate_alternative_plans",
			"params": map[string]interface{}{
				"user_id": userID,
			},
		}, nil
	case contains(message, "progress") || contains(message, "feedback"):
		// ✅ Fixed: Set needs_chaining flag for progress/feedback queries
		return map[string]interface{}{
			"tool": "get_user_progress",
			"params": map[string]interface{}{
				"user_id": userID,
				"message": message,
			},
			"needs_chaining": true, // Signal that feedback should follow
		}, nil
	default:
		// ✅ NEW: Fallback to AI-powered general query handler
		return map[string]interface{}{
			"tool": "handle_general_query",
			"params": map[string]interface{}{
				"user_id": userID,
				"message": message,
			},
		}, nil
	}
}

func reschedule_plan(params map[string]interface{}, repo *repository.PlanRepository) (map[string]interface{}, error) {
	userID, ok := params["user_id"].(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("user_id required")
	}

	message, ok := params["message"].(string)
	if !ok || message == "" {
		return nil, fmt.Errorf("message required")
	}

	delay := extractDelayDays(message)
	if delay == 0 {
		return nil, fmt.Errorf("no delay days found in message")
	}

	plan, err := repo.FindGoalByAI(userID, message)
	if err != nil {
		return nil, fmt.Errorf("goal not found: %v", err)
	}

	for i, task := range plan.Tasks {
		if !task.Deadline.IsZero() {
			plan.Tasks[i].Deadline = task.Deadline.AddDate(0, 0, delay)
		}
	}

	updatedPlan, err := repo.UpdatePlan(plan)
	if err != nil {
		return nil, fmt.Errorf("failed to update plan: %v", err)
	}

	return map[string]interface{}{
		"message": fmt.Sprintf("All tasks for goal '%s' shifted by %d days", updatedPlan.Goal, delay),
		"goal_id": updatedPlan.ID.Hex(),
		"tasks":   updatedPlan.Tasks,
	}, nil
}

func analyze_risks(params map[string]interface{}, repo *repository.PlanRepository) (map[string]interface{}, error) {
	userID, ok := params["user_id"].(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("user_id required")
	}

	threshold := 3
	if t, exists := params["threshold_days"]; exists {
		switch v := t.(type) {
		case int:
			threshold = v
		case float64:
			threshold = int(v)
		case string:
			if parsed, err := strconv.Atoi(v); err == nil {
				threshold = parsed
			}
		}
	}

	plans, err := repo.GetAllByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user plans: %v", err)
	}

	var risks []RiskTask
	now := time.Now()

	for _, plan := range plans {
		for _, task := range plan.Tasks {
			if !task.Deadline.IsZero() {
				daysLeft := int(task.Deadline.Sub(now).Hours() / 24)
				if daysLeft <= threshold {
					risks = append(risks, RiskTask{
						Goal:     plan.Goal,
						TaskName: task.Title,
						Deadline: task.Deadline.Format("2006-01-02"),
						DaysLeft: daysLeft,
					})
				}
			}
			checkSubTasks(plan.Goal, task.SubTasks, &risks, now, threshold)
		}
	}

	sort.SliceStable(risks, func(i, j int) bool {
		return risks[i].DaysLeft < risks[j].DaysLeft
	})

	return map[string]interface{}{
		"user_id":        userID,
		"risks":          risks,
		"count":          len(risks),
		"threshold_days": threshold,
	}, nil
}

func checkSubTasks(goal string, tasks []models.Task, risks *[]RiskTask, now time.Time, threshold int) {
	for _, t := range tasks {
		if !t.Deadline.IsZero() {
			daysLeft := int(t.Deadline.Sub(now).Hours() / 24)
			if daysLeft <= threshold {
				*risks = append(*risks, RiskTask{
					Goal:     goal,
					TaskName: t.Title,
					Deadline: t.Deadline.Format("2006-01-02"),
					DaysLeft: daysLeft,
				})
			}
		}
		checkSubTasks(goal, t.SubTasks, risks, now, threshold)
	}
}

func get_user_progress(params map[string]interface{}, repo *repository.PlanRepository) (map[string]interface{}, error) {
	userID, ok := params["user_id"].(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("user_id required")
	}

	message, _ := params["message"].(string)

	plans, err := repo.GetAllByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user plans: %v", err)
	}
	if len(plans) == 0 {
		return nil, fmt.Errorf("no plans found for this user")
	}

	var plan *models.Plan
	if message != "" {
		plan, err = repo.FindGoalByAI(userID, message)
		if err != nil {
			return nil, fmt.Errorf("failed to find matching plan: %v", err)
		}
	} else {
		plan = &plans[0]
	}

	totalTasks := 0
	completedTasks := 0

	var countTasks func(tasks []models.Task)
	countTasks = func(tasks []models.Task) {
		for _, t := range tasks {
			totalTasks++
			if strings.EqualFold(t.Status, "Completed") {
				completedTasks++
			}
			countTasks(t.SubTasks)
		}
	}
	countTasks(plan.Tasks)

	progress := 0
	if totalTasks > 0 {
		progress = (completedTasks * 100) / totalTasks
	}

	return map[string]interface{}{
		"user_id":               userID,
		"goal":                  plan.Goal,
		"completion_percentage": progress,
		"total_tasks":           totalTasks,
		"completed_tasks":       completedTasks,
	}, nil
}

func provide_feedback(params map[string]interface{}) (map[string]interface{}, error) {
	progressDataRaw, ok := params["progress_data"]
	if !ok || progressDataRaw == nil {
		return nil, fmt.Errorf("progress_data is required for provide_feedback")
	}

	progress, ok := progressDataRaw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("progress_data must be a map[string]interface{}")
	}

	pct, ok := progress["completion_percentage"].(int)
	if !ok {
		if f, ok := progress["completion_percentage"].(float64); ok {
			pct = int(f)
		} else {
			return nil, fmt.Errorf("completion_percentage missing or invalid")
		}
	}

	goal, _ := progress["goal"].(string)
	totalTasks, _ := progress["total_tasks"].(int)
	completedTasks, _ := progress["completed_tasks"].(int)
	remainingTasks := totalTasks - completedTasks

	var message, suggestion, tone string

	if pct == 0 {
		tone = "Let's get started! 🚀"
		message = fmt.Sprintf("You haven't started working on '%s' yet. The best time to begin is now!", goal)
		suggestion = "Start with the first task to build momentum. Small progress is still progress!"
	} else if pct < 25 {
		tone = "Good start! 💪"
		message = fmt.Sprintf("You've completed %d out of %d tasks for '%s'. You're %d%% of the way there!", 
			completedTasks, totalTasks, goal, pct)
		suggestion = "You're building momentum! Try to complete at least 2-3 tasks this week to stay on track."
	} else if pct < 50 {
		tone = "You're making progress! 📈"
		message = fmt.Sprintf("Great work! You're %d%% done with '%s'. %d tasks completed, %d to go!", 
			pct, goal, completedTasks, remainingTasks)
		suggestion = "You're doing well, but let's pick up the pace a bit. Try to complete more tasks this week to reach 50%!"
	} else if pct < 75 {
		tone = "Halfway there! 🎯"
		message = fmt.Sprintf("Excellent progress! You've crossed the halfway mark on '%s' with %d%% completion!", 
			goal, pct)
		suggestion = "Keep this momentum going! You're on the right track. Focus on consistency to finish strong."
	} else if pct < 90 {
		tone = "Almost there! 🔥"
		message = fmt.Sprintf("Outstanding! You're %d%% done with '%s'. Only %d tasks remaining!", 
			pct, goal, remainingTasks)
		suggestion = "You're so close to the finish line! Push through these last few tasks and you'll achieve your goal!"
	} else if pct < 100 {
		tone = "Final sprint! 🏃"
		message = fmt.Sprintf("Incredible work! You're at %d%% completion for '%s'. Just %d task(s) left!", 
			pct, goal, remainingTasks)
		suggestion = "You're almost done! Complete these last tasks and celebrate your achievement!"
	} else {
		tone = "Goal achieved! 🎉"
		message = fmt.Sprintf("Congratulations! You've completed all %d tasks for '%s'!", totalTasks, goal)
		suggestion = "Amazing work! Time to set a new goal and keep the momentum going!"
	}

	return map[string]interface{}{
		"feedback": map[string]interface{}{
			"tone":               tone,
			"message":            message,
			"suggestion":         suggestion,
			"progress_summary": map[string]interface{}{
				"goal":              goal,
				"completion_percentage": pct,
				"completed_tasks":   completedTasks,
				"remaining_tasks":   remainingTasks,
				"total_tasks":       totalTasks,
			},
		},
	}, nil
}

func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func extractDelayDays(message string) int {
	re := regexp.MustCompile(`(\d+)\s*day`)
	match := re.FindStringSubmatch(strings.ToLower(message))
	if len(match) > 1 {
		days, _ := strconv.Atoi(match[1])
		return days
	}
	return 0
}

func generate_alternative_plans(params map[string]interface{}) (map[string]interface{}, error) {
	goalID, ok := params["goal_id"].(string)
	if !ok || goalID == "" {
		return nil, fmt.Errorf("goal_id required")
	}

	options := []map[string]string{
		{"type": "speed", "description": "Focus on completing tasks faster, may reduce quality."},
		{"type": "balance", "description": "Balanced approach between speed and quality."},
		{"type": "quality", "description": "Focus on doing tasks with highest quality, may take longer."},
	}

	return map[string]interface{}{
		"goal_id": goalID,
		"options": options,
	}, nil
}

func handle_general_query(params map[string]interface{}, repo *repository.PlanRepository) (map[string]interface{}, error) {
	userID, ok := params["user_id"].(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("user_id required")
	}

	message, ok := params["message"].(string)
	if !ok || message == "" {
		return nil, fmt.Errorf("message required")
	}

	plans, err := repo.GetAllByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user plans: %v", err)
	}

	if len(plans) == 0 {
		return map[string]interface{}{
			"response": "You don't have any plans yet. Would you like to create one? Just tell me your goal!",
		}, nil
	}

	var contextBuilder strings.Builder
	contextBuilder.WriteString("User's current plans:\n")
	
	for i, plan := range plans {
		totalTasks := 0
		completedTasks := 0
		var countTasks func([]models.Task)
		countTasks = func(tasks []models.Task) {
			for _, t := range tasks {
				totalTasks++
				if strings.EqualFold(t.Status, "Completed") {
					completedTasks++
				}
				countTasks(t.SubTasks)
			}
		}
		countTasks(plan.Tasks)

		progress := 0
		if totalTasks > 0 {
			progress = (completedTasks * 100) / totalTasks
		}

		contextBuilder.WriteString(fmt.Sprintf("\n%d. Goal: %s\n", i+1, plan.Goal))
		contextBuilder.WriteString(fmt.Sprintf("   Progress: %d%% (%d/%d tasks completed)\n", progress, completedTasks, totalTasks))
		
		// Add upcoming deadlines
		now := time.Now()
		upcomingCount := 0
		for _, task := range plan.Tasks {
			if !task.Deadline.IsZero() && task.Deadline.After(now) && strings.EqualFold(task.Status, "Pending") {
				daysLeft := int(task.Deadline.Sub(now).Hours() / 24)
				if daysLeft <= 7 {
					contextBuilder.WriteString(fmt.Sprintf("   - %s (due in %d days)\n", task.Title, daysLeft))
					upcomingCount++
					if upcomingCount >= 3 {
						break
					}
				}
			}
		}
	}

	// Build AI prompt
	prompt := fmt.Sprintf(`You are a helpful task planning assistant. A user asked: "%s"

Here's what you know about the user:
%s

Provide a helpful, conversational response that:
1. Directly answers their question using the context above
2. Is encouraging and supportive
3. Suggests actionable next steps if relevant
4. Keep it concise (2-3 sentences max)

Response:`, message, contextBuilder.String())

	// Call OpenAI API
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return map[string]interface{}{
			"response": "I understand your question, but I'm having trouble generating a detailed response right now. Try asking about your progress, risks, or rescheduling tasks!",
		}, nil
	}

	aiResponse, err := CallOpenAIAPI(prompt)
	if err != nil {
		return map[string]interface{}{
			"response": "I understand your question, but I'm having trouble generating a detailed response right now. Try asking about your progress, risks, or rescheduling tasks!",
		}, nil
	}

	return map[string]interface{}{
		"response": strings.TrimSpace(aiResponse),
		"context_used": true,
	}, nil
}