package response

import (
	answerResponse "survey-backend/internal/answer/response"
)

type OptionCount struct {
	Option string `json:"option"`
	Count  int    `json:"count"`
}

// QuestionSummary is a server-side tally of one question's responses. Only the
// fields that apply to the question's type are populated.
type QuestionSummary struct {
	QuestionID    uint          `json:"question_id"`
	Text          string        `json:"text"`
	QuestionType  string        `json:"question_type"`
	ResponseCount int           `json:"response_count"`
	OptionCounts  []OptionCount `json:"option_counts"`
	RatingCounts  []int         `json:"rating_counts,omitempty"`
	AverageRating float64       `json:"average_rating,omitempty"`
	TextResponses []string      `json:"text_responses"`
}

type AnalyticDataResponse struct {
	Survey      SurveyResponse                  `json:"survey"`
	Answers     []answerResponse.AnswerResponse `json:"answers"`
	Summary     []QuestionSummary               `json:"summary"`
	AnswerTotal int                             `json:"answer_total"`
}

type AnswerDataDetailsResponse struct {
	Survey SurveyResponse                `json:"survey"`
	Answer answerResponse.AnswerResponse `json:"answer"`
}
