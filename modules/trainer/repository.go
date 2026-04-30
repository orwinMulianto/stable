package trainer

import (
	"errors"

	"stable/database/entities"
	"gorm.io/gorm"
)

type Repository interface {
	GetByUserID(userID uint) (*entities.TrainerProfile, error)
	GetAll() ([]entities.TrainerProfile, error)
	Create(profile *entities.TrainerProfile) (*entities.TrainerProfile, error)
	Update(profile *entities.TrainerProfile) (*entities.TrainerProfile, error)
	Delete(userID uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) GetByUserID(userID uint) (*entities.TrainerProfile, error) {
	var profile entities.TrainerProfile
	err := r.db.Preload("User").Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("trainer profile not found")
		}
		return nil, err
	}
	return &profile, nil
}

func (r *repository) GetAll() ([]entities.TrainerProfile, error) {
	var profiles []entities.TrainerProfile
	err := r.db.Preload("User").Find(&profiles).Error
	if err != nil {
		return nil, err
	}
	return profiles, nil
}

func (r *repository) Create(profile *entities.TrainerProfile) (*entities.TrainerProfile, error) {
	err := r.db.Create(profile).Error
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func (r *repository) Update(profile *entities.TrainerProfile) (*entities.TrainerProfile, error) {
	err := r.db.Save(profile).Error
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func (r *repository) Delete(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&entities.TrainerProfile{}).Error
}