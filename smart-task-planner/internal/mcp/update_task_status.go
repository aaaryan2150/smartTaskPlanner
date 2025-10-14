// mcp/update_task_status.go
package mcp

import (
	"context"
	"fmt"
	"smart-task-planner/internal/modules/plan/models"
	"smart-task-planner/internal/modules/plan/repository"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UpdateTaskStatus(taskID string, status string, repo *repository.PlanRepository) (*models.Task, error) {
	if taskID == "" || status == "" {
		return nil, fmt.Errorf("taskID and status required")
	}

	objID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return nil, fmt.Errorf("invalid task ID: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Update the task status
	filter := bson.M{"tasks._id": objID}
	update := bson.M{"$set": bson.M{"tasks.$.status": status}}

	res := repo.Collection.FindOneAndUpdate(ctx, filter, update)
	if res.Err() != nil {
		return nil, res.Err()
	}

	// Fetch updated task
	var plan models.Plan
	if err := repo.Collection.FindOne(ctx, filter).Decode(&plan); err != nil {
		return nil, err
	}

	for _, t := range plan.Tasks {
		if t.ID == objID {
			return &t, nil
		}
	}

	return nil, fmt.Errorf("task not found after update")
}
