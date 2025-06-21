package seeds

import (
	"skm/internal/models"

	"gorm.io/gorm"
)

// SeedEducations inserts default education levels into the educations table.
func SeedEducations(db *gorm.DB) error {
	educations := []models.Education{
		{Name: "SD Kebawah"},
		{Name: "SLTP"},
		{Name: "SLTA"},
		{Name: "S1-D3-D4"},
		{Name: "S-2 ke atas"},
	}

	for _, e := range educations {
		// FirstOrCreate memastikan tidak duplikat jika sudah ada
		if err := db.FirstOrCreate(&models.Education{}, e).Error; err != nil {
			return err
		}
	}
	return nil
}
