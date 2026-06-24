package repository

import (
	"recruitment-api/internal/models"
	"time"

	"gorm.io/gorm"
)

type VerificationRepository interface {
	FindActiveByEmail(email string) (*models.VerificationCode, error)
	FindValidCode(email, code string) (*models.VerificationCode, error)
	Create(vc *models.VerificationCode) error
	MarkAsUsed(id uint64) error
}

type verificationRepository struct {
	db *gorm.DB
}

func NewVerificationRepository(db *gorm.DB) VerificationRepository {
	return &verificationRepository{db: db}
}

// FindActiveByEmail finds a code that is not expired and not used
func (r *verificationRepository) FindActiveByEmail(email string) (*models.VerificationCode, error) {
	var vc models.VerificationCode
	err := r.db.Where("email = ? AND is_used = false AND expires_at > ?", email, time.Now()).
		First(&vc).Error
	if err != nil {
		return nil, err
	}
	return &vc, nil
}

// FindValidCode finds a matching, unexpired, unused code
func (r *verificationRepository) FindValidCode(email, code string) (*models.VerificationCode, error) {
	var vc models.VerificationCode
	err := r.db.Where("email = ? AND code = ? AND is_used = false AND expires_at > ?", email, code, time.Now()).
		First(&vc).Error
	if err != nil {
		return nil, err
	}
	return &vc, nil
}

func (r *verificationRepository) Create(vc *models.VerificationCode) error {
	return r.db.Create(vc).Error
}

func (r *verificationRepository) MarkAsUsed(id uint64) error {
	return r.db.Model(&models.VerificationCode{}).Where("id = ?", id).Update("is_used", true).Error
}
