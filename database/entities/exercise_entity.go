package entities

import (
	"time"
)

type Exercise struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"size:100;not null"`
	Difficulty  string    `gorm:"type:enum('EASY','MEDIUM','HARD');not null"`
	Description string    `gorm:"type:text"`
	CreatedBy   *uint
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Contents []ExerciseContent
}