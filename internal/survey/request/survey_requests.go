package request

type EditSurveyInfoRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Minutes     uint   `json:"minutes" binding:"required,min=1"`
	CategoryId  uint   `json:"category_id" binding:"required"`
}

type PublishSurveyRequest struct {
	ExpectedAnswerCounts uint `json:"expected_answer_counts" binding:"required,min=1"`
	PointsPerAnswer      uint `json:"points_per_answer" binding:"required,min=1"`
}

type PauseSurveyRequest struct {
	Paused *bool `json:"paused" binding:"required"`
}

type AddQuestionRequest struct {
	Text             string   `json:"text" binding:"required"`
	QuestionType     string   `json:"question_type" binding:"required"`
	AllowMultiAnswer bool     `json:"allow_multi_answer"`
	Options          []string `json:"options"`
}

type EditQuestionRequest struct {
	Text string `json:"text" binding:"required"`
	// Options is optional: a nil list leaves the existing choices untouched,
	// matching the editQuestion contract the Firestore data source exposed.
	Options []string `json:"options"`
}

type UserAnswerRequest struct {
	QuestionID uint     `json:"question_id" binding:"required"`
	Responses  []string `json:"responses"`
}

type SubmitAnswerRequest struct {
	Answers []UserAnswerRequest `json:"answers" binding:"required,min=1,dive"`
}
