package handlers

import (
	"net/http"
	"skm/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EducationHandler struct {
	db *gorm.DB
}

func EducationsHandler(db *gorm.DB) *EducationHandler {
	db.AutoMigrate(&models.Education{})
	return &EducationHandler{db: db}
}

func (h *EducationHandler) GetEducations(c *gin.Context) {
	var educations []models.Education
	if err := h.db.Order("id ASC").Find(&educations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data pendidikan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"educations": educations,
	})
}
