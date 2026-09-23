package survey

import (
	"survey-backend/pkg/database"
	"survey-backend/pkg/pagination"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreateSurveyAtRepo(survey *Survey) error {
	return database.DB.Create(survey).Error
}

func FindByID(id uint, s *Survey) error {
	return database.DB.First(s, id).Error
}

// FindByIDForUpdate loads a survey with a row lock, for use inside a
// transaction that is about to change its counters.
func FindByIDForUpdate(tx *gorm.DB, id uint, s *Survey) error {
	return tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(s, id).Error
}

func UpdateFields(id uint, fields map[string]any) error {
	return database.DB.Model(&Survey{}).Where("id = ?", id).Updates(fields).Error
}

// ListFilter describes a page of surveys to fetch. A nil pointer field means
// "do not filter on this".
type ListFilter struct {
	CreatorID    *uint
	States       []string
	CountryID    *uint
	RegionID     *uint
	ExcludeIDs   []uint
	OnlyOpen     bool // published and not yet at its expected answer count
	OrderByField string
	Descending   bool
}

func (f ListFilter) apply(db *gorm.DB) *gorm.DB {
	if f.CreatorID != nil {
		db = db.Where("creator_id = ?", *f.CreatorID)
	}
	if len(f.States) > 0 {
		db = db.Where("state IN ?", f.States)
	}
	if f.CountryID != nil {
		db = db.Where("country_id = ?", *f.CountryID)
	}
	if f.RegionID != nil {
		db = db.Where("region_id = ?", *f.RegionID)
	}
	if len(f.ExcludeIDs) > 0 {
		db = db.Where("id NOT IN ?", f.ExcludeIDs)
	}
	if f.OnlyOpen {
		db = db.Where("state = ?", StatePublished).
			Where("answer_count < expected_answer_counts")
	}
	return db
}

// List returns one page of surveys plus the total matching the filter, so the
// client can render a real page count instead of Firestore's cursor guesswork.
func List(filter ListFilter, page pagination.Query) ([]Survey, int64, error) {
	var (
		surveys []Survey
		total   int64
	)

	countQuery := filter.apply(database.DB.Model(&Survey{}))
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderField := filter.OrderByField
	if orderField == "" {
		orderField = "updated_at"
	}
	direction := "asc"
	if filter.Descending {
		direction = "desc"
	}

	query := filter.apply(database.DB.Model(&Survey{})).
		Order(orderField + " " + direction).
		Order("id desc").
		Offset(page.Offset()).
		Limit(page.Limit)

	if err := query.Find(&surveys).Error; err != nil {
		return nil, 0, err
	}

	return surveys, total, nil
}

func Count(filter ListFilter) (int64, error) {
	var total int64
	err := filter.apply(database.DB.Model(&Survey{})).Count(&total).Error
	return total, err
}

// CountsByState returns how many surveys a creator has in each state, in one
// query rather than the four count aggregations the Firestore version ran.
func CountsByState(creatorID uint) (map[string]int64, error) {
	var rows []struct {
		State string
		Total int64
	}

	err := database.DB.Model(&Survey{}).
		Select("state, count(*) as total").
		Where("creator_id = ?", creatorID).
		Group("state").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	counts := map[string]int64{
		StateDraft:     0,
		StatePublished: 0,
		StatePaused:    0,
		StateCompleted: 0,
	}
	for _, row := range rows {
		counts[row.State] = row.Total
	}

	return counts, nil
}

// FindByIDs loads a set of surveys and returns them in the order the ids were
// given, skipping any that no longer exist.
func FindByIDs(ids []uint) ([]Survey, error) {
	if len(ids) == 0 {
		return []Survey{}, nil
	}

	var surveys []Survey
	if err := database.DB.Where("id IN ?", ids).Find(&surveys).Error; err != nil {
		return nil, err
	}

	byID := make(map[uint]Survey, len(surveys))
	for _, s := range surveys {
		byID[s.ID] = s
	}

	ordered := make([]Survey, 0, len(ids))
	for _, id := range ids {
		if s, ok := byID[id]; ok {
			ordered = append(ordered, s)
		}
	}

	return ordered, nil
}
