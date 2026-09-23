package region

import (
	"survey-backend/internal/region/response"

	"gorm.io/gorm"
)

type Region struct {
	gorm.Model
	Name      string `json:"name"`
	Code      string `json:"code"`
	Active    bool   `gorm:"default:true" json:"active"`
	CountryID uint   `gorm:"not null" json:"country_id"`
}

// RegionSeed describes a region by the CODE of the country it belongs to.
//
// It deliberately does not carry a country id. Ids are assigned by the database
// at insert time and are not stable across engines: TiDB allocates
// auto-increment values in batches, which left a gap and shifted every country
// by one, silently attaching each region set to the wrong country.
type RegionSeed struct {
	Name        string
	Code        string
	CountryCode string
}

func GetAllRegions() []RegionSeed {
	return []RegionSeed{
		// Myanmar Regions (CountryCode: "MMR")
		{Name: "Yangon Region", Code: "YGN", CountryCode: "MMR"},
		{Name: "Mandalay Region", Code: "MDY", CountryCode: "MMR"},
		{Name: "Naypyidaw Union Territory", Code: "NYT", CountryCode: "MMR"},
		{Name: "Sagaing Region", Code: "SGG", CountryCode: "MMR"},
		{Name: "Bago Region", Code: "BGO", CountryCode: "MMR"},
		{Name: "Magway Region", Code: "MGW", CountryCode: "MMR"},
		{Name: "Ayeyarwady Region", Code: "AYY", CountryCode: "MMR"},
		{Name: "Mon State", Code: "MON", CountryCode: "MMR"},
		{Name: "Kayin State", Code: "KYN", CountryCode: "MMR"},
		{Name: "Kachin State", Code: "KCH", CountryCode: "MMR"},
		{Name: "Shan State", Code: "SHN", CountryCode: "MMR"},
		{Name: "Rakhine State", Code: "RKH", CountryCode: "MMR"},
		{Name: "Chin State", Code: "CHN", CountryCode: "MMR"},
		{Name: "Kayah State", Code: "KYH", CountryCode: "MMR"},
		{Name: "Tanintharyi Region", Code: "TNT", CountryCode: "MMR"},

		// Thailand Regions (CountryCode: "THA")
		{Name: "Bangkok", Code: "BKK", CountryCode: "THA"},
		{Name: "Chiang Mai", Code: "CMI", CountryCode: "THA"},
		{Name: "Phuket", Code: "HKT", CountryCode: "THA"},
		{Name: "Krabi", Code: "KBV", CountryCode: "THA"},
		{Name: "Pattaya", Code: "PTY", CountryCode: "THA"},
		{Name: "Surat Thani", Code: "URT", CountryCode: "THA"},
		{Name: "Chon Buri", Code: "CBI", CountryCode: "THA"},
		{Name: "Nakhon Ratchasima", Code: "NAK", CountryCode: "THA"},
		{Name: "Udon Thani", Code: "UDN", CountryCode: "THA"},
		{Name: "Hat Yai", Code: "HDY", CountryCode: "THA"},
		{Name: "Khon Kaen", Code: "KKN", CountryCode: "THA"},
		{Name: "Ubon Ratchathani", Code: "UBP", CountryCode: "THA"},
		{Name: "Nakhon Si Thammarat", Code: "NST", CountryCode: "THA"},
		{Name: "Rayong", Code: "RAY", CountryCode: "THA"},
		{Name: "Prachuap Khiri Khan", Code: "PKN", CountryCode: "THA"},
		{Name: "Samut Prakan", Code: "SPK", CountryCode: "THA"},
		{Name: "Nonthaburi", Code: "NON", CountryCode: "THA"},
		{Name: "Pathum Thani", Code: "PTE", CountryCode: "THA"},
		{Name: "Nakhon Pathom", Code: "NPT", CountryCode: "THA"},
		{Name: "Songkhla", Code: "SGK", CountryCode: "THA"},

		//India (CountryCode: "IND")
		{Name: "Andhra Pradesh", Code: "AP", CountryCode: "IND"},
		{Name: "Arunachal Pradesh", Code: "AR", CountryCode: "IND"},
		{Name: "Assam", Code: "AS", CountryCode: "IND"},
		{Name: "Bihar", Code: "BR", CountryCode: "IND"},
		{Name: "Chhattisgarh", Code: "CG", CountryCode: "IND"},
		{Name: "Goa", Code: "GA", CountryCode: "IND"},
		{Name: "Gujarat", Code: "GJ", CountryCode: "IND"},
		{Name: "Haryana", Code: "HR", CountryCode: "IND"},
		{Name: "Himachal Pradesh", Code: "HP", CountryCode: "IND"},
		{Name: "Jharkhand", Code: "JH", CountryCode: "IND"},
		{Name: "Karnataka", Code: "KA", CountryCode: "IND"},
		{Name: "Kerala", Code: "KL", CountryCode: "IND"},
		{Name: "Madhya Pradesh", Code: "MP", CountryCode: "IND"},
		{Name: "Maharashtra", Code: "MH", CountryCode: "IND"},
		{Name: "Manipur", Code: "MN", CountryCode: "IND"},
		{Name: "Meghalaya", Code: "ML", CountryCode: "IND"},
		{Name: "Mizoram", Code: "MZ", CountryCode: "IND"},
		{Name: "Nagaland", Code: "NL", CountryCode: "IND"},
		{Name: "Odisha", Code: "OD", CountryCode: "IND"},
		{Name: "Punjab", Code: "PB", CountryCode: "IND"},
		{Name: "Rajasthan", Code: "RJ", CountryCode: "IND"},
		{Name: "Sikkim", Code: "SK", CountryCode: "IND"},
		{Name: "Tamil Nadu", Code: "TN", CountryCode: "IND"},
		{Name: "Telangana", Code: "TG", CountryCode: "IND"},
		{Name: "Tripura", Code: "TR", CountryCode: "IND"},
		{Name: "Uttar Pradesh", Code: "UP", CountryCode: "IND"},
		{Name: "Uttarakhand", Code: "UK", CountryCode: "IND"},
		{Name: "West Bengal", Code: "WB", CountryCode: "IND"},
		{Name: "Andaman and Nicobar Islands", Code: "AN", CountryCode: "IND"},
		{Name: "Chandigarh", Code: "CH", CountryCode: "IND"},
		{Name: "Dadra and Nagar Haveli and Daman and Diu", Code: "DH", CountryCode: "IND"},
		{Name: "Delhi (National Capital Territory)", Code: "DL", CountryCode: "IND"},
		{Name: "Jammu and Kashmir", Code: "JK", CountryCode: "IND"},
		{Name: "Ladakh", Code: "LA", CountryCode: "IND"},
		{Name: "Lakshadweep", Code: "LD", CountryCode: "IND"},
		{Name: "Puducherry", Code: "PY", CountryCode: "IND"},
	}
}

func (r Region) ToResponse() response.RegionResponse {
	return response.RegionResponse{
		ID:   r.ID,
		Name: r.Name,
		Code: r.Code,
	}
}
