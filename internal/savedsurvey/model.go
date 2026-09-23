package savedsurvey

import "time"

// MaxSavedSurveys caps a respondent's reading list, the same limit the
// Firestore implementation enforced before writing the savedSurveys array.
const MaxSavedSurveys = 5

// SavedSurvey is a pure link row. It deliberately omits gorm.Model's DeletedAt:
// with a unique index on (user_id, survey_id), a soft-deleted row would keep
// squatting the index and a respondent could never re-save a survey they had
// removed from their list.
type SavedSurvey struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UserID    uint      `gorm:"not null;uniqueIndex:idx_saved_user_survey" json:"user_id"`
	SurveyID  uint      `gorm:"not null;uniqueIndex:idx_saved_user_survey" json:"survey_id"`
}
