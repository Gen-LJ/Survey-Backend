package response

type LoginResponse struct {
	Token        string       `json:"token"`
	UserResponse UserResponse `json:"user"`
}
