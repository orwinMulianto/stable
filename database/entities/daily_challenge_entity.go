package entities

import "time"

type DailyChallenge struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"not null;uniqueIndex:idx_user_challenge_date" json:"user_id"`
	ChallengeDate time.Time `gorm:"type:date;not null;uniqueIndex:idx_user_challenge_date" json:"challenge_date"`

	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	Difficulty  string `gorm:"type:varchar(50);default:'Daily'" json:"difficulty"`
	CreatedBy   uint   `gorm:"default:0" json:"created_by"`

	ChallengeKey string    `gorm:"type:varchar(50);not null" json:"challenge_key"`
	Repetitions  int       `gorm:"not null" json:"repetitions"`
	CompletedAt  time.Time `gorm:"not null" json:"completed_at"`
	CreatedAt    time.Time `json:"created_at"`
}