package data

import (
	"slices"
	"strings"

	"flick.io/internal/validator"
)

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafeList []string
}

// Metadata describes the pagination details returned alongside a filtered result set.
type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

func (f Filters) sortColumns() string {
	// Only allow sort fields from the caller-provided allowlist.
	if slices.Contains(f.SortSafeList, f.Sort) {
		return strings.TrimPrefix(f.Sort, "-")
	}

	panic("unsafe sort parameter: " + f.Sort)
}

func (f Filters) sortDirection() string {
	// A leading dash means descending order; otherwise sort ascending.
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

func (f Filters) limit() int {
	// The page size is the SQL LIMIT value.
	return f.PageSize
}

func (f Filters) offet() int {
	// Convert the 1-based page number into a zero-based SQL offset.
	return (f.Page - 1) * f.PageSize
}

func calculateMetadata(totalRecords, page, pageSize int) Metadata {
	// A zero-count result stays empty so callers can omit pagination details.
	if totalRecords == 0 {
		return Metadata{}
	}

	// Compute the last page from the total records and requested page size.
	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     (totalRecords + pageSize - 1) / pageSize,
		TotalRecords: totalRecords,
	}
}

func ValidateFilters(v *validator.Validator, f Filters) {
	// Keep page and page-size values in a bounded range before they reach SQL.
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")
	v.Check(validator.PermittedValues(f.Sort, f.SortSafeList...), "sort", "invalid sort value")
}
