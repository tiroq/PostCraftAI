package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tiroq/postcraftai/backend/handlers"
	"github.com/tiroq/postcraftai/backend/models"
)

func main() {
	// Pre-create admin user.
	models.PreCreateAdmin()

	// Use gin.Default() to include Logger and Recovery middleware.
	r := gin.Default()

	// Public endpoints.
	r.POST("/signup", handlers.Signup)
	r.POST("/login", handlers.Login)

	// Protected endpoints.
	protected := r.Group("/")
	protected.Use(handlers.AuthMiddleware())
	{
		protected.POST("/generate-post", handlers.GeneratePost)
		admin := protected.Group("/admin")
		admin.Use(handlers.AdminMiddleware())
		{
			admin.POST("/enable-user", handlers.AdminEnableUser)
			admin.GET("/list-users", handlers.AdminListUsers)
			admin.POST("/update-rate-limit", handlers.AdminUpdateRateLimit)
			admin.POST("/update-expiration", handlers.AdminUpdateExpiration)
			admin.GET("/request-stats", handlers.AdminRequestStats)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server running on :" + port)
	r.Run(":" + port)
}
