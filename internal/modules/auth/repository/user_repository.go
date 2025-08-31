package repository

import (
	"github.com/topboyasante/pitstop/internal/modules/auth/domain"
	"gorm.io/gorm"
)

// UserRepository handles user data operations
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *UserRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

// Delete soft deletes a user
func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&domain.User{}, id).Error
}

// UpdateFollowerCount updates the follower count for a user
func (r *UserRepository) UpdateFollowerCount(userID uint, increment bool) error {
	if increment {
		return r.db.Model(&domain.User{}).Where("id = ?", userID).
			UpdateColumn("follower_count", gorm.Expr("follower_count + ?", 1)).Error
	}
	return r.db.Model(&domain.User{}).Where("id = ?", userID).
		UpdateColumn("follower_count", gorm.Expr("follower_count - ?", 1)).Error
}

// UpdateFollowingCount updates the following count for a user
func (r *UserRepository) UpdateFollowingCount(userID uint, increment bool) error {
	if increment {
		return r.db.Model(&domain.User{}).Where("id = ?", userID).
			UpdateColumn("following_count", gorm.Expr("following_count + ?", 1)).Error
	}
	return r.db.Model(&domain.User{}).Where("id = ?", userID).
		UpdateColumn("following_count", gorm.Expr("following_count - ?", 1)).Error
}