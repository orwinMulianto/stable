package entities

import "time"

type UserStreak struct {
	ID                uint       `gorm:"primaryKey" json:"id"`
	UserID            uint       `gorm:"not null;uniqueIndex" json:"user_id"`
	CurrentStreak     int        `gorm:"not null;default:0" json:"current_streak"`
	LongestStreak     int        `gorm:"not null;default:0" json:"longest_streak"`
	LastCompletedDate *time.Time `gorm:"type:date" json:"last_completed_date"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}