package response

import (
	catResponse "survey-backend/internal/category/response"
	countResponse "survey-backend/internal/country/response"
)

type GetCreateSurveyFormResponse struct {
	CategoryList []catResponse.CategoryResponse  `json:"category_list"`
	CountryList  []countResponse.CountryResponse `json:"country_list"`
}
