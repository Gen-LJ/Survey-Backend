package category

import (
	"survey-backend/internal/category/response"

	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name string `gorm:"uniqueIndex;size:100;not null" json:"name"`
}

// Pre-defined category names as constants
const (
	EnergyEnvironment                 = "Energy & Environment"
	GenderInclusion                   = "Gender & Inclusion"
	ClimateChangeSustainability       = "Climate Change & Sustainability"
	HealthWellbeing                   = "Health & Well-being"
	EducationLiteracy                 = "Education & Literacy"
	EmploymentLaborMarket             = "Employment & Labor Market"
	IncomePoverty                     = "Income & Poverty"
	FoodSecurityAgriculture           = "Food Security & Agriculture"
	TechnologyDigitalAccess           = "Technology & Digital Access"
	HousingInfrastructure             = "Housing & Infrastructure"
	TransportMobility                 = "Transport & Mobility"
	GovernancePoliticalParticipation  = "Governance & Political Participation"
	DisasterPreparednessResponse      = "Disaster Preparedness & Response"
	SocialCohesionCommunityEngagement = "Social Cohesion & Community Engagement"
	ConsumerBehaviorMarketTrends      = "Consumer Behavior & Market Trends"
)

// Converts Category model to API response
func (c Category) ToResponse() response.CategoryResponse {
	return response.CategoryResponse{
		ID:   c.ID,
		Name: c.Name,
	}
}

// GetAllCategories returns all predefined categories for seeding
func GetAllCategories() []Category {
	return []Category{
		{Name: EnergyEnvironment},
		{Name: GenderInclusion},
		{Name: ClimateChangeSustainability},
		{Name: HealthWellbeing},
		{Name: EducationLiteracy},
		{Name: EmploymentLaborMarket},
		{Name: IncomePoverty},
		{Name: FoodSecurityAgriculture},
		{Name: TechnologyDigitalAccess},
		{Name: HousingInfrastructure},
		{Name: TransportMobility},
		{Name: GovernancePoliticalParticipation},
		{Name: DisasterPreparednessResponse},
		{Name: SocialCohesionCommunityEngagement},
		{Name: ConsumerBehaviorMarketTrends},
	}
}
