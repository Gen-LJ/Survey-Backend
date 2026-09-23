package response

import "time"

// AnsweredSurveyResponse is a survey summary plus when the caller answered it.
type AnsweredSurveyResponse struct {
	SurveySummaryResponse
	AnsweredAt time.Time `json:"answered_at"`
}

type InterviewerHomeResponse struct {
	Points          uint                    `json:"points"`
	DraftCount      int64                   `json:"draft_count"`
	PublishedCount  int64                   `json:"published_count"`
	PausedCount     int64                   `json:"paused_count"`
	CompletedCount  int64                   `json:"completed_count"`
	RecentPublished []SurveySummaryResponse `json:"recent_published"`
	RecentDrafts    []SurveySummaryResponse `json:"recent_drafts"`
}

type RespondentHomeResponse struct {
	Points         uint                    `json:"points"`
	AvailableCount int64                   `json:"available_count"`
	AnsweredCount  int64                   `json:"answered_count"`
	SavedCount     int64                   `json:"saved_count"`
	RecentSurveys  []SurveySummaryResponse `json:"recent_surveys"`
}
