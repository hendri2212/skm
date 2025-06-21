package seeds

import (
	"skm/internal/models"

	"gorm.io/gorm"
)

func SeedQuestions(db *gorm.DB) error {
	questions := []models.Question{
		{QuestionText: "Bagaimana pendapat Saudara tentang kesesuaian persyaratan pelayanan dengan jenis pelayanannya?"},
		{QuestionText: "Bagaimana pendapat saudara tentang kemudahan sistem, mekanisme dan prosedur pelayanan?."},
		{QuestionText: "Bagaimana pendapat saudara tentang kecepatan waktu dalam memberikan pelayanan?"},
		{QuestionText: "Bagaimana pendapat saudara tentang kewajaran biaya/tarif dalam pelayanan?"},
		{QuestionText: "Bagaimana pendapat saudara tentang kesesuaian produk pelayanan antara yang tercantum dalam standar pelayanan dengan hasil yang diberikan?"},
		{QuestionText: "Bagaimana pendapat saudara tentang kompetensi/kemampuan Petugas dalam memberikan pelayanan?"},
		{QuestionText: "Bagaimana pendapat saudara tentang kesopanan dan keramahan petugas dalam memberikan pelayanan?"},
		{QuestionText: "Bagaimana pendapat saudara tentang cara penanganan petugas terhadap pengaduan, saran dan masukan?"},
		{QuestionText: "Bagaimana pendapat saudara tentang kualitas sarana dan prasarana pelayanan?"},
	}

	for _, q := range questions {
		if err := db.FirstOrCreate(&models.Question{}, q).Error; err != nil {
			return err
		}
	}

	return nil
}
