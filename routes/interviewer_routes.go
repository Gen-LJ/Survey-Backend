package routes

import (
	"survey-backend/internal/survey"

	"github.com/gin-gonic/gin"
	"survey-backend/middleware"
)

func InterviewerRoutes(r *gin.RouterGroup) {
	r.Use(middleware.RoleMiddleware("interviewer"))
	r.POST("/survey/create", survey.CreateSurveyHandler)
}
