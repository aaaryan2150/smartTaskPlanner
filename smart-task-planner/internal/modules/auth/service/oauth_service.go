package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"smart-task-planner/internal/modules/auth/models"
	"smart-task-planner/internal/modules/auth/repository"
	"smart-task-planner/internal/utils"
	"sync"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig *oauth2.Config
	once              sync.Once
)

// InitGoogleOAuth initializes the Google OAuth config
func InitGoogleOAuth() {
	once.Do(func() {
		googleOauthConfig = &oauth2.Config{
			RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint:     google.Endpoint,
		}
	})
}

func GetGoogleLoginURL(state string) string {
	if googleOauthConfig == nil {
		InitGoogleOAuth()
	}
	
	url := googleOauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	fmt.Println("Generated OAuth URL:", url)
	return url
}

func HandleGoogleCallback(code string) (string, string, error) {
	if googleOauthConfig == nil {
		InitGoogleOAuth()
	}

	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return "", "", fmt.Errorf("token exchange failed: %w", err)
	}

	client := googleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return "", "", fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	var userInfo struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return "", "", fmt.Errorf("failed to decode user info: %w", err)
	}

	if userInfo.Email == "" {
		return "", "", errors.New("email not found in user info")
	}

	// Check if user exists or create
	user, err := repository.GetUserByEmail(userInfo.Email)
	if err != nil {
		// User doesn't exist, create new user
		newUser := models.User{
			Name:  userInfo.Name,
			Email: userInfo.Email,
		}
		createdUser, createErr := repository.CreateUser(newUser)
		if createErr != nil {
			return "", "", fmt.Errorf("failed to create user: %w", createErr)
		}
		user = createdUser
	}

	jwtToken, err := utils.GenerateToken(user.ID.Hex())
	if err != nil {
		return "", "", fmt.Errorf("failed to generate token: %w", err)
	}
	
	return jwtToken, user.Email, nil
}