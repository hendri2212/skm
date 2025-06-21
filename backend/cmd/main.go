package main

import (
	"log"
	"os"
	"skm/internal/models"
	"skm/internal/routes"
	"skm/internal/seeds"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}
	return db
}

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, relying on environment variables")
	}

	gin.SetMode(os.Getenv("GIN_MODE"))

	db := InitDB()

	// Auto migrate
	db.AutoMigrate(&models.User{}, &models.Question{}, &models.Occupation{}, &models.Education{}, &models.Choice{}, &models.Answer{})

	// Seed data
	seeds.SeedAll(db)

	router := gin.Default()
	routes.SetupRoutes(router, db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(router.Run(":" + port))
}
