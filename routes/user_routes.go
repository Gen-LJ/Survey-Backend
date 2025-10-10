package routes


import (
	"survey-backend/internal/region"
	"survey-backend/internal/user"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.RouterGroup) {
	r.GET("/me", user.GetUserHandler)
	r.GET("/regions/:country_id", region.GetRegionsByCountryHandler)
}
