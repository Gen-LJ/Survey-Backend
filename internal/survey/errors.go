package survey

import (
	"errors"
	"net/http"
)

// Sentinel errors let handlers map a failure to a status code without the
// service layer knowing anything about HTTP.
var (
	ErrNotFound            = errors.New("survey not found")
	ErrQuestionNotFound    = errors.New("question not found")
	ErrForbidden           = errors.New("you do not have access to this survey")
	ErrNotEditable         = errors.New("this survey can no longer be edited")
	ErrQuestionsLocked     = errors.New("questions cannot change once the survey has answers")
	ErrNoQuestions         = errors.New("add at least one question before publishing")
	ErrAlreadyPublished    = errors.New("this survey is already published")
	ErrNotPublished        = errors.New("this survey is not published")
	ErrInsufficientPoints  = errors.New("insufficient points")
	ErrSurveyClosed        = errors.New("this survey is no longer accepting answers")
	ErrAlreadyAnswered     = errors.New("you have already answered this survey")
	ErrNotAnswered         = errors.New("you have not answered this survey")
	ErrInvalidAnswer       = errors.New("invalid answer")
	ErrAlreadySaved        = errors.New("this survey is already saved")
	ErrNotSaved            = errors.New("this survey is not in your saved list")
	ErrSavedLimitReached   = errors.New("you can only save up to 5 surveys, answer a saved survey first")
	ErrInvalidState        = errors.New("invalid survey state")
	ErrInvalidQuestionType = errors.New("invalid question type")

	ErrOptionsRequired       = errors.New("a multiple choice question needs at least two options")
	ErrOptionsNotAllowed     = errors.New("this question type does not take options")
	ErrMultiAnswerNotAllowed = errors.New("only multiple choice questions can allow multiple answers")
	ErrUnknownCategory       = errors.New("unknown category")
	ErrRegionNotTargetable   = errors.New("surveys cannot target this country and region at the moment")
)

// StatusFor maps a service error to the HTTP status a handler should return.
func StatusFor(err error) int {
	switch {
	case errors.Is(err, ErrNotFound), errors.Is(err, ErrQuestionNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, ErrAlreadyAnswered),
		errors.Is(err, ErrAlreadySaved),
		errors.Is(err, ErrAlreadyPublished),
		errors.Is(err, ErrSavedLimitReached),
		errors.Is(err, ErrSurveyClosed),
		errors.Is(err, ErrInsufficientPoints),
		errors.Is(err, ErrNotEditable),
		errors.Is(err, ErrQuestionsLocked),
		errors.Is(err, ErrNotPublished):
		return http.StatusConflict
	case errors.Is(err, ErrNotSaved), errors.Is(err, ErrNotAnswered):
		return http.StatusNotFound
	case errors.Is(err, ErrInvalidAnswer),
		errors.Is(err, ErrNoQuestions),
		errors.Is(err, ErrInvalidState),
		errors.Is(err, ErrInvalidQuestionType),
		errors.Is(err, ErrOptionsRequired),
		errors.Is(err, ErrOptionsNotAllowed),
		errors.Is(err, ErrMultiAnswerNotAllowed),
		errors.Is(err, ErrUnknownCategory),
		errors.Is(err, ErrRegionNotTargetable):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
