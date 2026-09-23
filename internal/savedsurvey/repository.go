package savedsurvey

import (
	"survey-backend/pkg/database"

	"gorm.io/gorm"
)

// SurveyIDs returns the surveys a user has saved, newest first.
func SurveyIDs(userID uint) ([]uint, error) {
	var ids []uint
	err := database.DB.Model(&SavedSurvey{}).
		Where("user_id = ?", userID).
		Order("created_at desc").
		Pluck("survey_id", &ids).Error

	return ids, err
}

func Count(userID uint) (int64, error) {
	var total int64
	err := database.DB.Model(&SavedSurvey{}).
		Where("user_id = ?", userID).
		Count(&total).Error

	return total, err
}

func Exists(userID, surveyID uint) (bool, error) {
	var total int64
	err := database.DB.Model(&SavedSurvey{}).
		Where("user_id = ? AND survey_id = ?", userID, surveyID).
		Count(&total).Error

	return total > 0, err
}

func Create(saved *SavedSurvey) error {
	return database.DB.Create(saved).Error
}

func Delete(userID, surveyID uint) error {
	return database.DB.
		Where("user_id = ? AND survey_id = ?", userID, surveyID).
		Delete(&SavedSurvey{}).Error
}

// DeleteInTx removes one saved entry as part of a wider transaction - used when
// answering a survey clears it from the reading list.
func DeleteInTx(tx *gorm.DB, userID, surveyID uint) error {
	return tx.
		Where("user_id = ? AND survey_id = ?", userID, surveyID).
		Delete(&SavedSurvey{}).Error
}

// DeleteBySurveyID drops a survey from every reading list that holds it.
func DeleteBySurveyID(tx *gorm.DB, surveyID uint) error {
	return tx.Where("survey_id = ?", surveyID).Delete(&SavedSurvey{}).Error
}
