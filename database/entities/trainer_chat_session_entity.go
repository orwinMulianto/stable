package entities

import "time"

type TrainerChatSession struct {
	ID                    uint       `gorm:"primaryKey" json:"id"`
	UserID                uint       `gorm:"not null;index" json:"user_id"`
	TrainerID             uint       `gorm:"not null;index" json:"trainer_id"`
	TrainerName           string     `gorm:"type:varchar(120);not null" json:"trainer_name"`
	TrainerSpecialty      string     `gorm:"type:varchar(160);not null" json:"trainer_specialty"`
	DurationMinutes       int        `gorm:"not null;default:10" json:"duration_minutes"`
	Amount                int        `gorm:"not null" json:"amount"`
	Status                string     `gorm:"type:varchar(30);not null;default:'pending';index" json:"status"`
	OrderID               string     `gorm:"type:varchar(120);not null;uniqueIndex" json:"order_id"`
	SnapToken             string     `gorm:"type:varchar(255)" json:"snap_token"`
	RedirectURL           string     `gorm:"type:varchar(255)" json:"redirect_url"`
	MidtransTransactionID string     `gorm:"type:varchar(120)" json:"midtrans_transaction_id"`
	PaymentType           string     `gorm:"type:varchar(60)" json:"payment_type"`
	PaymentStatus         string     `gorm:"type:varchar(60)" json:"payment_status"`
	FraudStatus           string     `gorm:"type:varchar(60)" json:"fraud_status"`
	CustomerName string `gorm:"type:varchar(120)" json:"customer_name"`
	PaidAt                *time.Time `json:"paid_at"`
	StartedAt             *time.Time `json:"started_at"`
	ExpiresAt             *time.Time `json:"expires_at"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}