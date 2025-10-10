package routes

import (
	"survey-backend/internal/country"

	"github.com/gin-gonic/gin"
	"survey-backend/middleware"
)

func AdminRoutes(r *gin.RouterGroup) {
	r.Use(middleware.RoleMiddleware("admin"))
	r.POST("/country/toggle", country.ToggleCountryActiveHandler)
}
