package repository

import (
	"github.com/topboyasante/pitstop/internal/modules/user/domain"
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
func (r *UserRepository) GetByID(id string) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByProviderID retrieves a user by Provider and Provider ID
func (r *UserRepository) GetByProviderID(provider string, providerID string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("provider = ? AND provider_id = ?", provider, providerID).First(&user).Error
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

// GetAll retrieves all users with pagination
func (r *UserRepository) GetAll(page, limit int) ([]domain.User, int64, error) {
	var users []domain.User
	var totalCount int64

	offset := (page - 1) * limit

	// Get total count
	if err := r.db.Model(&domain.User{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Get users
	if err := r.db.Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, totalCount, nil
}

// Update updates a user
func (r *UserRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

// Delete soft deletes a user
func (r *UserRepository) Delete(id string) error {
	return r.db.Delete(&domain.User{}, id).Error
}
