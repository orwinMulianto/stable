package entities

type DailyChallengeExercise struct {
	ID               uint `gorm:"primaryKey"`
	DailyChallengeID uint `gorm:"not null"`
	ExerciseID       uint `gorm:"not null"`
	Sets             int  `gorm:"not null"`
	Reps             int  `gorm:"not null"`
	RestSeconds      int  `gorm:"not null"`

	DailyChallenge DailyChallenge `gorm:"constraint:OnDelete:CASCADE;"`
	Exercise       Exercise       `gorm:"constraint:OnDelete:CASCADE;"`
}