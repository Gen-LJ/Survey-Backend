package survey

import (
	"net/http"
	"survey-backend/internal/survey/request"
	surResponse "survey-backend/internal/survey/response"
	"survey-backend/internal/user"
	"survey-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

func CreateSurveyHandler(c *gin.Context) {
	r := &response.Response[Survey]{}

	var req request.CreateSurveyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, r.Fail(err.Error()))
		return
	}

	currentUser := c.MustGet("currentUser").(*user.User)

	survey := &Survey{
		Title:       req.Title,
		Description: req.Description,
		CategoryID:  req.CategoryId,
		CountryID:   req.CountryId,
		RegionID:    req.RegionId,
		CreatorID:   currentUser.ID,
	}

	err := CreateSurvey(survey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, r.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, r.OK(*survey))
}

func GetCreateSurveyFormHandler(c *gin.Context) {
	r := &response.Response[surResponse.GetCreateSurveyFormResponse]{}

	formData, err := GetCreateSurveyFormData()
	if err != nil {
		c.JSON(http.StatusNoContent, r.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, r.OK(*formData))

}
