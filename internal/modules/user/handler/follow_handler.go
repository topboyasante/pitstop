package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/topboyasante/pitstop/internal/core/logger"
	"github.com/topboyasante/pitstop/internal/core/response"
	"github.com/topboyasante/pitstop/internal/modules/user/service"
	utils "github.com/topboyasante/pitstop/internal/shared/utils"
)

// FollowHandler handles HTTP requests for follows
type FollowHandler struct {
	followService *service.FollowService
}

// NewFollowHandler creates a new follow handler instance
func NewFollowHandler(followService *service.FollowService) *FollowHandler {
	return &FollowHandler{
		followService: followService,
	}
}

// ToggleFollow toggles a follow relationship for a user
// @Summary Toggle follow on a user
// @Description Follow or unfollow a user
// @Tags follows
// @Accept json
// @Produce json
// @Param user_id path string true "User ID to follow/unfollow"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Security BearerAuth
// @Router /users/{user_id}/follow [post]
func (h *FollowHandler) ToggleFollow(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	if strings.TrimSpace(userID) == "" {
		return response.ValidationErrorJSON(c, "Invalid user ID", "User ID cannot be empty")
	}

	// Extract follower ID from JWT claims
	followerID, err := utils.ExtractUserIDFromContext(c)
	if err != nil {
		logger.Error("Failed to extract user ID from context", "error", err)
		return response.UnauthorizedJSON(c)
	}

	result, err := h.followService.ToggleFollow(followerID, userID)
	if err != nil {
		logger.Error("Failed to toggle follow", "follower_id", followerID, "following_id", userID, "error", err)
		if strings.Contains(err.Error(), "user not found") {
			return response.NotFoundJSON(c, "User")
		}
		if strings.Contains(err.Error(), "cannot follow themselves") {
			return response.ValidationErrorJSON(c, "Invalid operation", "Users cannot follow themselves")
		}
		return response.InternalErrorJSON(c, "Failed to toggle follow")
	}

	action := "followed"
	if !result.IsFollowing {
		action = "unfollowed"
	}

	logger.Info("Follow toggled successfully", "follower_id", followerID, "following_id", userID, "action", action)
	return response.SuccessJSON(c, result, "User "+action+" successfully")
}

// GetFollowers retrieves all followers for a user
// @Summary Get followers for a user
// @Description Retrieve all users who follow a specific user
// @Tags follows
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /users/{user_id}/followers [get]
func (h *FollowHandler) GetFollowers(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	if strings.TrimSpace(userID) == "" {
		return response.ValidationErrorJSON(c, "Invalid user ID", "User ID cannot be empty")
	}

	followers, err := h.followService.GetFollowers(userID)
	if err != nil {
		logger.Error("Failed to retrieve followers", "user_id", userID, "error", err)
		if strings.Contains(err.Error(), "user not found") {
			return response.NotFoundJSON(c, "User")
		}
		return response.InternalErrorJSON(c, "Failed to retrieve followers")
	}

	return response.SuccessJSON(c, followers, "Followers retrieved successfully")
}

// GetFollowing retrieves all users that a user is following
// @Summary Get users being followed
// @Description Retrieve all users that a specific user is following
// @Tags follows
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /users/{user_id}/following [get]
func (h *FollowHandler) GetFollowing(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	if strings.TrimSpace(userID) == "" {
		return response.ValidationErrorJSON(c, "Invalid user ID", "User ID cannot be empty")
	}

	following, err := h.followService.GetFollowing(userID)
	if err != nil {
		logger.Error("Failed to retrieve following", "user_id", userID, "error", err)
		if strings.Contains(err.Error(), "user not found") {
			return response.NotFoundJSON(c, "User")
		}
		return response.InternalErrorJSON(c, "Failed to retrieve following")
	}

	return response.SuccessJSON(c, following, "Following retrieved successfully")
}

// CheckFollowStatus checks if the current user is following a specific user
// @Summary Check follow status
// @Description Check if the authenticated user is following a specific user
// @Tags follows
// @Accept json
// @Produce json
// @Param user_id path string true "User ID to check"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Security BearerAuth
// @Router /users/{user_id}/follow/status [get]
func (h *FollowHandler) CheckFollowStatus(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	if strings.TrimSpace(userID) == "" {
		return response.ValidationErrorJSON(c, "Invalid user ID", "User ID cannot be empty")
	}

	// Extract follower ID from JWT claims
	followerID, err := utils.ExtractUserIDFromContext(c)
	if err != nil {
		logger.Error("Failed to extract user ID from context", "error", err)
		return response.UnauthorizedJSON(c)
	}

	isFollowing, err := h.followService.CheckUserFollowing(followerID, userID)
	if err != nil {
		logger.Error("Failed to check follow status", "follower_id", followerID, "following_id", userID, "error", err)
		return response.InternalErrorJSON(c, "Failed to check follow status")
	}

	return response.SuccessJSON(c, map[string]bool{"is_following": isFollowing}, "Follow status retrieved successfully")
}