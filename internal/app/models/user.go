package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Email         string       `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash  string       `gorm:"type:varchar(255);not null" json:"-"`
	FirstName     string       `gorm:"type:varchar(100)" json:"first_name"`
	LastName      string       `gorm:"type:varchar(100)" json:"last_name"`
	AvatarURL     string       `gorm:"type:text" json:"avatar_url,omitempty"`
	EmailVerified bool         `gorm:"default:false" json:"email_verified"`
	LastLoginAt   *time.Time   `gorm:"type:timestamp" json:"last_login_at,omitempty"`
	IsActive    bool         `gorm:"not null;default:true" json:"is_active"`
	Memberships []Membership `gorm:"foreignKey:UserID" json:"memberships,omitempty"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}
