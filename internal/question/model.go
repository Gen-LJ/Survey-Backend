package question

import (
	"survey-backend/internal/question/response"

	"gorm.io/gorm"
)

// Question types, mirroring the QuestionType enum on the client.
const (
	TypeTextInput      = "text_input"
	TypeMultipleChoice = "multiple_choice"
	TypeRating         = "rating"
)

func IsValidType(questionType string) bool {
	switch questionType {
	case TypeTextInput, TypeMultipleChoice, TypeRating:
		return true
	default:
		return false
	}
}

// RequiresOptions reports whether a question of this type is meaningless
// without a choice list.
func RequiresOptions(questionType string) bool {
	return questionType == TypeMultipleChoice
}

type Question struct {
	gorm.Model
	SurveyID         uint     `gorm:"index;not null" json:"survey_id"`
	Text             string   `gorm:"type:text;not null" json:"text"`
	QuestionType     string   `gorm:"size:32;not null" json:"question_type"`
	AllowMultiAnswer bool     `gorm:"default:false" json:"allow_multi_answer"`
	Position         uint     `gorm:"default:0" json:"position"`
	Options          []Option `gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE" json:"options"`
}

type Option struct {
	gorm.Model
	QuestionID uint   `gorm:"index;not null" json:"question_id"`
	Text       string `gorm:"type:text;not null" json:"text"`
	Position   uint   `gorm:"default:0" json:"position"`
}

func (o Option) ToResponse() response.OptionResponse {
	return response.OptionResponse{
		ID:       o.ID,
		Text:     o.Text,
		Position: o.Position,
	}
}

func (q Question) ToResponse() response.QuestionResponse {
	options := make([]response.OptionResponse, len(q.Options))
	for i, opt := range q.Options {
		options[i] = opt.ToResponse()
	}

	return response.QuestionResponse{
		ID:               q.ID,
		SurveyID:         q.SurveyID,
		Text:             q.Text,
		QuestionType:     q.QuestionType,
		AllowMultiAnswer: q.AllowMultiAnswer,
		Position:         q.Position,
		Options:          options,
	}
}

func ToResponseList(questions []Question) []response.QuestionResponse {
	list := make([]response.QuestionResponse, len(questions))
	for i, q := range questions {
		list[i] = q.ToResponse()
	}
	return list
}
