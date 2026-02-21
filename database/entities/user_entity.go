package entities

import (
	"time"
	"gorm.io/gorm"
	
)

type User struct {
	gorm.Model

	Email        string `gorm:"uniqueIndex;size:191;not null"`
	PasswordHash string `gorm:"size:255;not null"`

	Role string `gorm:"type:enum('USER','ADMIN');default:'USER';not null"`

	AuthProvider string `gorm:"type:enum('PASSWORD','GOOGLE');default:'PASSWORD';not null"`

	ProfileImage *string
	IsVerified   bool `gorm:"default:false"`

	ResetPasswordToken       *string `gorm:"index"`
	ResetPasswordTokenExpiry *time.Time

	Profile    *ProfileUser
	UserStreak *UserStreak
}