package response

// AdminCountryResponse carries the active flag, which the public
// CountryResponse deliberately does not: callers of that one only ever see
// countries that are already active.
type AdminCountryResponse struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Code   string `json:"code"`
	Active bool   `json:"active"`
}
