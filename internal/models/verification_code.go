package models

import "time"

type VerificationCode struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Email     string    `gorm:"type:varchar(255);not null" json:"email"`
	Code      string    `gorm:"type:varchar(6);not null" json:"code"`
	ExpiresAt time.Time `json:"expiresAt"`
	IsUsed    bool      `gorm:"default:false" json:"isUsed"`
	CreatedAt time.Time `json:"createdAt"`
}

func (VerificationCode) TableName() string {
	return "verification_codes"
}
