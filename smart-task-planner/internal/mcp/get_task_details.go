package mcp

import (
	"context"
	"fmt"
	"smart-task-planner/internal/modules/plan/models"
	"smart-task-planner/internal/modules/plan/repository"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
)

func GetTaskDetails(taskID string, repo *repository.PlanRepository) (*models.Task, error) {
	if taskID == "" {
		return nil, fmt.Errorf("taskID required")
	}

	objID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return nil, fmt.Errorf("invalid task ID: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var plan models.Plan
	err = repo.Collection.FindOne(ctx, bson.M{"tasks._id": objID}).Decode(&plan)
	if err != nil {
		return nil, err
	}

	for _, t := range plan.Tasks {
		if t.ID == objID {
			return &t, nil
		}
	}

	return nil, fmt.Errorf("task not found")
}
