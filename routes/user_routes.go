package routes

import (
	"survey-backend/internal/region"
	"survey-backend/internal/user"

	"github.com/gin-gonic/gin"
)

// UserRoutes are available to any signed-in user regardless of role.
//
// regions stays here rather than in PublicRoutes: signup reads its regions from
// /auth/register-form, so this is only for an account already in session
// changing their country.
func UserRoutes(r *gin.RouterGroup) {
	r.GET("/me", user.GetUserHandler)
	r.GET("/regions/:country_id", region.GetRegionsByCountryHandler)
}
