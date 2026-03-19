package users

import (
	"errors"
	"stable/database/entities"

	"gorm.io/gorm"
)

type userService struct {
	repo Repository
}

type Service interface {
	GetByID(id int) (*entities.User, error)
	// GetAll() ([]entities.User, error)
}

func NewService(repo Repository) Service {
	return &userService{repo: repo}
}

func (s *userService) GetByID(id int) (*entities.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	user.Password = "" 
	return &user, nil
}