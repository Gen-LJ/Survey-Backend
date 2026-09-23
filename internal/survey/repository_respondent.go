package survey

import (
	"time"

	"survey-backend/internal/answer"
	"survey-backend/pkg/database"
	"survey-backend/pkg/pagination"
)

// AnsweredSurvey pairs a survey with when this respondent answered it.
type AnsweredSurvey struct {
	Survey
	AnsweredAt time.Time
}

// ListAnsweredByRespondent returns a page of the surveys a respondent has
// answered. It joins through the answers table so paging happens in SQL rather
// than by fetching every answered id first.
func ListAnsweredByRespondent(respondentID uint, sortByAnsweredAt bool, page pagination.Query) ([]AnsweredSurvey, int64, error) {
	base := database.DB.Model(&Survey{}).
		Joins("JOIN answers ON answers.survey_id = surveys.id AND answers.deleted_at IS NULL").
		Where("answers.respondent_id = ?", respondentID)

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	order := "surveys.updated_at desc"
	if sortByAnsweredAt {
		order = "answers.created_at desc"
	}

	var rows []AnsweredSurvey
	err := database.DB.Model(&Survey{}).
		Select("surveys.*, answers.created_at as answered_at").
		Joins("JOIN answers ON answers.survey_id = surveys.id AND answers.deleted_at IS NULL").
		Where("answers.respondent_id = ?", respondentID).
		Order(order).
		Order("surveys.id desc").
		Offset(page.Offset()).
		Limit(page.Limit).
		Scan(&rows).Error
	if err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

// CountOpenForRespondent counts the surveys this respondent could answer right
// now: published, not yet full, matching their location, minus the ones they
// have already answered.
func CountOpenForRespondent(respondentID, countryID, regionID uint) (int64, error) {
	var total int64
	err := database.DB.Model(&Survey{}).
		Where("state = ?", StatePublished).
		Where("answer_count < expected_answer_counts").
		Where("country_id = ? AND region_id = ?", countryID, regionID).
		Where("id NOT IN (?)",
			database.DB.Model(&answer.Answer{}).
				Select("survey_id").
				Where("respondent_id = ? AND deleted_at IS NULL", respondentID),
		).
		Count(&total).Error

	return total, err
}
