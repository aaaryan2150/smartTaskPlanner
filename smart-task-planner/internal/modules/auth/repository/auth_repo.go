package repository

import (
	"context"
	"errors"
	"time"
	"smart-task-planner/internal/database"
	"smart-task-planner/internal/modules/auth/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func getUserCollection() *mongo.Collection {
	return database.DB.Collection("users")
}

// CreateUser inserts a new user
func CreateUser(user models.User) (models.User, error) {
	collection := getUserCollection()
	user.ID = primitive.NewObjectID()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, user)
	return user, err
}

// GetUserByEmail fetches a user by email
func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	collection := getUserCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return user, errors.New("user not found")
		}
		return user, err
	}
	return user, nil
}

// GetUserByID fetches user by ID
func GetUserByID(id string) (models.User, error) {
	var user models.User
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, errors.New("invalid user ID")
	}

	collection := getUserCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return user, errors.New("user not found")
		}
		return user, err
	}
	return user, nil
}
