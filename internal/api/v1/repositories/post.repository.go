package repositories

import "gorm.io/gorm"

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

func (pr *PostRepository) GetAll() ([]string, error) {
	// TODO: Implement actual post retrieval from database
	return []string{"Post 1", "Post 2", "Post 3"}, nil
}