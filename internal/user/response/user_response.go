package response

type UserResponse struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Role          string `json:"role"`
	Points        uint   `json:"points"`
	PendingPoints uint   `json:"pending_points"`
	CountryID     uint   `json:"country_id"`
	RegionID      uint   `json:"region_id"`
}
