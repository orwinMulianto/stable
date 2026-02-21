package entities

import(
	"time"
)

type UserDailyActivityLog struct {
	ID                uint      `gorm:"primaryKey"`
	UserID            uint      `gorm:"not null"`
	DailyChallengeID  uint      `gorm:"not null"`
	CompletedAt       time.Time `gorm:"autoCreateTime"`
	TotalDurationSec  *int
	CreatedAt         time.Time

	User           User
	DailyChallenge DailyChallenge

	// unique composite index
	_ struct{} `gorm:"uniqueIndex:idx_user_challenge,priority:1"`
}