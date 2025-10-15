package ai

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	openai "github.com/sashabaranov/go-openai"
)

// client will be initialized lazily
var (
	client     *openai.Client
	clientOnce sync.Once
)

// getClient initializes the OpenAI client safely
func getClient() *openai.Client {
	clientOnce.Do(func() {
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			panic("OPENAI_API_KEY is not set in the environment")
		}
		client = openai.NewClient(apiKey)
	})
	return client
}

// AskOpenAIForBestGoal asks OpenAI to match a user message to a goal
func AskOpenAIForBestGoal(message string, goals []string) (string, error) {
	if len(goals) == 0 {
		return "", fmt.Errorf("no goals provided")
	}

	ctx := context.Background()
prompt := fmt.Sprintf(`
You are a goal-matching assistant for a task planning app.

User message: "%s"

Here are the available goals:
%s

Rules:
- Choose exactly ONE goal from the list that best matches the user's message.
- Respond with the goal **exactly as written above**, no extra text.
- If none of the goals are relevant, reply with "NONE" (in uppercase, by itself).
`, message, strings.Join(goals, "\n"))

	resp, err := getClient().CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: "You are a precise AI that matches user intents to existing goals."},
			{Role: "user", Content: prompt},
		},
	})
	if err != nil {
		return "", fmt.Errorf("OpenAI API error: %v", err)
	}

	answer := strings.TrimSpace(resp.Choices[0].Message.Content)
	if strings.EqualFold(answer, "NONE") {
		return "", fmt.Errorf("AI could not determine a matching goal")
	}

	return answer, nil
}
