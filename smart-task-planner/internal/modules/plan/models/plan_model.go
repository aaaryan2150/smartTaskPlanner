package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Status      string             `bson:"status" json:"status"`
	Deadline    time.Time          `bson:"deadline" json:"deadline"`
}

type Plan struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID string             `bson:"user_id" json:"user_id"` // link to user
	Goal   string             `bson:"goal" json:"goal"`
	Tasks  []Task             `bson:"tasks" json:"tasks"`
}
