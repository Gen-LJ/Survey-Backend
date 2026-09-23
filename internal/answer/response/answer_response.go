package response

import "time"

type UserAnswerResponse struct {
	QuestionID uint     `json:"question_id"`
	Responses  []string `json:"responses"`
}

type AnswerResponse struct {
	ID           uint                 `json:"id"`
	SurveyID     uint                 `json:"survey_id"`
	RespondentID uint                 `json:"respondent_id"`
	PointsEarned uint                 `json:"points_earned"`
	AnsweredAt   time.Time            `json:"answered_at"`
	UserAnswers  []UserAnswerResponse `json:"user_answers"`
}
