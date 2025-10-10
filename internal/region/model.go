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

func GetAllRegions() []Region {
	return []Region{
		// Myanmar Regions (CountryID: 119)
		{Name: "Yangon Region", Code: "YGN", CountryID: 119, Active: true},
		{Name: "Mandalay Region", Code: "MDY", CountryID: 119, Active: true},
		{Name: "Naypyidaw Union Territory", Code: "NYT", CountryID: 119, Active: true},
		{Name: "Sagaing Region", Code: "SGG", CountryID: 119, Active: true},
		{Name: "Bago Region", Code: "BGO", CountryID: 119, Active: true},
		{Name: "Magway Region", Code: "MGW", CountryID: 119, Active: true},
		{Name: "Ayeyarwady Region", Code: "AYY", CountryID: 119, Active: true},
		{Name: "Mon State", Code: "MON", CountryID: 119, Active: true},
		{Name: "Kayin State", Code: "KYN", CountryID: 119, Active: true},
		{Name: "Kachin State", Code: "KCH", CountryID: 119, Active: true},
		{Name: "Shan State", Code: "SHN", CountryID: 119, Active: true},
		{Name: "Rakhine State", Code: "RKH", CountryID: 119, Active: true},
		{Name: "Chin State", Code: "CHN", CountryID: 119, Active: true},
		{Name: "Kayah State", Code: "KYH", CountryID: 119, Active: true},
		{Name: "Tanintharyi Region", Code: "TNT", CountryID: 119, Active: true},

		// Thailand Regions (CountryID: 173)
		{Name: "Bangkok", Code: "BKK", CountryID: 173, Active: true},
		{Name: "Chiang Mai", Code: "CMI", CountryID: 173, Active: true},
		{Name: "Phuket", Code: "HKT", CountryID: 173, Active: true},
		{Name: "Krabi", Code: "KBV", CountryID: 173, Active: true},
		{Name: "Pattaya", Code: "PTY", CountryID: 173, Active: true},
		{Name: "Surat Thani", Code: "URT", CountryID: 173, Active: true},
		{Name: "Chon Buri", Code: "CBI", CountryID: 173, Active: true},
		{Name: "Nakhon Ratchasima", Code: "NAK", CountryID: 173, Active: true},
		{Name: "Udon Thani", Code: "UDN", CountryID: 173, Active: true},
		{Name: "Hat Yai", Code: "HDY", CountryID: 173, Active: true},
		{Name: "Khon Kaen", Code: "KKN", CountryID: 173, Active: true},
		{Name: "Ubon Ratchathani", Code: "UBP", CountryID: 173, Active: true},
		{Name: "Nakhon Si Thammarat", Code: "NST", CountryID: 173, Active: true},
		{Name: "Rayong", Code: "RAY", CountryID: 173, Active: true},
		{Name: "Prachuap Khiri Khan", Code: "PKN", CountryID: 173, Active: true},
		{Name: "Samut Prakan", Code: "SPK", CountryID: 173, Active: true},
		{Name: "Nonthaburi", Code: "NON", CountryID: 173, Active: true},
		{Name: "Pathum Thani", Code: "PTE", CountryID: 173, Active: true},
		{Name: "Nakhon Pathom", Code: "NPT", CountryID: 173, Active: true},
		{Name: "Songkhla", Code: "SGK", CountryID: 173, Active: true},

		//India (CountryID: 75)
		{Name: "Andhra Pradesh", Code: "AP", CountryID: 75, Active: true},
		{Name: "Arunachal Pradesh", Code: "AR", CountryID: 75, Active: true},
		{Name: "Assam", Code: "AS", CountryID: 75, Active: true},
		{Name: "Bihar", Code: "BR", CountryID: 75, Active: true},
		{Name: "Chhattisgarh", Code: "CG", CountryID: 75, Active: true},
		{Name: "Goa", Code: "GA", CountryID: 75, Active: true},
		{Name: "Gujarat", Code: "GJ", CountryID: 75, Active: true},
		{Name: "Haryana", Code: "HR", CountryID: 75, Active: true},
		{Name: "Himachal Pradesh", Code: "HP", CountryID: 75, Active: true},
		{Name: "Jharkhand", Code: "JH", CountryID: 75, Active: true},
		{Name: "Karnataka", Code: "KA", CountryID: 75, Active: true},
		{Name: "Kerala", Code: "KL", CountryID: 75, Active: true},
		{Name: "Madhya Pradesh", Code: "MP", CountryID: 75, Active: true},
		{Name: "Maharashtra", Code: "MH", CountryID: 75, Active: true},
		{Name: "Manipur", Code: "MN", CountryID: 75, Active: true},
		{Name: "Meghalaya", Code: "ML", CountryID: 75, Active: true},
		{Name: "Mizoram", Code: "MZ", CountryID: 75, Active: true},
		{Name: "Nagaland", Code: "NL", CountryID: 75, Active: true},
		{Name: "Odisha", Code: "OD", CountryID: 75, Active: true},
		{Name: "Punjab", Code: "PB", CountryID: 75, Active: true},
		{Name: "Rajasthan", Code: "RJ", CountryID: 75, Active: true},
		{Name: "Sikkim", Code: "SK", CountryID: 75, Active: true},
		{Name: "Tamil Nadu", Code: "TN", CountryID: 75, Active: true},
		{Name: "Telangana", Code: "TG", CountryID: 75, Active: true},
		{Name: "Tripura", Code: "TR", CountryID: 75, Active: true},
		{Name: "Uttar Pradesh", Code: "UP", CountryID: 75, Active: true},
		{Name: "Uttarakhand", Code: "UK", CountryID: 75, Active: true},
		{Name: "West Bengal", Code: "WB", CountryID: 75, Active: true},
		{Name: "Andaman and Nicobar Islands", Code: "AN", CountryID: 75, Active: true},
		{Name: "Chandigarh", Code: "CH", CountryID: 75, Active: true},
		{Name: "Dadra and Nagar Haveli and Daman and Diu", Code: "DH", CountryID: 75, Active: true},
		{Name: "Delhi (National Capital Territory)", Code: "DL", CountryID: 75, Active: true},
		{Name: "Jammu and Kashmir", Code: "JK", CountryID: 75, Active: true},
		{Name: "Ladakh", Code: "LA", CountryID: 75, Active: true},
		{Name: "Lakshadweep", Code: "LD", CountryID: 75, Active: true},
		{Name: "Puducherry", Code: "PY", CountryID: 75, Active: true},
	}
}

func (r Region) ToResponse() response.RegionResponse {
	return response.RegionResponse{
		ID:   r.ID,
		Name: r.Name,
		Code: r.Code,
	}
}
