package pagination

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	DefaultPage  = 1
	DefaultLimit = 10
	MaxLimit     = 100
)

// Query holds a parsed offset/limit pagination request.
type Query struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// Parse reads ?page= and ?limit= from the request, falling back to defaults
// and clamping the limit so a client cannot ask for the whole table.
func Parse(c *gin.Context) Query {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = DefaultPage
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}

	return Query{Page: page, Limit: limit}
}

func (q Query) Offset() int {
	return (q.Page - 1) * q.Limit
}

// Meta is the pagination envelope returned next to a page of items.
type Meta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
}

func NewMeta(q Query, total int64) Meta {
	totalPages := 0
	if q.Limit > 0 {
		totalPages = int((total + int64(q.Limit) - 1) / int64(q.Limit))
	}

	return Meta{
		Page:       q.Page,
		Limit:      q.Limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    q.Page < totalPages,
	}
}

// Page pairs a slice of items with its pagination metadata.
type Page[T any] struct {
	Items []T  `json:"items"`
	Meta  Meta `json:"meta"`
}

func NewPage[T any](items []T, q Query, total int64) Page[T] {
	if items == nil {
		items = []T{}
	}
	return Page[T]{Items: items, Meta: NewMeta(q, total)}
}
