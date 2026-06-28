package trainerchat

import (
	"errors"
	"fmt"
	"stable/database/entities"
	"strings"
	"time"
)

var (
	ErrTrainerNotFound     = errors.New("trainer not found")
	ErrInvalidNotification = errors.New("invalid midtrans notification signature")
	ErrSessionNotAvailable = errors.New("chat session is not available")
)

type Service interface {
	ListTrainers() ([]TrainerCatalogItem, error)
	Checkout(request CheckoutRequest) (*CheckoutResponse, error)
	HandleNotification(notification MidtransNotification) error
	GetHistory(userID uint) ([]SessionResponse, error)
	GetSession(sessionID uint) (*SessionResponse, error)
	SendMessage(sessionID uint, request SendMessageRequest) (*MessageResponse, error)
	GetMyDashboard(userID uint) (*TrainerDashboardResponse, error)
	GetTrainerSessions(userID uint) ([]SessionResponse, error)
	UpdatePaymentDirect(sessionID uint, updates map[string]interface{}) error
	DevMarkPaid(sessionID uint) (*SessionResponse, error)
}

type service struct {
	repository Repository
	midtrans   *midtransClient
	now        func() time.Time
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
		midtrans:   newMidtransClient(),
		now:        time.Now,
	}
}

func (s *service) ListTrainers() ([]TrainerCatalogItem, error) {
	trainers, err := s.repository.ListTrainers()
	if err != nil {
		return nil, err
	}

	items := make([]TrainerCatalogItem, 0, len(trainers))
	for _, trainer := range trainers {
		items = append(items, trainerEntityToCatalogItem(trainer))
	}

	return items, nil
}

func (s *service) Checkout(request CheckoutRequest) (*CheckoutResponse, error) {
	trainerEntity, err := s.repository.FindTrainerByID(request.TrainerID)
	if err != nil {
		return nil, ErrTrainerNotFound
	}
	trainer := trainerEntityToCatalogItem(*trainerEntity)

	orderID := fmt.Sprintf("STABLE-CHAT-%d-%d", request.UserID, s.now().UnixNano())
	session := &entities.TrainerChatSession{
		UserID:           request.UserID,
		TrainerID:        trainer.ID,
		TrainerName:      trainer.Name,
		TrainerSpecialty: trainer.Specialty,
		DurationMinutes:  10,
		CustomerName: request.CustomerName,
		Amount:           trainer.Price,
		Status:           "pending",
		OrderID:          orderID,
	}

	if err := s.repository.CreateSession(session); err != nil {
		return nil, err
	}

	snap, err := s.midtrans.CreateSnapTransaction(session.ID, orderID, trainer, request)
	if err != nil {
		return nil, err
	}

	if err := s.repository.UpdateSnap(session.ID, snap.Token, snap.RedirectURL); err != nil {
		return nil, err
	}

	return &CheckoutResponse{
		SessionID:   session.ID,
		OrderID:     orderID,
		SnapToken:   snap.Token,
		RedirectURL: snap.RedirectURL,
		Amount:      trainer.Price,
	}, nil
}

func (s *service) GetTrainerSessions(userID uint) ([]SessionResponse, error) {
	trainer, err := s.repository.FindTrainerByUserID(userID)
	if err != nil {
		return nil, ErrTrainerNotFound
	}

	sessions, err := s.repository.ListSessionsByTrainerID(trainer.ID)
	if err != nil {
		return nil, err
	}

	responses := make([]SessionResponse, 0, len(sessions))
	for _, session := range sessions {
		messages, _ := s.repository.ListMessages(session.ID)
		responses = append(responses, sessionToResponse(session, messages))
	}

	return responses, nil
}

func (s *service) HandleNotification(notification MidtransNotification) error {
	if !s.midtrans.VerifyNotificationSignature(notification) {
		return ErrInvalidNotification
	}

	updates := map[string]interface{}{
		"midtrans_transaction_id": notification.TransactionID,
		"payment_type":            notification.PaymentType,
		"payment_status":          notification.TransactionStatus,
		"fraud_status":            notification.FraudStatus,
	}

	now := s.now()
	switch notification.TransactionStatus {
	case "settlement":
		startedAt := now
		expiresAt := startedAt.Add(10 * time.Minute)
		updates["status"] = "paid"
		updates["paid_at"] = now
		updates["started_at"] = startedAt
		updates["expires_at"] = expiresAt
	case "capture":
		if notification.FraudStatus == "accept" {
			startedAt := now
			expiresAt := startedAt.Add(10 * time.Minute)
			updates["status"] = "paid"
			updates["paid_at"] = now
			updates["started_at"] = startedAt
			updates["expires_at"] = expiresAt
		}
	case "pending":
		updates["status"] = "pending"
	case "expire":
		updates["status"] = "expired"
	case "cancel", "deny", "failure":
		updates["status"] = "canceled"
	}

	return s.repository.UpdatePaymentStatus(notification.OrderID, updates)
}

func (s *service) GetHistory(userID uint) ([]SessionResponse, error) {
	sessions, err := s.repository.ListSessionsByUser(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]SessionResponse, 0, len(sessions))
	for _, session := range sessions {
		messages, _ := s.repository.ListMessages(session.ID)
		responses = append(responses, sessionToResponse(session, messages))
	}

	return responses, nil
}

func (s *service) GetSession(sessionID uint) (*SessionResponse, error) {
	session, err := s.repository.FindSession(sessionID)
	if err != nil {
		return nil, err
	}

	messages, err := s.repository.ListMessages(sessionID)
	if err != nil {
		return nil, err
	}

	response := sessionToResponse(*session, messages)
	return &response, nil
}

func (s *service) SendMessage(sessionID uint, request SendMessageRequest) (*MessageResponse, error) {
	session, err := s.repository.FindSession(sessionID)
	if err != nil {
		return nil, err
	}
	if err := sessionCanChat(session, s.now()); err != nil {
		return nil, ErrSessionNotAvailable
	}

	sender := strings.ToLower(strings.TrimSpace(request.Sender))
	if sender != "user" && sender != "trainer" {
		sender = "user"
	}

	message := &entities.TrainerChatMessage{
		SessionID: sessionID,
		Sender:    sender,
		Message:   strings.TrimSpace(request.Message),
	}
	if err := s.repository.CreateMessage(message); err != nil {
		return nil, err
	}

	response := messageToResponse(*message)
	return &response, nil
}

func (s *service) GetMyDashboard(userID uint) (*TrainerDashboardResponse, error) {
	trainerEntity, err := s.repository.FindTrainerByUserID(userID)
	if err != nil {
		return nil, ErrTrainerNotFound
	}
	trainer := trainerEntityToCatalogItem(*trainerEntity)

	stats, err := s.repository.GetDashboardStatsByTrainerID(trainer.ID)
	if err != nil {
		return nil, err
	}

	clients, err := s.repository.ListRecentClientsByTrainerID(trainer.ID, 5)
	if err != nil {
		return nil, err
	}

	return &TrainerDashboardResponse{
		Trainer:       trainer,
		Stats:         stats,
		RecentClients: clients,
	}, nil
}

func (s *service) DevMarkPaid(sessionID uint) (*SessionResponse, error) {
	session, err := s.repository.FindSession(sessionID)
	if err != nil {
		return nil, err
	}

	now := s.now()
	expiresAt := now.Add(10 * time.Minute)

	err = s.repository.UpdatePaymentBySessionID(sessionID, map[string]interface{}{
		"status":     "paid",
		"paid_at":    now,
		"started_at": now,
		"expires_at": expiresAt,
	})
	if err != nil {
		return nil, err
	}

	session.Status = "paid"
	session.StartedAt = &now
	session.ExpiresAt = &expiresAt

	messages, _ := s.repository.ListMessages(sessionID)
	response := sessionToResponse(*session, messages)
	return &response, nil
}

func sessionToResponse(
	session entities.TrainerChatSession,
	messages []entities.TrainerChatMessage,
) SessionResponse {
	messageResponses := make([]MessageResponse, 0, len(messages))
	for _, message := range messages {
		messageResponses = append(messageResponses, messageToResponse(message))
	}

	return SessionResponse{
		ID:               session.ID,
		UserID:           session.UserID,
		TrainerID:        session.TrainerID,
		TrainerName:      session.TrainerName,
		TrainerSpecialty: session.TrainerSpecialty,
		Amount:           session.Amount,
		UserName: 		  session.CustomerName, 
		Status:           session.Status,
		OrderID:          session.OrderID,
		StartedAt:        session.StartedAt,
		ExpiresAt:        session.ExpiresAt,
		CreatedAt:        session.CreatedAt,
		Messages:         messageResponses,
	}
}

func messageToResponse(message entities.TrainerChatMessage) MessageResponse {
	return MessageResponse{
		ID:        message.ID,
		SessionID: message.SessionID,
		Sender:    message.Sender,
		Message:   message.Message,
		CreatedAt: message.CreatedAt,
	}
}

// Implementasi
func (s *service) UpdatePaymentDirect(sessionID uint, updates map[string]interface{}) error {
    return s.repository.UpdatePaymentBySessionID(sessionID, updates)
}
