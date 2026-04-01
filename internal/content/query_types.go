package content

import "GoHeadless/internal/domain"

// ListRecordsResult is the paginated response for GET /content/{slug}.
type ListRecordsResult struct {
	Data  []domain.Record `json:"data"`
	Total int64           `json:"total"`
	Page  int             `json:"page"`
	Limit int             `json:"limit"`
}

// ParsedFilter is one URL filter clause: filter[field][op]=value or filter[field]=value (eq).
type ParsedFilter struct {
	Field string
	Op    string // eq, gt, gte, lt, lte, ne, in, nin, contains
	Value string
}

// ParsedQuery is the output of QueryParser (URL → structured query).
type ParsedQuery struct {
	Search    string
	SortField string
	SortDesc  bool
	Page      int
	Limit     int
	Filters   []ParsedFilter
}
