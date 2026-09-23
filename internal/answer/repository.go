package answer

import (
	"survey-backend/pkg/database"

	"gorm.io/gorm"
)

// FindBySurveyID returns every submission for a survey, used to build analytics.
func FindBySurveyID(surveyID uint) ([]Answer, error) {
	var answers []Answer
	err := database.DB.
		Preload("UserAnswers").
		Where("survey_id = ?", surveyID).
		Order("created_at asc").
		Find(&answers).Error

	return answers, err
}

// FindBySurveyAndRespondent returns one respondent's submission for a survey.
func FindBySurveyAndRespondent(surveyID, respondentID uint, a *Answer) error {
	return database.DB.
		Preload("UserAnswers").
		Where("survey_id = ? AND respondent_id = ?", surveyID, respondentID).
		First(a).Error
}

// AnsweredSurveyIDs returns the ids of every survey this user has answered,
// replacing the answerList/{uid}.surveyList document.
func AnsweredSurveyIDs(respondentID uint) ([]uint, error) {
	var ids []uint
	err := database.DB.Model(&Answer{}).
		Where("respondent_id = ?", respondentID).
		Pluck("survey_id", &ids).Error

	return ids, err
}

func CountByRespondent(respondentID uint) (int64, error) {
	var total int64
	err := database.DB.Model(&Answer{}).
		Where("respondent_id = ?", respondentID).
		Count(&total).Error

	return total, err
}

func CountBySurveyID(surveyID uint) (int64, error) {
	var total int64
	err := database.DB.Model(&Answer{}).
		Where("survey_id = ?", surveyID).
		Count(&total).Error

	return total, err
}

func HasAnswered(surveyID, respondentID uint) (bool, error) {
	var total int64
	err := database.DB.Model(&Answer{}).
		Where("survey_id = ? AND respondent_id = ?", surveyID, respondentID).
		Count(&total).Error

	return total > 0, err
}

// CreateWithUserAnswers persists a submission and its per-question rows inside
// the caller's transaction, so it commits together with the survey counter and
// the point transfer.
func CreateWithUserAnswers(tx *gorm.DB, a *Answer) error {
	return tx.Create(a).Error
}

// DeleteBySurveyID removes every submission for a survey, used when the survey
// itself is deleted.
func DeleteBySurveyID(tx *gorm.DB, surveyID uint) error {
	var ids []uint
	if err := tx.Model(&Answer{}).
		Where("survey_id = ?", surveyID).
		Pluck("id", &ids).Error; err != nil {
		return err
	}

	if len(ids) > 0 {
		if err := tx.Unscoped().Where("answer_id IN ?", ids).Delete(&UserAnswer{}).Error; err != nil {
			return err
		}
	}

	return tx.Unscoped().Where("survey_id = ?", surveyID).Delete(&Answer{}).Error
}
