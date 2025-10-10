package routes

import (
	"survey-backend/internal/survey"

	"github.com/gin-gonic/gin"
	"survey-backend/middleware"
)

func RespondentRoutes(r *gin.RouterGroup) {
	r.Use(middleware.RoleMiddleware("respondent"))
	r.GET("/survey/form", survey.GetCreateSurveyFormHandler)
}
