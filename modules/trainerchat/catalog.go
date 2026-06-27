package trainerchat

import (
	"stable/database/entities"
	"strings"
)

type TrainerCatalogItem struct {
	ID         uint     `json:"id"`
	Name       string   `json:"name"`
	Specialty  string   `json:"specialty"`
	Categories []string `json:"categories"`
	Tags       []string `json:"tags"`
	Photo      string   `json:"photo"`
	Price      int      `json:"price"`
	Rating     string   `json:"rating"`
	Experience string   `json:"experience"`
	Bio        string   `json:"bio"`
	Online bool `json:"online"`
	IsOnline bool `json:"is_online"`
}

func trainerEntityToCatalogItem(trainer entities.TrainerProfile) TrainerCatalogItem {
	return TrainerCatalogItem{
		ID:         trainer.ID,
		Name:       trainer.Name,
		Specialty:  trainer.Specialty,
		Categories: splitCSV(trainer.Categories),
		Tags:       splitCSV(trainer.Tags),
		Photo:      trainer.Photo,
		Price:      trainer.Price,
		Rating:     trainer.Rating,
		Experience: trainer.Experience,
		Bio:        trainer.Bio,
		Online:    trainer.IsOnline,
		IsOnline:  trainer.IsOnline,
	}
}

func splitCSV(value string) []string {
	if strings.TrimSpace(value) == "" {
		return []string{}
	}

	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}