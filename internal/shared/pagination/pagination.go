package pagination

import "net/http"

const (
	DefaultPage    = 1
	DefaultPerPage = 20
	MaxPerPage     = 100
)

type Params struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

type Page[T any] struct {
	Items      []T    `json:"items"`
	Pagination Params `json:"pagination"`
	Total      int64  `json:"total"`
}

func NewParams(page, perPage int) Params {
	if page < 1 {
		page = DefaultPage
	}
	if perPage < 1 {
		perPage = DefaultPerPage
	}
	if perPage > MaxPerPage {
		perPage = MaxPerPage
	}
	return Params{Page: page, PerPage: perPage}
}

func FromRequest(r *http.Request) Params {
	query := r.URL.Query()
	return NewParams(parsePositiveInt(query.Get("page")), parsePositiveInt(query.Get("per_page")))
}

func (p Params) Offset() int {
	return (p.Page - 1) * p.PerPage
}

func (p Params) Limit() int {
	return p.PerPage
}

func parsePositiveInt(value string) int {
	var result int
	for _, r := range value {
		if r < '0' || r > '9' {
			return 0
		}
		result = result*10 + int(r-'0')
	}
	return result
}
