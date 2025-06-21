package seeds

import "gorm.io/gorm"

func SeedAll(db *gorm.DB) {
	SeedQuestions(db)
	SeedChoices(db)
	SeedEducations(db)
	SeedOccupations(db)
}
