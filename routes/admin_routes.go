package routes

import (
	"survey-backend/internal/country"
	"survey-backend/internal/user"
	"survey-backend/middleware"

	"github.com/gin-gonic/gin"
)

// AdminRoutes registers platform administration endpoints.
func AdminRoutes(r *gin.RouterGroup) {
	g := r.Group("/admin")
	g.Use(middleware.RoleMiddleware("admin"))

	g.POST("/country/toggle", country.ToggleCountryActiveHandler)
	g.POST("/points/grant", user.GrantPointsHandler)
}
