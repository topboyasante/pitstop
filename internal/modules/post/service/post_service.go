package service

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/topboyasante/pitstop/internal/core/logger"
	"github.com/topboyasante/pitstop/internal/modules/post/domain"
	"github.com/topboyasante/pitstop/internal/modules/post/dto"
	"github.com/topboyasante/pitstop/internal/modules/post/repository"
	"github.com/topboyasante/pitstop/internal/shared/events"
)

// PostService handles post business logic
type PostService struct {
	postRepo  *repository.PostRepository
	validator *validator.Validate
	eventBus  *events.EventBus
}

// NewPostService creates a new post service instance
func NewPostService(postRepo *repository.PostRepository, validator *validator.Validate, eventBus *events.EventBus) *PostService {
	return &PostService{
		postRepo:  postRepo,
		validator: validator,
		eventBus:  eventBus,
	}
}

// CreatePost creates a new post
func (s *PostService) CreatePost(req dto.CreatePostRequest) (*dto.PostResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	post := &domain.Post{
		UserID:  req.UserID,
		Content: req.Content,
	}

	if err := s.postRepo.Create(post); err != nil {
		logger.Error("Failed to create post", "error", err)
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	logger.Info("Post created successfully", "post_id", post.ID)

	return &dto.PostResponse{
		ID:        post.ID,
		UserID:    post.UserID,
		Content:   post.Content,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}, nil
}

// GetPostByID retrieves a post by ID
func (s *PostService) GetPostByID(id string) (*dto.PostResponse, error) {
	post, err := s.postRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("post not found: %w", err)
	}

	response := &dto.PostResponse{
		ID:           post.ID,
		UserID:       post.UserID,
		Content:      post.Content,
		CommentCount: post.CommentCount,
		LikeCount:    post.LikeCount,
		CreatedAt:    post.CreatedAt,
		UpdatedAt:    post.UpdatedAt,
	}

	if post.User != nil {
		response.User = &dto.PostUserResponse{
			Username:    post.User.Username,
			DisplayName: post.User.DisplayName,
			AvatarURL:   post.User.AvatarURL,
		}
	}

	return response, nil
}

// GetAllPosts retrieves all posts with pagination
func (s *PostService) GetAllPosts(page, limit int) (*dto.PostsResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	posts, totalCount, err := s.postRepo.GetAll(page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve posts: %w", err)
	}

	postResponses := make([]dto.PostResponse, len(posts))
	for i, post := range posts {
		postResponses[i] = dto.PostResponse{
			ID:           post.ID,
			UserID:       post.UserID,
			Content:      post.Content,
			CommentCount: post.CommentCount,
			LikeCount:    post.LikeCount,
			CreatedAt:    post.CreatedAt,
			UpdatedAt:    post.UpdatedAt,
		}

		if post.User != nil {
			postResponses[i].User = &dto.PostUserResponse{
				Username:    post.User.Username,
				DisplayName: post.User.DisplayName,
				AvatarURL:   post.User.AvatarURL,
			}
		}
	}

	hasNext := int64((page-1)*limit+len(posts)) < totalCount

	return &dto.PostsResponse{
		Posts:      postResponses,
		TotalCount: totalCount,
		Page:       page,
		Limit:      limit,
		HasNext:    hasNext,
	}, nil
}
