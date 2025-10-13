package service

import (
	"errors"
	"smart-task-planner/internal/modules/auth/dto"
	"smart-task-planner/internal/modules/auth/models"
	"smart-task-planner/internal/modules/auth/repository"
	"smart-task-planner/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

// Register a new user
func Register(req dto.RegisterRequest) (map[string]interface{}, error) {
	_, err := repository.GetUserByEmail(req.Email)
	if err == nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashed),
	}

	user, err = repository.CreateUser(user)
	if err != nil {
		return nil, err
	}

	token, _ := utils.GenerateToken(user.ID.Hex())

	return map[string]interface{}{
		"id":    user.ID.Hex(),
		"name":  user.Name,
		"email": user.Email,
		"token": token,
	}, nil
}

// Login user
func Login(req dto.LoginRequest) (map[string]interface{}, error) {
	user, err := repository.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, _ := utils.GenerateToken(user.ID.Hex())

	return map[string]interface{}{
		"id":    user.ID.Hex(),
		"name":  user.Name,
		"email": user.Email,
		"token": token,
	}, nil
}
