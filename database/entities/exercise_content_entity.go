package entities

import (
	"time"
)

type ExerciseContent struct {
	ID          uint      `gorm:"primaryKey"`
	ExerciseID  uint      `gorm:"not null"`
	ContentType string    `gorm:"type:enum('VIDEO','TEXT','IMAGE');not null"`
	ContentURL  string    `gorm:"size:255"`
	ContentText string    `gorm:"type:text"`
	CreatedAt   time.Time

	Exercise Exercise `gorm:"constraint:OnDelete:CASCADE;"`
}