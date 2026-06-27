package entities

import "time"

type TrainerProfile struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        *uint     `gorm:"index" json:"user_id"`
	Name          string    `gorm:"type:varchar(120);not null" json:"name"`
	Specialty     string    `gorm:"type:varchar(160);not null" json:"specialty"`
	Categories    string    `gorm:"type:varchar(255)" json:"categories"`
	Tags          string    `gorm:"type:varchar(255)" json:"tags"`
	Photo         string    `gorm:"type:varchar(255)" json:"photo"`
	Price         int       `gorm:"not null;default:0" json:"price"`
	Rating        string    `gorm:"type:varchar(10);default:'0.0'" json:"rating"`
	Experience    string    `gorm:"type:varchar(80)" json:"experience"`
	Bio           string    `gorm:"type:text" json:"bio"`
	IsActive      bool      `gorm:"not null;default:true" json:"is_active"`
	IsOnline 	  bool `gorm:"not null;default:false" json:"is_online"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
