package models

import (
	"strings"

	"github.com/arsenydubrovin/ad-submission/internal/validator"
)

type Filters struct {
	Page     int
	PageSize int
	Sort     string
}

// TODO: move to config if needed
var sortValidValues = []string{"id", "price", "date", "-id", "-price", "-date"}

func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be greater than 0")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 000 000")
	v.Check(f.PageSize > 0, "page_size", "must be greater than 0")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")

	var isValid bool
	for _, val := range sortValidValues {
		if f.Sort == val {
			isValid = true
		}
	}
	v.Check(isValid, "sort", "unknown key")
}

func (f *Filters) sortColumn() string {
	return strings.TrimPrefix(f.Sort, "-")
}

func (f *Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

func (f *Filters) limit() int {
	return f.PageSize
}

func (f *Filters) offset() int {
	return f.PageSize * (f.Page - 1)
}
