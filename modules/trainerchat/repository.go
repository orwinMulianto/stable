package trainerchat

import (
	"errors"
	"stable/database/entities"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	ListTrainers() ([]entities.TrainerProfile, error)
	FindTrainerByID(trainerID uint) (*entities.TrainerProfile, error)
	FindTrainerByUserID(userID uint) (*entities.TrainerProfile, error)
	CreateSession(session *entities.TrainerChatSession) error
	UpdateSnap(sessionID uint, snapToken string, redirectURL string) error
	FindSession(sessionID uint) (*entities.TrainerChatSession, error)
	FindSessionByOrderID(orderID string) (*entities.TrainerChatSession, error)
	UpdatePaymentStatus(orderID string, updates map[string]interface{}) error
	ListSessionsByUser(userID uint) ([]entities.TrainerChatSession, error)
	GetDashboardStatsByTrainerID(trainerID uint) (TrainerDashboardStats, error)
	ListRecentClientsByTrainerID(trainerID uint, limit int) ([]TrainerDashboardClient, error)
	CreateMessage(message *entities.TrainerChatMessage) error
	ListMessages(sessionID uint) ([]entities.TrainerChatMessage, error)
	ListSessionsByTrainerID(trainerID uint) ([]entities.TrainerChatSession, error)
	UpdatePaymentBySessionID(sessionID uint, updates map[string]interface{}) error
	
}



type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) ListTrainers() ([]entities.TrainerProfile, error) {
	var trainers []entities.TrainerProfile
	err := r.db.
		Where("is_active = ?", true).
		Order("name ASC").
		Find(&trainers).Error

	return trainers, err
}

func (r *repository) FindTrainerByID(trainerID uint) (*entities.TrainerProfile, error) {
	var trainer entities.TrainerProfile
	if err := r.db.
		Where("id = ? AND is_active = ?", trainerID, true).
		First(&trainer).Error; err != nil {
		return nil, err
	}

	return &trainer, nil
}

func (r *repository) ListSessionsByTrainerID(trainerID uint) ([]entities.TrainerChatSession, error) {
	var sessions []entities.TrainerChatSession
	err := r.db.
		Where("trainer_id = ? AND status = ?", trainerID, "paid").
		Order("created_at DESC").
		Find(&sessions).Error

	return sessions, err
}

func (r *repository) FindTrainerByUserID(userID uint) (*entities.TrainerProfile, error) {
	var trainer entities.TrainerProfile
	if err := r.db.
		Where("user_id = ? AND is_active = ?", userID, true).
		First(&trainer).Error; err != nil {
		return nil, err
	}

	return &trainer, nil
}

func (r *repository) CreateSession(session *entities.TrainerChatSession) error {
	return r.db.Create(session).Error
}

func (r *repository) UpdateSnap(sessionID uint, snapToken string, redirectURL string) error {
	return r.db.Model(&entities.TrainerChatSession{}).
		Where("id = ?", sessionID).
		Updates(map[string]interface{}{
			"snap_token":   snapToken,
			"redirect_url": redirectURL,
		}).Error
}

func (r *repository) FindSession(sessionID uint) (*entities.TrainerChatSession, error) {
	var session entities.TrainerChatSession
	if err := r.db.First(&session, sessionID).Error; err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *repository) FindSessionByOrderID(orderID string) (*entities.TrainerChatSession, error) {
	var session entities.TrainerChatSession
	if err := r.db.Where("order_id = ?", orderID).First(&session).Error; err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *repository) UpdatePaymentStatus(orderID string, updates map[string]interface{}) error {
	return r.db.Model(&entities.TrainerChatSession{}).
		Where("order_id = ?", orderID).
		Updates(updates).Error
}

func (r *repository) ListSessionsByUser(userID uint) ([]entities.TrainerChatSession, error) {
	var sessions []entities.TrainerChatSession
	err := r.db.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&sessions).Error

	return sessions, err
}

func (r *repository) GetDashboardStatsByTrainerID(trainerID uint) (TrainerDashboardStats, error) {
	var stats TrainerDashboardStats
	err := r.db.
		Table("trainer_chat_sessions").
		Select(`
			COALESCE(SUM(amount), 0) AS total_revenue,
			COUNT(DISTINCT user_id) AS total_clients,
			COUNT(*) AS paid_sessions
		`).
		Where("trainer_id = ? AND status = ?", trainerID, "paid").
		Scan(&stats).Error

	return stats, err
}

func (r *repository) ListRecentClientsByTrainerID(
	trainerID uint,
	limit int,
) ([]TrainerDashboardClient, error) {
	var clients []TrainerDashboardClient
	if limit <= 0 {
		limit = 5
	}

	err := r.db.
		Table("trainer_chat_sessions AS s").
		Select(`
			s.user_id,
			COALESCE(NULLIF(u.username, ''), NULLIF(u.email, ''), CONCAT('User #', s.user_id)) AS name,
			COALESCE(u.email, '') AS email,
			COUNT(*) AS total_sessions,
			COALESCE(SUM(s.amount), 0) AS total_paid,
			MAX(COALESCE(s.paid_at, s.created_at)) AS last_session_at
		`).
		Joins("LEFT JOIN users AS u ON u.id = s.user_id").
		Where("s.trainer_id = ? AND s.status = ?", trainerID, "paid").
		Group("s.user_id, u.username, u.email").
		Order("last_session_at DESC").
		Limit(limit).
		Scan(&clients).Error

	return clients, err
}

func (r *repository) CreateMessage(message *entities.TrainerChatMessage) error {
	return r.db.Create(message).Error
}

func (r *repository) ListMessages(sessionID uint) ([]entities.TrainerChatMessage, error) {
	var messages []entities.TrainerChatMessage
	err := r.db.
		Where("session_id = ?", sessionID).
		Order("created_at ASC").
		Find(&messages).Error

	return messages, err
}

func sessionCanChat(session *entities.TrainerChatSession, now time.Time) error {
	if session.Status != "paid" {
		return errors.New("chat session is not paid")
	}
	if session.ExpiresAt != nil && now.After(*session.ExpiresAt) {
		return errors.New("chat session has expired")
	}

	return nil
}

func (r *repository) UpdatePaymentBySessionID(sessionID uint, updates map[string]interface{}) error {
    return r.db.Model(&entities.TrainerChatSession{}).
        Where("id = ?", sessionID).
        Updates(updates).Error
}