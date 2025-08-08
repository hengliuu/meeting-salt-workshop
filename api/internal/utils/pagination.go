package utils

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaginationParams struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

func GetPaginationParams(c *gin.Context) PaginationParams {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}

	if limit < 1 || limit > 100 {
		limit = 10
	}

	return PaginationParams{
		Page:  page,
		Limit: limit,
	}
}

func GetOffset(page, limit int) int {
	return (page - 1) * limit
}

func CalculateTotalPages(total int64, limit int) int {
	return int(math.Ceil(float64(total) / float64(limit)))
}

func CreatePaginationMeta(page, limit int, total int64) PaginationMeta {
	totalPages := CalculateTotalPages(total, limit)
	
	return PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
}