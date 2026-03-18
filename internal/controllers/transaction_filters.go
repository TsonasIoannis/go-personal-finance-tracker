package controllers

import (
	"strconv"
	"time"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/apperrors"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/filters"
	"github.com/gin-gonic/gin"
)

func parseTransactionFilters(c *gin.Context) (filters.TransactionFilters, error) {
	transactionFilters := filters.TransactionFilters{}

	if transactionType := c.Query("type"); transactionType != "" {
		if transactionType != "income" && transactionType != "expense" {
			return filters.TransactionFilters{}, apperrors.Validation("invalid_transaction_type", "type must be either income or expense")
		}
		transactionFilters.Type = transactionType
	}

	if rawCategoryID := c.Query("category_id"); rawCategoryID != "" {
		categoryID, err := strconv.ParseUint(rawCategoryID, 10, 64)
		if err != nil || categoryID == 0 {
			return filters.TransactionFilters{}, apperrors.Validation("invalid_category_id", "category_id must be a positive integer")
		}

		parsedCategoryID := uint(categoryID)
		transactionFilters.CategoryID = &parsedCategoryID
	}

	if rawFrom := c.Query("from"); rawFrom != "" {
		from, err := parseTransactionFilterTime(rawFrom, false)
		if err != nil {
			return filters.TransactionFilters{}, apperrors.Validation("invalid_from", "from must be RFC3339 or YYYY-MM-DD")
		}
		transactionFilters.From = &from
	}

	if rawTo := c.Query("to"); rawTo != "" {
		to, err := parseTransactionFilterTime(rawTo, true)
		if err != nil {
			return filters.TransactionFilters{}, apperrors.Validation("invalid_to", "to must be RFC3339 or YYYY-MM-DD")
		}
		transactionFilters.To = &to
	}

	if transactionFilters.From != nil && transactionFilters.To != nil && transactionFilters.From.After(*transactionFilters.To) {
		return filters.TransactionFilters{}, apperrors.Validation("invalid_date_range", "from must be before or equal to to")
	}

	return transactionFilters, nil
}

func parseTransactionFilterTime(raw string, inclusiveEndOfDay bool) (time.Time, error) {
	if parsed, err := time.Parse(time.RFC3339, raw); err == nil {
		return parsed, nil
	}

	parsedDate, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return time.Time{}, err
	}

	if inclusiveEndOfDay {
		return parsedDate.Add(24*time.Hour - time.Nanosecond), nil
	}

	return parsedDate, nil
}
