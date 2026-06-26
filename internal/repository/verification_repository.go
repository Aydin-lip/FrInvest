package repository

import (
	"recruitment-api/internal/models"
	"time"

	"gorm.io/gorm"
)

type VerificationRepository interface {
	Create(vt *models.VerificationToken) error
	FindByToken(token string) (*models.VerificationToken, error)
	FindLatestValidByUserID(userID uint64) (*models.VerificationToken, error)
	FindLatestUnusedByUserID(userID uint64) (*models.VerificationToken, error)
	UpdateToken(id uint64, token string, expiresAt time.Time) error
	MarkAsUsed(id uint64, usedAt time.Time) error
}

type verificationRepository struct {
	db *gorm.DB
}

func NewVerificationRepository(db *gorm.DB) VerificationRepository {
	return &verificationRepository{db: db}
}

func (r *verificationRepository) Create(vt *models.VerificationToken) error {
	return r.db.Create(vt).Error
}

func (r *verificationRepository) FindByToken(token string) (*models.VerificationToken, error) {
	var vt models.VerificationToken
	err := r.db.Where("token = ?", token).First(&vt).Error
	if err != nil {
		return nil, err
	}
	return &vt, nil
}

func (r *verificationRepository) FindLatestValidByUserID(userID uint64) (*models.VerificationToken, error) {
	var vt models.VerificationToken
	err := r.db.Where("user_id = ? AND used_at IS NULL AND expires_at > ?", userID, time.Now()).
		Order("created_at DESC").
		First(&vt).Error
	if err != nil {
		return nil, err
	}
	return &vt, nil
}

func (r *verificationRepository) FindLatestUnusedByUserID(userID uint64) (*models.VerificationToken, error) {
	var vt models.VerificationToken
	err := r.db.Where("user_id = ? AND used_at IS NULL", userID).
		Order("created_at DESC").
		First(&vt).Error
	if err != nil {
		return nil, err
	}
	return &vt, nil
}

func (r *verificationRepository) UpdateToken(id uint64, token string, expiresAt time.Time) error {
	return r.db.Model(&models.VerificationToken{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"token":      token,
			"expires_at": expiresAt,
		}).Error
}

func (r *verificationRepository) MarkAsUsed(id uint64, usedAt time.Time) error {
	return r.db.Model(&models.VerificationToken{}).
		Where("id = ?", id).
		Update("used_at", usedAt).Error
}
