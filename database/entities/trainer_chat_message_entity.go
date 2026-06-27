package entities

import "time"

type TrainerChatMessage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SessionID uint      `gorm:"not null;index" json:"session_id"`
	Sender    string    `gorm:"type:varchar(20);not null" json:"sender"`
	Message   string    `gorm:"type:text;not null" json:"message"`
	CreatedAt time.Time `json:"created_at"`
}
