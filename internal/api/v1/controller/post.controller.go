package controller

import "github.com/gofiber/fiber/v2"

type PostController struct{}

func NewPostController() *PostController {
	return &PostController{}
}

// GetAllPosts retrieves all posts
// @Summary Get all posts
// @Description Retrieve a list of all posts
// @Tags posts
// @Accept json
// @Produce json
// @Success 200 {string} string "Get all posts"
// @Router /posts [get]
func (pc *PostController) GetAllPosts(c *fiber.Ctx) error {
	return c.SendString("Get all posts")
}