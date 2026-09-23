package survey

import (
	"errors"
	"strconv"
	"time"

	"survey-backend/internal/answer"
	answerResponse "survey-backend/internal/answer/response"
	"survey-backend/internal/question"
	"survey-backend/internal/savedsurvey"
	surResponse "survey-backend/internal/survey/response"
	"survey-backend/internal/user"
	"survey-backend/pkg/database"

	"gorm.io/gorm"
)

// SubmittedAnswer is one question's response as it arrives from a respondent.
type SubmittedAnswer struct {
	QuestionID uint
	Responses  []string
}

// SubmitAnswer records a respondent's submission. Everything that has to agree
// - the answer rows, the survey's counter and escrow, the respondent's balance,
// the saved list and the completion flip - happens in one transaction, with the
// survey row locked so two respondents cannot both take the last slot.
func SubmitAnswer(surveyID, respondentID uint, submitted []SubmittedAnswer) (*answerResponse.AnswerResponse, error) {
	var created answer.Answer

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var s Survey
		if err := FindByIDForUpdate(tx, surveyID, &s); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}

		if !s.IsOpen() {
			return ErrSurveyClosed
		}

		answered, err := answer.HasAnswered(s.ID, respondentID)
		if err != nil {
			return err
		}
		if answered {
			return ErrAlreadyAnswered
		}

		questions, err := question.FindBySurveyID(s.ID)
		if err != nil {
			return err
		}

		rows, err := buildUserAnswers(questions, submitted)
		if err != nil {
			return err
		}

		created = answer.Answer{
			SurveyID:     s.ID,
			RespondentID: respondentID,
			PointsEarned: s.PointsPerAnswer,
			UserAnswers:  rows,
		}
		if err := answer.CreateWithUserAnswers(tx, &created); err != nil {
			return err
		}

		// Draw the reward out of the survey's escrow and into the respondent's
		// balance, so points are never created or destroyed by answering.
		payout := s.PointsPerAnswer
		if payout > s.PendingPoints {
			payout = s.PendingPoints
		}
		if payout > 0 {
			if err := user.AddPoints(tx, respondentID, payout); err != nil {
				return err
			}
		}

		fields := map[string]any{
			"answer_count":   s.AnswerCount + 1,
			"pending_points": s.PendingPoints - payout,
		}
		if s.AnswerCount+1 >= s.ExpectedAnswerCounts {
			fields["state"] = StateCompleted
			fields["completed_at"] = time.Now()
		}

		if err := tx.Model(&Survey{}).Where("id = ?", s.ID).Updates(fields).Error; err != nil {
			return err
		}

		// Answering clears the survey off the respondent's reading list.
		return savedsurvey.DeleteInTx(tx, respondentID, s.ID)
	})
	if err != nil {
		return nil, err
	}

	result := created.ToResponse()
	return &result, nil
}

// buildUserAnswers validates a submission against the survey's questions and
// turns it into rows ready to insert. The Firestore version wrote whatever the
// client sent; this refuses anything the question cannot accept.
func buildUserAnswers(questions []question.Question, submitted []SubmittedAnswer) ([]answer.UserAnswer, error) {
	if len(questions) == 0 {
		return nil, ErrNoQuestions
	}

	byQuestionID := make(map[uint]question.Question, len(questions))
	for _, q := range questions {
		byQuestionID[q.ID] = q
	}

	seen := make(map[uint]bool, len(submitted))
	rows := make([]answer.UserAnswer, 0, len(submitted))

	for _, item := range submitted {
		q, ok := byQuestionID[item.QuestionID]
		if !ok {
			return nil, ErrInvalidAnswer
		}
		if seen[item.QuestionID] {
			return nil, ErrInvalidAnswer
		}
		seen[item.QuestionID] = true

		if err := validateResponses(q, item.Responses); err != nil {
			return nil, err
		}

		row := answer.UserAnswer{QuestionID: q.ID}
		if err := row.SetResponses(item.Responses); err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}

	// Every question must be accounted for; a partial submission would skew
	// the per-question analytics.
	if len(seen) != len(questions) {
		return nil, ErrInvalidAnswer
	}

	return rows, nil
}

func validateResponses(q question.Question, responses []string) error {
	if len(responses) == 0 {
		return ErrInvalidAnswer
	}

	switch q.QuestionType {
	case question.TypeTextInput:
		if len(responses) != 1 || responses[0] == "" {
			return ErrInvalidAnswer
		}

	case question.TypeRating:
		if len(responses) != 1 {
			return ErrInvalidAnswer
		}
		rating, err := strconv.Atoi(responses[0])
		if err != nil || rating < 1 || rating > 5 {
			return ErrInvalidAnswer
		}

	case question.TypeMultipleChoice:
		if !q.AllowMultiAnswer && len(responses) != 1 {
			return ErrInvalidAnswer
		}

		allowed := make(map[string]bool, len(q.Options))
		for _, opt := range q.Options {
			allowed[opt.Text] = true
		}

		picked := make(map[string]bool, len(responses))
		for _, r := range responses {
			if !allowed[r] || picked[r] {
				return ErrInvalidAnswer
			}
			picked[r] = true
		}

	default:
		return ErrInvalidQuestionType
	}

	return nil
}

// GetAnswerCount reports how many submissions a published survey has taken.
func GetAnswerCount(surveyID uint) (uint, error) {
	var s Survey
	if err := FindByID(surveyID, &s); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrNotFound
		}
		return 0, err
	}

	return s.AnswerCount, nil
}

// GetAnswerDataDetails returns one respondent's own submission alongside the
// survey it answered, for the "my answers" detail screen.
func GetAnswerDataDetails(surveyID, respondentID uint) (*surResponse.AnswerDataDetailsResponse, error) {
	var s Survey
	if err := FindByID(surveyID, &s); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	var a answer.Answer
	if err := answer.FindBySurveyAndRespondent(surveyID, respondentID, &a); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotAnswered
		}
		return nil, err
	}

	questions, err := question.FindBySurveyID(surveyID)
	if err != nil {
		return nil, err
	}

	return &surResponse.AnswerDataDetailsResponse{
		Survey: s.ToResponse(question.ToResponseList(questions)),
		Answer: a.ToResponse(),
	}, nil
}

// GetAnalyticData returns a survey, its questions and every submission, plus a
// per-question tally so the client does not have to aggregate thousands of rows
// itself the way it did against Firestore.
func GetAnalyticData(surveyID, creatorID uint) (*surResponse.AnalyticDataResponse, error) {
	s, err := loadOwned(surveyID, creatorID)
	if err != nil {
		return nil, err
	}

	questions, err := question.FindBySurveyID(surveyID)
	if err != nil {
		return nil, err
	}

	answers, err := answer.FindBySurveyID(surveyID)
	if err != nil {
		return nil, err
	}

	return &surResponse.AnalyticDataResponse{
		Survey:      s.ToResponse(question.ToResponseList(questions)),
		Answers:     answer.ToResponseList(answers),
		Summary:     summarizeAnswers(questions, answers),
		AnswerTotal: len(answers),
	}, nil
}

// summarizeAnswers tallies responses per question: option counts for choice
// questions, a 1-5 histogram plus average for ratings, and the raw text for
// free-text questions.
func summarizeAnswers(questions []question.Question, answers []answer.Answer) []surResponse.QuestionSummary {
	responsesByQuestion := make(map[uint][][]string, len(questions))
	for _, a := range answers {
		for _, ua := range a.UserAnswers {
			responsesByQuestion[ua.QuestionID] = append(responsesByQuestion[ua.QuestionID], ua.GetResponses())
		}
	}

	summaries := make([]surResponse.QuestionSummary, len(questions))
	for i, q := range questions {
		collected := responsesByQuestion[q.ID]

		summary := surResponse.QuestionSummary{
			QuestionID:    q.ID,
			Text:          q.Text,
			QuestionType:  q.QuestionType,
			ResponseCount: len(collected),
			OptionCounts:  []surResponse.OptionCount{},
			TextResponses: []string{},
		}

		switch q.QuestionType {
		case question.TypeMultipleChoice:
			tally := make(map[string]int, len(q.Options))
			for _, responses := range collected {
				for _, r := range responses {
					tally[r]++
				}
			}
			// Report in the author's option order, including unpicked options.
			for _, opt := range q.Options {
				summary.OptionCounts = append(summary.OptionCounts, surResponse.OptionCount{
					Option: opt.Text,
					Count:  tally[opt.Text],
				})
			}

		case question.TypeRating:
			histogram := make([]int, 5)
			total, rated := 0, 0
			for _, responses := range collected {
				if len(responses) == 0 {
					continue
				}
				rating, err := strconv.Atoi(responses[0])
				if err != nil || rating < 1 || rating > 5 {
					continue
				}
				histogram[rating-1]++
				total += rating
				rated++
			}
			summary.RatingCounts = histogram
			if rated > 0 {
				summary.AverageRating = float64(total) / float64(rated)
			}

		case question.TypeTextInput:
			for _, responses := range collected {
				summary.TextResponses = append(summary.TextResponses, responses...)
			}
		}

		summaries[i] = summary
	}

	return summaries
}
