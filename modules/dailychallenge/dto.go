package dailychallenge

import "time"

type CompleteChallengeRequest struct {
	UserID      uint `json:"user_id" binding:"required"`
	Repetitions int  `json:"repetitions" binding:"required,min=1"`
}

type ChallengeDTO struct {
	Key         string `json:"key"`
	Exercise    string `json:"exercise"`
	MuscleGroup string `json:"muscle_group"`
	Instruction string `json:"instruction"`
	TargetReps  int    `json:"target_reps"`
	XPReward    int    `json:"xp_reward"`
}

type TodayChallengeResponse struct {
	Date        string       `json:"date"`
	Challenge   ChallengeDTO `json:"challenge"`
	Completed   bool         `json:"completed"`
	CompletedAt *time.Time   `json:"completed_at,omitempty"`
	NextResetAt time.Time    `json:"next_reset_at"`
	Streak      StreakDTO    `json:"streak"`
}

type CompleteChallengeResponse struct {
	Message     string       `json:"message"`
	Challenge   ChallengeDTO `json:"challenge"`
	CompletedAt time.Time    `json:"completed_at"`
	XPReceived  int          `json:"xp_received"`
	Streak      StreakDTO    `json:"streak"`
}

type StreakDTO struct {
	Current int `json:"current"`
	Longest  int `json:"longest"`
}
