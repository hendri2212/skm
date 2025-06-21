package seeds

import (
	"skm/internal/models"

	"gorm.io/gorm"
)

// SeedOccupations inserts default occupation options into the occupations table.
func SeedOccupations(db *gorm.DB) error {
	occupations := []models.Occupation{
		{Name: "PNS/TNI/Polri"},
		{Name: "Pegawai Swasta"},
		{Name: "Wiraswastawan/Usahawan"},
		{Name: "Pelajar/Mahasiswa"},
		{Name: "Lainnya"},
	}

	for _, o := range occupations {
		// FirstOrCreate memastikan tidak duplikat jika sudah ada
		if err := db.FirstOrCreate(&models.Occupation{}, o).Error; err != nil {
			return err
		}
	}
	return nil
}
