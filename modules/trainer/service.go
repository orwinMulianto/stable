package trainer

import (
	"errors"

	"stable/database/entities" // sesuaikan import path
)

type Service interface {
	GetProfile(userID uint) (*TrainerProfileResponse, error)
	GetAllTrainers() ([]TrainerProfileResponse, error)
	CreateProfile(userID uint, req CreateTrainerProfileRequest) (*TrainerProfileResponse, error)
	UpdateProfile(userID uint, req UpdateTrainerProfileRequest) (*TrainerProfileResponse, error)
	DeleteProfile(userID uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
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

	var responses []TrainerProfileResponse
	for _, p := range profiles {
		responses = append(responses, *toResponse(&p))
	}
	return responses, nil
}

func (s *service) CreateProfile(userID uint, req CreateTrainerProfileRequest) (*TrainerProfileResponse, error) {
	// Cek apakah sudah ada
	existing, _ := s.repo.GetByUserID(userID)
	if existing != nil {
		return nil, errors.New("trainer profile already exists")
	}

	profile := &entities.TrainerProfile{
		UserID:         userID,
		Specialization: req.Specialization,
		Experience:     req.Experience,
		Certification:  req.Certification,
		Bio:            req.Bio,
	}

	created, err := s.repo.Create(profile)
	if err != nil {
		return nil, err
	}

	// Reload dengan preload User
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

	if req.Specialization != "" {
		profile.Specialization = req.Specialization
	}
	if req.Experience != "" {
		profile.Experience = req.Experience
	}
	if req.Certification != "" {
		profile.Certification = req.Certification
	}
	if req.Bio != "" {
		profile.Bio = req.Bio
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

// ── HELPER ──
func toResponse(p *entities.TrainerProfile) *TrainerProfileResponse {
	res := &TrainerProfileResponse{
		ID:             p.ID,
		UserID:         p.UserID,
		Specialization: p.Specialization,
		Experience:     p.Experience,
		Certification:  p.Certification,
		Rating:         p.Rating,
		TotalClients:   p.TotalClients,
		Bio:            p.Bio,
		CreatedAt:      p.CreatedAt,
	}

	if p.User != nil {
		res.Username     = p.User.Username
		res.Email        = p.User.Email
		res.ProfileImage = p.User.ProfileImage
	}

	return res
}