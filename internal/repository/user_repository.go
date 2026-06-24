package repository

import (
	"recruitment-api/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(email string) (*models.User, error)
	FindByID(id uint64) (*models.User, error)
	Create(user *models.User) error
	UpdateStatus(userID uint64, status int8) error
	GetAll() ([]models.User, error)
	GetStatusCounts() (map[int8]int64, error)
	GetTotalCount() (int64, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ? AND is_deleted = false", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id uint64) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ? AND is_deleted = false AND is_active = true", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) UpdateStatus(userID uint64, status int8) error {
	return r.db.Model(&models.User{}).
		Where("id = ? AND is_deleted = false", userID).
		Update("status", status).Error
}

func (r *userRepository) GetAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Where("is_deleted = false").Find(&users).Error
	return users, err
}

func (r *userRepository) GetStatusCounts() (map[int8]int64, error) {
	type StatusCount struct {
		Status int8
		Count  int64
	}

	var results []StatusCount
	err := r.db.Model(&models.User{}).
		Select("status, count(*) as count").
		Where("is_deleted = false").
		Group("status").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	counts := make(map[int8]int64)
	for _, r := range results {
		counts[r.Status] = r.Count
	}
	return counts, nil
}

func (r *userRepository) GetTotalCount() (int64, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("is_deleted = false").Count(&count).Error
	return count, err
}
