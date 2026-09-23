package survey

import (
	"errors"
	"time"

	"survey-backend/internal/question"
	questionResponse "survey-backend/internal/question/response"

	"gorm.io/gorm"
)

// loadEditable fetches an owned survey and refuses structural edits once a
// respondent has answered - changing the question set under existing answers
// would silently corrupt the analytics.
func loadEditable(surveyID, creatorID uint) (Survey, error) {
	s, err := loadOwned(surveyID, creatorID)
	if err != nil {
		return Survey{}, err
	}

	if s.State == StateCompleted {
		return Survey{}, ErrNotEditable
	}
	if s.AnswerCount > 0 {
		return Survey{}, ErrQuestionsLocked
	}

	return s, nil
}

// GetQuestions returns a survey's questions for its creator.
func GetQuestions(surveyID, creatorID uint) ([]questionResponse.QuestionResponse, error) {
	if _, err := loadOwned(surveyID, creatorID); err != nil {
		return nil, err
	}

	questions, err := question.FindBySurveyID(surveyID)
	if err != nil {
		return nil, err
	}

	return question.ToResponseList(questions), nil
}

// AddQuestion appends a question to a survey, keeping author order via position.
func AddQuestion(surveyID, creatorID uint, text, questionType string, allowMultiAnswer bool, options []string) (*questionResponse.QuestionResponse, error) {
	if _, err := loadEditable(surveyID, creatorID); err != nil {
		return nil, err
	}

	if !question.IsValidType(questionType) {
		return nil, ErrInvalidQuestionType
	}
	if question.RequiresOptions(questionType) && len(options) < 2 {
		return nil, ErrOptionsRequired
	}
	if !question.RequiresOptions(questionType) && allowMultiAnswer {
		return nil, ErrMultiAnswerNotAllowed
	}

	position, err := question.NextPosition(surveyID)
	if err != nil {
		return nil, err
	}

	optionRows := make([]question.Option, len(options))
	for i, text := range options {
		optionRows[i] = question.Option{Text: text, Position: uint(i)}
	}

	q := question.Question{
		SurveyID:         surveyID,
		Text:             text,
		QuestionType:     questionType,
		AllowMultiAnswer: allowMultiAnswer,
		Position:         position,
		Options:          optionRows,
	}

	if err := question.Create(&q); err != nil {
		return nil, err
	}

	// Refresh so the survey's updated_at moves, which is what the lists sort on.
	if err := touch(surveyID); err != nil {
		return nil, err
	}

	result := q.ToResponse()
	return &result, nil
}

// EditQuestion rewrites a question's text and, when options are supplied,
// replaces its choice list wholesale.
func EditQuestion(surveyID, creatorID, questionID uint, text string, options []string) error {
	if _, err := loadEditable(surveyID, creatorID); err != nil {
		return err
	}

	var q question.Question
	if err := question.FindByID(questionID, &q); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrQuestionNotFound
		}
		return err
	}
	if q.SurveyID != surveyID {
		return ErrQuestionNotFound
	}

	if options != nil {
		if question.RequiresOptions(q.QuestionType) && len(options) < 2 {
			return ErrOptionsRequired
		}
		if !question.RequiresOptions(q.QuestionType) && len(options) > 0 {
			return ErrOptionsNotAllowed
		}
	}

	if err := question.UpdateWithOptions(questionID, text, options); err != nil {
		return err
	}

	return touch(surveyID)
}

// DeleteQuestion removes a question and its options.
func DeleteQuestion(surveyID, creatorID, questionID uint) error {
	if _, err := loadEditable(surveyID, creatorID); err != nil {
		return err
	}

	var q question.Question
	if err := question.FindByID(questionID, &q); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrQuestionNotFound
		}
		return err
	}
	if q.SurveyID != surveyID {
		return ErrQuestionNotFound
	}

	if err := question.Delete(questionID); err != nil {
		return err
	}

	return touch(surveyID)
}

// touch bumps a survey's updated_at, the field the survey lists sort on. It
// stands in for the explicit lastModify write the Firestore code had to make.
func touch(surveyID uint) error {
	return UpdateFields(surveyID, map[string]any{"updated_at": time.Now()})
}
