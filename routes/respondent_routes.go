package routes

import (
	"survey-backend/internal/survey"
	"survey-backend/middleware"

	"github.com/gin-gonic/gin"
)

// RespondentRoutes registers the survey-answering API.
func RespondentRoutes(r *gin.RouterGroup) {
	g := r.Group("/respondent")
	g.Use(middleware.RoleMiddleware("respondent"))

	g.GET("/home", survey.RespondentHomeHandler)

	g.GET("/survey/list", survey.RespondentSurveyListHandler)
	g.GET("/survey/:id", survey.RespondentSurveyDetailHandler)
	g.POST("/survey/:id/answer", survey.SubmitAnswerHandler)
	g.GET("/survey/:id/answer-count", survey.AnswerCountHandler)

	g.GET("/saved", survey.SavedSurveysHandler)
	g.GET("/saved-ids", survey.SavedSurveyIDsHandler)
	g.POST("/saved/:id", survey.SaveSurveyHandler)
	g.DELETE("/saved/:id", survey.RemoveSavedSurveyHandler)

	g.GET("/completed", survey.RespondentCompletedHandler)
	g.GET("/completed/:id", survey.AnswerDetailsHandler)
}
