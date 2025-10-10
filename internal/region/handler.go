package region

import (
	"net/http"
	"strconv"
	regionResponse "survey-backend/internal/region/response"
	"survey-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

func GetRegionsByCountryHandler(c *gin.Context) {
	r := &response.Response[[]regionResponse.RegionResponse]{}

	// Get country ID from URL parameter
	countryIDStr := c.Param("country_id")
	countryID, err := strconv.ParseUint(countryIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, r.Fail("Invalid country ID"))
		return
	}

	// Get only active regions
	regions, err := FindActiveByCountryID(uint(countryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, r.Fail(err.Error()))
		return
	}

	// Convert to response DTOs
	regionResponses := make([]regionResponse.RegionResponse, len(regions))
	for i, reg := range regions {
		regionResponses[i] = reg.ToResponse()
	}

	c.JSON(http.StatusOK, r.OK(regionResponses))
}
