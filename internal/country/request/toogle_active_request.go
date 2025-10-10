package request

type ToggleActiveRequest struct {
	CountryID uint  `json:"id" binding:"required"`
	Active    *bool `json:"active" binding:"required"`
}
