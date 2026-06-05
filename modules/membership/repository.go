package membership

import (
	"stable/database/entities"

	"gorm.io/gorm"
)

type Repository interface {
	FindByUserID(userID uint) (*entities.UserMembership, error)
	Create(membership *entities.UserMembership) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) FindByUserID(userID uint) (*entities.UserMembership, error) {
	var membership entities.UserMembership

	err := r.db.Where("user_id = ?", userID).First(&membership).Error
	if err != nil {
		return nil, err
	}

	return &membership, nil
}

func (r *repository) Create(membership *entities.UserMembership) error {
	return r.db.Create(membership).Error
}