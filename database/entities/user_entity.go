package entities

import (
	"time"
	"gorm.io/gorm"
	
)

type User struct {
	gorm.Model
	Username     string `json:"username" form:"username" binding:"required"`
	Email        string `gorm:"uniqueIndex;size:191;not null"`
	PasswordHash string `gorm:"size:255;not null"`
	NoTelp       string `json:"noTelp" form:"noTelp"`
	JenisKelamin *bool  `json:"jenisKelamin" form:"jenisKelamin"`
	Role string `gorm:"type:enum('USER','ADMIN');default:'USER';not null"`

	AuthProvider string `gorm:"type:enum('PASSWORD','GOOGLE');default:'PASSWORD';not null"`

	ProfileImage string     `json:"profileImage" form:"profileImage"`
	IsVerified   bool `gorm:"default:false"`

	ResetPasswordToken       *string `gorm:"index"`
	ResetPasswordTokenExpiry *time.Time

	Profile    *ProfileUser
	UserStreak *UserStreak
}

type RegisterRequest struct {
    Username string `json:"username" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}