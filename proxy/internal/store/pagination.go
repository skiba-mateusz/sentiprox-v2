package store

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Paginated[T any] struct {
	Items T   `json:"items"`
	Total int `json:"total"`
}

type PaginatedQuery struct {
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	Search string `json:"search"`
	Since  string `json:"since"`
	To     string `json:"to"`
	Method string `json:"method"`
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

	search := qs.Get("search")
	if search != "" {
		q.Search = search
	}

	method := qs.Get("method")
	if method != "" {
		q.Method = method
	}

	since := qs.Get("since")
	if since != "" {
		q.Since = parseTime(since)
	}

	to := qs.Get("to")
	if to != "" {
		t, err := time.Parse(time.DateOnly, to)
		fmt.Println(err)
		if err == nil {
			q.To = t.AddDate(0, 0, 1).Format(time.DateTime)
		}
	}

	return q, nil
}

func parseTime(s string) string {
	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return ""
	}

	return t.Format(time.DateTime)
}
