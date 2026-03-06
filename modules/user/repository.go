package user

import (
	"stable/database/migrations"
	"stable/database/entities"
)

func FindByID(id int) (entities.User, error) {
	var user entities.User
	err := migrations.DB.First(&user, id).Error
	return user, err
}