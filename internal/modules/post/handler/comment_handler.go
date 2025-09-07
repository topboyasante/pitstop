package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/topboyasante/pitstop/internal/core/logger"
	"github.com/topboyasante/pitstop/internal/core/response"
	"github.com/topboyasante/pitstop/internal/modules/post/dto"
	"github.com/topboyasante/pitstop/internal/modules/post/service"
)

// CommentHandler handles comment HTTP requests
type CommentHandler struct {
	commentService *service.CommentService
}

// NewCommentHandler creates a new comment handler instance
func NewCommentHandler(commentService *service.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

// CreateComment handles comment creation
// @Summary Create a new comment
// @Description Create a new comment on a post
// @Tags Comments
// @Accept json
// @Produce json
// @Param post_id path string true "Post ID"
// @Param comment body dto.CreateCommentRequest true "Comment data"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /posts/{post_id}/comments [post]
func (h *CommentHandler) CreateComment(c *fiber.Ctx) error {
	postID := c.Params("post_id")
	userID := c.Locals("userID").(string)

	var req dto.CreateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error("Invalid request body", "error", err)
		return response.ValidationErrorJSON(c, "Invalid request body", err.Error())
	}

	// Basic validation
	if req.Content == "" {
		return response.ValidationErrorJSON(c, "Content is required", "")
	}

	comment, err := h.commentService.CreateComment(postID, userID, &req)
	if err != nil {
		if err.Error() == "post not found" {
			return response.NotFoundJSON(c, "Post not found")
		}
		logger.Error("Failed to create comment", "error", err)
		return response.InternalErrorJSON(c, "Failed to create comment")
	}

	return response.SuccessJSON(c, comment, "Comment created successfully")
}

// CreateReply handles reply creation
// @Summary Create a reply to a comment
// @Description Create a reply to an existing comment
// @Tags Comments
// @Accept json
// @Produce json
// @Param post_id path string true "Post ID"
// @Param parent_comment_id path string true "Parent Comment ID"
// @Param reply body dto.CreateCommentRequest true "Reply data"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /posts/{post_id}/comments/{parent_comment_id}/reply [post]
func (h *CommentHandler) CreateReply(c *fiber.Ctx) error {
	postID := c.Params("post_id")
	parentCommentID := c.Params("parent_comment_id")
	userID := c.Locals("userID").(string)

	var req dto.CreateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error("Invalid request body", "error", err)
		return response.ValidationErrorJSON(c, "Invalid request body", err.Error())
	}

	// Basic validation
	if req.Content == "" {
		return response.ValidationErrorJSON(c, "Content is required", "")
	}

	reply, err := h.commentService.CreateReply(postID, parentCommentID, userID, &req)
	if err != nil {
		if err.Error() == "post not found" || err.Error() == "parent comment not found" {
			return response.NotFoundJSON(c, err.Error())
		}
		if err.Error() == "parent comment doesn't belong to this post" {
			return response.ValidationErrorJSON(c, err.Error(), "")
		}
		logger.Error("Failed to create reply", "error", err)
		return response.InternalErrorJSON(c, "Failed to create reply")
	}

	return response.SuccessJSON(c, reply, "Reply created successfully")
}

// GetComments handles retrieving comments for a post
// @Summary Get comments for a post
// @Description Retrieve all comments for a specific post
// @Tags Comments
// @Produce json
// @Param post_id path string true "Post ID"
// @Success 200 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /posts/{post_id}/comments [get]
func (h *CommentHandler) GetComments(c *fiber.Ctx) error {
	postID := c.Params("post_id")

	comments, err := h.commentService.GetCommentsByPostID(postID)
	if err != nil {
		logger.Error("Failed to retrieve comments", "error", err)
		return response.InternalErrorJSON(c, "Failed to retrieve comments")
	}

	return response.SuccessJSON(c, comments, "Comments retrieved successfully")
}