package survey

import (
	"gorm.io/gorm"
)

type Survey struct {
	gorm.Model
	Title                string `json:"title"`
	Description          string `json:"description"`
	CategoryID           uint   `json:"category_id"`
	CountryID            uint   `json:"country_id"`
	RegionID             uint   `json:"region_id"`
	CreatorID            uint   `json:"creator_id"`
	ExpectedAnswerCounts uint   `gorm:"default:0" json:"expected_answer_counts"`
	PointsPerAnswer      uint   `gorm:"default:0" json:"points_per_answer"`
	PendingPoints        uint   `gorm:"default:0" json:"pending_points"`
	State                string `gorm:"default:'draft'" json:"state"`
}
