package response

type OptionResponse struct {
	ID       uint   `json:"id"`
	Text     string `json:"text"`
	Position uint   `json:"position"`
}

type QuestionResponse struct {
	ID               uint             `json:"id"`
	SurveyID         uint             `json:"survey_id"`
	Text             string           `json:"text"`
	QuestionType     string           `json:"question_type"`
	AllowMultiAnswer bool             `json:"allow_multi_answer"`
	Position         uint             `json:"position"`
	Options          []OptionResponse `json:"options"`
}
