package dailychallenge

import (
	"errors"
	"stable/database/entities"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	FindCompletion(userID uint, challengeDate time.Time) (*entities.DailyChallenge, error)
	FindStreak(userID uint) (*entities.UserStreak, error)
	CompleteAndUpdateStreak(challenge *entities.DailyChallenge) (*entities.UserStreak, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindCompletion(
	userID uint,
	challengeDate time.Time,
) (*entities.DailyChallenge, error) {
	var challenge entities.DailyChallenge

	err := r.db.
		Where("user_id = ? AND challenge_date = ?", userID, challengeDate).
		First(&challenge).Error
	if err != nil {
		return nil, err
	}

	return &challenge, nil
}

func (r *repository) FindStreak(userID uint) (*entities.UserStreak, error) {
	var streak entities.UserStreak

	if err := r.db.Where("user_id = ?", userID).First(&streak).Error; err != nil {
		return nil, err
	}

	return &streak, nil
}

func (r *repository) CompleteAndUpdateStreak(
	challenge *entities.DailyChallenge,
) (*entities.UserStreak, error) {
	var updatedStreak entities.UserStreak

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(challenge).Error; err != nil {
			if isDuplicateError(err) {
				return ErrAlreadyCompleted
			}
			return err
		}

		var streak entities.UserStreak
		err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", challenge.UserID).
			First(&streak).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			challengeDate := challenge.ChallengeDate
			streak = entities.UserStreak{
				UserID:            challenge.UserID,
				CurrentStreak:     1,
				LongestStreak:      1,
				LastCompletedDate: &challengeDate,
			}

			if err := tx.Create(&streak).Error; err != nil {
				return err
			}

			updatedStreak = streak
			return nil
		}
		if err != nil {
			return err
		}

		streak.CurrentStreak = calculateNextStreak(
			streak.CurrentStreak,
			streak.LastCompletedDate,
			challenge.ChallengeDate,
		)
		if streak.CurrentStreak > streak.LongestStreak {
			streak.LongestStreak = streak.CurrentStreak
		}

		challengeDate := challenge.ChallengeDate
		streak.LastCompletedDate = &challengeDate

		if err := tx.Save(&streak).Error; err != nil {
			return err
		}

		updatedStreak = streak
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &updatedStreak, nil
}

func calculateNextStreak(
	currentStreak int,
	lastCompletedDate *time.Time,
	challengeDate time.Time,
) int {
	if lastCompletedDate == nil {
		return 1
	}

	lastDate := normalizeDate(*lastCompletedDate)
	today := normalizeDate(challengeDate)

	if lastDate.Equal(today) {
		return currentStreak
	}
	if lastDate.Equal(today.AddDate(0, 0, -1)) {
		return currentStreak + 1
	}

	return 1
}

func normalizeDate(value time.Time) time.Time {
	year, month, day := value.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, jakartaLocation)
}

func isDuplicateError(err error) bool {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}

	message := strings.ToLower(err.Error())
	return strings.Contains(message, "duplicate") ||
		strings.Contains(message, "unique constraint")
}
