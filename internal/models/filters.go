package models

import (
	"math"
	"strings"

	"github.com/arsenydubrovin/ad-submission/internal/validator"
)

type Filters struct {
	Page     int
	PageSize int
	Sort     string
	Fields   []string
}

func ValidatePagination(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be greater than 0")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 000 000")
	v.Check(f.PageSize > 0, "page_size", "must be greater than 0")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")
}

// TODO: move to config if needed
var sortValidValues = []string{"id", "price", "date", "-id", "-price", "-date"}

func ValidateSorting(v *validator.Validator, f Filters) {
	var isValid bool
	for _, val := range sortValidValues {
		if f.Sort == val {
			isValid = true
		}
	}
	v.Check(isValid, "sort", "unknown key")
}

var fieldsValidValues = []string{"description", "allPhotos"}

func ValidateFields(v *validator.Validator, f Filters) {
	for _, field := range f.Fields {
		var isValid bool
		for _, val := range fieldsValidValues {
			if field == val {
				isValid = true
			}
		}
		v.Check(isValid, "fields", "unknown value")
	}
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

func (f *Filters) includeDescription() bool {
	for _, field := range f.Fields {
		if field == "description" {
			return true
		}
	}
	return false
}

func (f *Filters) includeAllPhotos() bool {
	for _, field := range f.Fields {
		if field == "allPhotos" {
			return true
		}
	}
	return false
}

type Info struct {
	CurrentPage  int `json:"currentPage,omitempty"`
	PageSize     int `json:"pageSize,omitempty"`
	FirstPage    int `json:"firstPage,omitempty"`
	LastPage     int `json:"lastPage,omitempty"`
	TotalRecords int `json:"totalRecords,omitempty"`
}

func (i *Info) calculate() {
	if i.TotalRecords == 0 {
		i = nil
	}
	i.FirstPage = 1
	i.LastPage = int(math.Ceil(float64(i.TotalRecords) / float64(i.PageSize)))
}
