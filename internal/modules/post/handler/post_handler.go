package handler

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/topboyasante/pitstop/internal/core/logger"
	"github.com/topboyasante/pitstop/internal/core/response"
	"github.com/topboyasante/pitstop/internal/modules/post/dto"
	"github.com/topboyasante/pitstop/internal/modules/post/service"
)

// PostHandler handles HTTP requests for posts
type PostHandler struct {
	postService *service.PostService
}

// NewPostHandler creates a new post handler instance
func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

// GetAllPosts retrieves all posts
// @Summary Get all posts
// @Description Retrieve a paginated list of posts
// @Tags posts
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Posts per page" default(20)
// @Success 200 {object} response.APIResponse
// @Router /posts [get]
func (h *PostHandler) GetAllPosts(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	posts, err := h.postService.GetAllPosts(page, limit)
	if err != nil {
		logger.Error("Failed to retrieve posts", "error", err)
		return response.InternalErrorJSON(c, "Failed to retrieve posts")
	}

	// Create pagination metadata
	meta := response.NewPaginationMeta(posts.Page, posts.Limit, posts.TotalCount, posts.HasNext)

	return response.SuccessJSONWithMeta(c, posts.Posts, "Posts retrieved successfully", meta)
}

// CreatePost creates a new post
// @Summary Create a new post
// @Description Create a new post
// @Tags posts
// @Accept json
// @Produce json
// @Param request body dto.CreatePostRequest true "Post details"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /posts [post]
func (h *PostHandler) CreatePost(c *fiber.Ctx) error {
	var req dto.CreatePostRequest
	if err := c.BodyParser(&req); err != nil {
		return response.ValidationErrorJSON(c, "Invalid request body", err.Error())
	}

	post, err := h.postService.CreatePost(req)
	if err != nil {
		logger.Error("Failed to create post", "error", err)
		return response.ValidationErrorJSON(c, "Failed to create post", err.Error())
	}

	logger.Info("Post created successfully", "post_id", post.ID)
	return response.CreatedJSON(c, post, "Post created successfully")
}

// GetPost retrieves a specific post by ID
// @Summary Get a post by ID
// @Description Retrieve a specific post
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /posts/{id} [get]
func (h *PostHandler) GetPost(c *fiber.Ctx) error {
	id := c.Params("id")
	if strings.TrimSpace(id) == "" {
		return response.ValidationErrorJSON(c, "Invalid post ID", "ID cannot be empty")
	}

	post, err := h.postService.GetPostByID(id)
	if err != nil {
		return response.NotFoundJSON(c, "Post")
	}

	return response.SuccessJSON(c, post, "Post retrieved successfully")
}
