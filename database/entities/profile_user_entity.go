package entities

import (
	"time"
)

type ProfileUser struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint `gorm:"uniqueIndex;not null"`


	BirthDate *time.Time
	HeightCm  *int
	WeightKg  *int
	AvatarURL string

	User *User `gorm:"constraint:OnDelete:CASCADE;"`
}