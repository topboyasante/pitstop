package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/topboyasante/pitstop/internal/core/logger"
	"github.com/topboyasante/pitstop/internal/core/response"
	"github.com/topboyasante/pitstop/internal/modules/post/service"
	utils "github.com/topboyasante/pitstop/internal/shared/utils"
)

// LikeHandler handles HTTP requests for likes
type LikeHandler struct {
	likeService *service.LikeService
}

// NewLikeHandler creates a new like handler instance
func NewLikeHandler(likeService *service.LikeService) *LikeHandler {
	return &LikeHandler{
		likeService: likeService,
	}
}

// ToggleLike toggles a like for a post
// @Summary Toggle like on a post
// @Description Like or unlike a post
// @Tags likes
// @Accept json
// @Produce json
// @Param post_id path string true "Post ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Security BearerAuth
// @Router /posts/{post_id}/like [post]
func (h *LikeHandler) ToggleLike(c *fiber.Ctx) error {
	postID := c.Params("post_id")
	if strings.TrimSpace(postID) == "" {
		return response.ValidationErrorJSON(c, "Invalid post ID", "Post ID cannot be empty")
	}

	// Extract user ID from JWT claims
	userID, err := utils.ExtractUserIDFromContext(c)
	if err != nil {
		logger.Error("Failed to extract user ID from context", "error", err)
		return response.UnauthorizedJSON(c)
	}

	result, err := h.likeService.ToggleLike(postID, userID)
	if err != nil {
		logger.Error("Failed to toggle like", "post_id", postID, "user_id", userID, "error", err)
		if strings.Contains(err.Error(), "post not found") {
			return response.NotFoundJSON(c, "Post")
		}
		return response.InternalErrorJSON(c, "Failed to toggle like")
	}

	action := "liked"
	if !result.Liked {
		action = "unliked"
	}
	
	logger.Info("Like toggled successfully", "post_id", postID, "user_id", userID, "action", action)
	return response.SuccessJSON(c, result, "Post "+action+" successfully")
}

// GetLikesByPost retrieves all likes for a post
// @Summary Get likes for a post
// @Description Retrieve all users who liked a specific post
// @Tags likes
// @Accept json
// @Produce json
// @Param post_id path string true "Post ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /posts/{post_id}/likes [get]
func (h *LikeHandler) GetLikesByPost(c *fiber.Ctx) error {
	postID := c.Params("post_id")
	if strings.TrimSpace(postID) == "" {
		return response.ValidationErrorJSON(c, "Invalid post ID", "Post ID cannot be empty")
	}

	likes, err := h.likeService.GetLikesByPost(postID)
	if err != nil {
		logger.Error("Failed to retrieve likes", "post_id", postID, "error", err)
		if strings.Contains(err.Error(), "post not found") {
			return response.NotFoundJSON(c, "Post")
		}
		return response.InternalErrorJSON(c, "Failed to retrieve likes")
	}

	return response.SuccessJSON(c, likes, "Likes retrieved successfully")
}

// CheckUserLiked checks if the current user has liked a specific post
// @Summary Check if user liked a post
// @Description Check if the authenticated user has liked a specific post
// @Tags likes
// @Accept json
// @Produce json
// @Param post_id path string true "Post ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Security BearerAuth
// @Router /posts/{post_id}/like/status [get]
func (h *LikeHandler) CheckUserLiked(c *fiber.Ctx) error {
	postID := c.Params("post_id")
	if strings.TrimSpace(postID) == "" {
		return response.ValidationErrorJSON(c, "Invalid post ID", "Post ID cannot be empty")
	}

	// Extract user ID from JWT claims
	userID, err := utils.ExtractUserIDFromContext(c)
	if err != nil {
		logger.Error("Failed to extract user ID from context", "error", err)
		return response.UnauthorizedJSON(c)
	}

	liked, err := h.likeService.CheckUserLiked(postID, userID)
	if err != nil {
		logger.Error("Failed to check user like status", "post_id", postID, "user_id", userID, "error", err)
		return response.InternalErrorJSON(c, "Failed to check like status")
	}

	return response.SuccessJSON(c, map[string]bool{"liked": liked}, "Like status retrieved successfully")
}