package membership

import (
	"errors"
	"stable/database/entities"
	"time"

	"gorm.io/gorm"
)

type Service interface {
	ClaimTrial(userID uint, plan string) (*entities.UserMembership, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository}
}

func (s *service) ClaimTrial(userID uint, plan string) (*entities.UserMembership, error) {
	if plan != "1_month" && plan != "3_month" && plan != "6_month" {
		return nil, errors.New("invalid membership plan")
	}

	existing, err := s.repository.FindByUserID(userID)

	if err == nil {
		if time.Now().Before(existing.ExpiredAt) {
			return nil, errors.New("trial is still active")
		}

		return nil, errors.New("trial already used")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	now := time.Now()

	membership := &entities.UserMembership{
		UserID:    userID,
		Plan:      plan,
		Status:    "trial",
		IsTrial:   true,
		StartedAt: now,
		ExpiredAt: now.AddDate(0, 0, 7),
	}

	if err := s.repository.Create(membership); err != nil {
		return nil, err
	}

	return membership, nil
}