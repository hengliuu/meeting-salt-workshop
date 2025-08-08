package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"api/internal/config"
	"api/internal/models"
	"api/internal/repositories"
	"api/internal/utils"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

type AuthService struct {
	userRepo   repositories.UserRepository
	jwtManager *utils.JWTManager
	config     *config.Config
	oauthConfig *oauth2.Config
}

type MicrosoftUser struct {
	ID          string `json:"id"`
	Email       string `json:"mail"`
	FirstName   string `json:"givenName"`
	LastName    string `json:"surname"`
	DisplayName string `json:"displayName"`
}

type LoginResponse struct {
	User  *models.User `json:"user"`
	Token string       `json:"token"`
}

func NewAuthService(userRepo repositories.UserRepository, config *config.Config) *AuthService {
	jwtManager := utils.NewJWTManager(config.Auth.JWTSecret)
	
	oauthConfig := &oauth2.Config{
		ClientID:     config.Auth.MicrosoftClientID,
		ClientSecret: config.Auth.MicrosoftClientSecret,
		RedirectURL:  config.Auth.MicrosoftRedirectURL,
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     microsoft.AzureADEndpoint(""),
	}

	return &AuthService{
		userRepo:    userRepo,
		jwtManager:  jwtManager,
		config:      config,
		oauthConfig: oauthConfig,
	}
}

func (s *AuthService) GetAuthURL(state string) string {
	return s.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (s *AuthService) HandleCallback(code string) (*LoginResponse, error) {
	token, err := s.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %v", err)
	}

	msUser, err := s.getMicrosoftUserInfo(token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}

	user, err := s.findOrCreateUser(msUser)
	if err != nil {
		return nil, fmt.Errorf("failed to find or create user: %v", err)
	}

	now := time.Now()
	user.LastLogin = &now
	user, err = s.userRepo.Update(user)
	if err != nil {
		return nil, fmt.Errorf("failed to update last login: %v", err)
	}

	jwtToken, err := s.jwtManager.GenerateToken(user.ID, user.Email, string(user.Role), user.MicrosoftID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	return &LoginResponse{
		User:  user,
		Token: jwtToken,
	}, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*utils.JWTClaims, error) {
	return s.jwtManager.ValidateToken(tokenString)
}

func (s *AuthService) RefreshToken(tokenString string) (string, error) {
	return s.jwtManager.RefreshToken(tokenString)
}

func (s *AuthService) getMicrosoftUserInfo(token *oauth2.Token) (*MicrosoftUser, error) {
	client := s.oauthConfig.Client(context.Background(), token)
	
	resp, err := client.Get("https://graph.microsoft.com/v1.0/me")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("microsoft graph API returned status %d", resp.StatusCode)
	}

	var msUser MicrosoftUser
	if err := json.NewDecoder(resp.Body).Decode(&msUser); err != nil {
		return nil, err
	}

	return &msUser, nil
}

func (s *AuthService) findOrCreateUser(msUser *MicrosoftUser) (*models.User, error) {
	user, err := s.userRepo.GetByMicrosoftID(msUser.ID)
	if err == nil {
		return user, nil
	}

	user, err = s.userRepo.GetByEmail(msUser.Email)
	if err == nil {
		user.MicrosoftID = msUser.ID
		return s.userRepo.Update(user)
	}

	newUser := &models.User{
		MicrosoftID: msUser.ID,
		Email:       msUser.Email,
		FirstName:   msUser.FirstName,
		LastName:    msUser.LastName,
		DisplayName: msUser.DisplayName,
		Role:        models.RoleEmployee,
		IsActive:    true,
	}

	return s.userRepo.Create(newUser)
}