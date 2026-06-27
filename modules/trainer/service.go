package trainer

import (
	"errors"
	"strconv"
	"strings"

	"stable/database/entities"

	"gorm.io/gorm"
)

var ErrTrainerProfileExists = errors.New("trainer profile already exists")

type Service interface {
	GetProfile(userID uint) (*TrainerProfileResponse, error)
	GetTrainerByID(id uint) (*TrainerProfileResponse, error)
	GetAllTrainers() ([]TrainerProfileResponse, error)
	CreateProfile(userID uint, req CreateTrainerProfileRequest) (*TrainerProfileResponse, error)
	UpdateProfile(userID uint, req UpdateTrainerProfileRequest) (*TrainerProfileResponse, error)
	DeleteProfile(userID uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetTrainerByID(id uint) (*TrainerProfileResponse, error) {
	trainer, err := s.repo.GetByUserID(id)
	if err != nil {
		return nil, err
	}
	return toResponse(trainer), nil
}

func (s *service) GetProfile(userID uint) (*TrainerProfileResponse, error) {
	profile, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	return toResponse(profile), nil
}

func (s *service) GetAllTrainers() ([]TrainerProfileResponse, error) {
	profiles, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	responses := make([]TrainerProfileResponse, 0, len(profiles))
	for i := range profiles {
		responses = append(responses, *toResponse(&profiles[i]))
	}

	return responses, nil
}

func (s *service) CreateProfile(userID uint, req CreateTrainerProfileRequest) (*TrainerProfileResponse, error) {
	existing, err := s.repo.GetByUserID(userID)
	if err == nil && existing != nil {
		return nil, ErrTrainerProfileExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	
	profile := &entities.TrainerProfile{
		UserID:        &userID,
		Specialty:    strings.TrimSpace(req.Specialization),
		Experience:  strings.TrimSpace(req.Experience),
		Bio:          strings.TrimSpace(req.Bio),
	}
	

	created, err := s.repo.Create(profile)
	if err != nil {
		return nil, err
	}

	result, err := s.repo.GetByUserID(userID)
	if err != nil {
		return toResponse(created), nil
	}

	return toResponse(result), nil
}

func (s *service) UpdateProfile(userID uint, req UpdateTrainerProfileRequest) (*TrainerProfileResponse, error) {
	profile, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(req.Specialization) != "" {
		profile.Specialty = strings.TrimSpace(req.Specialization)
	}
	if strings.TrimSpace(req.Experience) != "" {
		profile.Experience = strings.TrimSpace(req.Experience)
	}
	if strings.TrimSpace(req.Bio) != "" {
		profile.Bio = strings.TrimSpace(req.Bio)
	}
	if req.IsOnline != nil {
	profile.IsOnline = *req.IsOnline
}

	updated, err := s.repo.Update(profile)
	if err != nil {
		return nil, err
	}

	result, err := s.repo.GetByUserID(userID)
	if err != nil {
		return toResponse(updated), nil
	}

	return toResponse(result), nil
}

func (s *service) DeleteProfile(userID uint) error {
	return s.repo.Delete(userID)
}

func toResponse(p *entities.TrainerProfile) *TrainerProfileResponse {
	if p == nil {
		return nil
	}

	rating, _ := strconv.ParseFloat(p.Rating, 64)

	var userID uint
	if p.UserID != nil {
		userID = *p.UserID
	}

	return &TrainerProfileResponse{
		ID:             p.ID,
		UserID:         userID,
		Username:       p.Name,    
		Specialization: p.Specialty,
		Experience:     p.Experience,
		Rating:         rating,
		TotalClients:   0,
		Bio:            p.Bio,
		CreatedAt:      p.CreatedAt,
	}
}