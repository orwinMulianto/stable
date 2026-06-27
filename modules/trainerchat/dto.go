package trainerchat

import "time"

type CheckoutRequest struct {
	UserID        uint   `json:"user_id" binding:"required"`
	TrainerID     uint   `json:"trainer_id" binding:"required"`
	CustomerName  string `json:"customer_name"`
	CustomerEmail string `json:"customer_email"`
}

type CheckoutResponse struct {
	SessionID   uint   `json:"session_id"`
	OrderID     string `json:"order_id"`
	SnapToken   string `json:"snap_token"`
	RedirectURL string `json:"redirect_url"`
	Amount      int    `json:"amount"`
}

type SendMessageRequest struct {
	Sender  string `json:"sender" binding:"required"`
	Message string `json:"message" binding:"required"`
}

type MessageResponse struct {
	ID        uint      `json:"id"`
	SessionID uint      `json:"session_id"`
	Sender    string    `json:"sender"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

type SessionResponse struct {
    ID               uint              `json:"id"`
    UserID           uint              `json:"user_id"`
    UserName         string            `json:"user_name"` // ← tambahkan
    TrainerID        uint              `json:"trainer_id"`
    TrainerName      string            `json:"trainer_name"`
    TrainerSpecialty string            `json:"trainer_specialty"`
    Amount           int               `json:"amount"`
    Status           string            `json:"status"`
    OrderID          string            `json:"order_id"`
    StartedAt        *time.Time        `json:"started_at"`
    ExpiresAt        *time.Time        `json:"expires_at"`
    CreatedAt        time.Time         `json:"created_at"`
    Messages         []MessageResponse `json:"messages"`
}

type TrainerDashboardStats struct {
	TotalRevenue int64 `json:"total_revenue"`
	TotalClients int64 `json:"total_clients"`
	PaidSessions int64 `json:"paid_sessions"`
}

type TrainerDashboardClient struct {
	UserID        uint      `json:"user_id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	TotalSessions int64     `json:"total_sessions"`
	TotalPaid     int64     `json:"total_paid"`
	LastSessionAt time.Time `json:"last_session_at"`
}

type TrainerDashboardResponse struct {
	Trainer       TrainerCatalogItem       `json:"trainer"`
	Stats         TrainerDashboardStats    `json:"stats"`
	RecentClients []TrainerDashboardClient `json:"recent_clients"`
}

type MidtransNotification struct {
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	TransactionID     string `json:"transaction_id"`
	StatusMessage     string `json:"status_message"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
	PaymentType       string `json:"payment_type"`
	OrderID           string `json:"order_id"`
	MerchantID        string `json:"merchant_id"`
	GrossAmount       string `json:"gross_amount"`
	FraudStatus       string `json:"fraud_status"`
}

type WSMessage struct {
    Type      string `json:"type"`      // "message" | "ping"
    Sender    string `json:"sender"`    // "user" | "trainer" (diisi server)
    Content   string `json:"content"`
    Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}

