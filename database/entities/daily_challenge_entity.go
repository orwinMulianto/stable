package entities

import(
	"time"
)

type DailyChallenge struct {
	ID            uint      `gorm:"primaryKey"`
	ChallengeDate time.Time `gorm:"type:date;unique;not null"`
	Name          string    `gorm:"size:100;not null"`
	Description   string    `gorm:"type:text"`
	Difficulty    string    `gorm:"type:enum('EASY','MEDIUM','HARD');not null"`
	CreatedBy     *uint
	CreatedAt     time.Time

	Exercises []DailyChallengeExercise
}