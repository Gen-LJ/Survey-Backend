package country

import (
	"net/http"
	"strconv"
	"strings"

	"survey-backend/internal/country/request"
	countryResponse "survey-backend/internal/country/response"
	"survey-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// ListCountriesHandler lists countries with their active flag, so an admin can
// see what the platform accepts and find the id of one to turn on without
// opening the database.
//
// ?active=true narrows it to the active ones, ?active=false to the rest, and
// omitting it returns every country.
func ListCountriesHandler(c *gin.Context) {
	r := &response.Response[[]countryResponse.AdminCountryResponse]{}

	var filter *bool
	if raw := strings.TrimSpace(c.Query("active")); raw != "" {
		active, err := strconv.ParseBool(raw)
		if err != nil {
			c.JSON(http.StatusBadRequest, r.Fail("active must be true or false"))
			return
		}
		filter = &active
	}

	countries, err := FindFiltered(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, r.Fail(err.Error()))
		return
	}

	list := make([]countryResponse.AdminCountryResponse, len(countries))
	for i, country := range countries {
		list[i] = country.ToAdminResponse()
	}

	c.JSON(http.StatusOK, r.OK(list))
}

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
