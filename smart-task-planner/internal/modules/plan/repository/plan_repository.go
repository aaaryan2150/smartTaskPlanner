package repository

import (
	"context"
	"log"
	"smart-task-planner/internal/modules/plan/models"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PlanRepository struct {
	Collection *mongo.Collection
}

func NewPlanRepository(db *mongo.Database) *PlanRepository {
	return &PlanRepository{
		Collection: db.Collection("plans"),
	}
}

func (r *PlanRepository) Create(plan *models.Plan) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.Collection.InsertOne(ctx, plan)
	if err != nil {
		log.Println("Error creating plan:", err)
	}
	return err
}

// GetAllByUser fetches plans belonging to a specific user
func (r *PlanRepository) GetAllByUser(userID string) ([]models.Plan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.Collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var plans []models.Plan
	if err := cursor.All(ctx, &plans); err != nil {
		return nil, err
	}
	return plans, nil
}

func (r *PlanRepository) AddSubTasks(planID, taskID string, subtasks []models.Task) (*models.Plan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var plan models.Plan
	objectID, err := primitive.ObjectIDFromHex(planID)
	if err != nil {
		return nil, err
	}

	if err := r.Collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&plan); err != nil {
		return nil, err
	}

	// Recursively find the target task
	var addSubs func([]models.Task) []models.Task
	addSubs = func(tasks []models.Task) []models.Task {
		for i := range tasks {
			if tasks[i].ID.Hex() == taskID {
				tasks[i].SubTasks = append(tasks[i].SubTasks, subtasks...)
				return tasks
			}
			tasks[i].SubTasks = addSubs(tasks[i].SubTasks)
		}
		return tasks
	}

	plan.Tasks = addSubs(plan.Tasks)

	_, err = r.Collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{"tasks": plan.Tasks}},
	)
	if err != nil {
		return nil, err
	}

	return &plan, nil
}
