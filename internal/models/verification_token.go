package models

import "time"

type VerificationToken struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64     `gorm:"not null;index" json:"userId"`
	Token     string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"token"`
	ExpiresAt time.Time  `json:"expiresAt"`
	CreatedAt time.Time  `json:"createdAt"`
	UsedAt    *time.Time `json:"usedAt,omitempty"`
}

func (VerificationToken) TableName() string {
	return "verification_tokens"
}
