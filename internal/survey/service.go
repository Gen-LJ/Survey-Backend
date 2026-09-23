package survey

import (
	"errors"
	"time"

	"survey-backend/internal/answer"
	"survey-backend/internal/category"
	catResponse "survey-backend/internal/category/response"
	"survey-backend/internal/country"
	countryResponse "survey-backend/internal/country/response"
	"survey-backend/internal/question"
	"survey-backend/internal/region"
	"survey-backend/internal/savedsurvey"
	surResponse "survey-backend/internal/survey/response"
	"survey-backend/internal/user"
	"survey-backend/pkg/database"
	"survey-backend/pkg/pagination"

	"gorm.io/gorm"
)

// loadOwned fetches a survey and checks the caller created it. Every
// interviewer-side action starts here, so ownership is never assumed.
func loadOwned(surveyID, creatorID uint) (Survey, error) {
	var s Survey
	if err := FindByID(surveyID, &s); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Survey{}, ErrNotFound
		}
		return Survey{}, err
	}

	if s.CreatorID != creatorID {
		return Survey{}, ErrForbidden
	}

	return s, nil
}

// CreateSurvey stores a new draft after checking its category and location are
// ones the platform currently accepts.
func CreateSurvey(survey *Survey) error {
	if _, err := category.FindByID(survey.CategoryID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUnknownCategory
		}
		return err
	}

	regions, err := region.FindActiveByCountryID(survey.CountryID)
	if err != nil {
		return err
	}

	targeted := false
	for _, r := range regions {
		if r.ID == survey.RegionID {
			targeted = true
			break
		}
	}
	if !targeted {
		return ErrRegionNotTargetable
	}

	survey.State = StateDraft
	survey.AnswerCount = 0
	survey.PendingPoints = 0

	return CreateSurveyAtRepo(survey)
}

// GetSurveyForCreator returns one of the caller's own surveys, questions included.
func GetSurveyForCreator(surveyID, creatorID uint) (*surResponse.SurveyResponse, error) {
	s, err := loadOwned(surveyID, creatorID)
	if err != nil {
		return nil, err
	}

	questions, err := question.FindBySurveyID(s.ID)
	if err != nil {
		return nil, err
	}

	result := s.ToResponse(question.ToResponseList(questions))
	return &result, nil
}

// EditSurveyInfo updates the descriptive fields. A finished survey is frozen so
// its analytics keep describing what respondents actually saw.
func EditSurveyInfo(surveyID, creatorID uint, title, description string, minutes, categoryID uint) error {
	s, err := loadOwned(surveyID, creatorID)
	if err != nil {
		return err
	}

	if s.State == StateCompleted {
		return ErrNotEditable
	}

	if _, err := category.FindByID(categoryID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUnknownCategory
		}
		return err
	}

	return UpdateFields(s.ID, map[string]any{
		"title":       title,
		"description": description,
		"minutes":     minutes,
		"category_id": categoryID,
	})
}

// PublishSurvey moves a draft live and escrows the points it promises. The
// creator's balance is locked for the check-then-deduct so two concurrent
// publishes cannot both pass the affordability test.
func PublishSurvey(surveyID, creatorID uint, expectedAnswers, pointsPerAnswer uint) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		var s Survey
		if err := FindByIDForUpdate(tx, surveyID, &s); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}

		if s.CreatorID != creatorID {
			return ErrForbidden
		}
		if s.State != StateDraft {
			return ErrAlreadyPublished
		}

		questionCount, err := question.CountBySurveyID(s.ID)
		if err != nil {
			return err
		}
		if questionCount == 0 {
			return ErrNoQuestions
		}

		total := pointsPerAnswer * expectedAnswers

		var creator user.User
		if err := user.FindUserByIdForUpdate(tx, creatorID, &creator); err != nil {
			return err
		}
		if creator.Points < total {
			return ErrInsufficientPoints
		}

		if err := user.UpdatePoints(tx, creator.ID, creator.Points-total); err != nil {
			return err
		}

		now := time.Now()
		return tx.Model(&Survey{}).Where("id = ?", s.ID).Updates(map[string]any{
			"expected_answer_counts": expectedAnswers,
			"points_per_answer":      pointsPerAnswer,
			"pending_points":         total,
			"state":                  StatePublished,
			"published_at":           now,
		}).Error
	})
}

// SetPaused hides a live survey from respondents, or brings it back. A survey
// that filled up while paused comes back as completed rather than published.
func SetPaused(surveyID, creatorID uint, paused bool) error {
	s, err := loadOwned(surveyID, creatorID)
	if err != nil {
		return err
	}

	if paused {
		if s.State != StatePublished {
			return ErrNotPublished
		}
		return UpdateFields(s.ID, map[string]any{"state": StatePaused})
	}

	if s.State != StatePaused {
		return ErrNotPublished
	}

	next := StatePublished
	fields := map[string]any{}
	if s.AnswerCount >= s.ExpectedAnswerCounts {
		next = StateCompleted
		now := time.Now()
		fields["completed_at"] = now
	}
	fields["state"] = next

	return UpdateFields(s.ID, fields)
}

// DeleteSurvey removes a survey and everything hanging off it, refunding the
// escrow that was never paid out to respondents.
func DeleteSurvey(surveyID, creatorID uint) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		var s Survey
		if err := FindByIDForUpdate(tx, surveyID, &s); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}

		if s.CreatorID != creatorID {
			return ErrForbidden
		}

		if s.PendingPoints > 0 {
			if err := user.AddPoints(tx, s.CreatorID, s.PendingPoints); err != nil {
				return err
			}
		}

		if err := question.DeleteBySurveyID(tx, s.ID); err != nil {
			return err
		}
		if err := answer.DeleteBySurveyID(tx, s.ID); err != nil {
			return err
		}
		if err := savedsurvey.DeleteBySurveyID(tx, s.ID); err != nil {
			return err
		}

		return tx.Delete(&Survey{}, s.ID).Error
	})
}

// ListByCreator returns a page of the caller's surveys, newest activity first.
// States is optional; an empty slice means every state.
func ListByCreator(creatorID uint, states []string, page pagination.Query) (*pagination.Page[surResponse.SurveySummaryResponse], error) {
	for _, state := range states {
		if !IsValidState(state) {
			return nil, ErrInvalidState
		}
	}

	surveys, total, err := List(ListFilter{
		CreatorID:    &creatorID,
		States:       states,
		OrderByField: "updated_at",
		Descending:   true,
	}, page)
	if err != nil {
		return nil, err
	}

	summaries, err := summarize(surveys)
	if err != nil {
		return nil, err
	}

	result := pagination.NewPage(summaries, page, total)
	return &result, nil
}

// CountByCreator counts the caller's surveys in the given states.
func CountByCreator(creatorID uint, states []string) (int64, error) {
	for _, state := range states {
		if !IsValidState(state) {
			return 0, ErrInvalidState
		}
	}

	return Count(ListFilter{CreatorID: &creatorID, States: states})
}

// summarize attaches question counts to a page of surveys with one extra query.
func summarize(surveys []Survey) ([]surResponse.SurveySummaryResponse, error) {
	ids := make([]uint, len(surveys))
	for i, s := range surveys {
		ids[i] = s.ID
	}

	counts, err := question.CountBySurveyIDs(ids)
	if err != nil {
		return nil, err
	}

	return ToSummaryList(surveys, counts), nil
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
