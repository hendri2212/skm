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
			"https://skm.sipaktusarah.com",
		}
	} else {
		allowedOrigins = []string{
			"http://localhost:5173",
		}
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins: allowedOrigins,
		// AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	router.OPTIONS("/*path", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNoContent)
	})

	// userHandler := handlers.UsersHandler(db)
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
		// api.GET("/users", userHandler.GetUsers)
		// api.POST("/users", userHandler.CreateUser)
		// api.GET("/users/:id", userHandler.GetUserByID)
		// api.PUT("/users/:id", userHandler.UpdateUser)
		// api.DELETE("/users/:id", userHandler.DeleteUser)

		api.GET("/questions", questionsHandler.GetQuestions)
		api.POST("/answers", questionsHandler.SubmitAnswers)
		api.GET("/educations", educationsHandler.GetEducations)
		api.GET("/occupations", occupationsHandler.GetOccupations)
	}
}
