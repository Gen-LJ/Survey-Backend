package user

import (
	"survey-backend/internal/user/response"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name          string `json:"name"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	Role          string `json:"role"`
	Points        uint   `json:"points"`
	PendingPoints uint   `json:"pending_points"`
	CountryID     uint   `json:"country_id"`
	RegionID      uint   `json:"region_id"`
}

// Converts full User model to safe API response
func (u User) ToResponse() response.UserResponse {
	return response.UserResponse{
		Name:          u.Name,
		Email:         u.Email,
		Role:          u.Role,
		Points:        u.Points,
		PendingPoints: u.PendingPoints,
		CountryID:     u.CountryID,
		RegionID:      u.RegionID,
	}
}
