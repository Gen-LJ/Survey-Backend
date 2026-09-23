package survey

import (
	"net/http"

	answerResponse "survey-backend/internal/answer/response"
	"survey-backend/internal/survey/request"
	surResponse "survey-backend/internal/survey/response"
	"survey-backend/pkg/pagination"
	"survey-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

func RespondentHomeHandler(c *gin.Context) {
	r := &response.Response[surResponse.RespondentHomeResponse]{}

	result, err := GetRespondentHome(currentUser(c))
	if err != nil {
		c.JSON(StatusFor(err), r.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, r.OK(*result))
}

func RespondentSurveyListHandler(c *gin.Context) {
	r := &response.Response[pagination.Page[surResponse.SurveySummaryResponse]]{}

	result, err := ListForRespondent(currentUser(c), pagination.Parse(c))
	if err != nil {
		c.JSON(StatusFor(err), r.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, r.OK(*result))
}

func RespondentSurveyDetailHandler(c *gin.Context) {
	r := &response.Response[surResponse.SurveyResponse]{}

	surveyID, ok := idParam(c, "id")
	if !ok {
		return
	}

	result, err := GetSurveyForRespondent(surveyID, currentUser(c))
	if err != nil {
		c.JSON(StatusFor(err), r.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, r.OK(*result))
}

func SubmitAnswerHandler(c *gin.Context) {
	r := &response.Response[answerResponse.AnswerResponse]{}

	surveyID, ok := idParam(c, "id")
	if !ok {
		return
	}

	var req request.SubmitAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, r.Fail(err.Error()))
		return
	}

	submitted := make([]SubmittedAnswer, len(req.Answers))
	for i, item := range req.Answers {
		submitted[i] = SubmittedAnswer{QuestionID: item.QuestionID, Responses: item.Responses}
	}

	result, err := SubmitAnswer(surveyID, currentUser(c).ID, submitted)
	if err != nil {
		c.JSON(StatusFor(err), r.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, r.OK(*result))
}

func AnswerCountHandler(c *gin.Context) {
	r := &response.Response[map[string]uint]{}

	surveyID, ok := idParam(c, "id")
	if !ok {
		return
	}

	count, err := GetAnswerCount(surveyID)
	if err != nil {
		c.JSON(StatusFor(err), r.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, r.OK(map[string]uint{"answer_count": count}))
}

func SaveSurveyHandler(c *gin.Context) {
	r := &response.Response[struct{}]{}

	surveyID, ok := idParam(c, "id")
	if !ok {
		return
	}

	if err := SaveSurvey(currentUser(c).ID, surveyID); err != nil {
		c.JSON(StatusFor(err), r.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, r.StatusOK("Survey saved successfully"))
}

func RemoveSavedSurveyHandler(c *gin.Context) {
	r := &response.Response[struct{}]{}

	surveyID, ok := idParam(c, "id")
	if !ok {
		return
	}

	if err := RemoveSavedSurvey(currentUser(c).ID, surveyID); err != nil {
		c.JSON(StatusFor(err), r.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, r.StatusOK("Survey removed from saved list"))
}

func SavedSurveyIDsHandler(c *gin.Context) {
	r := &response.Response[[]uint]{}

	ids, err := GetSavedSurveyIDs(currentUser(c).ID)
	if err != nil {
		c.JSON(StatusFor(err), r.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, r.OK(ids))
}

func SavedSurveysHandler(c *gin.Context) {
	r := &response.Response[[]surResponse.SurveySummaryResponse]{}

	surveys, err := GetSavedSurveys(currentUser(c).ID)
	if err != nil {
		c.JSON(StatusFor(err), r.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, r.OK(surveys))
}

func RespondentCompletedHandler(c *gin.Context) {
	r := &response.Response[pagination.Page[surResponse.AnsweredSurveyResponse]]{}

	// sort=last_modified falls back to the survey's own activity, matching the
	// sortLastAnswerAt toggle the Flutter data source exposed.
	sortByAnsweredAt := c.Query("sort") != "last_modified"

	result, err := ListCompletedForRespondent(currentUser(c).ID, sortByAnsweredAt, pagination.Parse(c))
	if err != nil {
		c.JSON(StatusFor(err), r.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, r.OK(*result))
}

func AnswerDetailsHandler(c *gin.Context) {
	r := &response.Response[surResponse.AnswerDataDetailsResponse]{}

	surveyID, ok := idParam(c, "id")
	if !ok {
		return
	}

	result, err := GetAnswerDataDetails(surveyID, currentUser(c).ID)
	if err != nil {
		c.JSON(StatusFor(err), r.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, r.OK(*result))
}
