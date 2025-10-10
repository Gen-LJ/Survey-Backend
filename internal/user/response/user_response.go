package response

type UserResponse struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	Role          string `json:"role"`
	Points        uint   `json:"points"`
	PendingPoints uint   `json:"pending_points"`
}
