package response

type CountryWithRegions struct {
	ID      uint          `json:"id"`
	Name    string        `json:"name"`
	Code    string        `json:"code"`
	Regions []RegionBrief `json:"regions"`
}

type RegionBrief struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}
