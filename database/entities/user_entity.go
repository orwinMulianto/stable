package entities

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string `json:"username" form:"username" binding:"required"`
	Email        string `gorm:"uniqueIndex;size:191;not null"`
	Password string `gorm:"size:255;not null"`
	NoTelp       string `json:"noTelp" form:"noTelp"`
	JenisKelamin *bool  `json:"jenisKelamin" form:"jenisKelamin"`
	Role string `gorm:"type:enum('USER','ADMIN');default:'USER';not null"`

	AuthProvider string `gorm:"type:enum('PASSWORD','GOOGLE');default:'PASSWORD';not null"`

	ProfileImage string     `json:"profileImage" form:"profileImage"`
	IsVerified   bool `gorm:"default:false"`
	VerificationCode         *string    `json:"verificationCode,omitempty" gorm:"size:100"`
	ResetPasswordToken       *string `gorm:"index"`
	ResetPasswordTokenExpiry *time.Time

	Profile    *ProfileUser
	UserStreak *UserStreak
}
