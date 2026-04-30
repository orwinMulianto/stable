package trainer

import "time"

type CreateTrainerProfileRequest struct {
	Specialization string  `json:"specialization" binding:"required"`
	Experience     string  `json:"experience" binding:"required"`
	Certification  string  `json:"certification"`
	Bio            string  `json:"bio"`
}

type UpdateTrainerProfileRequest struct {
	Specialization string  `json:"specialization"`
	Experience     string  `json:"experience"`
	Certification  string  `json:"certification"`
	Bio            string  `json:"bio"`
}

type TrainerProfileResponse struct {
	ID             uint      `json:"id"`
	UserID         uint      `json:"userId"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	ProfileImage   string    `json:"profileImage"`
	Specialization string    `json:"specialization"`
	Experience     string    `json:"experience"`
	Certification  string    `json:"certification"`
	Rating         float64   `json:"rating"`
	TotalClients   int       `json:"totalClients"`
	Bio            string    `json:"bio"`
	CreatedAt      time.Time `json:"createdAt"`
}