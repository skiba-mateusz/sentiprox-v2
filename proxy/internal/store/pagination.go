package store

import (
	"net/http"
	"strconv"
)

type PaginatedQuery struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

func (q PaginatedQuery) Parse(r *http.Request) (PaginatedQuery, error) {
	qs := r.URL.Query()

	limit := qs.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return q, err
		}

		q.Limit = l
	}

	offset := qs.Get("offset")
	if offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return q, err
		}

		q.Offset = o
	}

	return q, nil
}
