package question

import (
	"survey-backend/pkg/database"

	"gorm.io/gorm"
)

// FindBySurveyID returns a survey's questions in author-defined order,
// options included.
func FindBySurveyID(surveyID uint) ([]Question, error) {
	var questions []Question
	err := database.DB.
		Preload("Options", func(db *gorm.DB) *gorm.DB {
			return db.Order("options.position asc, options.id asc")
		}).
		Where("survey_id = ?", surveyID).
		Order("position asc, id asc").
		Find(&questions).Error

	return questions, err
}

// FindBySurveyIDs loads the questions for several surveys at once, grouped by
// survey id, so a list endpoint does not issue one query per survey.
func FindBySurveyIDs(surveyIDs []uint) (map[uint][]Question, error) {
	grouped := make(map[uint][]Question, len(surveyIDs))
	if len(surveyIDs) == 0 {
		return grouped, nil
	}

	var questions []Question
	err := database.DB.
		Preload("Options", func(db *gorm.DB) *gorm.DB {
			return db.Order("options.position asc, options.id asc")
		}).
		Where("survey_id IN ?", surveyIDs).
		Order("position asc, id asc").
		Find(&questions).Error
	if err != nil {
		return nil, err
	}

	for _, q := range questions {
		grouped[q.SurveyID] = append(grouped[q.SurveyID], q)
	}

	return grouped, nil
}

func FindByID(id uint, q *Question) error {
	return database.DB.
		Preload("Options", func(db *gorm.DB) *gorm.DB {
			return db.Order("options.position asc, options.id asc")
		}).
		First(q, id).Error
}

// CountBySurveyIDs returns a question count per survey id.
func CountBySurveyIDs(surveyIDs []uint) (map[uint]int64, error) {
	counts := make(map[uint]int64, len(surveyIDs))
	if len(surveyIDs) == 0 {
		return counts, nil
	}

	var rows []struct {
		SurveyID uint
		Total    int64
	}
	err := database.DB.Model(&Question{}).
		Select("survey_id, count(*) as total").
		Where("survey_id IN ?", surveyIDs).
		Group("survey_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		counts[row.SurveyID] = row.Total
	}

	return counts, nil
}

func CountBySurveyID(surveyID uint) (int64, error) {
	var total int64
	err := database.DB.Model(&Question{}).Where("survey_id = ?", surveyID).Count(&total).Error
	return total, err
}

// NextPosition returns the position to give the next question appended to a
// survey, preserving creation order the way Firestore's createdAt did.
func NextPosition(surveyID uint) (uint, error) {
	var max *uint
	err := database.DB.Model(&Question{}).
		Select("max(position)").
		Where("survey_id = ?", surveyID).
		Scan(&max).Error
	if err != nil || max == nil {
		return 0, err
	}
	return *max + 1, nil
}

func Create(q *Question) error {
	return database.DB.Create(q).Error
}

// UpdateWithOptions rewrites a question's text and, when options is non-nil,
// replaces its option list in a single transaction.
func UpdateWithOptions(questionID uint, text string, options []string) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Question{}).
			Where("id = ?", questionID).
			Update("text", text).Error; err != nil {
			return err
		}

		if options == nil {
			return nil
		}

		if err := tx.Unscoped().Where("question_id = ?", questionID).Delete(&Option{}).Error; err != nil {
			return err
		}

		if len(options) == 0 {
			return nil
		}

		rows := make([]Option, len(options))
		for i, text := range options {
			rows[i] = Option{QuestionID: questionID, Text: text, Position: uint(i)}
		}

		return tx.Create(&rows).Error
	})
}

// Delete removes a question and its options.
func Delete(questionID uint) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().Where("question_id = ?", questionID).Delete(&Option{}).Error; err != nil {
			return err
		}
		return tx.Unscoped().Delete(&Question{}, questionID).Error
	})
}

// DeleteBySurveyID removes every question (and option) belonging to a survey.
// It takes a transaction because it runs as part of deleting the survey.
func DeleteBySurveyID(tx *gorm.DB, surveyID uint) error {
	var ids []uint
	if err := tx.Model(&Question{}).
		Where("survey_id = ?", surveyID).
		Pluck("id", &ids).Error; err != nil {
		return err
	}

	if len(ids) > 0 {
		if err := tx.Unscoped().Where("question_id IN ?", ids).Delete(&Option{}).Error; err != nil {
			return err
		}
	}

	return tx.Unscoped().Where("survey_id = ?", surveyID).Delete(&Question{}).Error
}
