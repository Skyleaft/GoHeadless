package content

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

const (
	DefaultPage  = 1
	DefaultLimit = 10
	MaxLimit     = 100
)

var (
	filterEqRE  = regexp.MustCompile(`^filter\[([^\]]+)\]$`)
	filterOpRE  = regexp.MustCompile(`^filter\[([^\]]+)\]\[([^\]]+)\]$`)
	allowedSort = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_.]*$`)
)

// QueryParser converts raw URL query strings into ParsedQuery (single responsibility: parsing only).
type QueryParser struct{}

func NewQueryParser() *QueryParser {
	return &QueryParser{}
}

// Parse decodes the raw query string (without leading '?').
func (p *QueryParser) Parse(rawQuery string) (ParsedQuery, error) {
	out := ParsedQuery{
		Page:  DefaultPage,
		Limit: DefaultLimit,
	}

	if rawQuery == "" {
		return out, nil
	}

	values, err := url.ParseQuery(rawQuery)
	if err != nil {
		return out, fmt.Errorf("invalid query string: %w", err)
	}

	for key, vals := range values {
		if len(vals) == 0 {
			continue
		}
		val := vals[len(vals)-1]

		switch key {
		case "search":
			out.Search = strings.TrimSpace(val)
			continue
		case "sort":
			s := strings.TrimSpace(val)
			if s != "" {
				if strings.HasPrefix(s, "-") {
					out.SortDesc = true
					out.SortField = strings.TrimPrefix(s, "-")
				} else {
					out.SortField = s
				}
				if out.SortField != "" && !allowedSort.MatchString(out.SortField) {
					return out, fmt.Errorf("invalid sort field")
				}
			}
			continue
		case "page":
			n, err := strconv.Atoi(val)
			if err != nil || n < 1 {
				return out, fmt.Errorf("invalid page")
			}
			out.Page = n
			continue
		case "limit":
			n, err := strconv.Atoi(val)
			if err != nil || n < 1 {
				return out, fmt.Errorf("invalid limit")
			}
			if n > MaxLimit {
				n = MaxLimit
			}
			out.Limit = n
			continue
		}

		if m := filterOpRE.FindStringSubmatch(key); len(m) == 3 {
			op := strings.ToLower(m[2])
			if !isAllowedFilterOp(op) {
				return out, fmt.Errorf("unsupported filter operator: %s", op)
			}
			out.Filters = append(out.Filters, ParsedFilter{
				Field: m[1],
				Op:    op,
				Value: val,
			})
			continue
		}

		if m := filterEqRE.FindStringSubmatch(key); len(m) == 2 {
			out.Filters = append(out.Filters, ParsedFilter{
				Field: m[1],
				Op:    "eq",
				Value: val,
			})
			continue
		}
	}

	return out, nil
}

func isAllowedFilterOp(op string) bool {
	switch op {
	case "eq", "ne", "gt", "gte", "lt", "lte", "in", "nin", "contains":
		return true
	default:
		return false
	}
}
