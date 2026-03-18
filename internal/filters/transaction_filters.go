package filters

import "time"

type TransactionFilters struct {
	Type       string
	CategoryID *uint
	From       *time.Time
	To         *time.Time
}
