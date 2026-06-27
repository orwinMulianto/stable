package dailychallenge

import (
	"errors"
	"stable/database/entities"
	"time"

	"gorm.io/gorm"
)

var (
	ErrAlreadyCompleted = errors.New("daily challenge already completed today")
	ErrTargetNotReached = errors.New("target repetitions have not been reached")
)

var jakartaLocation = time.FixedZone("Asia/Jakarta", 7*60*60)

var challengeRotation = []ChallengeDTO{
	{
		Key: "bicep-curl", Exercise: "Bicep Curl", MuscleGroup: "Biceps",
		Instruction: "Selesaikan 12 repetisi dengan gerakan terkontrol.",
		TargetReps: 12, XPReward: 20,
	},
	{
		Key: "tricep-dips", Exercise: "Tricep Dips", MuscleGroup: "Triceps",
		Instruction: "Selesaikan 10 repetisi dan jaga siku tetap mengarah ke belakang.",
		TargetReps: 10, XPReward: 20,
	},
	{
		Key: "bodyweight-squat", Exercise: "Bodyweight Squat", MuscleGroup: "Legs",
		Instruction: "Selesaikan 15 repetisi dengan posisi lutut tetap stabil.",
		TargetReps: 15, XPReward: 20,
	},
	{
		Key: "push-up", Exercise: "Push Up", MuscleGroup: "Chest",
		Instruction: "Selesaikan 10 repetisi dengan tubuh membentuk garis lurus.",
		TargetReps: 10, XPReward: 20,
	},
	{
		Key: "reverse-lunge", Exercise: "Reverse Lunge", MuscleGroup: "Legs",
		Instruction: "Selesaikan 12 repetisi total secara bergantian.",
		TargetReps: 12, XPReward: 20,
	},
	{
		Key: "shoulder-press", Exercise: "Shoulder Press", MuscleGroup: "Shoulders",
		Instruction: "Selesaikan 12 repetisi tanpa melengkungkan punggung.",
		TargetReps: 12, XPReward: 20,
	},
	{
		Key: "sit-up", Exercise: "Sit Up", MuscleGroup: "Core",
		Instruction: "Selesaikan 15 repetisi dengan tempo yang stabil.",
		TargetReps: 15, XPReward: 20,
	},
}

type Service interface {
	GetToday(userID uint) (*TodayChallengeResponse, error)
	CompleteToday(userID uint, repetitions int) (*CompleteChallengeResponse, error)
}

type service struct {
	repository Repository
	now        func() time.Time
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
		now:        time.Now,
	}
}

func (s *service) GetToday(userID uint) (*TodayChallengeResponse, error) {
	now := s.now().In(jakartaLocation)
	today := startOfDay(now)
	challenge := challengeForDate(today)
	nextReset := today.AddDate(0, 0, 1)

	response := &TodayChallengeResponse{
		Date:        today.Format("2006-01-02"),
		Challenge:   challenge,
		Completed:   false,
		NextResetAt: nextReset,
	}

	streak, streakErr := s.repository.FindStreak(userID)
	if streakErr == nil {
		response.Streak = effectiveStreak(streak, today)
	} else if !errors.Is(streakErr, gorm.ErrRecordNotFound) {
		return nil, streakErr
	}

	completion, err := s.repository.FindCompletion(userID, today)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return response, nil
	}
	if err != nil {
		return nil, err
	}

	completedAt := completion.CompletedAt
	response.Completed = true
	response.CompletedAt = &completedAt

	return response, nil
}

func (s *service) CompleteToday(
	userID uint,
	repetitions int,
) (*CompleteChallengeResponse, error) {
	now := s.now().In(jakartaLocation)
	today := startOfDay(now)
	challenge := challengeForDate(today)

	if repetitions < challenge.TargetReps {
		return nil, ErrTargetNotReached
	}

	_, err := s.repository.FindCompletion(userID, today)
	if err == nil {
		return nil, ErrAlreadyCompleted
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	completedAt := now

	dailyChallenge := &entities.DailyChallenge{
		UserID:        userID,
		ChallengeDate: today,
		Name:          challenge.Exercise,
		Description:   challenge.Instruction,
		Difficulty:    "Daily",
		CreatedBy:     userID,
		ChallengeKey:  challenge.Key,
		Repetitions:   repetitions,
		CompletedAt:   completedAt,
	}

	streak, err := s.repository.CompleteAndUpdateStreak(dailyChallenge)
	if err != nil {
		return nil, err
	}

	return &CompleteChallengeResponse{
		Message:     "daily challenge completed",
		Challenge:   challenge,
		CompletedAt: completedAt,
		XPReceived:  challenge.XPReward,
		Streak: StreakDTO{
			Current: streak.CurrentStreak,
			Longest:  streak.LongestStreak,
		},
	}, nil
}

func effectiveStreak(streak *entities.UserStreak, today time.Time) StreakDTO {
	result := StreakDTO{Longest: streak.LongestStreak}

	if streak.LastCompletedDate == nil {
		return result
	}

	lastDate := normalizeDate(*streak.LastCompletedDate)
	if lastDate.Equal(today) || lastDate.Equal(today.AddDate(0, 0, -1)) {
		result.Current = streak.CurrentStreak
	}

	return result
}

func startOfDay(value time.Time) time.Time {
	year, month, day := value.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, jakartaLocation)
}

func challengeForDate(date time.Time) ChallengeDTO {
	dayNumber := int(date.Unix() / int64(24*time.Hour/time.Second))
	index := dayNumber % len(challengeRotation)

	if index < 0 {
		index += len(challengeRotation)
	}

	return challengeRotation[index]
}