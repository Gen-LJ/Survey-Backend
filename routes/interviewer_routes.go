package routes

import (
	"survey-backend/internal/survey"
	"survey-backend/middleware"

	"github.com/gin-gonic/gin"
)

// InterviewerRoutes registers the survey-authoring API.
//
// The group is created with Group("") rather than calling Use on the shared
// parent: Use would append the role check to the parent's handler chain, so
// every group registered afterwards would inherit it too.
func InterviewerRoutes(r *gin.RouterGroup) {
	g := r.Group("/interviewer")
	g.Use(middleware.RoleMiddleware("interviewer"))

	g.GET("/home", survey.InterviewerHomeHandler)
	g.GET("/completed", survey.InterviewerCompletedHandler)

	g.GET("/survey/form", survey.GetCreateSurveyFormHandler)
	g.POST("/survey/create", survey.CreateSurveyHandler)
	g.GET("/survey/list", survey.ListSurveysHandler)
	g.GET("/survey/count", survey.SurveyCountHandler)

	g.GET("/survey/:id", survey.GetSurveyHandler)
	g.PUT("/survey/:id", survey.EditSurveyInfoHandler)
	g.DELETE("/survey/:id", survey.DeleteSurveyHandler)
	g.POST("/survey/:id/publish", survey.PublishSurveyHandler)
	g.POST("/survey/:id/pause", survey.PauseSurveyHandler)
	g.GET("/survey/:id/analytics", survey.AnalyticsHandler)

	g.GET("/survey/:id/questions", survey.GetQuestionsHandler)
	g.POST("/survey/:id/questions", survey.AddQuestionHandler)
	g.PUT("/survey/:id/questions/:question_id", survey.EditQuestionHandler)
	g.DELETE("/survey/:id/questions/:question_id", survey.DeleteQuestionHandler)
}
