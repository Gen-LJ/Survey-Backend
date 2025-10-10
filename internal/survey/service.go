package survey

import (
	"survey-backend/internal/category"
	catResponse "survey-backend/internal/category/response"
	"survey-backend/internal/country"
	countryResponse "survey-backend/internal/country/response"
	surResponse "survey-backend/internal/survey/response"
)

func CreateSurvey(
	survey *Survey,
) error {
	return CreateSurveyAtRepo(survey)
}

func GetCreateSurveyFormData() (*surResponse.GetCreateSurveyFormResponse, error) {
	// Get categories
	categories, err := category.FindAll()
	if err != nil {
		return nil, err
	}

	// Get countries
	countries, err := country.FindActive()
	if err != nil {
		return nil, err
	}

	// Convert categories to response DTOs
	categoryResponses := make([]catResponse.CategoryResponse, len(categories))
	for i, cat := range categories {
		categoryResponses[i] = cat.ToResponse()
	}

	// Convert countries to response DTOs
	countryResponses := make([]countryResponse.CountryResponse, len(countries))
	for i, cntry := range countries {
		countryResponses[i] = cntry.ToResponse()
	}

	return &surResponse.GetCreateSurveyFormResponse{
		CategoryList: categoryResponses,
		CountryList:  countryResponses,
	}, nil
}
