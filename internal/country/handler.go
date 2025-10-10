package country

import (
	"net/http"
	"survey-backend/internal/country/request"
	"survey-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

func ToggleCountryActiveHandler(c *gin.Context) {
	r := &response.Response[struct{}]{}

	var req request.ToggleActiveRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, r.Fail(err.Error()))
		return
	}

	if err := ToggleActiveByID(req.CountryID, *req.Active); err != nil {
		c.JSON(http.StatusInternalServerError, r.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, r.StatusOK(
		"Country status updated successfully",
	))
}
