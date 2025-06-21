package handlers

import (
	"net/http"
	"skm/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OccupationHandler struct {
	db *gorm.DB
}

func OccupationsHandler(db *gorm.DB) *OccupationHandler {
	db.AutoMigrate(&models.Occupation{})
	return &OccupationHandler{db: db}
}

func (h *OccupationHandler) GetOccupations(c *gin.Context) {
	var occupations []models.Occupation
	if err := h.db.Order("id ASC").Find(&occupations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data pendidikan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      "success",
		"occupations": occupations,
	})
}
