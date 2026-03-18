package historydto

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)

type PaginationQuery struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

type PaginatedResponse[T any] struct {
	Items    []T `json:"items"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}

func (q PaginationQuery) Offset() int {
	if q.Page <= 1 {
		return 0
	}

	return (q.Page - 1) * q.PageSize
}

func ParsePaginationQuery(values url.Values) (PaginationQuery, error) {
	query := PaginationQuery{
		Page:     DefaultPage,
		PageSize: DefaultPageSize,
	}

	if rawPage := strings.TrimSpace(values.Get("page")); rawPage != "" {
		page, err := strconv.Atoi(rawPage)
		if err != nil || page < 1 {
			return PaginationQuery{}, fmt.Errorf("page must be a positive integer")
		}
		query.Page = page
	}

	if rawPageSize := strings.TrimSpace(values.Get("pageSize")); rawPageSize != "" {
		pageSize, err := strconv.Atoi(rawPageSize)
		if err != nil || pageSize < 1 {
			return PaginationQuery{}, fmt.Errorf("pageSize must be a positive integer")
		}
		if pageSize > MaxPageSize {
			pageSize = MaxPageSize
		}
		query.PageSize = pageSize
	}

	return query, nil
}

func NewPaginatedResponse[T any](query PaginationQuery, items []T, total int) PaginatedResponse[T] {
	if items == nil {
		items = make([]T, 0)
	}

	return PaginatedResponse[T]{
		Items:    items,
		Page:     query.Page,
		PageSize: query.PageSize,
		Total:    total,
	}
}
