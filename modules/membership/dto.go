package membership

import "time"

type ClaimTrialRequest struct {
	UserID uint   `json:"user_id"`
	Plan   string `json:"plan" binding:"required"`
}

type MembershipResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Plan      string    `json:"plan"`
	Status    string    `json:"status"`
	IsTrial   bool      `json:"is_trial"`
	StartedAt time.Time `json:"started_at"`
	ExpiredAt time.Time `json:"expired_at"`
}