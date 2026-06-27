package profile

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

const usernameChangeCooldown = 14 * 24 * time.Hour

type UsernameCooldownError struct {
	CanChangeAt time.Time
}

func (err UsernameCooldownError) Error() string {
	return fmt.Sprintf("username can be changed again at %s", err.CanChangeAt.Format(time.RFC3339))
}

type Repository interface {
	FindByUserID(userID uint) (*ProfileResponse, error)
	Update(userID uint, request UpdateProfileRequest) (*ProfileResponse, error)
	UpdateProfileImage(userID uint, imageURL string) (*ProfileResponse, error)
}

type repository struct {
	db *gorm.DB
}

type profileRow struct {
	ID                uint       `gorm:"column:id"`
	Username          string     `gorm:"column:username"`
	Email             string     `gorm:"column:email"`
	Role              string     `gorm:"column:role"`
	UsernameChangedAt *time.Time `gorm:"column:username_changed_at"`
	Gender            *bool      `gorm:"column:gender"`
	BirthDate         *time.Time `gorm:"column:birth_date"`
	ProfileImage      *string    `gorm:"column:profile_image"`
	Weight            *int       `gorm:"column:weight"`
	Height            *int       `gorm:"column:height"`
	CurrentStreak     *int       `gorm:"column:current_streak"`
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindByUserID(userID uint) (*ProfileResponse, error) {
	var row profileRow

	result := r.db.
		Table("users AS u").
		Select(`
			u.id,
			u.username,
			u.email,
			u.role,
			u.username_changed_at,
			u.jenis_kelamin AS gender,
			p.birth_date,
			p.avatar_url AS profile_image,
			p.weight_kg AS weight,
			p.height_cm AS height,
			COALESCE(s.current_streak, 0) AS current_streak
		`).
		Joins("LEFT JOIN profile_users AS p ON p.user_id = u.id").
		Joins("LEFT JOIN user_streaks AS s ON s.user_id = u.id").
		Where("u.id = ?", userID).
		Scan(&row)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return row.toResponse(), nil
}

func (r *repository) Update(
	userID uint,
	request UpdateProfileRequest,
) (*ProfileResponse, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		userUpdates := map[string]interface{}{}

		if request.Username != nil {
			username := strings.TrimSpace(*request.Username)
			if username != "" {
				var current struct {
					Username          string     `gorm:"column:username"`
					UsernameChangedAt *time.Time `gorm:"column:username_changed_at"`
				}

				if err := tx.Table("users").
					Select("username, username_changed_at").
					Where("id = ?", userID).
					First(&current).Error; err != nil {
					return err
				}

				if username != current.Username {
					if current.UsernameChangedAt != nil {
						canChangeAt := current.UsernameChangedAt.Add(usernameChangeCooldown)
						if time.Now().Before(canChangeAt) {
							return UsernameCooldownError{CanChangeAt: canChangeAt}
						}
					}

					userUpdates["username"] = username
					userUpdates["username_changed_at"] = time.Now()
				}
			}
		}

		if request.Gender != nil {
			userUpdates["jenis_kelamin"] = *request.Gender
		}

		if len(userUpdates) > 0 {
			if err := tx.Table("users").
				Where("id = ?", userID).
				Updates(userUpdates).Error; err != nil {
				return err
			}
		}

		profileUpdates := map[string]interface{}{}
		if request.ProfileImage != nil {
			profileUpdates["avatar_url"] = strings.TrimSpace(*request.ProfileImage)
		}
		if request.Weight != nil {
			profileUpdates["weight_kg"] = *request.Weight
		}
		if request.Height != nil {
			profileUpdates["height_cm"] = *request.Height
		}

		if len(profileUpdates) == 0 {
			return nil
		}

		var existing struct {
			ID uint `gorm:"column:id"`
		}
		err := tx.Table("profile_users").
			Select("id").
			Where("user_id = ?", userID).
			First(&existing).Error

		if err == nil {
			return tx.Table("profile_users").
				Where("user_id = ?", userID).
				Updates(profileUpdates).Error
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		profileUpdates["user_id"] = userID
		return tx.Table("profile_users").Create(profileUpdates).Error
	})
	if err != nil {
		return nil, err
	}

	return r.FindByUserID(userID)
}

func (r *repository) UpdateProfileImage(
	userID uint,
	imageURL string,
) (*ProfileResponse, error) {
	profileUpdates := map[string]interface{}{
		"avatar_url": strings.TrimSpace(imageURL),
	}

	if err := r.upsertProfile(userID, profileUpdates); err != nil {
		return nil, err
	}

	return r.FindByUserID(userID)
}

func (r *repository) upsertProfile(
	userID uint,
	profileUpdates map[string]interface{},
) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var existing struct {
			ID uint `gorm:"column:id"`
		}

		err := tx.Table("profile_users").
			Select("id").
			Where("user_id = ?", userID).
			First(&existing).Error

		if err == nil {
			return tx.Table("profile_users").
				Where("user_id = ?", userID).
				Updates(profileUpdates).Error
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		profileUpdates["user_id"] = userID
		return tx.Table("profile_users").Create(profileUpdates).Error
	})
}

func (row profileRow) toResponse() *ProfileResponse {
	var canChangeUsernameAt *time.Time
	if row.UsernameChangedAt != nil {
		nextChange := row.UsernameChangedAt.Add(usernameChangeCooldown)
		canChangeUsernameAt = &nextChange
	}

	return &ProfileResponse{
		ID:                  row.ID,
		Username:            row.Username,
		Email:               row.Email,
		Role:                row.Role,
		UsernameChangedAt:   row.UsernameChangedAt,
		CanChangeUsernameAt: canChangeUsernameAt,
		BirthDate:           row.BirthDate,
		Gender:              row.Gender,
		ProfileImage:        stringValue(row.ProfileImage),
		Weight:              row.Weight,
		Height:              row.Height,
		CurrentStreak:       intValue(row.CurrentStreak),
	}
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}

func intValue(value *int) int {
	if value == nil {
		return 0
	}

	return *value
}
