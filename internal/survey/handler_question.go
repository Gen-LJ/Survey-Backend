package survey

import (
	"net/http"

	questionResponse "survey-backend/internal/question/response"
	"survey-backend/internal/survey/request"
	"survey-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

func GetQuestionsHandler(c *gin.Context) {
	r := &response.Response[[]questionResponse.QuestionResponse]{}

	surveyID, ok := idParam(c, "id")
	if !ok {
		return
	}

	questions, err := GetQuestions(surveyID, currentUser(c).ID)
	if err != nil {
		c.JSON(StatusFor(err), r.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, r.OK(questions))
}

func AddQuestionHandler(c *gin.Context) {
	r := &response.Response[questionResponse.QuestionResponse]{}

	surveyID, ok := idParam(c, "id")
	if !ok {
		return
	}

	var req request.AddQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, r.Fail(err.Error()))
		return
	}

	created, err := AddQuestion(surveyID, currentUser(c).ID, req.Text, req.QuestionType, req.AllowMultiAnswer, req.Options)
	if err != nil {
		c.JSON(StatusFor(err), r.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, r.OK(*created))
}

func EditQuestionHandler(c *gin.Context) {
	r := &response.Response[struct{}]{}

	surveyID, ok := idParam(c, "id")
	if !ok {
		return
	}
	questionID, ok := idParam(c, "question_id")
	if !ok {
		return
	}

	var req request.EditQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, r.Fail(err.Error()))
		return
	}

	if err := EditQuestion(surveyID, currentUser(c).ID, questionID, req.Text, req.Options); err != nil {
		c.JSON(StatusFor(err), r.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, r.StatusOK("Question updated successfully"))
}

func DeleteQuestionHandler(c *gin.Context) {
	r := &response.Response[struct{}]{}

	surveyID, ok := idParam(c, "id")
	if !ok {
		return
	}
	questionID, ok := idParam(c, "question_id")
	if !ok {
		return
	}

	if err := DeleteQuestion(surveyID, currentUser(c).ID, questionID); err != nil {
		c.JSON(StatusFor(err), r.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, r.StatusOK("Question deleted successfully"))
}
