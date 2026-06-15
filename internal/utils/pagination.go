package utils

import (
	"strconv"

	"twitter-clone-go/internal/constants"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

func ParsePagination(c *gin.Context) Pagination {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(constants.DefaultPageLimit)))
	if err != nil || limit <= 0 || limit > constants.MaxPageLimit {
		limit = constants.DefaultPageLimit
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}
	return Pagination{Limit: limit, Offset: offset}
}

type PaginatedResponse[T any] struct {
	Items  []T   `json:"items"`
	Total  int64 `json:"total"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
}

func NewPaginatedResponse[T any](items []T, total int64, p Pagination) PaginatedResponse[T] {
	return PaginatedResponse[T]{
		Items:  items,
		Total:  total,
		Limit:  p.Limit,
		Offset: p.Offset,
	}
}
