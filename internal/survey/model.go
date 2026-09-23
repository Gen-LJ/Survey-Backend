package survey

import (
	"time"

	questionResponse "survey-backend/internal/question/response"
	"survey-backend/internal/survey/response"

	"gorm.io/gorm"
)

// Survey lifecycle. These replace the Firestore publish/pause boolean pair,
// which could not express "finished" without a second collection.
const (
	StateDraft     = "draft"
	StatePublished = "published"
	StatePaused    = "paused"
	StateCompleted = "completed"
)

func IsValidState(state string) bool {
	switch state {
	case StateDraft, StatePublished, StatePaused, StateCompleted:
		return true
	default:
		return false
	}
}

type Survey struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	CategoryID  uint   `gorm:"index" json:"category_id"`
	CountryID   uint   `gorm:"index" json:"country_id"`
	RegionID    uint   `gorm:"index" json:"region_id"`
	CreatorID   uint   `gorm:"index" json:"creator_id"`
	Minutes     uint   `gorm:"default:0" json:"minutes"`

	ExpectedAnswerCounts uint `gorm:"default:0" json:"expected_answer_counts"`
	PointsPerAnswer      uint `gorm:"default:0" json:"points_per_answer"`
	// PendingPoints is the escrow still owed to respondents: it is funded from
	// the creator's balance at publish time and drawn down one answer at a time.
	PendingPoints uint `gorm:"default:0" json:"pending_points"`
	AnswerCount   uint `gorm:"default:0" json:"answer_count"`

	State       string     `gorm:"size:16;default:'draft';index" json:"state"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// TotalPoints is the full escrow a survey costs to publish.
func (s Survey) TotalPoints() uint {
	return s.PointsPerAnswer * s.ExpectedAnswerCounts
}

// IsOpen reports whether a respondent may still answer this survey.
func (s Survey) IsOpen() bool {
	return s.State == StatePublished && s.AnswerCount < s.ExpectedAnswerCounts
}

// Editable reports whether the survey's questions and info may still change.
// Once points are escrowed, respondents may already have answered the current
// question set, so structural edits are refused.
func (s Survey) Editable() bool {
	return s.State == StateDraft
}

func (s Survey) ToSummary() response.SurveySummaryResponse {
	return response.SurveySummaryResponse{
		ID:                   s.ID,
		Title:                s.Title,
		Description:          s.Description,
		CategoryID:           s.CategoryID,
		CountryID:            s.CountryID,
		RegionID:             s.RegionID,
		CreatorID:            s.CreatorID,
		Minutes:              s.Minutes,
		ExpectedAnswerCounts: s.ExpectedAnswerCounts,
		PointsPerAnswer:      s.PointsPerAnswer,
		PendingPoints:        s.PendingPoints,
		TotalPoints:          s.TotalPoints(),
		AnswerCount:          s.AnswerCount,
		State:                s.State,
		PublishedAt:          s.PublishedAt,
		CompletedAt:          s.CompletedAt,
		LastModifiedAt:       s.UpdatedAt,
	}
}

func (s Survey) ToResponse(questions []questionResponse.QuestionResponse) response.SurveyResponse {
	if questions == nil {
		questions = []questionResponse.QuestionResponse{}
	}

	summary := s.ToSummary()
	summary.QuestionCount = len(questions)

	return response.SurveyResponse{
		SurveySummaryResponse: summary,
		Questions:             questions,
	}
}

func ToSummaryList(surveys []Survey, questionCounts map[uint]int64) []response.SurveySummaryResponse {
	list := make([]response.SurveySummaryResponse, len(surveys))
	for i, s := range surveys {
		summary := s.ToSummary()
		summary.QuestionCount = int(questionCounts[s.ID])
		list[i] = summary
	}
	return list
}
