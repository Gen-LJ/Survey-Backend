package survey

import (
	"survey-backend/pkg/database"
)

func CreateSurveyAtRepo(
	survey *Survey,
) error {
	return database.DB.Create(survey).Error
}
