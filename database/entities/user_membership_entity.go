package entities

import "time"

type UserMembership struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;uniqueIndex" json:"user_id"`
	Plan      string    `gorm:"type:varchar(20);not null" json:"plan"`
	Status    string    `gorm:"type:varchar(20);default:'trial'" json:"status"`
	IsTrial   bool      `gorm:"default:true" json:"is_trial"`
	StartedAt time.Time `json:"started_at"`
	ExpiredAt time.Time `json:"expired_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}