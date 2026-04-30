package entities
	
import "gorm.io/gorm"

type TrainerProfile struct {
    gorm.Model
    UserID         uint    `gorm:"uniqueIndex;not null"`
    Specialization string  `json:"specialization"`
    Experience     string  `json:"experience"`
    Certification  string  `json:"certification"`
    Rating         float64 `json:"rating" gorm:"default:0"`
    TotalClients   int     `json:"totalClients" gorm:"default:0"`
    Bio            string  `json:"bio"`

    User *User `gorm:"constraint:OnDelete:CASCADE;"`
}