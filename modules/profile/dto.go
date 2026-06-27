package profile

import "time"

type ProfileResponse struct {
	ID                  uint       `json:"id"`
	Username            string     `json:"username"`
	Email               string     `json:"email"`
	Role                string     `json:"role"`
	UsernameChangedAt   *time.Time `json:"username_changed_at,omitempty"`
	CanChangeUsernameAt *time.Time `json:"can_change_username_at,omitempty"`
	BirthDate           *time.Time `json:"birth_date,omitempty"`
	Gender              *bool      `json:"gender"`
	ProfileImage        string     `json:"profile_image"`
	Weight              *int       `json:"weight"`
	Height              *int       `json:"height"`
	FitnessLevel        string     `json:"fitness_level"`
	MainGoal            string     `json:"main_goal"`
	WorkoutDays         *int       `json:"workout_days"`
	CurrentStreak       int        `json:"current_streak"`
}

type UpdateProfileRequest struct {
	Username     *string  `json:"username"`
	Gender       *bool    `json:"gender"`
	ProfileImage *string  `json:"profile_image"`
	Weight       *int     `json:"weight"`
	Height       *int     `json:"height"`
	FitnessLevel *string  `json:"fitness_level"`
	MainGoal     *string  `json:"main_goal"`
	WorkoutDays  *int     `json:"workout_days"`
}
