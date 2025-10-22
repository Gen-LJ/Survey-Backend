package main

import (
	"fmt"
	"log"
	"os"

	"survey-backend/internal/category"
	"survey-backend/internal/country"
	"survey-backend/internal/region"
	"survey-backend/internal/user"
	"survey-backend/middleware"
	"survey-backend/pkg/database"
	"survey-backend/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error Loading .env File")
	}

	// database
	database.ConnectDB()
	database.DB.AutoMigrate(&user.User{}, &category.Category{}, &country.Country{}, &region.Region{})

	// seed
	if err := category.SeedCategories(); err != nil {
		log.Fatal("Failed to seed categories:", err)
	}
	if err := country.SeedCountries(); err != nil {
		log.Fatal("Failed to seed countries:", err)
	}
	if err := region.SeedRegions(); err != nil {
		log.Fatal("Failed to seed regions:", err)
	}

	// auth routes
	r := gin.Default()
	r.POST("/auth/register", user.RegisterHandler)
	r.POST("/auth/login", user.LoginHandler)
	r.GET("/auth/register-form",user.GetRegisterFormHandler)

	// protected routes
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	
	// Attach role-based groups
	routes.UserRoutes(protected)
	routes.RespondentRoutes(protected)
	routes.InterviewerRoutes(protected)
	routes.AdminRoutes(protected)

	port := os.Getenv("PORT")

	fmt.Println("Server runnning on port", port)
	r.Run(":" + port)
}
