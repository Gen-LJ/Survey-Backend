package routes

import (
	"survey-backend/internal/user"

	"github.com/gin-gonic/gin"
)

// PublicRoutes need no token.
//
// register-form is the one the signup screen depends on: it returns the active
// countries with their regions nested, so a visitor can fill both pickers
// before they have an account.
func PublicRoutes(r *gin.RouterGroup) {
	r.POST("/auth/register", user.RegisterHandler)
	r.POST("/auth/login", user.LoginHandler)
	r.GET("/auth/register-form", user.GetRegisterFormHandler)
}
