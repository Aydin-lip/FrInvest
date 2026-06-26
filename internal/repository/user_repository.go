package repository

import (
	"recruitment-api/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(email string) (*models.User, error)
	FindByID(id uint64) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	SetVerified(userID uint64) error
	UpdateStatus(userID uint64, status int8) error
	GetVerifiedUsers() ([]models.User, error)
	GetVerifiedStatusCounts() (map[int8]int64, error)
	GetVerifiedTotalCount() (int64, error)
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
	err := r.db.Where("id = ? AND is_deleted = false", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) SetVerified(userID uint64) error {
	result := r.db.Model(&models.User{}).
		Where("id = ? AND is_deleted = false", userID).
		UpdateColumn("verify", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *userRepository) UpdateStatus(userID uint64, status int8) error {
	return r.db.Model(&models.User{}).
		Where("id = ? AND is_deleted = false AND verify = true", userID).
		Update("status", status).Error
}

func (r *userRepository) GetVerifiedUsers() ([]models.User, error) {
	var users []models.User
	err := r.db.Where("verify = true AND is_deleted = false").Find(&users).Error
	return users, err
}

func (r *userRepository) GetVerifiedStatusCounts() (map[int8]int64, error) {
	type StatusCount struct {
		Status int8
		Count  int64
	}

	var results []StatusCount
	err := r.db.Model(&models.User{}).
		Select("status, count(*) as count").
		Where("verify = true AND is_deleted = false").
		Group("status").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	counts := make(map[int8]int64)
	for _, row := range results {
		counts[row.Status] = row.Count
	}
	return counts, nil
}

func (r *userRepository) GetVerifiedTotalCount() (int64, error) {
	var count int64
	err := r.db.Model(&models.User{}).
		Where("verify = true AND is_deleted = false").
		Count(&count).Error
	return count, err
}
