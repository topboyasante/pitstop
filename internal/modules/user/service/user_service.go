package service

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/topboyasante/pitstop/internal/core/logger"
	"github.com/topboyasante/pitstop/internal/modules/user/domain"
	"github.com/topboyasante/pitstop/internal/modules/user/dto"
	"github.com/topboyasante/pitstop/internal/modules/user/repository"
	"github.com/topboyasante/pitstop/internal/shared/events"
)

// UserService handles user business logic
type UserService struct {
	userRepo  *repository.UserRepository
	validator *validator.Validate
	eventBus  *events.EventBus
}

// NewUserService creates a new user service instance
func NewUserService(userRepo *repository.UserRepository, validator *validator.Validate, eventBus *events.EventBus) *UserService {
	return &UserService{
		userRepo:  userRepo,
		validator: validator,
		eventBus:  eventBus,
	}
}

// CreateUser creates a user from OAuth provider data
func (s *UserService) CreateUser(req dto.CreateUserRequest) (*dto.UserResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		logger.Error("OAuth user validation failed", 
			"event", "user.oauth_validation_failed",
			"provider", req.Provider,
			"error", err)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if user already exists by provider ID
	existingUser, err := s.userRepo.GetByProviderID(req.Provider, req.ProviderID)
	if err == nil {
		logger.Info("Existing OAuth user found",
			"event", "user.oauth_login",
			"provider", req.Provider,
			"user_id", existingUser.ID)
		return s.mapUserToResponse(existingUser), nil
	}

	// Create new user from OAuth data
	user := &domain.User{
		ProviderID: req.ProviderID,
		Provider:   req.Provider,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Email:      req.Email,
		AvatarURL:  req.AvatarURL,
		Locale:     req.Locale,
	}

	if err := s.userRepo.Create(user); err != nil {
		logger.Error("Failed to create OAuth user",
			"event", "user.oauth_create_failed",
			"provider", req.Provider,
			"error", err)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	logger.Info("OAuth user created successfully",
		"event", "user.oauth_created",
		"provider", req.Provider,
		"user_id", user.ID)

	// Publish user created event
	// s.eventBus.Publish("UserRegistered", &events.UserRegistered{
	// 	UserID:   user.ID,
	// 	Email:    user.Email,
	// 	Provider: user.Provider,
	// })

	return s.mapUserToResponse(user), nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(id string) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return s.mapUserToResponse(user), nil
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(email string) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return s.mapUserToResponse(user), nil
}

// UpdateUser updates user profile
func (s *UserService) UpdateUser(userID string, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Update fields if provided
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}

	if err := s.userRepo.Update(user); err != nil {
		logger.Error("Failed to update user profile",
			"event", "user.update_failed",
			"user_id", userID,
			"error", err)
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	logger.Info("User profile updated successfully",
		"event", "user.updated",
		"user_id", userID)

	return s.mapUserToResponse(user), nil
}

// GetAllUsers retrieves all users with pagination
func (s *UserService) GetAllUsers(page, limit int) (*dto.UsersResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	users, totalCount, err := s.userRepo.GetAll(page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %w", err)
	}

	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = *s.mapUserToResponse(&user)
	}

	hasNext := int64((page-1)*limit+len(users)) < totalCount

	return &dto.UsersResponse{
		Users:      userResponses,
		TotalCount: totalCount,
		Page:       page,
		Limit:      limit,
		HasNext:    hasNext,
	}, nil
}

// mapUserToResponse converts domain User to UserResponse DTO
func (s *UserService) mapUserToResponse(user *domain.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        user.ID,
		Provider:  user.Provider,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
		FullName:  user.FullName(),
		IsOAuth:   user.IsOAuthUser(),
		CreatedAt: user.CreatedAt,
	}
}
