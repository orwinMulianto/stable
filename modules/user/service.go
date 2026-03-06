package user

import(
	"stable/database/entities"
	"errors"
	"gorm.io/gorm"
)

func GetByID(id int) (*entities.User, error) {
	user, err := FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	user.PasswordHash = "" 
	return &user, nil
}