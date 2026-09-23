package response

import (
	"time"

	questionResponse "survey-backend/internal/question/response"
)

// SurveySummaryResponse is the list-view shape: everything a card needs,
// without the question bodies.
type SurveySummaryResponse struct {
	ID                   uint       `json:"id"`
	Title                string     `json:"title"`
	Description          string     `json:"description"`
	CategoryID           uint       `json:"category_id"`
	CountryID            uint       `json:"country_id"`
	RegionID             uint       `json:"region_id"`
	CreatorID            uint       `json:"creator_id"`
	Minutes              uint       `json:"minutes"`
	ExpectedAnswerCounts uint       `json:"expected_answer_counts"`
	PointsPerAnswer      uint       `json:"points_per_answer"`
	PendingPoints        uint       `json:"pending_points"`
	TotalPoints          uint       `json:"total_points"`
	AnswerCount          uint       `json:"answer_count"`
	QuestionCount        int        `json:"question_count"`
	State                string     `json:"state"`
	PublishedAt          *time.Time `json:"published_at,omitempty"`
	CompletedAt          *time.Time `json:"completed_at,omitempty"`
	LastModifiedAt       time.Time  `json:"last_modified_at"`
}

// SurveyResponse is the detail-view shape, questions included.
type SurveyResponse struct {
	SurveySummaryResponse
	Questions []questionResponse.QuestionResponse `json:"questions"`
}
