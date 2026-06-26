package models

import "time"

type User struct {
	ID          uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName   string     `gorm:"type:varchar(100);not null" json:"firstName"`
	LastName    *string    `gorm:"type:varchar(100)" json:"lastName"`
	PhoneNumber *string    `gorm:"type:varchar(11)" json:"phoneNumber"`
	Email       string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Status      int8       `gorm:"type:tinyint;not null;default:0" json:"status"`
	Verify      bool       `gorm:"column:verify;default:false" json:"verify"`
	IsActive    bool       `gorm:"default:true" json:"isActive"`
	IsDeleted   bool       `gorm:"default:false" json:"isDeleted"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}

func (User) TableName() string {
	return "users"
}
