package survey

import (
	"errors"

	"survey-backend/internal/answer"
	"survey-backend/internal/question"
	"survey-backend/internal/savedsurvey"
	surResponse "survey-backend/internal/survey/response"
	"survey-backend/internal/user"
	"survey-backend/pkg/pagination"

	"gorm.io/gorm"
)

// RecentSurveyLimit is how many survey cards each home screen shows.
const RecentSurveyLimit = 10

// ListForRespondent returns a page of surveys this user can answer: published,
// still taking answers, targeted at their country and region, and not already
// answered by them.
func ListForRespondent(u user.User, page pagination.Query) (*pagination.Page[surResponse.SurveySummaryResponse], error) {
	answeredIDs, err := answer.AnsweredSurveyIDs(u.ID)
	if err != nil {
		return nil, err
	}

	countryID, regionID := u.CountryID, u.RegionID

	surveys, total, err := List(ListFilter{
		OnlyOpen:     true,
		CountryID:    &countryID,
		RegionID:     &regionID,
		ExcludeIDs:   answeredIDs,
		OrderByField: "published_at",
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

// GetSurveyForRespondent returns a survey with its questions, for answering.
// It refuses surveys that are closed or that the caller already answered, so
// the answer screen never opens on something that cannot be submitted.
func GetSurveyForRespondent(surveyID uint, u user.User) (*surResponse.SurveyResponse, error) {
	var s Survey
	if err := FindByID(surveyID, &s); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if !s.IsOpen() {
		return nil, ErrSurveyClosed
	}

	answered, err := answer.HasAnswered(s.ID, u.ID)
	if err != nil {
		return nil, err
	}
	if answered {
		return nil, ErrAlreadyAnswered
	}

	questions, err := question.FindBySurveyID(s.ID)
	if err != nil {
		return nil, err
	}

	result := s.ToResponse(question.ToResponseList(questions))
	return &result, nil
}

// SaveSurvey adds a survey to the respondent's reading list, capped at
// MaxSavedSurveys the way the Firestore implementation was.
func SaveSurvey(userID, surveyID uint) error {
	var s Survey
	if err := FindByID(surveyID, &s); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}

	if !s.IsOpen() {
		return ErrSurveyClosed
	}

	exists, err := savedsurvey.Exists(userID, surveyID)
	if err != nil {
		return err
	}
	if exists {
		return ErrAlreadySaved
	}

	count, err := savedsurvey.Count(userID)
	if err != nil {
		return err
	}
	if count >= savedsurvey.MaxSavedSurveys {
		return ErrSavedLimitReached
	}

	return savedsurvey.Create(&savedsurvey.SavedSurvey{UserID: userID, SurveyID: surveyID})
}

// RemoveSavedSurvey drops a survey from the respondent's reading list.
func RemoveSavedSurvey(userID, surveyID uint) error {
	exists, err := savedsurvey.Exists(userID, surveyID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrNotSaved
	}

	return savedsurvey.Delete(userID, surveyID)
}

// GetSavedSurveyIDs returns just the ids, for marking cards as saved in a list.
func GetSavedSurveyIDs(userID uint) ([]uint, error) {
	ids, err := savedsurvey.SurveyIDs(userID)
	if err != nil {
		return nil, err
	}
	if ids == nil {
		ids = []uint{}
	}
	return ids, nil
}

// GetSavedSurveys returns the saved surveys themselves. Entries that are no
// longer answerable are dropped from the list as they are found, which is what
// the Firestore version did for paused surveys.
func GetSavedSurveys(userID uint) ([]surResponse.SurveySummaryResponse, error) {
	ids, err := savedsurvey.SurveyIDs(userID)
	if err != nil {
		return nil, err
	}

	surveys, err := FindByIDs(ids)
	if err != nil {
		return nil, err
	}

	open := make([]Survey, 0, len(surveys))
	found := make(map[uint]bool, len(surveys))
	for _, s := range surveys {
		found[s.ID] = true
		if s.IsOpen() {
			open = append(open, s)
			continue
		}
		// Stale entry: prune it so the 5-survey cap is not held by dead rows.
		if err := savedsurvey.Delete(userID, s.ID); err != nil {
			return nil, err
		}
	}

	// Surveys deleted outright leave a dangling saved row behind.
	for _, id := range ids {
		if !found[id] {
			if err := savedsurvey.Delete(userID, id); err != nil {
				return nil, err
			}
		}
	}

	return summarize(open)
}

// ListCompletedForRespondent returns a page of the surveys this user answered.
func ListCompletedForRespondent(userID uint, sortByAnsweredAt bool, page pagination.Query) (*pagination.Page[surResponse.AnsweredSurveyResponse], error) {
	rows, total, err := ListAnsweredByRespondent(userID, sortByAnsweredAt, page)
	if err != nil {
		return nil, err
	}

	surveys := make([]Survey, len(rows))
	for i, row := range rows {
		surveys[i] = row.Survey
	}

	summaries, err := summarize(surveys)
	if err != nil {
		return nil, err
	}

	items := make([]surResponse.AnsweredSurveyResponse, len(summaries))
	for i, summary := range summaries {
		items[i] = surResponse.AnsweredSurveyResponse{
			SurveySummaryResponse: summary,
			AnsweredAt:            rows[i].AnsweredAt,
		}
	}

	result := pagination.NewPage(items, page, total)
	return &result, nil
}

// GetRespondentHome builds the respondent dashboard in one call.
func GetRespondentHome(u user.User) (*surResponse.RespondentHomeResponse, error) {
	availableCount, err := CountOpenForRespondent(u.ID, u.CountryID, u.RegionID)
	if err != nil {
		return nil, err
	}

	answeredCount, err := answer.CountByRespondent(u.ID)
	if err != nil {
		return nil, err
	}

	savedCount, err := savedsurvey.Count(u.ID)
	if err != nil {
		return nil, err
	}

	recent, err := ListForRespondent(u, pagination.Query{Page: 1, Limit: RecentSurveyLimit})
	if err != nil {
		return nil, err
	}

	return &surResponse.RespondentHomeResponse{
		Points:         u.Points,
		AvailableCount: availableCount,
		AnsweredCount:  answeredCount,
		SavedCount:     savedCount,
		RecentSurveys:  recent.Items,
	}, nil
}

// GetInterviewerHome builds the interviewer dashboard: one grouped count query
// plus the two recent lists, in place of Firestore's four count aggregations.
func GetInterviewerHome(u user.User) (*surResponse.InterviewerHomeResponse, error) {
	counts, err := CountsByState(u.ID)
	if err != nil {
		return nil, err
	}

	recentPublished, err := ListByCreator(u.ID, []string{StatePublished}, pagination.Query{Page: 1, Limit: RecentSurveyLimit})
	if err != nil {
		return nil, err
	}

	recentDrafts, err := ListByCreator(u.ID, []string{StateDraft}, pagination.Query{Page: 1, Limit: RecentSurveyLimit})
	if err != nil {
		return nil, err
	}

	return &surResponse.InterviewerHomeResponse{
		Points:          u.Points,
		DraftCount:      counts[StateDraft],
		PublishedCount:  counts[StatePublished],
		PausedCount:     counts[StatePaused],
		CompletedCount:  counts[StateCompleted],
		RecentPublished: recentPublished.Items,
		RecentDrafts:    recentDrafts.Items,
	}, nil
}
