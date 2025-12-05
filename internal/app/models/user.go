package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	gorm.Model
	Name          string       `gorm:"size:100;not null" json:"name"`
	Email         string       `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash  string       `gorm:"type:varchar(255);not null" json:"-"`
	Age           int          `gorm:"default:0" json:"age"`
	AvatarURL     string       `gorm:"type:text" json:"avatar_url,omitempty"`
	EmailVerified bool         `gorm:"default:false" json:"email_verified"`
	LastLoginAt   *time.Time   `gorm:"type:timestamp" json:"last_login_at,omitempty"`
	IsActive      bool         `gorm:"not null;default:true" json:"is_active"`
	Memberships   []Membership `gorm:"foreignKey:UserID" json:"memberships,omitempty"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}
