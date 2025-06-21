package seeds

import (
	"skm/internal/models"

	"gorm.io/gorm"
)

func SeedChoices(db *gorm.DB) error {
	choices := []models.Choice{
		{QuestionID: 1, ChoiceText: "Tidak sesuai", Points: 1},
		{QuestionID: 1, ChoiceText: "Kurang sesuai", Points: 2},
		{QuestionID: 1, ChoiceText: "Sesuai", Points: 3},
		{QuestionID: 1, ChoiceText: "Sangat sesuai", Points: 4},

		{QuestionID: 2, ChoiceText: "Tidak mudah", Points: 1},
		{QuestionID: 2, ChoiceText: "Kurang mudah", Points: 2},
		{QuestionID: 2, ChoiceText: "Mudah", Points: 3},
		{QuestionID: 2, ChoiceText: "Sangat mudah", Points: 4},

		{QuestionID: 3, ChoiceText: "Tidak cepat", Points: 1},
		{QuestionID: 3, ChoiceText: "Kurang cepat", Points: 2},
		{QuestionID: 3, ChoiceText: "Cepat", Points: 3},
		{QuestionID: 3, ChoiceText: "Sangat cepat", Points: 4},

		{QuestionID: 4, ChoiceText: "Sangat mahal", Points: 1},
		{QuestionID: 4, ChoiceText: "Cukup mahal", Points: 2},
		{QuestionID: 4, ChoiceText: "Murah", Points: 3},
		{QuestionID: 4, ChoiceText: "Gratis", Points: 4},

		{QuestionID: 5, ChoiceText: "Tidak sesuai", Points: 1},
		{QuestionID: 5, ChoiceText: "Kurang sesuai", Points: 2},
		{QuestionID: 5, ChoiceText: "Sesuai", Points: 3},
		{QuestionID: 5, ChoiceText: "Sangat sesuai", Points: 4},

		{QuestionID: 6, ChoiceText: "Tidak mampu", Points: 1},
		{QuestionID: 6, ChoiceText: "Kurang mampu", Points: 2},
		{QuestionID: 6, ChoiceText: "Mampu", Points: 3},
		{QuestionID: 6, ChoiceText: "Sangat mampu", Points: 4},

		{QuestionID: 7, ChoiceText: "Tidak sopan dan tidak ramah", Points: 1},
		{QuestionID: 7, ChoiceText: "Kurang sopan dan kurang ramah", Points: 2},
		{QuestionID: 7, ChoiceText: "Sopan dan ramah", Points: 3},
		{QuestionID: 7, ChoiceText: "Sangat sopan dan sangat ramah", Points: 4},

		{QuestionID: 8, ChoiceText: "Tidak ada", Points: 1},
		{QuestionID: 8, ChoiceText: "Ada, tetapi tidak berfungsi", Points: 2},
		{QuestionID: 8, ChoiceText: "Berfungsi, tetapi kurang maksimal", Points: 3},
		{QuestionID: 8, ChoiceText: "Dikelola dengan baik", Points: 4},

		{QuestionID: 9, ChoiceText: "Buruk", Points: 1},
		{QuestionID: 9, ChoiceText: "Cukup", Points: 2},
		{QuestionID: 9, ChoiceText: "Baik", Points: 3},
		{QuestionID: 9, ChoiceText: "Sangat baik", Points: 4},
	}

	for _, c := range choices {
		if err := db.FirstOrCreate(&models.Choice{}, c).Error; err != nil {
			return err
		}
	}

	return nil
}
