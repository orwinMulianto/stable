package users

import (
	"stable/database/entities"
	"log"
	"gorm.io/gorm"
	// "stable/database/migrations"
)

type repository struct {
	db *gorm.DB
}

// FindByVerificationCode implements [Repository].
func (r *repository) FindByVerificationCode(code string) (entities.User, error) {
	panic("unimplemented")
}


type Repository interface {
	// FindAll() ([]entities.User, error)
	FindByID(id int) (entities.User, error)
	FindByEmail(email string) (entities.User, error)
	// FindByVerificationCode(code string) (entities.User, error)
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

func (r *repository) Update(user entities.User) error {
    log.Println("[UPDATE] User ID:", user.ID)
    log.Println("[UPDATE] VerificationCode:", user.VerificationCode)
    
    result := r.db.Model(&entities.User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
        "username":                    user.Username,
        "email":                       user.Email,
        "password":                    user.Password,
        "no_telp":                     user.NoTelp,
        "jenis_kelamin":               user.JenisKelamin,
        "role":                        user.Role,
        "auth_provider":               user.AuthProvider,
        "profile_image":               user.ProfileImage,
        "is_verified":                 user.IsVerified,
        "verification_code":           user.VerificationCode,
        "reset_password_token":        user.ResetPasswordToken,
        "reset_password_token_expiry": user.ResetPasswordTokenExpiry,
    })
    
    log.Println("[UPDATE] Rows affected:", result.RowsAffected)
    log.Println("[UPDATE] Error:", result.Error)
    
    return result.Error
}