package main

import (
	"fmt"
	"log"
	"os"

	"survey-backend/internal/category"
	"survey-backend/internal/country"
	"survey-backend/internal/region"
	"survey-backend/internal/survey"
	"survey-backend/internal/user"
	"survey-backend/middleware"
	"survey-backend/pkg/database"

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
		log.Fatal("Failed to seed countries:", err)
	}

	// auth routes
	r := gin.Default()
	r.POST("/auth/register", user.RegisterHandler)
	r.POST("/auth/login", user.LoginHandler)

	// protected routes
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		//user
		protected.GET("/me", user.GetUserHandler)
		protected.GET("/survey/form", survey.GetCreateSurveyFormHandler)
		protected.GET("/regions/:country_id", region.GetRegionsByCountryHandler)
		protected.POST("/survey/create", survey.CreateSurveyHandler)
		// add other protected routes here

		//Admin
		protected.POST("/country/toggle", country.ToggleCountryActiveHandler)

	}

	port := os.Getenv("PORT")

	fmt.Println("Server runnning on port", port)
	r.Run(":" + port)
}
