package request

type CreateSurveyRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	CategoryId  uint   `json:"category_id" binding:"required"`
	CountryId   uint   `json:"country_id" binding:"required"`
	RegionId    uint   `json:"region_id" binding:"required"`
}
