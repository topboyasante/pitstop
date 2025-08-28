package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/topboyasante/pitstop/internal/api/v1/services"
)

type PostController struct {
	postService *services.PostService
}

func NewPostController(postService *services.PostService) *PostController {
	return &PostController{
		postService: postService,
	}
}

// GetAllPosts retrieves all posts
// @Summary Get all posts
// @Description Retrieve a list of all posts
// @Tags posts
// @Accept json
// @Produce json
// @Success 200 {array} string "List of posts"
// @Router /posts [get]
func (pc *PostController) GetAllPosts(c *fiber.Ctx) error {
	posts, err := pc.postService.GetAllPosts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve posts",
		})
	}
	
	return c.JSON(fiber.Map{
		"data": posts,
	})
}