package controllers

import (
	"strconv"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/apperrors"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/pagination"
	"github.com/gin-gonic/gin"
)

func parsePaginationParams(c *gin.Context) (pagination.Params, error) {
	page := pagination.DefaultPage
	pageSize := pagination.DefaultPageSize

	if rawPage := c.Query("page"); rawPage != "" {
		parsedPage, err := strconv.Atoi(rawPage)
		if err != nil || parsedPage < 1 {
			return pagination.Params{}, apperrors.Validation("invalid_page", "page must be a positive integer")
		}
		page = parsedPage
	}

	if rawPageSize := c.Query("page_size"); rawPageSize != "" {
		parsedPageSize, err := strconv.Atoi(rawPageSize)
		if err != nil || parsedPageSize < 1 || parsedPageSize > pagination.MaxPageSize {
			return pagination.Params{}, apperrors.Validation("invalid_page_size", "page_size must be between 1 and 100")
		}
		pageSize = parsedPageSize
	}

	return pagination.New(page, pageSize), nil
}
