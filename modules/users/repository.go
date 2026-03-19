package users

import (
	"gorm.io/gorm"
	"stable/database/entities"
	// "stable/database/migrations"
)

type repository struct {
	db *gorm.DB
}

// FindByVerificationCode implements [Repository].
func (r *repository) FindByVerificationCode(code string) (entities.User, error) {
	panic("unimplemented")
}

// Update implements [Repository].
func (r *repository) Update(user entities.User) error {
	panic("unimplemented")
}

type Repository interface {
	// FindAll() ([]entities.User, error)
	FindByID(id int) (entities.User, error)
	FindByEmail(email string) (entities.User, error)
	FindByVerificationCode(code string) (entities.User, error)
	// FindByResetToken(token string) (entities.User, error)
	Create(user *entities.User) error
	Update(user entities.User) error
	// Delete(id int) error
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindByID(id int) (entities.User, error) {
	var user entities.User
	err := r.db.First(&user, id).Error
	return user, err
}

func (r *repository) FindByEmail(email string) (entities.User, error) {
	var user entities.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return user, err
}

func (r *repository) Create(user *entities.User) error {
	return r.db.Create(user).Error
}
