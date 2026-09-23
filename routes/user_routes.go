package routes

import (
	"survey-backend/internal/region"
	"survey-backend/internal/user"

	"github.com/gin-gonic/gin"
)

// UserRoutes are available to any signed-in user regardless of role.
func UserRoutes(r *gin.RouterGroup) {
	r.GET("/me", user.GetUserHandler)
	r.GET("/regions/:country_id", region.GetRegionsByCountryHandler)
}
