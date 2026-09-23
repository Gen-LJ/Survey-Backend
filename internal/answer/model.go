package answer

import (
	"encoding/json"

	"survey-backend/internal/answer/response"

	"gorm.io/gorm"
)

// Answer is one respondent's submission for one survey. The composite unique
// index is what enforces "a user answers a survey at most once" - the rule the
// Firestore version had to check by hand inside a transaction.
type Answer struct {
	gorm.Model
	SurveyID     uint         `gorm:"not null;uniqueIndex:idx_answer_survey_respondent" json:"survey_id"`
	RespondentID uint         `gorm:"not null;uniqueIndex:idx_answer_survey_respondent" json:"respondent_id"`
	PointsEarned uint         `gorm:"default:0" json:"points_earned"`
	UserAnswers  []UserAnswer `gorm:"foreignKey:AnswerID;constraint:OnDelete:CASCADE" json:"user_answers"`
}

// UserAnswer is the response to a single question. Responses is stored as a
// JSON array so one shape covers text input, single choice, multi choice and
// rating alike.
type UserAnswer struct {
	gorm.Model
	AnswerID   uint   `gorm:"index;not null" json:"answer_id"`
	QuestionID uint   `gorm:"index;not null" json:"question_id"`
	Responses  string `gorm:"type:text" json:"-"`
}

func (ua *UserAnswer) SetResponses(responses []string) error {
	if responses == nil {
		responses = []string{}
	}

	encoded, err := json.Marshal(responses)
	if err != nil {
		return err
	}

	ua.Responses = string(encoded)
	return nil
}

func (ua UserAnswer) GetResponses() []string {
	if ua.Responses == "" {
		return []string{}
	}

	var responses []string
	if err := json.Unmarshal([]byte(ua.Responses), &responses); err != nil {
		// A row written outside this API; surface it rather than dropping it.
		return []string{ua.Responses}
	}

	return responses
}

func (ua UserAnswer) ToResponse() response.UserAnswerResponse {
	return response.UserAnswerResponse{
		QuestionID: ua.QuestionID,
		Responses:  ua.GetResponses(),
	}
}

func (a Answer) ToResponse() response.AnswerResponse {
	userAnswers := make([]response.UserAnswerResponse, len(a.UserAnswers))
	for i, ua := range a.UserAnswers {
		userAnswers[i] = ua.ToResponse()
	}

	return response.AnswerResponse{
		ID:           a.ID,
		SurveyID:     a.SurveyID,
		RespondentID: a.RespondentID,
		PointsEarned: a.PointsEarned,
		AnsweredAt:   a.CreatedAt,
		UserAnswers:  userAnswers,
	}
}

func ToResponseList(answers []Answer) []response.AnswerResponse {
	list := make([]response.AnswerResponse, len(answers))
	for i, a := range answers {
		list[i] = a.ToResponse()
	}
	return list
}
