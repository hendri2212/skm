package routes

import (
	"net/http"
	"skm/internal/handlers"

	"gorm.io/gorm"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	router.Static("/uploads", "./uploads")

	var allowedOrigins []string
	if gin.Mode() == gin.ReleaseMode {
		allowedOrigins = []string{
			"http://sipaktusarah.com", "https://skm.sipaktusarah.com",
		}
	} else {
		allowedOrigins = []string{
			"http://127.0.0.1:8000",
		}
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-CSRF-Token", "X-XSRF-TOKEN", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	router.OPTIONS("/*path", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNoContent)
	})

	userHandler := handlers.UsersHandler(db)
	questionsHandler := handlers.QuestionsHandler(db)
	educationsHandler := handlers.EducationsHandler(db)
	occupationsHandler := handlers.OccupationsHandler(db)

	api := router.Group("/api")
	{
		// api.POST("/login", userHandler.LoginUser)
		api.POST("/logout", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
		})

		// api.Use(middlewares.AuthMiddleware())

		// api.GET("/me", userHandler.Me)
		api.GET("/countAge", userHandler.CountAge)
		api.GET("/countGender", userHandler.CountGender)
		api.GET("/countEducation", userHandler.CountEducation)
		api.GET("/countOccupation", userHandler.CountOccupation)
		api.GET("/users", userHandler.GetUsers)
		api.GET("/user-answer", userHandler.GetUserAnswers)
		api.GET("/user-answer-all", userHandler.GetUserAnswerAll)
		api.GET("/user-answer/:id", userHandler.GetUserAnswerByID)

		api.GET("/questions", questionsHandler.GetQuestions)
		api.POST("/answers", questionsHandler.SubmitAnswers)
		api.GET("/educations", educationsHandler.GetEducations)
		api.GET("/occupations", occupationsHandler.GetOccupations)
	}
}
