package request

type RegisterRequest struct {
	Name      string `json:"name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	Role      string `json:"role" binding:"required"`
	CountryID uint   `json:"country_id" binding:"required"`
	RegionID  uint   `json:"region_id" binding:"required"`
}
