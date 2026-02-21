package entities

import(
	"time"
)

type UserStreak struct {
	ID                uint      `gorm:"primaryKey"`
	UserID            uint      `gorm:"uniqueIndex;not null"`
	CurrentStreak     int       `gorm:"default:0"`
	LongestStreak     int       `gorm:"default:0"`
	LastCompletedDate *time.Time
	UpdatedAt         time.Time

	User User `gorm:"constraint:OnDelete:CASCADE;"`
}