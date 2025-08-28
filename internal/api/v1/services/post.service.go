package services

import (
	"github.com/go-playground/validator/v10"
	"github.com/topboyasante/pitstop/internal/api/v1/repositories"
)

type PostService struct {
	postRepo  *repositories.PostRepository
	validator *validator.Validate
}

func NewPostService(postRepo *repositories.PostRepository, validator *validator.Validate) *PostService {
	return &PostService{
		postRepo:  postRepo,
		validator: validator,
	}
}

func (ps *PostService) GetAllPosts() ([]string, error) {
	return ps.postRepo.GetAll()
}