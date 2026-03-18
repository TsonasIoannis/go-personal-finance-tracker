package pagination

const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)

type Params struct {
	Page     int
	PageSize int
}

func New(page, pageSize int) Params {
	if page < 1 {
		page = DefaultPage
	}

	if pageSize < 1 {
		pageSize = DefaultPageSize
	}

	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	return Params{
		Page:     page,
		PageSize: pageSize,
	}
}

func (p Params) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func TotalPages(total int64, pageSize int) int {
	if total == 0 {
		return 0
	}

	return int((total + int64(pageSize) - 1) / int64(pageSize))
}
