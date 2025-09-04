package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/topboyasante/pitstop/internal/modules/post/dto"
	"github.com/topboyasante/pitstop/internal/modules/post/service"
	"github.com/topboyasante/pitstop/internal/core/logger"
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
// @Success 200 {object} dto.PostsResponse
// @Router /posts [get]
func (h *PostHandler) GetAllPosts(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	posts, err := h.postService.GetAllPosts(page, limit)
	if err != nil {
		logger.Error("Failed to retrieve posts", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve posts",
		})
	}

	return c.JSON(posts)
}

// CreatePost creates a new post
// @Summary Create a new post
// @Description Create a new post
// @Tags posts
// @Accept json
// @Produce json
// @Param request body dto.CreatePostRequest true "Post details"
// @Success 201 {object} dto.PostResponse
// @Failure 400 {object} map[string]string
// @Router /posts [post]
func (h *PostHandler) CreatePost(c *fiber.Ctx) error {
	var req dto.CreatePostRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	post, err := h.postService.CreatePost(req)
	if err != nil {
		logger.Error("Failed to create post", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	logger.Info("Post created successfully", "post_id", post.ID)
	return c.Status(fiber.StatusCreated).JSON(post)
}

// GetPost retrieves a specific post by ID
// @Summary Get a post by ID
// @Description Retrieve a specific post
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} dto.PostResponse
// @Failure 404 {object} map[string]string
// @Router /posts/{id} [get]
func (h *PostHandler) GetPost(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid post ID",
		})
	}

	post, err := h.postService.GetPostByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Post not found",
		})
	}

	return c.JSON(post)
}
