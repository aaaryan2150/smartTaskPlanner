package repository

import (
	"context"
	"log"
	"smart-task-planner/internal/modules/plan/models"
	"time"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"smart-task-planner/internal/modules/plan/ai"
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

func (r *PlanRepository) FindGoalByAI(userID, message string) (*models.Plan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1️⃣ Fetch all plans for the user
	cursor, err := r.Collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user plans: %v", err)
	}
	defer cursor.Close(ctx)

	var plans []models.Plan
	if err := cursor.All(ctx, &plans); err != nil {
		return nil, fmt.Errorf("failed to decode plans: %v", err)
	}

	if len(plans) == 0 {
		return nil, fmt.Errorf("no plans found for this user")
	}

	// 2️⃣ Create a list of plan goals
	var goalList []string
	for _, p := range plans {
		goalList = append(goalList, p.Goal)
	}

	// 3️⃣ Ask the AI which goal best matches the user message
	bestGoal, err := ai.AskOpenAIForBestGoal(message, goalList)
	if err != nil {
		return nil, err
	}

	// 4️⃣ Find the plan object corresponding to the chosen goal
	for _, p := range plans {
		if strings.EqualFold(p.Goal, bestGoal) {
			return &p, nil
		}
	}

	return nil, fmt.Errorf("AI suggested goal not found: %s", bestGoal)
}

func (r *PlanRepository) UpdatePlan(plan *models.Plan) (*models.Plan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.Collection.UpdateOne(
		ctx,
		bson.M{"_id": plan.ID},             // match by plan ID
		bson.M{"$set": bson.M{"tasks": plan.Tasks}}, // update tasks
	)
	if err != nil {
		return nil, err
	}

	return plan, nil
}
