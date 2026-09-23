package survey

import (
	"net/http"
	"strconv"
	"strings"

	"survey-backend/internal/survey/request"
	surResponse "survey-backend/internal/survey/response"
	"survey-backend/internal/user"
	"survey-backend/pkg/pagination"
	"survey-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// currentUser reads the user the auth middleware stored. The middleware sets a
// user.User value, so this asserts to the value type, not a pointer.
func currentUser(c *gin.Context) user.User {
	return c.MustGet("currentUser").(user.User)
}

// idParam parses a uint path parameter, writing the 400 itself when it fails.
func idParam(c *gin.Context, name string) (uint, bool) {
	r := &response.Response[struct{}]{}

	id, err := strconv.ParseUint(c.Param(name), 10, 32)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, r.Fail("invalid "+name))
		return 0, false
	}

	return uint(id), true
}

// statesParam reads a repeated or comma separated ?state= filter.
func statesParam(c *gin.Context) []string {
	var states []string
	for _, raw := range c.QueryArray("state") {
		for _, state := range strings.Split(raw, ",") {
			if state = strings.TrimSpace(state); state != "" {
				states = append(states, state)
			}
		}
	}
	return states
}

func CreateSurveyHandler(c *gin.Context) {
	r := &response.Response[surResponse.SurveyResponse]{}

	var req request.CreateSurveyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, r.Fail(err.Error()))
		return
	}

	creator := currentUser(c)

	survey := &Survey{
		Title:       req.Title,
		Description: req.Description,
		CategoryID:  req.CategoryId,
		CountryID:   req.CountryId,
		RegionID:    req.RegionId,
		Minutes:     req.Minutes,
		CreatorID:   creator.ID,
	}

	if err := CreateSurvey(survey); err != nil {
		c.JSON(StatusFor(err), r.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, r.OK(survey.ToResponse(nil)))
}

func GetSurveyHandler(c *gin.Context) {
	r := &response.Response[surResponse.SurveyResponse]{}

	surveyID, ok := idParam(c, "id")
	if !ok {
		return
	}

	result, err := GetSurveyForCreator(surveyID, currentUser(c).ID)
	if err != nil {
		c.JSON(StatusFor(err), r.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, r.OK(*result))
}

func EditSurveyInfoHandler(c *gin.Context) {
	r := &response.Response[struct{}]{}

	surveyID, ok := idParam(c, "id")
	if !ok {
		return
	}

	var req request.EditSurveyInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, r.Fail(err.Error()))
		return
	}

	err := EditSurveyInfo(surveyID, currentUser(c).ID, req.Title, req.Description, req.Minutes, req.CategoryId)
	if err != nil {
		c.JSON(StatusFor(err), r.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, r.StatusOK("Survey updated successfully"))
}

func PublishSurveyHandler(c *gin.Context) {
	r := &response.Response[struct{}]{}

	surveyID, ok := idParam(c, "id")
	if !ok {
		return
	}

	var req request.PublishSurveyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, r.Fail(err.Error()))
		return
	}

	err := PublishSurvey(surveyID, currentUser(c).ID, req.ExpectedAnswerCounts, req.PointsPerAnswer)
	if err != nil {
		c.JSON(StatusFor(err), r.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, r.StatusOK("Survey published successfully"))
}

func PauseSurveyHandler(c *gin.Context) {
	r := &response.Response[struct{}]{}

	surveyID, ok := idParam(c, "id")
	if !ok {
		return
	}

	var req request.PauseSurveyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, r.Fail(err.Error()))
		return
	}

	if err := SetPaused(surveyID, currentUser(c).ID, *req.Paused); err != nil {
		c.JSON(StatusFor(err), r.Fail(err.Error()))
		return
	}

	message := "Survey resumed successfully"
	if *req.Paused {
		message = "Survey paused successfully"
	}

	c.JSON(http.StatusOK, r.StatusOK(message))
}

func DeleteSurveyHandler(c *gin.Context) {
	r := &response.Response[struct{}]{}

	surveyID, ok := idParam(c, "id")
	if !ok {
		return
	}

	if err := DeleteSurvey(surveyID, currentUser(c).ID); err != nil {
		c.JSON(StatusFor(err), r.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, r.StatusOK("Survey deleted successfully"))
}

func ListSurveysHandler(c *gin.Context) {
	r := &response.Response[pagination.Page[surResponse.SurveySummaryResponse]]{}

	result, err := ListByCreator(currentUser(c).ID, statesParam(c), pagination.Parse(c))
	if err != nil {
		c.JSON(StatusFor(err), r.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, r.OK(*result))
}

func SurveyCountHandler(c *gin.Context) {
	r := &response.Response[map[string]int64]{}

	total, err := CountByCreator(currentUser(c).ID, statesParam(c))
	if err != nil {
		c.JSON(StatusFor(err), r.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, r.OK(map[string]int64{"count": total}))
}

func GetCreateSurveyFormHandler(c *gin.Context) {
	r := &response.Response[surResponse.GetCreateSurveyFormResponse]{}

	formData, err := GetCreateSurveyFormData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, r.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, r.OK(*formData))
}

func InterviewerHomeHandler(c *gin.Context) {
	r := &response.Response[surResponse.InterviewerHomeResponse]{}

	result, err := GetInterviewerHome(currentUser(c))
	if err != nil {
		c.JSON(StatusFor(err), r.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, r.OK(*result))
}

func InterviewerCompletedHandler(c *gin.Context) {
	r := &response.Response[pagination.Page[surResponse.SurveySummaryResponse]]{}

	result, err := ListByCreator(currentUser(c).ID, []string{StateCompleted}, pagination.Parse(c))
	if err != nil {
		c.JSON(StatusFor(err), r.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, r.OK(*result))
}

func AnalyticsHandler(c *gin.Context) {
	r := &response.Response[surResponse.AnalyticDataResponse]{}

	surveyID, ok := idParam(c, "id")
	if !ok {
		return
	}

	result, err := GetAnalyticData(surveyID, currentUser(c).ID)
	if err != nil {
		c.JSON(StatusFor(err), r.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, r.OK(*result))
}
