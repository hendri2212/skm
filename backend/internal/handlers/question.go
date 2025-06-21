package handlers

import (
	"net/http"
	"skm/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type QuestionHandler struct {
	db *gorm.DB
}

func QuestionsHandler(db *gorm.DB) *QuestionHandler {
	db.AutoMigrate(&models.Question{})
	return &QuestionHandler{db: db}
}

func (h *QuestionHandler) GetQuestions(c *gin.Context) {
	var questions []models.Question
	// Preload Choices
	if err := h.db.Preload("Choices").Find(&questions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, questions)
}

func (h *QuestionHandler) SubmitAnswers(c *gin.Context) {
	// Bind and validate payload from models.SubmitAnswersPayload
	var payload models.SubmitAnswersPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse birth date string to time.Time
	birthDate, err := time.Parse("2006-01-02", payload.User.BirthDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format birthDate harus YYYY-MM-DD"})
		return
	}

	// Start transaction
	tx := h.db.Begin()

	// Create user record
	user := models.User{
		FullName:         payload.User.Name,
		PlaceOfBirth:     payload.User.BirthPlace,
		DateOfBirth:      birthDate,
		IsMale:           payload.User.IsMale,
		LastEducationID:  payload.User.EducationID,
		MainOccupationID: payload.User.OccupationID,
	}
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan user: " + err.Error()})
		return
	}

	// Create answer records
	for _, a := range payload.Answers {
		ans := models.Answer{
			UserID:     user.ID,
			QuestionID: a.QuestionID,
			ChoiceID:   a.ChoiceID,
		}
		if err := tx.Create(&ans).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan jawaban: " + err.Error()})
			return
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal commit transaksi"})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, gin.H{
		"message": "Data berhasil disimpan",
		"user_id": user.ID,
	})
}
