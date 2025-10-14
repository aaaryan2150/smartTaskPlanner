package mcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
	"smart-task-planner/internal/modules/plan/models"
)

// type OpenAIResponse struct {
// 	Choices []struct {
// 		Message struct {
// 			Content string `json:"content"`
// 		} `json:"message"`
// 	} `json:"choices"`
// }

// Temp struct for parsing OpenAI JSON (deadline as string)
type AITask struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	DeadlineStr string `json:"deadline"`
}

type TaskPlan struct {
	Tasks []models.Task
}

func CreateTaskPlan(params map[string]interface{}) (TaskPlan, error) {
	goal, ok := params["goal"].(string)
	if !ok || goal == "" {
		return TaskPlan{}, fmt.Errorf("goal required")
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return TaskPlan{}, fmt.Errorf("OPENAI_API_KEY not set")
	}

	prompt := fmt.Sprintf(`You are an expert AI task planner.
Generate at least 10 actionable, detailed tasks for this goal:
"%s"
Each task must have:
- title
- description
- deadline (YYYY-MM-DD, evenly distributed across goal duration)
Return ONLY valid JSON:
[
  {"title": "...", "description": "...", "deadline": "..."}
]`, goal)

	url := "https://api.openai.com/v1/chat/completions"
	payload := map[string]interface{}{
		"model": "gpt-4.1-mini",
		"messages": []map[string]string{
			{"role": "system", "content": "You are an expert AI task planner."},
			{"role": "user", "content": prompt},
		},
		"temperature": 0.7,
	}

	data, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return TaskPlan{}, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return TaskPlan{}, fmt.Errorf("OpenAI error: %s", string(body))
	}

	var openAIResp OpenAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return TaskPlan{}, err
	}

	if len(openAIResp.Choices) == 0 {
		return TaskPlan{}, fmt.Errorf("no response from OpenAI")
	}

	raw := strings.TrimSpace(openAIResp.Choices[0].Message.Content)
	// Clean markdown
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimSuffix(raw, "```")
	raw = strings.TrimSpace(raw)

	var aiTasks []AITask
	if err := json.Unmarshal([]byte(raw), &aiTasks); err != nil {
		return TaskPlan{}, fmt.Errorf("failed to parse AI JSON: %v\nRaw: %s", err, raw)
	}

	// Convert to models.Task
	var tasks []models.Task
	now := time.Now()
	for i, t := range aiTasks {
		deadline, err := time.Parse("2006-01-02", t.DeadlineStr)
		if err != nil {
			deadline = now.Add(time.Duration(i+1) * 24 * time.Hour)
		}
		tasks = append(tasks, models.Task{
			Title:       t.Title,
			Description: t.Description,
			Status:      "Pending",
			Deadline:    deadline,
		})
	}

	return TaskPlan{Tasks: tasks}, nil
}
