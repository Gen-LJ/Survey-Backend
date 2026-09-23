package request

type GrantPointsRequest struct {
	Email  string `json:"email" binding:"required,email"`
	Points uint   `json:"points" binding:"required,min=1"`
}
